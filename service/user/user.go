package main

import (
	"context"
	"douyin/database"
	pbuser "douyin/grpc/user"
	"douyin/model"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/package/util"
	"douyin/response"
	"errors"
	"strconv"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	pbuser.UnimplementedUserServer
}

func (s *UserService) Register(ctx context.Context, req *pbuser.RegisterRequest) (*pbuser.RegisterResponse, error) {
	// 判断用户名是否合法
	if len(req.Username) <= 0 || len(req.Username) > 32 {
		return nil, errors.New(constant.BadParaRequest)
	}
	if len(req.Password) < 6 || len(req.Password) > 32 {
		return nil, errors.New(constant.SecretFormatError)
	}
	// TODO 复杂度判断 可以使用正则 记得去除常数
	if req.Password == constant.EasySecret {
		return nil, errors.New(constant.SecretFormatEasy)
	}
	//先判断用户存不存在 有唯一索引 其实可以不判断
	_, err := database.SelectUserByName(req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error(constant.DatabaseError, zap.Error(err))
		return nil, err
	}
	if err != gorm.ErrRecordNotFound {
		return nil, errors.New(constant.UserDepulicate)
	}
	// 对密码进行加密并存储
	encryptedPassword := util.BcryptHash(req.Password)
	user := &model.User{
		Username:        req.Username,
		Password:        encryptedPassword,
		Avatar:          util.GenerateAvatar(),
		BackgroundImage: util.GenerateImage(),
		Signature:       util.GenerateSignatrue(),
	}
	userID, err := database.CreateUser(user)
	if err != nil {
		zap.L().Error(constant.DatabaseError, zap.Error(err))
		return nil, err
	}
	// 1. 缓存用户的个人信息
	// 2. 缓存关注和粉丝列表  这个刚关注肯定没有
	// 3. 缓存发布视频和喜欢的视频
	// 注册和登录之后是一样的
	go func() {
		// 将用户ID加入到布隆过滤器里  对抗缓存穿透
		cache.UserIDBloomFilter.AddString(strconv.FormatUint(userID, 10))
		err = cache.SetUserInfo(user)
		if err != nil {
			zap.L().Sugar().Error(err)
		}
		err = cache.SetFavoriteSet(userID, []uint64{})
		if err != nil {
			zap.L().Sugar().Error(err)
		}
		// 用0值维护 redis key 的存在
		err = cache.SetFollowUserIDSet(userID, []uint64{})
		if err != nil {
			zap.L().Sugar().Error(err)
		}
	}()
	return &pbuser.RegisterResponse{
		StatusCode: response.Success,
		StatusMsg:  response.RegisterSuccess,
		UserId:     userID,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *pbuser.LoginRequest) (*pbuser.LoginResponse, error) {
	// 使用redis 限制用户一定时间的登录次数
	loginKey := constant.LoginCounterPrefix + req.Username
	logintimes, err := cache.UserRedisClient.Get(loginKey).Result()
	var logintimesInt int
	if err != nil {
		// 说明没有这个键 初始化键的登录次数
		cache.UserRedisClient.Set(loginKey, 0, constant.MaxloginInernal)
	} else {
		logintimesInt, _ = strconv.Atoi(logintimes)
		if logintimesInt >= constant.MaxLoginTime {
			return nil, errors.New(constant.FrequentLogin)
		}
	}
	// 无论登录成功还是失败 这里redis记录的数据都+1
	go cache.UserRedisClient.Set(loginKey, logintimesInt+1, constant.MaxloginInernal)
	//先判断用户存不存在
	user, err := database.SelectUserByName(req.Username)
	if err != nil {
		zap.L().Error(constant.DatabaseError, zap.Error(err))
		return nil, errors.New(constant.UserNoExist)
	}
	if user.ID == 0 {
		return nil, errors.New(constant.UserNoExist)
	}
	if !util.BcryptCheck(req.Password, user.Password) {
		return nil, errors.New(constant.SecretError)
	}
	// redis预热 用户要查看个人信息 发布的视频 喜欢的视频
	go func() {
		// 个人的用户信息
		err = cache.SetUserInfo(user)
		if err != nil {
			zap.L().Sugar().Error(err)
		}
		// 喜欢的视频列表
		favoriteIDs, err := database.SelectFavoriteVideoByUserID(user.ID)
		if err != nil {
			zap.L().Error(constant.DatabaseError, zap.Error(err))
		} else {
			err = cache.SetFavoriteSet(user.ID, favoriteIDs)
			if err != nil {
				zap.L().Sugar().Error(err)
			}
		}
		// 关注列表
		followUserIDSet, err := database.SelectFollowingByUserID(user.ID)
		if err != nil {
			zap.L().Error(constant.DatabaseError, zap.Error(err))
			return
		}
		err = cache.SetFollowUserIDSet(user.ID, followUserIDSet)
		if err != nil {
			zap.L().Sugar().Error(err)
		}
	}()
	return &pbuser.LoginResponse{
		StatusCode: response.Success,
		StatusMsg:  response.LoginSucess,
		UserId:     user.ID,
	}, nil
}

func (s *UserService) Info(ctx context.Context, req *pbuser.InfoRequest) (*pbuser.InfoResponse, error) {
	// 使用布隆过滤器判断用户ID是否存在
	if !cache.UserIDBloomFilter.TestString(strconv.FormatUint(req.UserID, 10)) {
		err := errors.New(constant.BloomFilterRejected)
		zap.L().Sugar().Error(err)
		return nil, err
	}
	// 去redis里查询用户信息 这是热点数据 redis缓存确实快了很多
	user, err := cache.GetUserInfo(req.UserID)
	// 缓存未命中再去查数据库
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		user, err = database.SelectUserByID(req.UserID)
		if err != nil {
			zap.L().Error(constant.DatabaseError, zap.Error(err))
			return nil, err
		}
		// 设置缓存
		go func() {
			err = cache.SetUserInfo(user)
			if err != nil {
				zap.L().Error(constant.SetCacheError, zap.Error(err))
			}
		}()
	}
	// 判断是否是关注用户
	var isFollow bool
	// 用户未登录
	if req.LoginUserID == 0 {
		isFollow = false
	} else if req.LoginUserID == req.UserID { // 自己查自己 当然是关注了的
		isFollow = true
	} else {
		isFollow, err = cache.IsFollow(req.LoginUserID, req.UserID)
		// 缓存未命中 查询数据库
		if err != nil {
			zap.L().Warn(constant.CacheMiss, zap.Error(err))
			isFollow, err = database.IsFollowed(req.LoginUserID, req.UserID)
			if err != nil {
				zap.L().Error(constant.DatabaseError, zap.Error(err))
				return nil, err
			}
			go func() {
				// 关注列表
				followUserIDSet, err := database.SelectFollowingByUserID(req.LoginUserID)
				if err != nil {
					zap.L().Error(constant.DatabaseError, zap.Error(err))
					return
				}
				err = cache.SetFollowUserIDSet(user.ID, followUserIDSet)
				if err != nil {
					zap.L().Sugar().Error(err)
				}
			}()
		}
	}
	return &pbuser.InfoResponse{
		StatusCode: response.Success,
		User:       response.UserInfo(user, isFollow),
	}, nil
}

package main

import (
	"context"
	"douyin/config"
	pbrelation "douyin/grpc/relation"
	pbuser "douyin/grpc/user"
	"douyin/model"
	"douyin/package/cache"
	"douyin/package/constant"
	"errors"
	"strconv"

	"go.uber.org/zap"
)

type UserService struct {
	pbuser.UnimplementedUserServer
}

func (s *UserService) Register(ctx context.Context, req *pbuser.RegisterRequest) (*pbuser.RegisterResponse, error) {
	err := isUsernameOK(req.Username)
	if err != nil {
		return nil, err
	}
	err = isPasswordOK(req.Password)
	if err != nil {
		return nil, err
	}
	// 对密码进行加密并存储
	encryptedPassword := bcryptHash(req.Password)
	user := &model.User{
		Username:        req.Username,
		Password:        encryptedPassword,
		Avatar:          generateAvatar(),
		BackgroundImage: generateImage(),
		Signature:       generateSignatrue(),
	}
	userID, err := CreateUser(user)
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
		err = SetUserInfo(user)
		if err != nil {
			zap.L().Sugar().Error(err)
		}
		// FIXME 这里预热应该
		// err = cache.SetFavoriteSet(userID, []uint64{})
		// if err != nil {
		// 	zap.L().Sugar().Error(err)
		// }
		// // 用0值维护 redis key 的存在
		// err = cache.SetFollowUserIDSet(userID, []uint64{})
		// if err != nil {
		// 	zap.L().Sugar().Error(err)
		// }
	}()
	return &pbuser.RegisterResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.RegisterSuccess,
		UserId:     userID,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *pbuser.LoginRequest) (*pbuser.LoginResponse, error) {
	// 使用redis 限制用户一定时间的登录次数
	loginKey := constant.LoginCounterPrefix + req.Username
	logintimes, err := userRedis.Get(loginKey).Result()
	var logintimesInt int
	if err != nil {
		// 说明没有这个键 初始化键的登录次数
		userRedis.Set(loginKey, 0, constant.MaxloginInernal)
	} else {
		logintimesInt, _ = strconv.Atoi(logintimes)
		if logintimesInt >= constant.MaxLoginTime {
			return nil, errors.New(constant.FrequentLogin)
		}
	}
	// 无论登录成功还是失败 这里redis记录的数据都+1
	go userRedis.Set(loginKey, logintimesInt+1, constant.MaxloginInernal)
	//先判断用户存不存在
	user, err := SelectUserByName(req.Username)
	if err != nil {
		zap.L().Error(constant.DatabaseError, zap.Error(err))
		return nil, errors.New(constant.UserNoExist)
	}
	if user.ID == 0 {
		return nil, errors.New(constant.UserNoExist)
	}
	if !bcryptCheck(req.Password, user.Password) {
		return nil, errors.New(constant.SecretError)
	}
	// redis预热 用户要查看个人信息 发布的视频 喜欢的视频
	go func() {
		// 个人的用户信息
		err = SetUserInfo(user)
		if err != nil {
			zap.L().Sugar().Error(err)
		}
		// FIXME 这一部分缓存预热 必须掉接口 跨服务调用
		// // 喜欢的视频列表
		// favoriteIDs, err := SelectFavoriteVideoByUserID(user.ID)
		// if err != nil {
		// 	zap.L().Error(constant.DatabaseError, zap.Error(err))
		// } else {
		// 	err = cache.SetFavoriteSet(user.ID, favoriteIDs)
		// 	if err != nil {
		// 		zap.L().Sugar().Error(err)
		// 	}
		// }
		// 关注列表
		// followUserIDSet, err := SelectFollowingByUserID(user.ID)
		// if err != nil {
		// 	zap.L().Error(constant.DatabaseError, zap.Error(err))
		// 	return
		// }
		// err = cache.SetFollowUserIDSet(user.ID, followUserIDSet)
		// if err != nil {
		// 	zap.L().Sugar().Error(err)
		// }
	}()
	return &pbuser.LoginResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.LoginSuccess,
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
	user, err := GetUserInfo(req.UserID)
	// 缓存未命中再去查数据库
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		user, err = SelectUserByID(req.UserID)
		if err != nil {
			zap.L().Error(constant.DatabaseError, zap.Error(err))
			return nil, err
		}
		// 设置缓存
		go func() {
			err = SetUserInfo(user)
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
		isFollowRes, err := relationClient.IsFollow(ctx, &pbrelation.ListRequest{
			UserID:      req.UserID,
			LoginUserID: req.LoginUserID,
		})
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		isFollow = isFollowRes.Result
	}
	return &pbuser.InfoResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.InfoSuccess,
		User:       userResponse(user, isFollow),
	}, nil
}

func userResponse(user *model.User, isFollowed bool) *pbuser.UserInfo {
	return &pbuser.UserInfo{
		Avatar:          config.System.Qiniu.OssDomain + "/" + user.Avatar,
		BackgroundImage: config.System.Qiniu.OssDomain + "/" + user.BackgroundImage,
		FavoriteCount:   user.FavoriteCount,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		Id:              user.ID,
		IsFollow:        isFollowed,
		Name:            user.Username,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
	}
}

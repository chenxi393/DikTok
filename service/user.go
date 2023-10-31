package service

import (
	"douyin/database"
	"douyin/model"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/package/util"
	"douyin/response"
	"errors"
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	// 密码，最长32个字符
	Password string `query:"password"`
	// 注册用户名，最长32个字符
	Username string `query:"username"`
	// 用户鉴权token
	Token string `query:"token"`
	// 用户id 注意上面token会带一个userID
	UserID uint64 `query:"user_id"`
}

func (service *UserService) RegisterService() (*response.UserRegisterOrLogin, error) {
	// 判断用户名是否合法
	if len(service.Username) <= 0 || len(service.Username) > 32 {
		return nil, errors.New(constant.BadParaRequest)
	}
	// TODO 复杂度判断 可以使用正则 记得去除常数
	if len(service.Password) < 6 || len(service.Password) > 32 {
		return nil, fmt.Errorf("密码长度错误")
	}
	if service.Password == "123456" {
		return nil, fmt.Errorf("密码太简单")
	}
	//先判断用户存不存在
	_, err := database.SelectUserByName(service.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error(constant.DatabaseError, zap.Error(err))
		return nil, err
	}
	if err != gorm.ErrRecordNotFound {
		err := fmt.Errorf("用户名已被注册")
		zap.L().Info(err.Error())
		return nil, err
	}
	// 对密码进行加密并存储
	encryptedPassword := util.BcryptHash(service.Password)
	user := &model.User{
		Username: service.Username,
		Password: encryptedPassword,
	}
	userID, err := database.CreateUser(user)
	if err != nil {
		zap.L().Error(constant.DatabaseError, zap.Error(err))
		return nil, err
	}
	// 签发token
	token, err := util.SignToken(userID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 将用户ID加入到布隆过滤器里  对抗缓存穿透
	cache.UserIDBloomFilter.AddString(strconv.FormatUint(userID, 10))
	// 1. 缓存用户的个人信息
	// 2. 缓存关注和粉丝列表  这个刚关注肯定没有
	// 3. 缓存发布视频和喜欢的视频
	// 注册和登录之后是一样的
	go func() {
		err = cache.SetUserInfo(user)
		if err != nil {
			zap.L().Sugar().Error(err)
		}
		err = cache.SetFavoriteSet(userID, []uint64{0})
		if err != nil {
			zap.L().Sugar().Error(err)
		}
		// 用0值维护 redis key 的存在
		err = cache.SetFollowUserIDSet(userID, []uint64{0})
		if err != nil {
			zap.L().Sugar().Error(err)
		}
	}()
	return &response.UserRegisterOrLogin{
		StatusCode: response.Success,
		StatusMsg:  response.RegisterSuccess,
		Token:      &token,
		UserID:     &userID,
	}, nil
}

func (service *UserService) LoginService() (*response.UserRegisterOrLogin, error) {
	// 使用redis 限制用户一定时间的登录次数
	loginKey := constant.LoginCounterPrefix + service.Username
	logintimes, err := cache.UserRedisClient.Get(loginKey).Result()
	var logintimesInt int
	if err != nil {
		// 说明没有这个键 初始化键的登录次数
		cache.UserRedisClient.Set(loginKey, 0, constant.MaxloginInernal)
	} else {
		logintimesInt, _ = strconv.Atoi(logintimes)
		if logintimesInt >= constant.MaxLoginTime {
			err := fmt.Errorf("登录次数过多 5分钟后再试")
			zap.L().Error(err.Error())
			return nil, err
		}
	}
	// 无论登录成功还是失败 这里redis记录的数据都+1
	cache.UserRedisClient.Set(loginKey, logintimesInt+1, constant.MaxloginInernal)
	//先判断用户存不存在
	user, err := database.SelectUserByName(service.Username)
	if err != nil {
		zap.L().Error(constant.DatabaseError, zap.Error(err))
		return nil, err
	}
	if user.ID == 0 {
		err := fmt.Errorf("用户不存在")
		zap.L().Error(err.Error())
		return nil, err
	}
	if !util.BcryptCheck(service.Password, user.Password) {
		err := fmt.Errorf("用户密码错误")
		zap.L().Error(err.Error())
		return nil, err
	}
	// 签发token
	token, err := util.SignToken(user.ID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
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
			if len(favoriteIDs) == 0 {
				// 缓存空值
				favoriteIDs = append(favoriteIDs, 0)
			}
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
		// 0 关注用户加一个数来维持 redis key 的存在
		if len(followUserIDSet) == 0 {
			followUserIDSet = append(followUserIDSet, 0)
		}
		err = cache.SetFollowUserIDSet(user.ID, followUserIDSet)
		if err != nil {
			zap.L().Sugar().Error(err)
		}
	}()
	return &response.UserRegisterOrLogin{
		StatusCode: response.Success,
		StatusMsg:  response.LoginSucess,
		Token:      &token,
		UserID:     &user.ID,
	}, nil
}

func (service *UserService) InfoService(loginUserID uint64) (*response.InfoResponse, error) {
	// 使用布隆过滤器判断用户ID是否存在
	if !cache.UserIDBloomFilter.TestString(strconv.FormatUint(loginUserID, 10)) {
		err := fmt.Errorf("布隆过滤器拦截 用户ID不存在")
		zap.L().Sugar().Error(err)
		return nil, err
	}
	// 去redis里查询用户信息 这是热点数据 redis缓存确实快了很多
	user, err := cache.GetUserInfo(service.UserID)
	// 缓存未命中再去查数据库
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		user, err = database.SelectUserByID(service.UserID)
		if err != nil {
			zap.L().Error(constant.DatabaseError, zap.Error(err))
			return nil, err
		}
		// 设置缓存
		go func() {
			err = cache.SetUserInfo(user)
			if err != nil {
				zap.L().Error("设置缓存失败", zap.Error(err))
			}
		}()
	}
	// 判断是否是关注用户
	var isFollow bool
	// 自己查自己 当然是关注了的
	if loginUserID == service.UserID {
		isFollow = true
	} else {
		isFollow, err = cache.IsFollow(loginUserID, service.UserID)
		// 缓存未命中 查询数据库
		if err != nil {
			zap.L().Warn(constant.CacheMiss, zap.Error(err))
			isFollow, err = database.IsFollowed(loginUserID, service.UserID)
			if err != nil {
				zap.L().Error(constant.DatabaseError, zap.Error(err))
				return nil, err
			}
			go func() {
				// 关注列表
				followUserIDSet, err := database.SelectFollowingByUserID(loginUserID)
				if err != nil {
					zap.L().Error(constant.DatabaseError, zap.Error(err))
					return
				}
				// 0 关注用户加一个数来维持 redis key 的存在
				if len(followUserIDSet) == 0 {
					followUserIDSet = append(followUserIDSet, 0)
				}
				err = cache.SetFollowUserIDSet(user.ID, followUserIDSet)
				if err != nil {
					zap.L().Sugar().Error(err)
				}
			}()
		}
	}
	return &response.InfoResponse{
		StatusCode: response.Success,
		User:       response.UserInfo(user, isFollow),
	}, nil
}

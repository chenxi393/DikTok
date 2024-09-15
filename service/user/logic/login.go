package logic

import (
	"context"
	"strconv"

	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/util"
	"diktok/service/user/storage"

	"go.uber.org/zap"
)

func Login(ctx context.Context, req *pbuser.LoginRequest) (*pbuser.LoginResponse, error) {
	resp := &pbuser.LoginResponse{}
	// 使用redis 限制用户一定时间的登录次数
	loginKey := constant.LoginCounterPrefix + req.Username
	logintimes, err := storage.UserRedis.Get(loginKey).Result()
	var logintimesInt int
	if err != nil {
		// 说明没有这个键 初始化键的登录次数
		storage.UserRedis.Set(loginKey, 0, constant.MaxloginInernal)
	} else {
		logintimesInt, _ = strconv.Atoi(logintimes)
		if logintimesInt >= constant.MaxLoginTime {
			resp.StatusCode = -1
			resp.StatusMsg = constant.FrequentLogin
			return resp, nil
		}
	}
	// 无论登录成功还是失败 这里redis记录的数据都+1
	go storage.UserRedis.Set(loginKey, logintimesInt+1, constant.MaxloginInernal)
	//先判断用户存不存在
	user, err := storage.SelectUserByName(req.Username)
	if err != nil {
		zap.L().Error(constant.DatabaseError, zap.Error(err))
		resp.StatusCode = -1
		resp.StatusMsg = constant.UserNoExist
		return resp, nil
	}
	if user.ID == 0 {

		resp.StatusCode = -1
		resp.StatusMsg = constant.UserNoExist
		return resp, nil
	}
	if !util.BcryptCheck(req.Password, user.Password) {
		resp.StatusCode = -1
		resp.StatusMsg = constant.SecretError
		return resp, nil
	}
	// redis预热 用户要查看个人信息 发布的视频 喜欢的视频
	go func() {
		// 个人的用户信息
		err = storage.SetUserInfo(user)
		if err != nil {
			zap.L().Sugar().Error(err)
		}
	}()
	resp.StatusCode = constant.Success
	resp.StatusMsg = constant.LoginSuccess
	resp.UserId = user.ID
	return resp, nil
}

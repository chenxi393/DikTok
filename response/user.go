package response

import (
	"douyin/config"
	pbuser "douyin/grpc/user"
	"douyin/model"
)

type UserRegisterOrLogin struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 用户id
	UserID *uint64 `json:"user_id"`
}

type InfoResponse struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg *string `json:"status_msg"`
	// 用户信息
	User *User `json:"user"`
}

// User
type User struct {
	// 用户头像
	Avatar string `json:"avatar"`
	// 用户个人页顶部大图
	BackgroundImage string `json:"background_image"`
	// 喜欢数
	FavoriteCount int64 `json:"favorite_count"`
	// 关注总数
	FollowCount int64 `json:"follow_count"`
	// 粉丝总数
	FollowerCount int64 `json:"follower_count"`
	// 用户id
	ID uint64 `json:"id"`
	// true-已关注，false-未关注
	IsFollow bool `json:"is_follow"`
	// 用户名称
	Username string `json:"name" gorm:"column:username"` //done
	// 个人简介
	Signature string `json:"signature"`
	// 获赞数量
	TotalFavorited int64 `json:"total_favorited"`
	// 作品数
	WorkCount int64 `json:"work_count"`
}

func UserInfo(user *model.User, isFollowed bool) *pbuser.UserInfo {
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

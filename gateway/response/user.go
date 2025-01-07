package response

import (
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
)

type UserRegisterOrLogin struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int32 `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 用户id
	UserID *int64 `json:"user_id"`
	Token  string `json:"token"`
}

type InfoResponse struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
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
	ID int64 `json:"id"`
	// true-已关注，false-未关注
	IsFollow bool `json:"is_follow"`
	// 用户名称
	Name string `json:"name"` //done
	// 个人简介
	Signature string `json:"signature"`
	// 获赞数量
	TotalFavorited int64 `json:"total_favorited"`
	// 作品数
	WorkCount int64 `json:"work_count"`
}

func BuildInfoResponse(u *pbuser.UserInfo) *InfoResponse {
	return &InfoResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.LoadSuccess,
		User:       BuildUser(u),
	}
}

func BuildUser(user *pbuser.UserInfo) *User {
	return &User{
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		FavoriteCount:   user.FavoriteCount,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		ID:              user.Id,
		IsFollow:        user.IsFollow,
		Name:            user.Name,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
	}
}

func BuildUserMap(userList *pbuser.ListResp) map[int64]*User {
	if userList == nil {
		return nil
	}
	mp := make(map[int64]*User, len(userList.User))
	for _, v := range userList.User {
		mp[v.Id] = BuildUser(v)
	}
	return mp
}

func BuildLoginRes(loginRes *pbuser.LoginResponse) *UserRegisterOrLogin {
	return &UserRegisterOrLogin{
		StatusCode: loginRes.StatusCode,
		StatusMsg:  loginRes.StatusMsg,
		UserID:     &loginRes.UserId,
		Token:      loginRes.Token,
	}
}

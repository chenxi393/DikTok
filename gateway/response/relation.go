package response

import (
	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"
)

type RelationListResponse struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 用户信息列表
	UserList []User `json:"user_list"`
}

type FriendResponse struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 用户信息列表
	UserList []*FriendUser `json:"user_list"`
}
type FriendUser struct {
	User
	Message string `json:"message"`
	MsgType int32  `json:"msgType"`
}

func BuildFrindsRes(frinedRes *pbrelation.FriendsResponse) *FriendResponse {
	friendData := frinedRes.UserList
	res := &FriendResponse{
		UserList:   make([]*FriendUser, 0, len(friendData)),
		StatusCode: constant.Success,
		StatusMsg:  constant.FeedSuccess,
	}
	for _, v := range friendData {
		if v != nil {
			res.UserList = append(res.UserList, BuildFriendUser(v))
		}
	}
	return res
}

func BuildFriendUser(item *pbrelation.FriendInfo) *FriendUser {
	return &FriendUser{
		User: User{
			Avatar:          item.Avatar,
			BackgroundImage: item.BackgroundImage,
			FavoriteCount:   item.FavoriteCount,
			FollowCount:     item.FollowCount,
			FollowerCount:   item.FollowerCount,
			ID:              item.Id,
			IsFollow:        item.IsFollow,
			Name:            item.Name,
			Signature:       item.Signature,
			TotalFavorited:  item.TotalFavorited,
			WorkCount:       item.WorkCount},
		Message: item.Message,
		MsgType: item.MsgType,
	}
}

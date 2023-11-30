package response

import (
	"douyin/config"
	"time"
)

type FeedResponse struct {
	// 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	NextTime int64 `json:"next_time"`
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 视频列表
	VideoList []Video `json:"video_list"`
}

type Video struct {
	// 视频作者信息
	Author User `json:"author"`
	// 视频的评论总数
	CommentCount int64 `json:"comment_count"`
	// 视频封面地址
	CoverURL string `json:"cover_url"`
	// 视频的点赞总数
	FavoriteCount int64 `json:"favorite_count"`
	// 视频唯一标识
	ID uint64 `json:"id"`
	// true-已点赞，false-未点赞
	IsFavorite bool `json:"is_favorite"`
	// 视频播放地址
	PlayURL string `json:"play_url"`
	// 视频标题
	Title string `json:"title"`
	// 新增返回视频发布时间
	PublishTime string `json:"publish_time"`
	// 新增视频话题
	Topic string `json:"topic"`
}

type VideoData struct {
	User               `json:"author"`
	VideoID            uint64 `gorm:"column:vid" json:"id"`
	PlayURL            string `gorm:"column:play_url" json:"play_url"`
	CoverURL           string `gorm:"column:cover_url" json:"cover_url"`
	VideoFavoriteCount int64  `gorm:"column:vfavorite_count" json:"favorite_count"`
	CommentCount       int64  `gorm:"column:comment_count" json:"comment_count"`
	Title              string `gorm:"column:title" json:"title"`
	IsFavorite         bool   `json:"is_favorite"`
	PublishTime        time.Time
	Topic              string // 11.3 新增字段
}

type VideoListResponse struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 用户发布的视频列表
	VideoList []Video `json:"video_list"`
}

func VideoDataInfo(data []VideoData) []Video {
	items := make([]Video, 0, len(data))
	for _, item := range data {
		v := Video{
			Author:        *addUserDomain(&item.User),
			ID:            item.VideoID,
			PlayURL:       config.System.Qiniu.OssDomain + "/" + item.PlayURL,
			CoverURL:      config.System.Qiniu.OssDomain + "/" + item.CoverURL,
			FavoriteCount: item.VideoFavoriteCount,
			CommentCount:  item.CommentCount,
			Title:         item.Title,
			PublishTime:   item.PublishTime.Format("2006-01-02 15:04"),
			Topic:         item.Topic,
		}
		items = append(items, v)
	}
	return items
}

func addUserDomain(user *User) *User {
	return &User{
		Avatar:          config.System.Qiniu.OssDomain + "/" + user.Avatar,
		BackgroundImage: config.System.Qiniu.OssDomain + "/" + user.BackgroundImage,
		FavoriteCount:   user.FavoriteCount,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		ID:              user.ID,
		IsFollow:        user.IsFollow,
		Name:            user.Name,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
	}
}

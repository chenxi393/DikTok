package response

import (
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
)

type FeedResponse struct {
	// 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	NextTime int64 `json:"next_time"`
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 视频列表
	VideoList []*Video `json:"video_list"`
}

type Video struct {
	// 视频作者信息
	Author *User `json:"author"`
	// 视频的评论总数
	CommentCount int64 `json:"comment_count"`
	// 视频封面地址
	CoverURL string `json:"cover_url"`
	// 视频的点赞总数
	FavoriteCount int64 `json:"favorite_count"`
	// 视频唯一标识
	ID int64 `json:"id"`
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

type VideoListResponse struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 用户发布的视频列表
	VideoList []*Video `json:"video_list"`
}

type UploadTokenResponse struct {
	UploadToken string `json:"upload_token"`
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	FileName  string `json:"file_name"`
}

func BuildVideoList(videoData []*pbvideo.VideoData) *VideoListResponse {
	res := &VideoListResponse{
		VideoList:  make([]*Video, 0, len(videoData)),
		StatusCode: constant.Success,
		StatusMsg:  constant.FeedSuccess,
	}
	for _, v := range videoData {
		if v != nil {
			res.VideoList = append(res.VideoList, BuildVideo(v))
		}
	}
	return res
}

func BuildFeed(feedData *pbvideo.FeedResponse) *FeedResponse {
	videoData := feedData.VideoList
	res := &FeedResponse{
		VideoList:  make([]*Video, 0, len(videoData)),
		StatusCode: constant.Success,
		StatusMsg:  constant.FeedSuccess,
		NextTime:   feedData.NextTime,
	}
	for _, v := range videoData {
		if v != nil {
			res.VideoList = append(res.VideoList, BuildVideo(v))
		}
	}
	return res
}

// 仅透传字段 video服务已打包好
func BuildVideo(item *pbvideo.VideoData) *Video {
	return &Video{
		Author:        BuildUser(item.GetAuthor()),
		ID:            item.Id,
		PlayURL:       item.PlayUrl,
		CoverURL:      item.CoverUrl,
		IsFavorite:    item.IsFavorite,
		FavoriteCount: item.FavoriteCount,
		CommentCount:  item.CommentCount,
		Title:         item.Title,
		PublishTime:   item.PublishTime,
		Topic:         item.Topic,
	}
}

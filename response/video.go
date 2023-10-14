package response

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
}

func VideoDataInfo(data []VideoData) []Video {
	items := make([]Video, 0, len(data))
	for _, item := range data {
		v := Video{
			Author:        item.User,
			ID:            item.VideoID,
			PlayURL:       item.PlayURL,
			CoverURL:      item.CoverURL,
			FavoriteCount: item.VideoFavoriteCount,
			CommentCount:  item.CommentCount,
			Title:         item.Title,
		}
		items = append(items, v)
	}
	return items
}

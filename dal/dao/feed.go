package dao

import (
	"douyin/dal/model"
	"douyin/response"
	"time"
)

// 这里要获取视频的信息 还要获取视频作者的信息
// 抖音面板有是否关注 且用户很有可能点进去

// 这里下面必须使用scan scan对传入的个数有严格限制
func SelectFeedVideoList(numberVideos int, lastTime *int64) ([]response.Video, error) {
	if lastTime == nil || *lastTime == 0 { //很多时候这种写法会有空指针的问题
		// 所以得产生一个新的变量 杜绝空指针问题 内存回收交给gc
		// 待验证 客户端传入的是毫秒还是秒 文档说是秒
		currentTime := time.Now().UnixMilli()
		lastTime = &currentTime
	}
	res := make([]response.Video, 0, 30)
	// 这里使用外连接 双表联查 可以考虑单次拿出30个video 再构造一个map（包含30个uid）去数据库里查
	// FIX 这里有问题 两表联查有重复字段  需要重新开一个结构体 手动select而不是* 然后对应字段
	// 或者干脆一次只查一个表 然后走批量查询
	err := global_db.Model(&model.User{}).Select("*").Joins(
		"right join video on video.author_id = user.id").Where("video.publish_time <= ? ",
		time.UnixMilli(*lastTime)).Limit(numberVideos).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SelectPublishTimeByVideoID(vid uint64) (time.Time, error) {
	var publishTime time.Time
	err := global_db.Model(&model.Video{}).Select("publish_time").Where("id = ?", vid).First(&publishTime).Error
	return publishTime, err
}

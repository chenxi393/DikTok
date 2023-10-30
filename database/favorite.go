package database

import (
	"douyin/model"
	"douyin/package/cache"

	"gorm.io/gorm"
)

// cnt=1表示 点赞 cnt=-1 表示取消赞
func FavoriteVideo(userID, videoID uint64, cnt int64) error {
	favorite := model.Favorite{
		UserID:  userID,
		VideoID: videoID,
	}
	// 一般输入流程 在是事务里 使用tx而不是db 返回任何错误都会回滚事务
	return global_db.Transaction(func(tx *gorm.DB) error {
		var err error
		if cnt == 1 {
			err = tx.Model(&model.Favorite{}).Create(&favorite).Error
		} else {
			err = tx.Model(&model.Favorite{}).Where("user_id = ? AND video_id = ? ", userID, videoID).Delete(&favorite).Error
		}
		if err != nil {
			return err
		}
		// 视频表增加该视频的点赞
		video := model.Video{ID: videoID} // model里需要有主键
		// 否则favorite.go:33 WHERE conditions required
		err = tx.Model(&model.Video{}).Select("favorite_count, author_id").Where("id = ?", videoID).First(&video).Error
		if err != nil {
			return err
		}
		err = tx.Model(&video).Update("favorite_count", video.FavoriteCount+cnt).Error
		if err != nil {
			return err
		}
		// 增加视频作者被点赞数
		author := model.User{ID: video.AuthorID} // 同理上面 可以看Model函数的说明
		err = tx.Model(&author).Select("total_favorited").First(&author).Error
		if err != nil {
			return err
		}
		err = tx.Model(&author).Update("total_favorited", author.TotalFavorited+cnt).Error
		if err != nil {
			return err
		}
		// 增加用户的点赞数
		user := model.User{ID: userID} // 同理上面
		err = tx.Model(&user).Select("favorite_count").First(&user).Error
		if err != nil {
			return err
		}
		err = tx.Model(&user).Update("favorite_count", user.FavoriteCount+cnt).Error
		if err != nil {
			return err
		}
		return cache.FavoriteAction(userID, author.ID, videoID, cnt)
	})
}

func SelectFavoriteVideoByUserID(userID uint64) ([]uint64, error) {
	res := make([]uint64, 0)
	err := global_db.Model(&model.Favorite{}).Select("video_id").Where("user_id = ?", userID).Find(&res).Error
	return res, err
}

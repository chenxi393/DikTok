package main

import (
	"errors"

	"diktok/package/constant"
	"diktok/storage/database"
	"diktok/storage/database/model"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// cnt=1表示 点赞 cnt=-1 表示取消赞
func FavoriteVideo(userID, videoID int64, cnt int64) error {
	favorite := model.Favorite{
		UserID:  userID,
		VideoID: videoID,
	}
	// 一般输入流程 在是事务里 使用tx而不是db 返回任何错误都会回滚事务
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 先看有没有点赞过
		var isFavorite int64
		err := tx.Model(&model.Favorite{}).Where("user_id = ? AND video_id = ?", userID, videoID).Count(&isFavorite).Error
		if err != nil {
			return err
		}
		if cnt == 1 && isFavorite == 0 {
			err = tx.Model(&model.Favorite{}).Create(&favorite).Error
		} else if cnt == -1 && isFavorite == 1 {
			err = tx.Model(&model.Favorite{}).Where("user_id = ? AND video_id = ? ", userID, videoID).Delete(&favorite).Error
		} else {
			err = errors.New(constant.BadParaRequest)
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
		return FavoriteAction(userID, author.ID, videoID, cnt)
	})
}

func SelectFavoriteVideoByUserID(userID int64) ([]int64, error) {
	res := make([]int64, 0)
	err := database.DB.Model(&model.Favorite{}).Select("video_id").Where("user_id = ?", userID).Order("id desc").Find(&res).Error
	return res, err
}

func GetFavoriteNumByVideoIDFromMaster(videoID int64) (int64, error) {
	var cnt int64
	err := database.DB.Clauses(dbresolver.Write).Model(&model.Favorite{}).Where("video_id = ?", videoID).Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

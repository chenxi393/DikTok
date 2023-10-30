package database

import (
	"douyin/model"
	"douyin/package/cache"

	"gorm.io/gorm"
)

func CommentAdd(com *model.Comment) error {
	err := global_db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.Comment{}).Create(&com).Error
		if err != nil {
			return err
		}
		video := model.Video{ID: com.VideoID}
		err = tx.Model(&video).Select("comment_count").First(&video).Error
		if err != nil {
			return err
		}
		err = tx.Model(&video).Update("comment_count", video.CommentCount+1).Error
		if err != nil {
			return err
		}
		return cache.CommentAdd(com)
	})
	if err != nil {
		return err
	}
	return nil
}

func CommentDelete(commentID *string, videoID, userID uint64) (*model.Comment, error) {
	comment := model.Comment{}
	// delete 不会回写到comment里  Clauses(clause.Returning{}) 这个才会回写
	err := global_db.Transaction(func(tx *gorm.DB) error {
		// 删除要先检查里面有没有啊
		err := tx.Where("id = ? AND video_id = ? AND user_id = ?", commentID, videoID, userID).First(&comment).Error
		if err != nil || comment.ID == 0 {
			return err
		}
		err = tx.Delete(&comment).Error
		if err != nil {
			return err
		}
		video := model.Video{ID: videoID}
		err = tx.Model(&video).Select("comment_count").First(&video).Error
		if err != nil {
			return err
		}
		err = tx.Model(&video).Update("comment_count", video.CommentCount-1).Error
		if err != nil {
			return err
		}
		return cache.CommentDelete(&comment)
	})
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func GetCommentsByVideoID(videoID uint64) ([]*model.Comment, error) {
	videos := make([]*model.Comment, 0)
	err := global_db.Model(&model.Comment{}).Where("video_id = ?", videoID).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

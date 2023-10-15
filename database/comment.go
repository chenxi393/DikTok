package database

import (
	"douyin/model"
	"time"

	"gorm.io/gorm"
)

func CommentAdd(comment *string, videoID, userID uint64) (*model.Comment, error) {
	com := model.Comment{
		VideoID:     videoID,
		UserID:      userID,
		Content:     *comment,
		CreatedTime: time.Now(),
	}
	// 还是需要使用事务 先更新评论表 再更新视频表
	err := global_db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.Comment{}).Create(&com).Error
		if err != nil {
			return err
		}
		video := model.Video{ID: videoID}
		err = tx.Model(&video).Select("comment_count").First(&video).Error
		if err != nil {
			return err
		}
		err = tx.Model(&video).Update("comment_count", video.CommentCount+1).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &com, nil
}

func CommentDelete(commentID *string, videoID, userID uint64) (*model.Comment, error) {
	// 还是需要使用事务 先更新评论表 再更新视频表
	comment := model.Comment{}
	// fix 这里返回的评论没有东西 一点也没有
	err := global_db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.Comment{}).Where("id = ? AND video_id = ? AND user_id = ?",
			commentID, videoID, userID).Delete(&comment).Error
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
		return nil
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

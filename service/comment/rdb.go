package main

import (
	"douyin/model"
	"douyin/package/database"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func CreateComment(com *model.Comment) error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 先查询videoID是否存在
		video := model.Video{ID: com.VideoID}
		err := tx.First(&video).Error
		if err != nil {
			return err
		}
		err = tx.Model(&model.Comment{}).Create(&com).Error
		if err != nil {
			return err
		}
		err = tx.Model(&video).Update("comment_count", video.CommentCount+1).Error
		if err != nil {
			return err
		}
		return CommentAdd(com)
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteComment(commentID, videoID, userID uint64) (*model.Comment, error) {
	comment := model.Comment{}
	// delete 不会回写到comment里  Clauses(clause.Returning{}) 这个才会回写
	err := database.DB.Transaction(func(tx *gorm.DB) error {
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
		return CommentDelete(&comment)
	})
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func GetCommentsByVideoIDRDB(videoID uint64) ([]*model.Comment, error) {
	videos := make([]*model.Comment, 0)
	err := database.DB.Model(&model.Comment{}).Where("video_id = ?", videoID).Order("created_time desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func GetCommentsByVideoIDFromMaster(videoID uint64, offset, count int) ([]*model.Comment, error) {
	videos := make([]*model.Comment, 0)
	err := database.DB.Clauses(dbresolver.Write).Model(&model.Comment{}).Where("video_id = ?", videoID).Order("created_time desc").Offset(int(offset)).Limit(count).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func GetCommentsNumByVideoIDFromMaster(videoID uint64) (int64, error) {
	var cnt int64
	err := database.DB.Clauses(dbresolver.Write).Model(&model.Comment{}).Where("video_id = ?", videoID).Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

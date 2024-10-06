package storage

import (
	"context"
	"diktok/storage/database"
	"diktok/storage/database/model"
	"diktok/storage/database/query"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

// 这个主题表 后续看有没有必要 变成 评论总数表 count计数
// User not found, create it with give conditions and Attrs
// func FirstOrCreateSubjectByObj(ctx context.Context, objType int32, objID int64) (*model.CommentSubject, error) {
// 	so := query.Use(database.DB.Clauses(dbresolver.Read)).CommentSubject
// 	return so.WithContext(ctx).Attrs(field.Attrs(&model.CommentSubject{})).Where(so.ObjType.Eq(objType), so.ObjID.Eq(objID)).FirstOrCreate()
// }

func CreateCommentContent(ctx context.Context, meta *model.CommentMetum, content *model.CommentContent) error {
	q := query.Use(database.DB)
	return q.Transaction(func(tx *query.Query) error {
		if err := tx.CommentContent.WithContext(ctx).Create(content); err != nil {
			return err
		}
		if err := tx.CommentMetum.WithContext(ctx).Create(meta); err != nil {
			return err
		}
		return nil
	})
}

func DeleteCommentContent(ctx context.Context, commentID int64) error {
	q := query.Use(database.DB)
	var commentsID []int64
	return q.Transaction(func(tx *query.Query) error {
		_, err := tx.CommentContent.WithContext(ctx).Where(tx.CommentContent.ID.Eq(commentID)).Delete()
		if err != nil {
			return err
		}
		_, err = tx.CommentMetum.WithContext(ctx).Where(tx.CommentMetum.CommentID.Eq(commentID)).Delete()
		if err != nil {
			return err
		}
		// 子评论也需要都置为删除
		err = tx.CommentMetum.WithContext(ctx).Select(q.CommentMetum.CommentID).Where(tx.CommentMetum.ParentID.Eq(commentID)).Scan(&commentsID)
		if err != nil {
			return err
		}
		if len(commentsID) <= 0 {
			return nil
		}
		_, err = tx.CommentContent.WithContext(ctx).Where(tx.CommentContent.ID.In(commentsID...)).Delete()
		if err != nil {
			return err
		}
		_, err = tx.CommentMetum.WithContext(ctx).Where(tx.CommentMetum.CommentID.In(commentsID...)).Delete()
		if err != nil {
			return err
		}
		// 如果后续有计数表
		// 需要再事务里 计数-1
		return nil
	})
}

func MGetCommentsByCond(ctx context.Context, offset, limit int, conds []gen.Condition, order ...field.Expr) ([]*model.CommentMetum, error) {
	return query.Use(database.DB).CommentMetum.WithContext(ctx).Where(conds...).Order(order...).Offset(offset).Limit(limit).Find()
}

func CountByCond(ctx context.Context, conds []gen.Condition) (int64, error) {
	return query.Use(database.DB).CommentMetum.WithContext(ctx).Where(conds...).Count()
}

func GetContentByIDs(ctx context.Context, commentIDs []int64) ([]*model.CommentContent, error) {
	q := query.Use(database.DB)
	return q.CommentContent.WithContext(ctx).Where(q.CommentContent.ID.In(commentIDs...)).Find()
}

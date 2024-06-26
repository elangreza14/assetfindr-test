package repository

import (
	"context"

	"github.com/elangreza14/assetfindr-test/model"
	"gorm.io/gorm"
)

type (
	PostRepository struct {
		db *gorm.DB
	}
)

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db}
}

func (pr *PostRepository) GetPosts(ctx context.Context) ([]model.Post, error) {
	res := []model.Post{}
	err := pr.db.WithContext(ctx).Model(&model.Post{}).Preload("Tags").Order("id desc").Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (pr *PostRepository) CreatePost(ctx context.Context, req model.Post) error {
	err := pr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		post := model.Post{
			Title:   req.Title,
			Content: req.Content,
		}

		err := tx.WithContext(ctx).Create(&post).Error
		if err != nil {
			return err
		}

		for _, tag := range req.Tags {
			reqTag := model.Tag{}
			err := tx.Where(model.Tag{Label: tag.Label}).FirstOrCreate(&reqTag).Error
			if err != nil {
				return err
			}

			err = tx.Exec(`INSERT INTO post_tags ("post_id","tag_id") VALUES ($1, $2) ON CONFLICT DO NOTHING`, post.ID, reqTag.ID).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (pr *PostRepository) GetPost(ctx context.Context, id int) (*model.Post, error) {
	res := model.Post{}
	err := pr.db.WithContext(ctx).Model(&model.Post{}).Preload("Tags").First(&res, id).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (pr *PostRepository) UpdatePost(ctx context.Context, req model.Post, tagsToBeDeleted ...int) error {
	err := pr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Exec(`DELETE FROM post_tags WHERE post_id=? and tag_id IN ?;`, req.ID, tagsToBeDeleted).Error
		if err != nil {
			return err
		}

		post := model.Post{
			ID:      req.ID,
			Title:   req.Title,
			Content: req.Content,
		}

		err = tx.WithContext(ctx).Updates(&post).Error
		if err != nil {
			return err
		}

		for _, tag := range req.Tags {
			reqTag := model.Tag{}
			err := tx.Where(model.Tag{Label: tag.Label}).FirstOrCreate(&reqTag).Error
			if err != nil {
				return err
			}

			err = tx.Exec(`INSERT INTO post_tags ("post_id","tag_id") VALUES ($1, $2) ON CONFLICT DO NOTHING`, post.ID, reqTag.ID).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (pr *PostRepository) DeletePost(ctx context.Context, req model.Post) error {
	err := pr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(`DELETE FROM post_tags WHERE post_id=$1;`, req.ID).Error
		if err != nil {
			return err
		}

		err = tx.WithContext(ctx).Delete(&req).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

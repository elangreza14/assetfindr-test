package service

//go:generate mockgen -source $GOFILE -destination ../mock/service/mock_$GOFILE -package $GOPACKAGE

import (
	"context"
	"errors"

	"github.com/elangreza14/assetfindr-test/dto"
	"github.com/elangreza14/assetfindr-test/model"
	"gorm.io/gorm"
)

type (
	IPostRepository interface {
		GetPosts(ctx context.Context) ([]model.Post, error)
		CreatePost(ctx context.Context, req model.Post) error
		GetPost(ctx context.Context, id int) (*model.Post, error)
		UpdatePost(ctx context.Context, req model.Post, tagsToBeDeleted ...int) error
		DeletePost(ctx context.Context, req model.Post) error
	}

	PostService struct {
		postRepository IPostRepository
	}
)

func NewPostService(postRepository IPostRepository) *PostService {
	return &PostService{
		postRepository: postRepository,
	}
}

func (ps *PostService) GetPosts(ctx context.Context) ([]dto.GetPostResponse, error) {
	posts, err := ps.postRepository.GetPosts(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]dto.GetPostResponse, len(posts))
	for i, post := range posts {
		tags := make([]string, len(post.Tags))
		for j, tag := range post.Tags {
			tags[j] = tag.Label
		}

		res[i] = dto.GetPostResponse{
			ID:      post.ID,
			Title:   post.Title,
			Content: post.Content,
			Tags:    tags,
		}
	}

	return res, nil
}

func (ps *PostService) CreatePost(ctx context.Context, req dto.CreateOrUpdatePostRequest) error {
	tags := make([]*model.Tag, len(req.Tags))
	for i, tag := range req.Tags {
		tags[i] = &model.Tag{
			Label: tag,
		}
	}

	err := ps.postRepository.CreatePost(ctx, model.Post{
		Title:   req.Title,
		Content: req.Content,
		Tags:    tags,
	})
	if err != nil {
		return err
	}

	return nil
}

func (ps *PostService) GetPost(ctx context.Context, id int) (*dto.GetPostResponse, error) {
	post, err := ps.postRepository.GetPost(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrorNotFound{
				EntityName: "post",
				EntityID:   id,
			}
		}
		return nil, err
	}

	tags := make([]string, len(post.Tags))
	for j, tag := range post.Tags {
		tags[j] = tag.Label
	}

	return &dto.GetPostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Tags:    tags,
	}, nil
}

func (ps *PostService) UpdatePost(ctx context.Context, req dto.CreateOrUpdatePostRequest, id int) error {
	post, err := ps.postRepository.GetPost(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrorNotFound{
				EntityName: "post",
				EntityID:   id,
			}
		}
		return err
	}

	newTagsToBeSave := make([]*model.Tag, len(req.Tags))
	for i, tag := range req.Tags {
		newTagsToBeSave[i] = &model.Tag{
			Label: tag,
		}
	}

	prevTags := make(map[string]*model.Tag)
	for _, tag := range post.Tags {
		prevTags[tag.Label] = tag
	}

	for _, reqTag := range req.Tags {
		if _, ok := prevTags[reqTag]; !ok {
			continue
		}

		delete(prevTags, reqTag)
	}

	prevTagsIDToBeDelete := make([]int, 0)
	for _, tag := range prevTags {
		prevTagsIDToBeDelete = append(prevTagsIDToBeDelete, tag.ID)
	}

	err = ps.postRepository.UpdatePost(ctx, model.Post{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
		Tags:    newTagsToBeSave,
	}, prevTagsIDToBeDelete...)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PostService) DeletePost(ctx context.Context, id int) error {
	post, err := ps.postRepository.GetPost(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrorNotFound{
				EntityName: "post",
				EntityID:   id,
			}
		}
		return err
	}
	err = ps.postRepository.DeletePost(ctx, *post)
	if err != nil {
		return err
	}

	return nil
}

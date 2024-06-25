package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/elangreza14/assetfindr-test/dto"
	gomockService "github.com/elangreza14/assetfindr-test/mock/service"
	"github.com/elangreza14/assetfindr-test/model"
	. "github.com/elangreza14/assetfindr-test/service"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type TestPostServiceSuite struct {
	suite.Suite

	MockPostRepo      *gomockService.MockIPostRepository
	MockCreatePostReq dto.CreateOrUpdatePostRequest
	Cs                *PostService
	Ctrl              *gomock.Controller
}

func (suite *TestPostServiceSuite) SetupSuite() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockPostRepo = gomockService.NewMockIPostRepository(suite.Ctrl)
	suite.MockCreatePostReq = dto.CreateOrUpdatePostRequest{
		Title:   "test",
		Content: "test",
		Tags:    []string{"test1", "test2"},
	}
	suite.Cs = NewPostService(suite.MockPostRepo)

}

func (suite *TestPostServiceSuite) TearDownSuite() {
	suite.Ctrl.Finish()
}

func TestPostServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TestPostServiceSuite))
}

func (suite *TestPostServiceSuite) TestPostService_CreatePost() {
	suite.Run("error when create", func() {
		suite.MockPostRepo.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(errors.New("err from db"))

		err := suite.Cs.CreatePost(context.Background(), suite.MockCreatePostReq)
		suite.Error(err)
		suite.Equal(err.Error(), "err from db")
	})

	suite.Run("success", func() {
		suite.MockPostRepo.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(nil)

		err := suite.Cs.CreatePost(context.Background(), suite.MockCreatePostReq)
		suite.NoError(err)
	})
}

func (suite *TestPostServiceSuite) TestPostService_DeletePost() {
	suite.Run("error when get post", func() {
		suite.MockPostRepo.EXPECT().GetPost(gomock.Any(), gomock.Any()).Return(nil, errors.New("err from db"))

		err := suite.Cs.DeletePost(context.Background(), 1)
		suite.Error(err)
		suite.Equal(err.Error(), "err from db")
	})

	suite.Run("error not found when get post", func() {
		suite.MockPostRepo.EXPECT().GetPost(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		err := suite.Cs.DeletePost(context.Background(), 1)
		suite.Error(err)
		suite.Equal(err.Error(), "cannot find post with id 1")
	})

	suite.Run("error delete post", func() {
		suite.MockPostRepo.EXPECT().GetPost(gomock.Any(), gomock.Any()).Return(&model.Post{
			ID:      1,
			Title:   "test",
			Content: "test",
			Tags: []*model.Tag{{
				ID:    1,
				Label: "test label",
			}, {
				ID:    1,
				Label: "test 2",
			}},
		}, nil)
		suite.MockPostRepo.EXPECT().DeletePost(gomock.Any(), gomock.Any()).Return(errors.New("error delete"))

		err := suite.Cs.DeletePost(context.Background(), 1)
		suite.Error(err)
		suite.Equal(err.Error(), "error delete")
	})

	suite.Run("success delete", func() {
		suite.MockPostRepo.EXPECT().GetPost(gomock.Any(), gomock.Any()).Return(&model.Post{
			ID:      1,
			Title:   "test",
			Content: "test",
			Tags: []*model.Tag{{
				ID:    1,
				Label: "1",
			}},
		}, nil)
		suite.MockPostRepo.EXPECT().DeletePost(gomock.Any(), gomock.Any()).Return(nil)

		err := suite.Cs.DeletePost(context.Background(), 1)
		suite.NoError(err)
	})
}

func (suite *TestPostServiceSuite) TestPostService_UpdatePost() {
	suite.Run("error when get post", func() {
		suite.MockPostRepo.EXPECT().GetPost(gomock.Any(), gomock.Any()).Return(nil, errors.New("err from db"))

		err := suite.Cs.UpdatePost(context.Background(), suite.MockCreatePostReq, 1)
		suite.Error(err)
		suite.Equal(err.Error(), "err from db")
	})

	suite.Run("error not found when get post", func() {
		suite.MockPostRepo.EXPECT().GetPost(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		err := suite.Cs.UpdatePost(context.Background(), suite.MockCreatePostReq, 1)
		suite.Error(err)
		suite.Equal(err.Error(), "cannot find post with id 1")
	})

	suite.Run("error update post", func() {
		suite.MockPostRepo.EXPECT().GetPost(gomock.Any(), gomock.Any()).Return(&model.Post{
			ID:      1,
			Title:   "test",
			Content: "test",
			Tags: []*model.Tag{{
				ID:    1,
				Label: "test label",
			}, {
				ID:    2,
				Label: "test2",
			}, {
				ID:    3,
				Label: "test 3",
			}},
		}, nil)
		suite.MockPostRepo.EXPECT().UpdatePost(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error update"))

		err := suite.Cs.UpdatePost(context.Background(), suite.MockCreatePostReq, 1)
		suite.Error(err)
		suite.Equal(err.Error(), "error update")
	})

	suite.Run("success", func() {
		suite.MockPostRepo.EXPECT().GetPost(gomock.Any(), gomock.Any()).Return(&model.Post{
			ID:      1,
			Title:   "test",
			Content: "test",
			Tags: []*model.Tag{{
				ID:    1,
				Label: "test label",
			}, {
				ID:    1,
				Label: "test 2",
			}},
		}, nil)
		suite.MockPostRepo.EXPECT().UpdatePost(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := suite.Cs.UpdatePost(context.Background(), suite.MockCreatePostReq, 1)
		suite.NoError(err)
	})
}

func (suite *TestPostServiceSuite) TestPostService_GetPost() {
	suite.Run("error when get post", func() {
		suite.MockPostRepo.EXPECT().GetPost(gomock.Any(), gomock.Any()).Return(nil, errors.New("err from db"))

		res, err := suite.Cs.GetPost(context.Background(), 1)
		suite.Error(err)
		suite.Nil(res)
		suite.Equal(err.Error(), "err from db")
	})

	suite.Run("error not found when get post", func() {
		suite.MockPostRepo.EXPECT().GetPost(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		res, err := suite.Cs.GetPost(context.Background(), 1)
		suite.Error(err)
		suite.Nil(res)
		suite.Equal(err.Error(), "cannot find post with id 1")
	})

	suite.Run("success", func() {
		suite.MockPostRepo.EXPECT().GetPost(gomock.Any(), gomock.Any()).Return(&model.Post{
			ID:      1,
			Title:   "test",
			Content: "test",
			Tags: []*model.Tag{{
				ID:    1,
				Label: "test tag",
			}},
		}, nil)

		res, err := suite.Cs.GetPost(context.Background(), 1)
		suite.NoError(err)
		suite.NotNil(res)
	})
}

func (suite *TestPostServiceSuite) TestPostService_GetPosts() {
	suite.Run("error when get posts", func() {
		suite.MockPostRepo.EXPECT().GetPosts(gomock.Any()).Return(nil, errors.New("err from db"))

		res, err := suite.Cs.GetPosts(context.Background())
		suite.Error(err)
		suite.Nil(res)
		suite.Equal(err.Error(), "err from db")
	})

	suite.Run("success", func() {
		suite.MockPostRepo.EXPECT().GetPosts(gomock.Any()).Return([]model.Post{{
			ID:      1,
			Title:   "test",
			Content: "test",
			Tags: []*model.Tag{{
				ID:    1,
				Label: "test tag",
			}}},
		}, nil)

		res, err := suite.Cs.GetPosts(context.Background())
		suite.NoError(err)
		suite.NotNil(res)
	})
}

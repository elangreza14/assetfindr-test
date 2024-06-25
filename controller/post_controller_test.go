package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elangreza14/assetfindr-test/cmd/http/routes"
	"github.com/elangreza14/assetfindr-test/controller"
	"github.com/elangreza14/assetfindr-test/dto"
	PostController "github.com/elangreza14/assetfindr-test/mock/controller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type TestPostControllerSuite struct {
	suite.Suite

	Ctrl            *gomock.Controller
	MockPostService *PostController.MockIPostService
}

func (suite *TestPostControllerSuite) SetupSuite() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockPostService = PostController.NewMockIPostService(suite.Ctrl)
}

func (suite *TestPostControllerSuite) TearDownSuite() {
	suite.Ctrl.Finish()
}

func TestPostControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TestPostControllerSuite))
}

func (suite *TestPostControllerSuite) TestPostController_CreatePost() {
	postController := controller.NewPostController(suite.MockPostService)

	router := gin.Default()
	apiGroup := router.Group("/api")
	routes.PostRoute(apiGroup, postController)

	suite.Run("error from validation", func() {
		errRequestBody := dto.CreateOrUpdatePostRequest{
			Title:   "",
			Content: "",
			Tags:    []string{},
		}
		errPayload, _ := json.Marshal(errRequestBody)

		bodyReader := bytes.NewReader(errPayload)
		req, _ := http.NewRequest(http.MethodPost, "/api/posts", bodyReader)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"errors","error":[{"field":"Title","message":"This field is required"},{"field":"Content","message":"This field is required"},{"field":"Tags","message":"Should be greater than 0"}]}`, string(responseData))
		suite.Equal(http.StatusBadRequest, w.Code)
	})

	successBody := dto.CreateOrUpdatePostRequest{
		Title:   "test",
		Content: "test",
		Tags:    []string{"test"},
	}
	payload, _ := json.Marshal(successBody)

	suite.Run("error from service", func() {
		bodyReader := bytes.NewReader(payload)
		suite.MockPostService.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(errors.New("test error from service"))
		req, _ := http.NewRequest(http.MethodPost, "/api/posts", bodyReader)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"error","error":"test error from service"}`, string(responseData))
		suite.Equal(http.StatusInternalServerError, w.Code)
	})

	suite.Run("success", func() {
		bodyReader := bytes.NewReader(payload)
		suite.MockPostService.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(nil)
		req, _ := http.NewRequest(http.MethodPost, "/api/posts", bodyReader)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"created"}`, string(responseData))
		suite.Equal(http.StatusCreated, w.Code)
	})
}

func (suite *TestPostControllerSuite) TestPostController_GetPosts() {
	postController := controller.NewPostController(suite.MockPostService)

	router := gin.Default()
	apiGroup := router.Group("/api")
	routes.PostRoute(apiGroup, postController)

	suite.Run("error from service", func() {

		suite.MockPostService.EXPECT().GetPosts(gomock.Any()).Return(nil, errors.New("test error from service"))
		req, _ := http.NewRequest(http.MethodGet, "/api/posts", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"error","error":"test error from service"}`, string(responseData))
		suite.Equal(http.StatusInternalServerError, w.Code)
	})

	suite.Run("success", func() {

		suite.MockPostService.EXPECT().GetPosts(gomock.Any()).Return([]dto.GetPostResponse{{
			ID:      1,
			Title:   "test",
			Content: "test",
			Tags:    []string{"test"},
		}}, nil)
		req, _ := http.NewRequest(http.MethodGet, "/api/posts", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"data":[{"id":1,"title":"test","content":"test","tags":["test"]}],"result":"ok"}`, string(responseData))
		suite.Equal(http.StatusOK, w.Code)
	})
}

func (suite *TestPostControllerSuite) TestPostController_UpdatePost() {
	postController := controller.NewPostController(suite.MockPostService)

	router := gin.Default()
	apiGroup := router.Group("/api")
	routes.PostRoute(apiGroup, postController)

	suite.Run("error from uri id", func() {
		errRequestBody := dto.CreateOrUpdatePostRequest{
			Title:   "",
			Content: "",
			Tags:    []string{},
		}
		errPayload, _ := json.Marshal(errRequestBody)

		bodyReader := bytes.NewReader(errPayload)
		req, _ := http.NewRequest(http.MethodPut, "/api/posts/1212aasas", bodyReader)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"error","error":"strconv.ParseInt: parsing \"1212aasas\": invalid syntax"}`, string(responseData))
		suite.Equal(http.StatusBadRequest, w.Code)
	})

	suite.Run("error from body", func() {
		errRequestBody := dto.CreateOrUpdatePostRequest{
			Title:   "",
			Content: "",
			Tags:    []string{},
		}
		errPayload, _ := json.Marshal(errRequestBody)

		bodyReader := bytes.NewReader(errPayload)
		req, _ := http.NewRequest(http.MethodPut, "/api/posts/1212", bodyReader)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"errors","error":[{"field":"Title","message":"This field is required"},{"field":"Content","message":"This field is required"},{"field":"Tags","message":"Should be greater than 0"}]}`, string(responseData))
		suite.Equal(http.StatusBadRequest, w.Code)
	})

	successBody := dto.CreateOrUpdatePostRequest{
		Title:   "test",
		Content: "test",
		Tags:    []string{"test"},
	}
	payload, _ := json.Marshal(successBody)

	suite.Run("error internal from service", func() {
		bodyReader := bytes.NewReader(payload)
		suite.MockPostService.EXPECT().UpdatePost(gomock.Any(), gomock.Any(), 1).Return(errors.New("test error from service"))
		req, _ := http.NewRequest(http.MethodPut, "/api/posts/1", bodyReader)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"error","error":"test error from service"}`, string(responseData))
		suite.Equal(http.StatusInternalServerError, w.Code)
	})

	suite.Run("error not found from service", func() {
		bodyReader := bytes.NewReader(payload)
		suite.MockPostService.EXPECT().UpdatePost(gomock.Any(), gomock.Any(), 3).Return(dto.ErrorNotFound{EntityName: "post", EntityID: 3})
		req, _ := http.NewRequest(http.MethodPut, "/api/posts/3", bodyReader)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"error","error":"cannot find post with id 3"}`, string(responseData))
		suite.Equal(http.StatusNotFound, w.Code)
	})

	suite.Run("success", func() {
		bodyReader := bytes.NewReader(payload)
		suite.MockPostService.EXPECT().UpdatePost(gomock.Any(), gomock.Any(), 2).Return(nil)
		req, _ := http.NewRequest(http.MethodPut, "/api/posts/2", bodyReader)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"updated"}`, string(responseData))
		suite.Equal(http.StatusOK, w.Code)
	})
}

func (suite *TestPostControllerSuite) TestPostController_DeletePost() {
	postController := controller.NewPostController(suite.MockPostService)

	router := gin.Default()
	apiGroup := router.Group("/api")
	routes.PostRoute(apiGroup, postController)

	suite.Run("error from uri id", func() {

		req, _ := http.NewRequest(http.MethodDelete, "/api/posts/1212aasas", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"error","error":"strconv.ParseInt: parsing \"1212aasas\": invalid syntax"}`, string(responseData))
		suite.Equal(http.StatusBadRequest, w.Code)
	})

	suite.Run("error internal from service", func() {

		suite.MockPostService.EXPECT().DeletePost(gomock.Any(), 1).Return(errors.New("test error from service"))
		req, _ := http.NewRequest(http.MethodDelete, "/api/posts/1", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"error","error":"test error from service"}`, string(responseData))
		suite.Equal(http.StatusInternalServerError, w.Code)
	})

	suite.Run("error not found from service", func() {

		suite.MockPostService.EXPECT().DeletePost(gomock.Any(), 3).Return(dto.ErrorNotFound{EntityName: "post", EntityID: 3})
		req, _ := http.NewRequest(http.MethodDelete, "/api/posts/3", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"error","error":"cannot find post with id 3"}`, string(responseData))
		suite.Equal(http.StatusNotFound, w.Code)
	})

	suite.Run("success", func() {

		suite.MockPostService.EXPECT().DeletePost(gomock.Any(), 2).Return(nil)
		req, _ := http.NewRequest(http.MethodDelete, "/api/posts/2", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"deleted"}`, string(responseData))
		suite.Equal(http.StatusOK, w.Code)
	})
}

func (suite *TestPostControllerSuite) TestPostController_GetPost() {
	postController := controller.NewPostController(suite.MockPostService)

	router := gin.Default()
	apiGroup := router.Group("/api")
	routes.PostRoute(apiGroup, postController)

	suite.Run("error from uri id", func() {

		req, _ := http.NewRequest(http.MethodGet, "/api/posts/1212aasas", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"error","error":"strconv.ParseInt: parsing \"1212aasas\": invalid syntax"}`, string(responseData))
		suite.Equal(http.StatusBadRequest, w.Code)
	})

	suite.Run("error internal from service", func() {

		suite.MockPostService.EXPECT().GetPost(gomock.Any(), 1).Return(nil, errors.New("test error from service"))
		req, _ := http.NewRequest(http.MethodGet, "/api/posts/1", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"error","error":"test error from service"}`, string(responseData))
		suite.Equal(http.StatusInternalServerError, w.Code)
	})

	suite.Run("error not found from service", func() {

		suite.MockPostService.EXPECT().GetPost(gomock.Any(), 3).Return(nil, dto.ErrorNotFound{EntityName: "post", EntityID: 3})
		req, _ := http.NewRequest(http.MethodGet, "/api/posts/3", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"result":"error","error":"cannot find post with id 3"}`, string(responseData))
		suite.Equal(http.StatusNotFound, w.Code)
	})

	suite.Run("success", func() {

		suite.MockPostService.EXPECT().GetPost(gomock.Any(), 2).Return(&dto.GetPostResponse{
			ID:      1,
			Title:   "test",
			Content: "test",
			Tags:    []string{"test"},
		}, nil)
		req, _ := http.NewRequest(http.MethodGet, "/api/posts/2", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		suite.Equal(`{"data":{"id":1,"title":"test","content":"test","tags":["test"]},"result":"ok"}`, string(responseData))
		suite.Equal(http.StatusOK, w.Code)
	})
}

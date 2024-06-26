package repository_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elangreza14/assetfindr-test/model"
	. "github.com/elangreza14/assetfindr-test/repository"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupDbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatal(err)
	}

	return sqlDB, gormDB, mock
}

type TestPostRepositorySuite struct {
	suite.Suite

	sqlDB    *sql.DB
	gormDB   *gorm.DB
	mock     sqlmock.Sqlmock
	Ctrl     *gomock.Controller
	postRepo *PostRepository
}

func (suite *TestPostRepositorySuite) SetupSuite() {
	suite.Ctrl = gomock.NewController(suite.T())
	sqlDB, gormDB, mock := setupDbMock(suite.T())

	suite.sqlDB = sqlDB
	suite.gormDB = gormDB
	suite.mock = mock
	suite.postRepo = NewPostRepository(gormDB)
}

func (suite *TestPostRepositorySuite) TearDownSuite() {
	suite.Ctrl.Finish()
	suite.sqlDB.Close()
}

func TestPostRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TestPostRepositorySuite))
}

func (suite *TestPostRepositorySuite) TestPostRepository_GetPost() {
	suite.Run("err", func() {
		suite.mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "posts" WHERE "posts"."id" = $1 ORDER BY "posts"."id" LIMIT $2`)).
			WithArgs(1, 1).WillReturnError(gorm.ErrInvalidData)

		res, err := suite.postRepo.GetPost(context.Background(), 1)
		suite.Error(err)
		suite.True(errors.Is(err, gorm.ErrInvalidData))
		suite.Nil(res)
	})

	suite.Run("success", func() {
		suite.mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "posts" WHERE "posts"."id" = $1 ORDER BY "posts"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "tags"}).AddRow(1, 1, 1, 1))

		suite.mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "post_tags" WHERE "post_tags"."post_id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"post_id", "tag_id"}).AddRow(1, 1))

		suite.mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "tags" WHERE "tags"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "label"}).AddRow(1, 1))

		res, err := suite.postRepo.GetPost(context.Background(), 1)
		suite.NoError(err)
		suite.NotNil(res)
	})
}

func (suite *TestPostRepositorySuite) TestPostRepository_GetPosts() {

	suite.Run("err", func() {
		suite.mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "posts" ORDER BY id desc`)).
			WillReturnError(gorm.ErrRecordNotFound)

		res, err := suite.postRepo.GetPosts(context.Background())
		suite.Error(err)
		suite.Nil(res)
	})

	suite.Run("success", func() {
		suite.mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "posts" ORDER BY id desc`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "tags"}).AddRow(1, 1, 1, 1))

		suite.mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "post_tags" WHERE "post_tags"."post_id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"post_id", "tag_id"}).AddRow(1, 1))

		suite.mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "tags" WHERE "tags"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "label"}).AddRow(1, 1))

		res, err := suite.postRepo.GetPosts(context.Background())
		suite.NoError(err)
		suite.NotNil(res)
	})
}

func (suite *TestPostRepositorySuite) TestPostRepository_CreatePost() {
	testReq := model.Post{
		Title:   "test",
		Content: "test",
		Tags: []*model.Tag{{
			ID:    1,
			Label: "test",
		}},
	}

	suite.Run("success", func() {

		suite.mock.ExpectBegin()
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "posts" ("title","content") VALUES ($1,$2) RETURNING "id"`)).
			WithArgs("test", "test").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		suite.mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "tags" WHERE "tags"."label" = $1 ORDER BY "tags"."id" LIMIT $2`)).
			WithArgs("test", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "label"}).AddRow(1, "test"))
		suite.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO post_tags ("post_id","tag_id") VALUES ($1, $2) ON CONFLICT DO NOTHING`)).
			WithArgs(1, 1).WillReturnResult(driver.ResultNoRows)

		suite.mock.ExpectCommit()

		err := suite.postRepo.CreatePost(context.Background(), testReq)
		suite.NoError(err)
	})
}

func (suite *TestPostRepositorySuite) TestPostRepository_UpdatePost() {
	testReq := model.Post{
		ID:      1,
		Title:   "test",
		Content: "test",
		Tags: []*model.Tag{{
			ID:    1,
			Label: "test 1",
		}},
	}

	suite.Run("success", func() {

		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM post_tags WHERE post_id=$1 and tag_id IN ($2);`)).
			WithArgs(1, 1).WillReturnResult(driver.ResultNoRows)
		updUserSQL := "UPDATE \"posts\" SET .+"
		suite.mock.ExpectExec(updUserSQL).
			WithArgs("test", "test", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		suite.mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "tags" WHERE "tags"."label" = $1 ORDER BY "tags"."id" LIMIT $2`)).
			WithArgs("test 1", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "label"}).AddRow(1, "test 1"))
		suite.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO post_tags ("post_id","tag_id") VALUES ($1, $2) ON CONFLICT DO NOTHING`)).
			WithArgs(1, 1).WillReturnResult(driver.ResultNoRows)

		suite.mock.ExpectCommit()

		err := suite.postRepo.UpdatePost(context.Background(), testReq, 1)
		suite.NoError(err)
	})
}

func (suite *TestPostRepositorySuite) TestPostRepository_DeletePost() {
	testReq := model.Post{
		ID: 1,
	}

	suite.Run("success", func() {

		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM post_tags WHERE post_id=$1;`)).
			WithArgs(1).WillReturnResult(driver.ResultNoRows)
		suite.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "posts" WHERE "posts"."id" = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		suite.mock.ExpectCommit()

		err := suite.postRepo.DeletePost(context.Background(), testReq)
		suite.NoError(err)
	})
}

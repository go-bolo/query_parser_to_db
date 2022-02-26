package query_parser_to_db

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGormDBOperations(t *testing.T) {
	assert := assert.New(t)

	db := GetFakeGormDB()
	err := db.AutoMigrate(&ContentModelStub{})
	assert.Nil(err)

	urlString := "https://example.com/example"
	parsedURL, _ := url.Parse(urlString)

	q := NewQuery(50)
	err = q.ParseFromURLValues(parsedURL.Query())
	assert.Nil(err)

	t.Run("Should generate a valid equal query", func(t *testing.T) {
		fieldName, value := "title", "Lov"
		db.DryRun = true

		queryI, err := gormDBOperations["equal"](fieldName, value, db, q)
		assert.Nil(err)

		query := queryI.(*gorm.DB)
		query.Find(&[]ContentModelStub{})

		assert.Equal("SELECT * FROM `content_model_stubs` WHERE title = ?", query.Statement.SQL.String())
		assert.Equal([]interface{}{"Lov"}, query.Statement.Vars)

		db.DryRun = false
	})

	t.Run("Should generate a valid not-equal query", func(t *testing.T) {
		fieldName, value := "body", "NotLovi"
		db.DryRun = true

		queryI, err := gormDBOperations["not-equal"](fieldName, value, db, q)
		assert.Nil(err)

		query := queryI.(*gorm.DB)
		query.Find(&[]ContentModelStub{})

		assert.Equal("SELECT * FROM `content_model_stubs` WHERE body != ?", query.Statement.SQL.String())
		assert.Equal([]interface{}{"NotLovi"}, query.Statement.Vars)

		db.DryRun = false
	})
	t.Run("Should generate a valid is-null query", func(t *testing.T) {
		fieldName := "title"
		db.DryRun = true

		queryI, err := gormDBOperations["is-null"](fieldName, "", db, q)
		assert.Nil(err)

		query := queryI.(*gorm.DB)
		query.Find(&[]ContentModelStub{})

		assert.Equal("SELECT * FROM `content_model_stubs` WHERE title IS NULL", query.Statement.SQL.String())
		assert.Equal([]interface{}{}, query.Statement.Vars)

		db.DryRun = false
	})

	t.Run("Should generate a valid is-not-null query", func(t *testing.T) {
		fieldName := "title"
		db.DryRun = true

		queryI, err := gormDBOperations["is-not-null"](fieldName, "", db, q)
		assert.Nil(err)

		query := queryI.(*gorm.DB)
		query.Find(&[]ContentModelStub{})

		assert.Equal("SELECT * FROM `content_model_stubs` WHERE title IS NOT NULL", query.Statement.SQL.String())
		assert.Equal([]interface{}{}, query.Statement.Vars)

		db.DryRun = false
	})

	t.Run("Should generate a valid starts-with query", func(t *testing.T) {
		fieldName, value := "title", "Lo"
		db.DryRun = true

		queryI, err := gormDBOperations["starts-with"](fieldName, value, db, q)
		assert.Nil(err)

		query := queryI.(*gorm.DB)
		query.Find(&[]ContentModelStub{})

		assert.Equal("SELECT * FROM `content_model_stubs` WHERE title LIKE ?", query.Statement.SQL.String())
		assert.Equal([]interface{}{"Lo%"}, query.Statement.Vars)

		db.DryRun = false
	})

	t.Run("Should generate a valid not-starts-with query", func(t *testing.T) {
		fieldName, value := "title", "Lo"
		db.DryRun = true

		queryI, err := gormDBOperations["not-starts-with"](fieldName, value, db, q)
		assert.Nil(err)

		query := queryI.(*gorm.DB)
		query.Find(&[]ContentModelStub{})

		assert.Equal("SELECT * FROM `content_model_stubs` WHERE title NOT LIKE ?", query.Statement.SQL.String())
		assert.Equal([]interface{}{"Lo%"}, query.Statement.Vars)

		db.DryRun = false
	})

	t.Run("Should generate a valid ends-with query", func(t *testing.T) {
		fieldName, value := "title", "ve"
		db.DryRun = true

		queryI, err := gormDBOperations["ends-with"](fieldName, value, db, q)
		assert.Nil(err)

		query := queryI.(*gorm.DB)
		query.Find(&[]ContentModelStub{})

		assert.Equal("SELECT * FROM `content_model_stubs` WHERE title LIKE ?", query.Statement.SQL.String())
		assert.Equal([]interface{}{"%ve"}, query.Statement.Vars)

		db.DryRun = false
	})

	t.Run("Should generate a valid not-ends-with query", func(t *testing.T) {
		fieldName, value := "body", "ve"
		db.DryRun = true

		queryI, err := gormDBOperations["not-ends-with"](fieldName, value, db, q)
		assert.Nil(err)

		query := queryI.(*gorm.DB)
		query.Find(&[]ContentModelStub{})

		assert.Equal("SELECT * FROM `content_model_stubs` WHERE body NOT LIKE ?", query.Statement.SQL.String())
		assert.Equal([]interface{}{"%ve"}, query.Statement.Vars)

		db.DryRun = false
	})

	t.Run("Should generate a valid contains query", func(t *testing.T) {
		fieldName, value := "body", "ov"
		db.DryRun = true

		queryI, err := gormDBOperations["contains"](fieldName, value, db, q)
		assert.Nil(err)

		query := queryI.(*gorm.DB)
		query.Find(&[]ContentModelStub{})

		assert.Equal("SELECT * FROM `content_model_stubs` WHERE body LIKE ?", query.Statement.SQL.String())
		assert.Equal([]interface{}{"%ov%"}, query.Statement.Vars)

		db.DryRun = false
	})

	t.Run("Should generate a valid not-contains query", func(t *testing.T) {
		fieldName, value := "body", "hate"
		db.DryRun = true

		queryI, err := gormDBOperations["not-contains"](fieldName, value, db, q)
		assert.Nil(err)

		query := queryI.(*gorm.DB)
		query.Find(&[]ContentModelStub{})

		assert.Equal("SELECT * FROM `content_model_stubs` WHERE body NOT LIKE ?", query.Statement.SQL.String())
		assert.Equal([]interface{}{"%hate%"}, query.Statement.Vars)

		db.DryRun = false
	})
}

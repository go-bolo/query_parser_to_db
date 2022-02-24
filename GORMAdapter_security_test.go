package query_parser

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGORMAdapterSecurity(t *testing.T) {
	assert := assert.New(t)

	db := GetFakeGormDB()
	err := db.AutoMigrate(&ContentModelStub{})
	assert.Nil(err)

	t.Run("Should ignore contains params if has 'drop table users;' dryRun", func(t *testing.T) {
		urlString := "https://example.com/example?title_contains='he;drop table users;'&limit=5&page=1"
		parsedURL, _ := url.Parse(urlString)
		// rawParamName
		q := NewQuery(50)
		err := q.ParseFromURLValues(parsedURL.Query())
		assert.Nil(err)

		query := GetFakeGormDB()
		query.DryRun = true

		query2, err := q.SetDatabaseQueryForModel(query, &ContentModelStub{})
		assert.Nil(err)
		query = query2.(*gorm.DB)

		records := []ContentModelStub{}

		r := query.Find(&records)

		assert.Nil(r.Error)
		assert.Equal("SELECT * FROM `content_model_stubs`", r.Statement.SQL.String())
		assert.Equal(0, len(r.Statement.Vars))

		query.DryRun = false
	})

	t.Run("Should ignore = params if has 'drop table users;' dryRun", func(t *testing.T) {
		urlString := "https://example.com/example?title='he;drop table users;'&limit=5&page=1"
		parsedURL, _ := url.Parse(urlString)
		// rawParamName
		q := NewQuery(50)
		err := q.ParseFromURLValues(parsedURL.Query())
		assert.Nil(err)

		query := GetFakeGormDB()
		query.DryRun = true

		query2, err := q.SetDatabaseQueryForModel(query, &ContentModelStub{})
		assert.Nil(err)
		query = query2.(*gorm.DB)

		records := []ContentModelStub{}

		r := query.Find(&records)

		assert.Nil(r.Error)
		assert.Equal("SELECT * FROM `content_model_stubs`", r.Statement.SQL.String())
		assert.Equal(0, len(r.Statement.Vars))

		query.DryRun = false
	})

}

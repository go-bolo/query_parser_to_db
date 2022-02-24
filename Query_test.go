package query_parser

import (
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

var db *gorm.DB

func GetFakeGormDB() *gorm.DB {
	if db != nil {
		return db
	}

	d, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: gorm_logger.Default.LogMode(gorm_logger.Warn),
	})
	db = d

	return db
}

type ContentModelStub struct {
	ID         uint64 `json:"id" filter:"param:id;type:string"`
	Title      string `json:"title" filter:"param:title;type:string"`
	Body       string `json:"body" filter:"type:string"`
	Published  bool   `json:"published" filter:"param:published;type:bool"`
	ClickCount int64  `json:"clickCount" filter:"param:clickCount;type:number"`
	Secret     string `json:"-"`
	Email      string `json:"email"`
	Email2     string `json:"email2" filter:""`
	PrivateBio string `json:"-" filter:"-"`
}

func GetContentModelStub() ContentModelStub {
	return ContentModelStub{
		// ID:         gofakeit.Uint64(),
		Title:      gofakeit.Paragraph(1, 4, 4, " "),
		Body:       gofakeit.Paragraph(1, 3, 5, " "),
		Published:  true,
		Secret:     gofakeit.Word(),
		Email:      gofakeit.Email(),
		Email2:     gofakeit.Email(),
		PrivateBio: gofakeit.Paragraph(1, 4, 4, ""),
	}
}

type QueryConfiguration struct {
	// string, number ...etc
	Type string
}

func TestQueryParser(t *testing.T) {
	assert := assert.New(t)

	db := GetFakeGormDB()
	err := db.AutoMigrate(&ContentModelStub{})
	assert.Nil(err)

	t.Run("Should parse and load data from valid url query params", func(t *testing.T) {
		urlString := "https://example.com/example?content_contains=He&id=10&limit=3&page=2"
		parsedURL, _ := url.Parse(urlString)

		// rawParamName
		q := NewQuery(50)
		err := q.ParseFromURLValues(parsedURL.Query())
		assert.Nil(err)

		assert.Equal(int64(3), q.GetLimit())
		assert.Equal(int64(2), q.GetPage())
		assert.Equal("10", q.GetParamValue("id"))
		assert.Equal("He", q.GetParamValue("content"))
	})
	t.Run("Should parse and generate a DryRun GORM sql", func(t *testing.T) {
		urlString := "https://example.com/example?title_contains=Hello&limit=5&page=3"
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
		assert.Equal("SELECT * FROM `content_model_stubs` WHERE title LIKE ? LIMIT 5 OFFSET 10", r.Statement.SQL.String())

		query.DryRun = false
	})
	t.Run("Should parse and run a valid sql query", func(t *testing.T) {
		urlString := "https://example.com/example?title_contains=Hello&limit=3"
		parsedURL, _ := url.Parse(urlString)

		// rawParamName
		q := NewQuery(50)
		err := q.ParseFromURLValues(parsedURL.Query())
		assert.Nil(err)

		db := GetFakeGormDB()

		// add 3 records
		recordsToPreSave := []ContentModelStub{
			GetContentModelStub(),
			GetContentModelStub(),
			GetContentModelStub(),
		}
		recordsToPreSave[2].Title = "Hellou World"

		err = db.Create(recordsToPreSave).Error
		assert.Nil(err)

		query := GetFakeGormDB()

		query2, err := q.SetDatabaseQueryForModel(query, &ContentModelStub{})
		assert.Nil(err)
		query = query2.(*gorm.DB)

		records := []ContentModelStub{}

		r := query.Find(&records)
		assert.Nil(r.Error)
		assert.Equal(1, len(records))
		assert.Equal("Hellou World", records[0].Title)
	})

}

func TestQueryParserGetSetMethods(t *testing.T) {
	assert := assert.New(t)

	t.Run("Should set and get Limit", func(t *testing.T) {
		urlString := "https://example.com/example?title_contains=Hello&limit=3"
		parsedURL, _ := url.Parse(urlString)

		// rawParamName
		q := NewQuery(50)
		err := q.ParseFromURLValues(parsedURL.Query())
		assert.Nil(err)

		assert.Equal(int64(3), q.GetLimit())
		q.SetLimit(20)
		assert.Equal(int64(20), q.GetLimit())
		// should not pass the limitMax value
		q.SetLimit(2000000)
		assert.Equal(int64(50), q.GetLimit())

		// should not be negative
		q.SetLimit(-10)
		assert.Equal(int64(10), q.GetLimit())
	})

	t.Run("Should set and get page", func(t *testing.T) {
		urlString := "https://example.com/example?title_contains=Hello&limit=10&page=2"
		parsedURL, _ := url.Parse(urlString)

		// rawParamName
		q := NewQuery(50)
		err := q.ParseFromURLValues(parsedURL.Query())
		assert.Nil(err)

		assert.Equal(int64(2), q.GetPage())
		q.SetPage(5)
		assert.Equal(int64(5), q.GetPage())
		// should not be negative
		q.SetPage(-10)
		assert.Equal(int64(0), q.GetPage())
	})
}

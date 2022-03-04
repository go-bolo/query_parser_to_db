# Go Catupiri Query Parser to DB

Query parser with support for build database operations / params from query params.

For now we only have supports for GORM queries.

Examples: https://github.com/go-catupiry/query_parser_to_db_examples

## Usage:

```go
  // GORM Model configuration
  // use filter param to allow or disable any query param parsing:
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

  // get a new instance of the query parser:
  q := query_parser_to_db.NewQuery(50)
  // parse url query params and its operations:
	q.ParseFromURLValues(req.URL.Query())
  
  // db = *gorm.DB, that will parse 
	queryInterface, _ := q.SetDatabaseQueryForModel(db, &ContentModelStub{})
  
  // get the query, a *gorm.DB var
	query := queryInterface.(*gorm.DB)
  
  // execute the query as any gorm query:
	records := []ContentModelStub{}
	dbResultTX := query.Find(&records)
```

## Operations:

Will accept this query params as filters:

- 'get /post?id=[id]'
- 'get /post?id_equal=[id]'
- 'get /post?id_is-null=true'
- 'get /post?id_is-null=true'
- 'get /post?id_not-is-null=true'
- 'get /post?id_between=10-20'
- 'get /post?id_not-between=10-30'
- 'get /post?id_gt=2'
- 'get /post?id_gte=2'
- 'get /post?id_lt=20'
- 'get /post?id_lte=20'
- 'get /post?title=Oi mundo'
- 'get /post?title_equal=Oi mundo'
- 'get /post?title_is-null=true'
- 'get /post?title_not-is-null=true'
- 'get /post?title_starts-with=Oi'
- 'get /post?title_not-starts-with=Oi'
- 'get /post?title_ends-with=Mundo'
- 'get /post?title_not-ends-with=Mundo'
- 'get /post?title_contains=Mundo'
- 'get /post?title_not-contains=Mundo'
- 'get /post?body=Something'
- 'get /post?body_equal=Something'
- 'get /post?body_equal=Something'

## Roadmap

- Improve to allows database adapter extension with interfaces
  - Create one mongoDB adapter
- Add github CI tests and coverage
- Add support for between operation
- Add support for advanced JSON operations



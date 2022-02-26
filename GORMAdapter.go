package query_parser_to_db

import "gorm.io/gorm"

var gormDBOperations = DBOperations{
	"equal": func(fieldName, value string, q interface{}, r QueryInterface) (interface{}, error) {
		query := q.(*gorm.DB)
		query = query.Where(fieldName+" = ?", value)

		return query, nil
	},
	"not-equal": func(fieldName, value string, q interface{}, r QueryInterface) (interface{}, error) {
		query := q.(*gorm.DB)
		query = query.Where(fieldName+" != ?", value)

		return query, nil
	},
	"is-null": func(fieldName, value string, q interface{}, r QueryInterface) (interface{}, error) {
		query := q.(*gorm.DB)
		query = query.Where(fieldName + " IS NULL")

		return query, nil
	},
	"is-not-null": func(fieldName, value string, q interface{}, r QueryInterface) (interface{}, error) {
		query := q.(*gorm.DB)
		query = query.Where(fieldName + " IS NOT NULL")

		return query, nil
	},
	"starts-with": func(fieldName, value string, q interface{}, r QueryInterface) (interface{}, error) {
		query := q.(*gorm.DB)
		query = query.Where(fieldName+" LIKE ?", value+"%")

		return query, nil
	},
	"not-starts-with": func(fieldName, value string, q interface{}, r QueryInterface) (interface{}, error) {
		query := q.(*gorm.DB)
		query = query.Where(fieldName+" NOT LIKE ?", value+"%")

		return query, nil
	},
	"ends-with": func(fieldName, value string, q interface{}, r QueryInterface) (interface{}, error) {
		query := q.(*gorm.DB)
		query = query.Where(fieldName+" LIKE ?", "%"+value)

		return query, nil
	},
	"not-ends-with": func(fieldName, value string, q interface{}, r QueryInterface) (interface{}, error) {
		query := q.(*gorm.DB)
		query = query.Where(fieldName+" NOT LIKE ?", "%"+value)

		return query, nil
	},
	"contains": func(fieldName, value string, q interface{}, r QueryInterface) (interface{}, error) {
		query := q.(*gorm.DB)
		query = query.Where(fieldName+" LIKE ?", "%"+value+"%")

		return query, nil
	},
	"not-contains": func(fieldName, value string, q interface{}, r QueryInterface) (interface{}, error) {
		query := q.(*gorm.DB)
		query = query.Where(fieldName+" NOT LIKE ?", "%"+value+"%")

		return query, nil
	},
}

func NewGORMDBAdapter() DBAdapter {
	GORMDBAdapter = DBAdapter{
		"default": {
			"equal":       gormDBOperations["equal"],
			"not-equal":   gormDBOperations["not-equal"],
			"is-null":     gormDBOperations["is-null"],
			"is-not-null": gormDBOperations["is-not-null"],
		},
		"string": {
			"equal":           gormDBOperations["equal"],
			"not-equal":       gormDBOperations["not-equal"],
			"is-null":         gormDBOperations["is-null"],
			"is-not-null":     gormDBOperations["is-not-null"],
			"starts-with":     gormDBOperations["starts-with"],
			"not-starts-with": gormDBOperations["not-starts-with"],
			"ends-with":       gormDBOperations["ends-with"],
			"not-ends-with":   gormDBOperations["not-ends-with"],
			"contains":        gormDBOperations["contains"],
			"not-contains":    gormDBOperations["not-contains"],
		},
		"bool": {
			"equal":       gormDBOperations["equal"],
			"not-equal":   gormDBOperations["not-equal"],
			"is-null":     gormDBOperations["is-null"],
			"is-not-null": gormDBOperations["is-not-null"],
		},
		"number": {
			"equal":       gormDBOperations["equal"],
			"not-equal":   gormDBOperations["not-equal"],
			"is-null":     gormDBOperations["is-null"],
			"is-not-null": gormDBOperations["is-not-null"],
		},
		"pagination": {
			"pager": func(fieldName, value string, dbQuery interface{}, r QueryInterface) (interface{}, error) {
				query := dbQuery.(*gorm.DB)
				query = query.Limit(int(r.GetLimit())).Offset(r.GetOffset())
				return query, nil
			},
		},
	}
	// text and blob here will have same operations like string:
	GORMDBAdapter["text"] = GORMDBAdapter["string"]
	GORMDBAdapter["blob"] = GORMDBAdapter["string"]
	// TODO! add a bette support for JSON
	GORMDBAdapter["json"] = GORMDBAdapter["default"]
	// Date for all date like formats:
	GORMDBAdapter["time"] = GORMDBAdapter["date"]
	GORMDBAdapter["dateOnly"] = GORMDBAdapter["date"]

	return GORMDBAdapter
}

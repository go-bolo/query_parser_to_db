package query_parser_to_db

import "strings"

// [fieldType][queryType]function
type DBAdapter map[string]DBOperations

func (r DBAdapter) Run(fieldType, operator, column, value string, dbQuery interface{}, q *Query) (interface{}, error) {
	if r[fieldType] == nil || r[fieldType][operator] == nil {
		return nil, nil
	}

	return r[fieldType][operator](column, strings.TrimSpace(value), dbQuery, q)
}

type DBOperations map[string]func(column, value string, dbQuery interface{}, q QueryInterface) (interface{}, error)

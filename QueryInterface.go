package query_parser

import "net/url"

type QueryInterface interface {
	ParseFromURLValues(query url.Values) error
	AddQueryParamFromRaw(paramName string, values []string) error
	AddQueryString(paramName string, values []string)
	GetQueryString(paramName string) string
	GetParamValue(paramName string) string
	GetParam(paramName string) *QueryAttr
	// Get limit query param
	GetLimit() int64
	SetLimit(v int64)
	// Get page query param
	GetPage() int64
	SetPage(v int64)
	GetOffset() int
	SetDatabaseQueryForModel(query interface{}, model interface{}) (interface{}, error)
}

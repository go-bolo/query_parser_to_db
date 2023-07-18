package query_parser_to_db

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// [modelName][fieldName]ModelFieldTagConfig
var modelSearchTagsCache map[string]map[string]*ModelFieldTagConfig

var GORMDBAdapter DBAdapter

type ModelFieldTagConfig struct {
	Param       string
	Type        string
	DBFieldName string
}

type QueryAttr struct {
	Operator   string
	Values     []string
	IsMultiple bool
	ParamName  string
}

type Query struct {
	Fields      []QueryAttr
	Limit       int64
	LimitMax    int64
	Page        int64
	QueryString string
}

func init() {
	modelSearchTagsCache = make(map[string]map[string]*ModelFieldTagConfig)

	GORMDBAdapter = NewGORMDBAdapter()
}

func (r *Query) ParseFromURLValues(query url.Values) error {
	for key, param := range query {
		// get limit with max value for security:
		if key == "limit" && len(param) == 1 {
			queryLimit, err := strconv.ParseInt(param[0], 10, 64)
			if err != nil {
				return ErrInvalidQueryOperator
			}
			if queryLimit > 0 && queryLimit < r.LimitMax {
				r.Limit = queryLimit
			}
		}
		// page for build offset on queries:
		if key == "page" && len(param) == 1 {
			page, _ := strconv.ParseInt(param[0], 10, 64)
			r.Page = page
			continue
		}

		r.AddQueryParamFromRaw(key, param)
	}

	return nil
}

func (r *Query) AddQueryParamFromRaw(paramName string, values []string) error {
	if len(values) == 0 {
		return nil
	}

	if paramName == "page" {
		return nil
	}

	r.AddQueryString(paramName, values)

	var qAttr QueryAttr

	if len(values) > 1 {
		qAttr.IsMultiple = true
	}

	if !strings.Contains(paramName, "_") {
		qAttr.Values = values
		qAttr.ParamName = paramName
		qAttr.Operator = "equal"
		r.Fields = append(r.Fields, qAttr)
		return nil
	}

	for op := range gormDBOperations {
		if strings.HasSuffix(paramName, "_"+op) {
			qAttr.Values = values
			qAttr.ParamName = strings.Replace(paramName, "_"+op, "", 1)
			qAttr.Operator = op
			r.Fields = append(r.Fields, qAttr)
			return nil
		}
	}

	return nil
}

func (r *Query) AddQueryString(paramName string, values []string) {
	if r.QueryString != "" {
		r.QueryString += "&"
	}

	if len(values) > 1 {
		for i := range values {
			r.QueryString += paramName + "[]=" + values[i]
		}
	} else {
		r.QueryString += paramName + "=" + values[0]
	}
}

func (r *Query) GetQueryString(paramName string) string {
	for i := range r.Fields {
		if r.Fields[i].ParamName == paramName {
			if len(r.Fields[i].Values) == 0 {
				return ""
			} else if len(r.Fields[i].Values) == 1 {
				return paramName + `=` + r.Fields[i].Values[0]
			} else {
				var results []string
				for vi := range r.Fields[i].Values {
					result := paramName + "[]=" + r.Fields[i].Values[vi]
					results = append(results, result)
				}

				return strings.Join(results, "&")
			}
		}
	}

	return ""
}

func (r *Query) GetParam(paramName string) *QueryAttr {
	for i := range r.Fields {
		if r.Fields[i].ParamName == paramName {
			return &r.Fields[i]
		}
	}

	return nil
}

func (r *Query) GetParamValue(paramName string) string {
	for i := range r.Fields {
		if r.Fields[i].ParamName == paramName {
			if len(r.Fields[i].Values) != 0 {
				return r.Fields[i].Values[0]
			}
		}
	}

	return ""
}

func (r *Query) GetLimit() int64 {
	return r.Limit
}

func (r *Query) SetLimit(v int64) {
	if v > r.LimitMax {
		r.Limit = r.LimitMax
		return
	}

	if v < 0 {
		r.Limit = 10
		return
	}

	r.Limit = v
}

func (r *Query) GetPage() int64 {
	return r.Page
}

func (r *Query) SetPage(v int64) {
	if v < 0 {
		r.Page = 0
		return
	}

	r.Page = v
}

func (r *Query) GetOffset() int {
	page := int(r.Page)

	if page < 2 {
		return 0
	}

	limit := int(r.Limit)
	return limit * (page - 1)
}

func (r *Query) SetDatabaseQueryForModel(query interface{}, model interface{}) (interface{}, error) {

	modelType := reflect.TypeOf(model).String()

	if modelSearchTagsCache[modelType] == nil {
		err := parseAndCacheModel(model)
		if err != nil {
			return query, fmt.Errorf("query parser: model parse error: %w", err)
		}
	}

	if modelSearchTagsCache[modelType] == nil {
		return query, nil
	}

	modelCfg := modelSearchTagsCache[modelType]
	// each model field:
	for i := range modelCfg {
		p := r.GetParam(modelCfg[i].Param)
		if p == nil {
			continue
		}

		query, _ = GORMDBAdapter.Run(modelCfg[i].Type, p.Operator, p.ParamName, p.Values[0], query, r)
	}

	GORMDBAdapter["pagination"]["pager"]("", "", query, r)

	return query, nil
}

func parseAndCacheModel(model interface{}) error {
	modelCfg := make(map[string]*ModelFieldTagConfig)
	modelType := reflect.TypeOf(model).String()

	ut := reflect.TypeOf(model).Elem()
	for i := 0; i < ut.NumField(); i++ {
		field := ut.Field(i)

		if filterConfig, ok := field.Tag.Lookup("filter"); ok {
			if filterConfig == "-" {
				continue
			}

			cfg := ModelFieldTagConfig{
				// default name is the struct field name:
				Param: field.Name,
				Type:  "default",
			}

			if filterConfig == "" {
				modelCfg[cfg.Param] = &cfg
				continue
			}

			rawTagData := field.Tag.Get("filter")
			tagDataLine := strings.Split(rawTagData, ";")

			for _, v := range tagDataLine {
				tagData := strings.Split(v, ":")

				if tagData[0] == "param" && tagData[1] != "" {
					cfg.Param = tagData[1]
				}

				if tagData[0] == "type" && tagData[1] != "" {
					cfg.Type = tagData[1]
				}
			}

			modelCfg[cfg.Param] = &cfg
		}
	}

	modelSearchTagsCache[modelType] = modelCfg

	return nil
}

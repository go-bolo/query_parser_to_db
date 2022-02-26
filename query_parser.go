package query_parser_to_db

func NewQuery(limitMax int64) QueryInterface {
	q := Query{
		LimitMax: limitMax,
	}

	return &q
}

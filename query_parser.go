package query_parser

func NewQuery(limitMax int64) QueryInterface {
	q := Query{
		LimitMax: limitMax,
	}

	return &q
}

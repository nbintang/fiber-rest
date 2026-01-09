package httpx

import "errors"

type PaginationQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (q PaginationQuery) Normalize(defaultLimit, maxLimit int) PaginationQuery {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = defaultLimit
	}

	if maxLimit > 0 && q.Limit > maxLimit {
		q.Limit = maxLimit
	}
	return q
}

func (q PaginationQuery) Validate() error {
	if q.Page <= 0 {
		return errors.New("page must be greater than 0")
	}
	if q.Limit <= 0 {
		return errors.New("limit must be greater than 0")
	}
	return nil
}

func (q PaginationQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}

package models

import "math"

var (
	// DefaultPerPage is the default number of results to return in a page result set
	DefaultPerPage = 25
)

// Pagination is used to hold pagination values of a query params
type Pagination struct {
	Page    int `query:"page" json:"page"`
	PerPage int `query:"per_page" json:"per_page"`
}

func NewPagination() *Pagination {
	return &Pagination{Page: 1, PerPage: DefaultPerPage}
}

// ApplyLimits sets pagination limit values
func (q *Pagination) ApplyLimits() {
	// min value allowed 1 and max 100
	q.PerPage = int(math.Max(math.Min(float64(q.PerPage), 100), 1))
	// min value allowed 1
	q.Page = int(math.Max(1, float64(q.Page)))
}

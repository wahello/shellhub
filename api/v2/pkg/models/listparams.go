package models

// ListParams holds pagination, filtering and sorting data used for list queries
type ListParams struct {
	Pagination `json:",inline"`
	Filters    FilterList `json:"filters" query:"filters"`
}

// IsValid check if is valid, including its filters
func (l *ListParams) IsValid() error {
	for _, filter := range l.Filters {
		if err := filter.IsValid(); err != nil {
			return err
		}
	}

	return nil
}

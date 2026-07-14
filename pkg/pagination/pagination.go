package pagination

type Params struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func DefaultParams() Params {
	return Params{Page: 1, PerPage: 25}
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func (p Params) Limit() int {
	if p.PerPage < 1 || p.PerPage > 100 {
		return 25
	}
	return p.PerPage
}

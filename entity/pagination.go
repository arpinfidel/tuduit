package entity

import "github.com/arpinfidel/tuduit/pkg/db"

type Pagination struct {
	Page     int    `form:"page"      json:"page"      validate:"gt=0"` // 1-based
	PageSize int    `form:"page_size" json:"page_size" validate:"gt=0"`
	Sort     string `form:"sort"      json:"sort"     `
	SortDesc bool   `form:"sort_desc" json:"sort_desc"`

	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

func (p *Pagination) SetDefault() Pagination {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	return *p
}

func (p *Pagination) Limit() int {
	return p.PageSize
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) SetTotal(totalData int) Pagination {
	p.TotalData = totalData
	p.TotalPage = (totalData + p.PageSize - 1) / p.PageSize
	return *p
}

func (p *Pagination) QBPaginate() *db.Pagination {
	return &db.Pagination{
		Limit:  p.PageSize,
		Offset: p.Offset(),
	}
}

func (p *Pagination) QBSort() []db.Sort {
	if p.Sort == "" {
		return nil
	}
	return []db.Sort{
		{
			Field: p.Sort,
			Asc:   !p.SortDesc,
		},
	}
}

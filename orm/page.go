package orm

import (
	"gorm.io/gorm"
)

type Page struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
}

func (p *Page) Apply(db *gorm.DB) {
	if p.PageSize > 100 {
		p.PageSize = 20
	}
	*db = *db.Offset((p.PageNo - 1) * p.PageSize)
	*db = *db.Limit(p.PageSize)
}

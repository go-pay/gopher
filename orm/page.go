package orm

import (
	"github.com/jinzhu/gorm"
	"xorm.io/xorm"
)

type Page struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
}

func (p *Page) ApplyGorm(db *gorm.DB) {
	if p.PageSize > 100 {
		p.PageSize = 20
	}
	*db = *db.Offset((p.PageNo - 1) * p.PageSize)
	*db = *db.Limit(p.PageSize)
}

func (p *Page) ApplyXorm(db *xorm.Session) {
	if p.PageSize > 100 {
		p.PageSize = 20
	}
	*db = *db.Limit(p.PageSize, (p.PageNo-1)*p.PageSize)
}

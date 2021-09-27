package orm

import "xorm.io/xorm"

type Page struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
}

func (p *Page) Apply(db *xorm.Session) {
	if p.PageSize > 100 {
		p.PageSize = 20
	}
	*db = *db.Limit(p.PageSize, (p.PageNo-1)*p.PageSize)
}

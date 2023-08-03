package shared

import (
	"gorm.io/gorm"
	"xmlt/internal/shared/enum"
)

// Page 分页查询 辅助结构
type Page struct {
	Page     int   // 当前页
	PageSize int   // 单页数量
	Total    int64 // 数据总量
	Pages    int64 // 总页数
}

func NewPage(page int, pageSize int) *Page {
	p := &Page{Page: page, PageSize: pageSize}
	// 过滤 当前页、单页数量； 计算总页数
	if p.Page < 1 {
		p.Page = 1
	}
	switch {
	case p.PageSize > 100:
		p.PageSize = enum.MaxPageSize
	case p.PageSize <= 0:
		p.PageSize = enum.MinPageSize
	}
	return p
}

// Paginate 分页
func (p *Page) Paginate(table interface{}) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		// TODO 后续添加匹配规则，目前是每次查询都会下发。
		// 拼接Count
		countDb := d.Session(&gorm.Session{NewDB: true})
		countDb.Model(table).Count(&p.Total)

		// 计算总页数 Total / PageSize = Pages
		p.Pages = p.Total / int64(p.PageSize)
		// 如果还有余那么也可以查询
		if p.Total%int64(p.PageSize) != 0 {
			p.Pages++
		}
		// 拼接分页
		d.Offset(p.offset()).Limit(p.PageSize)
		return d
	}
}

func (p *Page) offset() int {
	return (p.Page - 1) * p.PageSize
}

// RangeBy ZSet 中 Score 的范围
type RangeBy struct {
	Start int64
	Stop  int64
	Order uint8 // 0 从小到大，1 从大到小
}

func NewRange(p *Page) *RangeBy {
	r := RangeBy{}
	r.Start = int64((p.Page - 1) * (p.PageSize - 1))
	r.Stop = int64(p.PageSize - 1)
	return &r
}

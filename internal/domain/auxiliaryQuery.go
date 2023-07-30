package domain

import (
	"gorm.io/gorm"
)

// Page 分页查询 辅助结构
type Page struct {
	Page     int   // 当前页
	PageSize int   // 单页数量
	Total    int64 // 数据总量
	Pages    int64 // 总页数
	//TableName   string // 本次查询表名
}

func NewPage(page int, pageSize int) *Page {
	return &Page{
		Page:     page,
		PageSize: pageSize,
	}
}

// Paginate 分页
func (p *Page) Paginate() func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		// 过滤 当前页、单页数量； 计算总页数
		if p.Page < 0 {
			p.Page = 0
		}

		switch {
		case p.PageSize > 100:
			p.PageSize = 100
		case p.PageSize <= 0:
			p.PageSize = 10
		}

		// 拼接分页
		d.Offset(p.offset()).Limit(p.PageSize)
		return d
	}
}

// GetPagingStruct 获取分页构造
func (p *Page) GetPagingStruct() func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		// 拼接Count
		d.Count(&p.Total)
		// 计算总页数 Total / PageSize = Pages
		p.Pages = p.Total / int64(p.PageSize)
		// 如果还有余那么也可以查询
		if p.Total%int64(p.PageSize) != 0 {
			p.Pages++
		}
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

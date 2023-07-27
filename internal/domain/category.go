package domain

type Category struct {
	ID           uint64
	Name         string
	Description  string
	ArticleCount uint64
	State        bool
	Ctime        int64
	Utime        int64
}

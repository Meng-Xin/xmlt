package domain

type Category struct {
	ID           uint64
	Name         string
	Description  string
	ArticleCount uint64
	State        bool
	Articles     []Article
	Ctime        int64
	Utime        int64
}

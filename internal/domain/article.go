package domain

import "time"

type Article struct {
	ID      uint64
	Title   string
	Content string
	Ctime   time.Time
	Utime   time.Time

	Author Author
}

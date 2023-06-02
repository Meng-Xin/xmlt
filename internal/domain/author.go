package domain

// Author 是一个值对象
// 在用户的领域内应该是叫做 User
// 但是在写帖子这里，没有用户这个概念
// 只有对应的 Author 的概念
// 本质上是 User 转化过来的一个值对象
type Author struct {
	ID   uint64
	Name string
}

package model

type Post struct {
	ID      int `gorm:"primaryKey"`
	Title   string
	Content string
	Tags    []*Tag `gorm:"many2many:post_tags;"`
}

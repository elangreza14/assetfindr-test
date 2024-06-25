package model

type Tag struct {
	ID    int     `gorm:"primaryKey"`
	Label string  `gorm:"uniqueIndex:idx_label_tag"`
	Posts []*Post `gorm:"many2many:post_tags;"`
}

package model

import "gorm.io/gorm"

type Post struct {
	ID      int `gorm:"primaryKey"`
	Title   string
	Content string
	Tags    []*Tag `gorm:"many2many:post_tags;"`
}

func (p *Post) BeforeDelete(tx *gorm.DB) error {
	err := tx.Exec(`DELETE FROM post_tags WHERE post_id=$1;`, p.ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *Post) BeforeSave(tx *gorm.DB) error {
	for i, tag := range p.Tags {
		reqTag := Tag{}
		err := tx.Where(Tag{Label: tag.Label}).FirstOrCreate(&reqTag).Error
		if err != nil {
			return err
		}

		p.Tags[i] = &reqTag
	}

	return nil
}

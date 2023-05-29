package models

type Tag struct {
	ID       int        `gorm:"primaryKey" json:"id"`
	Title    string     `gorm:"size:50;unique" json:"title"`
	Articles []*Article `gorm:"many2many:articles_tag" json:"articles"`
}

// CreateNewTags is function to create new tag.
func CreateNewTags(tag *Tag) error {
	return db.Create(tag).Error
}

// GetTagByID is function to get tag by id.
func GetTagByID(id int) (Tag, error) {
	var tag Tag
	err := db.Model(&Tag{}).Where("id = ?", id).Preload("Articles").First(&tag).Error
	return tag, err
}

// GetTagByTitle
func GetTagByTitle(title string) (Tag, error) {
	var tag Tag
	err := db.Model(&Tag{}).Where("title = ?", title).Preload("Articles").First(&tag).Error
	return tag, err
}

// GetAllTags is function to get all tags.
func GetAllTags() []Tag {
	var tags []Tag
	db.Model(&Tag{}).Preload("Articles").Find(&tags)
	return tags
}

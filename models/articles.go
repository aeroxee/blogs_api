package models

import "time"

type Article struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Tags        []*Tag    `gorm:"many2many:articles_tag" json:"tags"`
	UserID      int       `json:"user_id"`
	Title       string    `gorm:"size:50;unique" json:"title"`
	Slug        string    `gorm:"size:60;unique" json:"slug"`
	Logo        string    `gorm:"size:255;null" json:"logo"`
	Description string    `gorm:"size:255" json:"description"`
	Content     string    `json:"content"`
	Views       int       `gorm:"default:0" json:"views"`
	Status      string    `gorm:"default:DRAFTED;size:9" json:"status"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateNewArticle is function to create new article.
func CreateNewArticle(article *Article) error {
	return db.Create(article).Error
}

// GetArticleByID is function to get article by id.
func GetArticleByID(id int) (Article, error) {
	var article Article
	err := db.Model(&Article{}).Where("id = ?", id).Preload("Tags").First(&article).Error
	return article, err
}

// GetArticleBySlug is function to get article by slug.
func GetArticleBySlug(slug string) (Article, error) {
	var article Article
	err := db.Model(&Article{}).Where("slug = ?", slug).Preload("Tags").First(&article).Error
	return article, err
}

// GetAllArticles is function to get all article.
func GetAllArticles(offset, limit int, sorted, status string) []Article {
	var articles []Article
	db.Model(&Article{}).Offset(offset).Limit(limit).Preload("Tags").
		Where("status = ?", status).Find(&articles)
	return articles
}

// GetArticleFilterBySlug is function to get article by filtering title.
func GetArticleFilterBySlug(slug string) []Article {
	var articles []Article
	db.Where("slug LIKE ?", "%"+slug+"%").Preload("Tags").Find(&articles)
	return articles
}

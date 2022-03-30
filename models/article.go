package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	Model
	TagID      int    `json:"tag_id" gorm:"index"`
	Tag        Tag    `json:"tag"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func (a *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (a *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}

func GetActile(articleId int) (acticle Article) {
	db.Where("id = ?", articleId).First(&acticle)
	db.Model(&acticle).Related(&acticle.Tag)
	return
}

func GetActiles(pageNum int, pageSize int, condition interface{}) (articles []Article) {
	db.Preload("Tag").Where(condition).Offset(pageNum).Limit(pageSize).Find(&articles)
	return
}

func GetArticlesTotal(condition interface{}) (total int) {
	db.Model(&Article{}).Where(condition).Count(&total)
	return
}

func AddArticle(condition map[string]interface{}) bool {
	db.Create(&Article{
		TagID:     condition["tag_id"].(int),
		Title:     condition["title"].(string),
		Desc:      condition["desc"].(string),
		Content:   condition["content"].(string),
		CreatedBy: condition["created_by"].(string),
		State:     condition["state"].(int),
	})
	return true
}

func EditArticle(id int, condition interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Update(condition)
	return true
}

func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(&Article{})
	return true
}

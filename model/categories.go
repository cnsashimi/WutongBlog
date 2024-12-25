package model

import (
	"imzixuan/config"
)

type (
	CategoriesModel struct {
		ID    int64  `gorm:"column:id" db:"id" json:"id" form:"id" `
		Name  string `gorm:"column:name" db:"name" json:"name" form:"name"`
		Count int64  `gorm:"column:count"  `
	}
)

func (CategoriesModel) TableName() string {
	return "blog.categories"
}
func CountCategories() int64 {
	var CategoriesModel CategoriesModel
	var count int64
	config.Connect.Model(&CategoriesModel).Count(&count)
	return count
}

func DeleteCategories(where string, key string) bool {
	var CategoriesModel CategoriesModel

	config.Connect.Where(where, key).Delete(&CategoriesModel)
	return true
}
func GetCategories(PageNo int, PageSize int) []CategoriesModel {
	var CategoriesModels []CategoriesModel

	offset := (PageNo - 1) * PageSize

	config.Connect.Order("id desc").Limit(PageSize).Offset(offset).Find(&CategoriesModels)

	return CategoriesModels
}

func GetAllCates() []CategoriesModel {
	var CategoriesModels []CategoriesModel
	config.Connect.Select("categories.id,categories.name, (select count(*) from blog.blog where cid = categories.id) as count").Find(&CategoriesModels)
	return CategoriesModels
}
func GetCategorie(id string) (CategoriesModel, error) {
	var CategoriesModel CategoriesModel
	err := config.Connect.First(&CategoriesModel, id)
	if err.Error != nil {
		return CategoriesModel, err.Error
	}
	return CategoriesModel, nil
}
func Updatecate(cate CategoriesModel) bool {
	if cate.Name == "" {
		return false
	}

	config.Connect.Select("Name").Save(&cate)
	return true

}

func Insertcate(cate CategoriesModel) bool {
	if cate.Name == "" {
		return false
	}
	config.Connect.Select("Name").Create(&cate)
	return true
}

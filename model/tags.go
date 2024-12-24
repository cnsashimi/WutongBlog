package model

import (
	"imzixuan/config"
)

type (
	TagsModel struct {
		ID   int64  `gorm:"column:id" db:"id" json:"id" form:"id" `
		Name string `gorm:"column:name" db:"name" json:"name" form:"name"`
	}
)

func (TagsModel) TableName() string {
	return "blog.tags"
}
func GetAllTags() []TagsModel {
	var TagsModels []TagsModel
	config.Connect.Find(&TagsModels)
	return TagsModels
}
func CountTags() int64 {
	var TagsModel TagsModel
	var count int64
	config.Connect.Model(&TagsModel).Count(&count)
	return count
}

func GetTagss(PageNo int, PageSize int) []TagsModel {
	var TagsModels []TagsModel

	offset := (PageNo - 1) * PageSize

	config.Connect.Order("id desc").Limit(PageSize).Offset(offset).Find(&TagsModels)

	return TagsModels
}

func DeleteTags(where string, key string) bool {
	var TagsModel TagsModel

	config.Connect.Where(where, key).Delete(&TagsModel)
	return true
}
func GetTagsbyid(ids []int64) []TagsModel {
	var TagsModels []TagsModel
	config.Connect.Where("id IN ?", ids).Find(&TagsModels)
	return TagsModels
}

func Gettags(id string) (TagsModel, error) {
	var TagsModel TagsModel
	err := config.Connect.First(&TagsModel, id)

	if err.Error != nil {
		return TagsModel, err.Error
	}
	return TagsModel, nil
}
func Updatetags(tagsm TagsModel) bool {
	if tagsm.Name == "" {
		return false
	}

	config.Connect.Select("Name").Save(&tagsm)
	return true

}

func GetTagsbyname(tagname string) TagsModel {
	var TagsModel TagsModel
	config.Connect.Where("name like ?", tagname).First(&TagsModel)
	return TagsModel
}
func GetOrCreateTag(tagname string) int64 {
	var TagsModel TagsModel

	result := config.Connect.Where("name like ?", tagname).First(&TagsModel)

	if result.Error != nil {

		TagsModel.Name = tagname
		config.Connect.Create(&TagsModel)
	}

	return TagsModel.ID
}

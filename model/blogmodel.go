package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"imzixuan/config"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	BlogModel struct {
		ID         int64       `gorm:"column:id" db:"id" json:"id" form:"id" `
		Title      string      `gorm:"column:title" db:"title" json:"title" form:"title" validate:"required,min=1,max=50"`
		Text       string      `gorm:"column:text" db:"text" json:"text" form:"text" validate:"required"`
		Addtime    time.Time   `gorm:"column:addtime" db:"addtime" json:"addtime" form:"addtime" validate:"required"`
		Tags       string      `gorm:"column:tags" db:"tags" json:"tags" form:"tags"`
		Tagss      []TagsModel `json:"tagss" gorm:"-"`
		Totop      int         `gorm:"column:totop" db:"totop" json:"totop" form:"totop" `
		By         int64       `gorm:"column:by" db:"by" json:"by" form:"by"`
		Cid        int64       `gorm:"column:cid" db:"cid" json:"cid" form:"cid"`
		Categories string      `gorm:"column:categories"`
		Images     string      `gorm:"column:images" db:"images" json:"images" form:"images"`

		Imagess  []string `json:"imagess" gorm:"-"`
		Textmini string   `json:"textmini" gorm:"-"`

		Upimgs string `gorm:"-"`

		AddtimeFormat string `gorm:"-"`
		Author        string `gorm:"-"`
	}

	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}
	XValidator struct {
		validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

func (BlogModel) TableName() string {
	return "blog.blog"
}
func Addblog(blogmodel BlogModel) string {
	blogmodel.Addtime = time.Now()
	error := ""

	myValidator := &XValidator{
		validator: validate,
	}
	if errs := myValidator.Validate(blogmodel); len(errs) > 0 && errs[0].Error {
		//errMsgs := make([]string, 0)

		for _, err := range errs {
			//errMsgs = append(errMsgs, fmt.Sprintf(
			//	"[%s]: '%v' | Needs to implement '%s'",
			//	err.FailedField,
			//	err.Value,
			//	err.Tag,
			//))

			error = error + fmt.Sprintf(
				"[%s] ",
				err.FailedField,
				err.Value,
				err.Tag,
			)
		}

	}
	fmt.Println("error::", error)

	if error == "" {

		config.Connect.Select("Title", "Text", "Addtime", "Tags", "By", "Cid", "Images", "Totop").Create(&blogmodel)
	}

	return error

}
func Updateblog(blogmodel BlogModel, id string) string {
	blogmodel.Addtime = time.Now()
	error := ""

	myValidator := &XValidator{
		validator: validate,
	}
	if errs := myValidator.Validate(blogmodel); len(errs) > 0 && errs[0].Error {
		//errMsgs := make([]string, 0)

		for _, err := range errs {
			//errMsgs = append(errMsgs, fmt.Sprintf(
			//	"[%s]: '%v' | Needs to implement '%s'",
			//	err.FailedField,
			//	err.Value,
			//	err.Tag,
			//))

			error = error + fmt.Sprintf(
				"[%s] ",
				err.FailedField,
				err.Value,
				err.Tag,
			)
		}

	}
	fmt.Println("error::", error)
	if error == "" {
		blogmodel.ID, _ = strconv.ParseInt(id, 10, 64)
		config.Connect.Where("id = ?", id).Select("Title", "Text", "Addtime", "Tags", "By", "Cid", "Images", "Totop").Save(&blogmodel)
	}
	return error

}
func GetBlog(id string) (BlogModel, error) {
	var blogmodel BlogModel
	dberr := config.Connect.Select("blog.*, categories.name as categories").Joins("left join blog.categories on categories.id = blog.cid").First(&blogmodel, id)
	if dberr.Error != nil {
		return blogmodel, dberr.Error
	}

	err := json.Unmarshal([]byte(blogmodel.Images), &blogmodel.Imagess)
	if err != nil {
		//fmt.Println("Error JSON: %v", err)
	}

	var tagids []int64
	err = json.Unmarshal([]byte(blogmodel.Tags), &tagids)
	if err != nil {
		//fmt.Println("Error Tags JSON: %v", err)
	} else {
		blogmodel.Tagss = GetTagsbyid(tagids)
	}

	blogmodel.Author = GetUsername(blogmodel.By)
	blogmodel.AddtimeFormat = blogmodel.Addtime.Format("2006-01-02 15:04:05") // 你可以根据需要更改格式

	return blogmodel, nil
}
func GetBlogs(PageNo int, PageSize int, where string, value string) []BlogModel {
	var BlogModel []BlogModel

	offset := (PageNo - 1) * PageSize

	query := config.Connect.Model(&BlogModel)
	if where != "" {
		query = query.Where(where, value)
	}

	query.Order("blog.totop DESC,blog.id desc").Select("blog.*, categories.name as categories").Limit(PageSize).Joins("left join blog.categories on categories.id = blog.cid").Offset(offset).Find(&BlogModel)

	for i, _ := range BlogModel {

		err := json.Unmarshal([]byte(BlogModel[i].Images), &BlogModel[i].Imagess)
		if err != nil {
			fmt.Println("Error Images JSON: %v", err)
		}

		var tagids []int64
		err = json.Unmarshal([]byte(BlogModel[i].Tags), &tagids)
		if err != nil {
			fmt.Println("Error Tags JSON: %v", err)
		} else {
			BlogModel[i].Tagss = GetTagsbyid(tagids)
		}

		BlogModel[i].Textmini = truncateChineseString(trimHtml(BlogModel[i].Text), 30)
		BlogModel[i].Author = GetUsername(BlogModel[i].By)
		BlogModel[i].AddtimeFormat = BlogModel[i].Addtime.Format("2006-01-02 15:04:05") // 你可以根据需要更改格式

	}

	return BlogModel
}
func GetBlogTop3() []BlogModel {
	var BlogModel []BlogModel
	config.Connect.Order("id desc").Where("JSON_LENGTH(images) > 0  ").Select("blog.images,blog.id ,blog.addtime,blog.title, categories.name as categories").Limit(3).Joins("left join blog.categories on categories.id = blog.cid").Offset(0).Find(&BlogModel)
	for i, _ := range BlogModel {

		err := json.Unmarshal([]byte(BlogModel[i].Images), &BlogModel[i].Imagess)
		if err != nil {
			fmt.Println("Error Images JSON: %v", err)
		}

		BlogModel[i].AddtimeFormat = BlogModel[i].Addtime.Format("2006-01-02 15:04:05") // 你可以根据需要更改格式
	}
	return BlogModel
}

func Deleteblog(where string, key string) bool {
	var blogmodel BlogModel
	config.Connect.Where(where, key).Delete(&blogmodel)
	return true
}
func CountBlogs(where string, value string) int64 {
	var blog BlogModel
	var count int64
	query := config.Connect.Model(&blog)
	if where != "" {
		query = query.Where(where, value)
	}
	query.Count(&count)
	return count
}

var validate = validator.New()

func (v XValidator) Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
func isChineseRune(r rune) bool {
	return r >= '\u4e00' && r <= '\u9fff' // 基本汉字范围（包括扩展A、B、C、D、E、F和补充汉字）
	// 注意：这里只检查了基本汉字范围，如果需要包括更多汉字范围（如罕见字、异体字等），
	// 需要扩展这个范围。例如，可以添加'\u3400'到'\u4DBF'（CJK统一表意文字扩展A）等范围。
}
func truncateChineseString(s string, size int) string {
	var count int
	var runes []rune
	for _, r := range s {
		if isChineseRune(r) {
			count++
		}
		runes = append(runes, r)
		if count > size {

			return string(runes[:len(runes)-1]) + "..." // 去掉最后一个字符（可能是非中文字符），然后添加"..."
		}
	}
	return string(runes) // 没有超过size个中文字符，则返回原始字符串
}
func trimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

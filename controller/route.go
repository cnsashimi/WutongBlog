package controller

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/mojocn/base64Captcha"
	"html/template"
	"imzixuan/config"
	"imzixuan/model"
	"imzixuan/task"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type Loginfrom struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}
type Jsonback struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
}
type UpfileJsonback struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}
type Page struct {
	PageNo    int   `json:"page_no"`
	ItemSize  int64 `json:"item_size"`
	PageCount int64 `json:"page_count"`
	PageSize  int   `json:"page_size"`
}
type TableCols struct {
	Title string `json:"title"`
	Field string `json:"field"`
	Width int    `json:"width"`
}

func InitController() {

	engine := html.New("./views", ".html")
	engine.AddFunc(
		// 定义一个  raw 做模板代码不转义使用
		"raw", func(s string) template.HTML {
			return template.HTML(s)
		},
	)
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	store := session.New()
	//TODO::public 根目录
	app.Static("/", "./public")

	app.Get("/test", func(c *fiber.Ctx) error {
		jsonsrt := `["asdfasdf","d800","ddddd"]`

		return c.Render("test", fiber.Map{
			"jsonsrt": jsonsrt,
		})
	})

	//TODO::首页
	app.Get("/", func(c *fiber.Ctx) error {

		count := model.CountBlogs("", "")
		var err error
		var page Page
		page.PageSize = 10
		page.ItemSize = count
		page.PageCount = int64(math.Ceil(float64(count / int64(page.PageSize))))
		page.PageNo, err = strconv.Atoi(c.Query("page", "1"))
		if err != nil {
			page.PageNo = 1
		}
		if page.PageNo == 0 {
			page.PageNo = 1
		}

		blogs := model.GetBlogs(page.PageNo, page.PageSize, "", "")

		categories := config.RedisGet("categories")
		if categories == "" {
			categories = task.Cachecategories()
		}

		var categorieslist []model.CategoriesModel
		err = json.Unmarshal([]byte(categories), &categorieslist)
		if err != nil {
			fmt.Println("categories转换出错：", err)
		}

		tags := config.RedisGet("tags")
		if tags == "" {
			tags = task.CacheTags()
		}

		var taglist []model.TagsModel
		err = json.Unmarshal([]byte(tags), &taglist)
		if err != nil {
			fmt.Println("tags转换出错：", err)
		}

		top3blog := config.RedisGet("top3blog")
		if top3blog == "" {
			top3blog = task.CacheTopblog()
		}

		var top3bloglist []model.BlogModel
		err = json.Unmarshal([]byte(top3blog), &top3bloglist)
		if err != nil {
			fmt.Println("top3blog转换出错：", err)
		}

		return c.Render("index", fiber.Map{
			"setting":    config.Getsetting(),
			"page":       page,
			"blogs":      blogs,
			"categories": categorieslist,
			"tags":       taglist,
			"top3blog":   top3bloglist,
		})

	})

	//TODO::类别列表
	app.Get("/cate/:id<int>/", func(c *fiber.Ctx) error {

		cate, err := model.GetCategorie(c.Params("id", "0"))

		if err != nil {
			return c.Status(404).Render("404", fiber.Map{})
		}

		count := model.CountBlogs("cid = ?", c.Params("id", "0"))

		var page Page
		page.PageSize = 10
		page.ItemSize = count
		page.PageCount = int64(math.Ceil(float64(count / int64(page.PageSize))))
		page.PageNo, err = strconv.Atoi(c.Query("page", "1"))
		if err != nil {
			page.PageNo = 1
		}
		if page.PageNo <= 0 {
			page.PageNo = 1
		}

		blogs := model.GetBlogs(page.PageNo, page.PageSize, "cid = ?", c.Params("id", "0"))

		return c.Render("list", fiber.Map{
			"setting": config.Getsetting(),
			"title":   "类型:" + cate.Name,
			"page":    page,
			"blogs":   blogs,
		})

	})
	//TODO::tags列表
	app.Get("/tags/:id<int>/", func(c *fiber.Ctx) error {

		tags, err := model.Gettags(c.Params("id", "0"))

		if err != nil {
			return c.Status(404).Render("404", fiber.Map{})
		}

		count := model.CountBlogs("JSON_CONTAINS(blog.tags->'$[*]' ,?)", c.Params("id", "0"))

		var page Page
		page.PageSize = 10
		page.ItemSize = count
		page.PageCount = int64(math.Ceil(float64(count / int64(page.PageSize))))
		page.PageNo, err = strconv.Atoi(c.Query("page", "1"))
		if err != nil {
			page.PageNo = 1
		}
		if page.PageNo <= 0 {
			page.PageNo = 1
		}

		blogs := model.GetBlogs(page.PageNo, page.PageSize, "JSON_CONTAINS(blog.tags->'$[*]' ,?)", c.Params("id", "0"))

		return c.Render("list", fiber.Map{
			"setting": config.Getsetting(),
			"title":   "标签:" + tags.Name,
			"page":    page,
			"blogs":   blogs,
		})

	})
	//TODO::作者列表
	app.Get("/author/:id<int>/", func(c *fiber.Ctx) error {

		user, err := model.GetUserbyid(c.Params("id", "0"))

		if err != nil {
			return c.Status(404).Render("404", fiber.Map{})
		}

		count := model.CountBlogs("blog.by = ?", c.Params("id", "0"))

		var page Page
		page.PageSize = 10
		page.ItemSize = count
		page.PageCount = int64(math.Ceil(float64(count / int64(page.PageSize))))
		page.PageNo, err = strconv.Atoi(c.Query("page", "1"))
		if err != nil {
			page.PageNo = 1
		}
		if page.PageNo <= 0 {
			page.PageNo = 1
		}

		blogs := model.GetBlogs(page.PageNo, page.PageSize, "blog.by = ?", c.Params("id", "0"))

		return c.Render("author", fiber.Map{
			"setting": config.Getsetting(),
			"title":   "作者：" + user.Nickname + "文章列表",
			"user":    user,
			"page":    page,
			"blogs":   blogs,
		})

	})
	//TODO::搜索列表
	app.Get("/search", func(c *fiber.Ctx) error {

		count := model.CountBlogs("title like ?", "%"+c.Query("s", "0")+"%")
		var err error
		var page Page
		page.PageSize = 10
		page.ItemSize = count
		page.PageCount = int64(math.Ceil(float64(count / int64(page.PageSize))))
		page.PageNo, err = strconv.Atoi(c.Query("page", "1"))
		if err != nil {
			page.PageNo = 1
		}
		if page.PageNo <= 0 {
			page.PageNo = 1
		}

		blogs := model.GetBlogs(page.PageNo, page.PageSize, "title like ?", "%"+c.Query("s", "0")+"%")

		return c.Render("list", fiber.Map{
			"setting": config.Getsetting(),
			"title":   "搜索关键字:" + c.Query("s", "0"),
			"page":    page,
			"blogs":   blogs,
		})

	})

	//TODO::about
	app.Get("/about", func(c *fiber.Ctx) error {

		return c.Render("about", fiber.Map{
			"setting": config.Getsetting(),
		})

	})

	//TODO::blog文章详情页
	app.Get("/detail_:id<int>.html", func(c *fiber.Ctx) error {
		var blog model.BlogModel
		blog, err := model.GetBlog(c.Params("id", "0"))

		if err != nil {
			return c.Status(404).Render("404", fiber.Map{})
		}
		categories := config.RedisGet("categories")

		if categories == "" {
			categories = task.Cachecategories()
		}

		var categorieslist []model.CategoriesModel

		err = json.Unmarshal([]byte(categories), &categorieslist)
		if err != nil {
			fmt.Println("categories转换出错：", err)
		}

		return c.Render("detail", fiber.Map{
			"setting": config.Getsetting(),
			"blog":    blog,
		})

	})

	//TODO::登陆页
	app.Get("/login", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			fmt.Println(err.Error())
		}

		if sess.Get("iflogin") == true {
			return c.Redirect("/admin/")
		}

		return c.Render("login", fiber.Map{
			"setting": config.Getsetting(),
		})

	})
	//TODO::验证码
	app.Get("/captcha", func(c *fiber.Ctx) error {
		var id string
		var b64s string
		var answer string
		var err error
		sess, err := store.Get(c)
		if err != nil {
			fmt.Println(err.Error())
		}

		var store = base64Captcha.DefaultMemStore
		driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
		captcha := base64Captcha.NewCaptcha(driver, store)

		id, b64s, answer, err = captcha.Generate()

		sess.Set("captchaid", id)
		sess.Set("captchaanswer", answer)
		if err := sess.Save(); err != nil {
			return err
		}

		c.Set("Content-Type", "image/png")

		imgData, err := base64.StdEncoding.DecodeString(b64s[strings.Index(b64s, ",")+1:])
		if err != nil {
			fmt.Println("Error decoding base64:", err)
			//return
		}

		_, err = c.Write(imgData)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("发送出错")
		}

		return nil

	})
	//TODO::loginsubmit登陆提交处理
	app.Post("/loginsubmit", func(c *fiber.Ctx) error {
		loginfrom := new(Loginfrom)
		sess, err := store.Get(c)
		if err != nil {
			fmt.Println(err.Error())
		}
		//name := sess.Get("captchaanswer")

		capt := sess.Get("captchaanswer")

		if err := c.BodyParser(loginfrom); err != nil {
			fmt.Println("error = ", err)
			sess.Delete("iflogin")
			sess.Delete("loginname")
			sess.Delete("loginid")
			if err := sess.Save(); err != nil {
				return err
			}
			return c.JSON(Jsonback{
				Code: 1,
				Msg:  "提交出错",
				Data: nil,
			})
		}

		if capt != loginfrom.Captcha {
			sess.Delete("iflogin")
			sess.Delete("loginname")
			sess.Delete("loginid")
			if err := sess.Save(); err != nil {
				return err
			}
			return c.JSON(Jsonback{
				Code: 1,
				Msg:  "验证码出错",
				Data: nil,
			})

		}

		user, err := model.GetUser(loginfrom.Username)
		if err != nil {
			sess.Delete("iflogin")
			sess.Delete("loginname")
			sess.Delete("loginid")
			if err := sess.Save(); err != nil {
				return err
			}
			return c.JSON(Jsonback{
				Code: 1,
				Msg:  "密码或用户名错误",
				Data: nil,
			})
		}

		if user.Password != Mmd5(loginfrom.Password) {
			sess.Delete("iflogin")
			sess.Delete("loginname")
			sess.Delete("loginid")
			if err := sess.Save(); err != nil {
				return err
			}
			return c.JSON(Jsonback{
				Code: 1,
				Msg:  "密码或用户名错误",
				Data: nil,
			})

		}
		sess.Set("iflogin", true)
		sess.Set("loginname", loginfrom.Username)
		sess.Set("loginid", user.Id)

		if err := sess.Save(); err != nil {
			return err
		}

		return c.JSON(Jsonback{
			Code: 0,
			Msg:  "",
			Data: loginfrom,
		})

	})

	//admin-------------------
	RouteAdmin(app, store)
	//TODO:: 错误页 一定放在最后面
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("404", fiber.Map{})
	})

	app.Listen(":3000")

}

func Mmd5(str string) string {
	password := md5.Sum([]byte(str))
	strString := hex.EncodeToString(password[:])
	return strString
}

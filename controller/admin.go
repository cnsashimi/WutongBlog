package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"imzixuan/config"
	"imzixuan/model"
	"imzixuan/task"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func RouteAdmin(app *fiber.App, store *session.Store) {

	//TODO::后台 - 中间件
	adminMiddleware := func(store *session.Store) fiber.Handler {
		return func(c *fiber.Ctx) error {
			sess, err := store.Get(c)
			if err != nil {
				fmt.Println(err.Error())
			}
			if sess.Get("iflogin") != true {
				return c.Redirect("/login")
			}
			return c.Next()
		}
	}
	admincontroller := app.Group("/admin", adminMiddleware(store))
	//admincontroller.Use(AdminSecurity)
	admincontroller.Get("/", func(c *fiber.Ctx) error {
		return c.Render("admin/admin", fiber.Map{
			"Title": config.Getsetting().Name,
		})

	})
	//TODO::后台uploadimg - 上传图片
	admincontroller.Post("uploadimg", func(c *fiber.Ctx) error {
		if form, err := c.MultipartForm(); err == nil {
			files := form.File["edit"]
			for _, file := range files {
				fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
				if file.Header["Content-Type"][0] != "image/jpeg" && file.Header["Content-Type"][0] != "image/png" {
					return c.JSON(
						Jsonback{
							Code: 1,
							Msg:  "只能上传Jpg 或 Png 文件",
							Data: nil,
						})

				}

				if file.Size > 1024000 {
					return c.JSON(
						Jsonback{
							Code: 1,
							Msg:  "文件上传不能大于 1M",
							Data: nil,
						})

				}

				now := time.Now()

				fileformat := strings.Split(file.Header["Content-Type"][0], "/")
				filename := strconv.FormatInt(now.UnixNano(), 10) + "." + fileformat[1]
				urlpath := "/uploads/" + now.Format("2006/01/02/")
				// 构造目录和文件名
				dir := filepath.Join("./public/", urlpath)
				filepath := filepath.Join(dir, filename)

				// 确保目录存在
				if err := os.MkdirAll(dir, 0755); err != nil {
					return c.JSON(
						Jsonback{
							Code: 1,
							Msg:  "上传出错:" + err.Error(),
							Data: nil,
						})
				}

				if err := c.SaveFile(file, filepath); err != nil {

					return c.JSON(
						Jsonback{
							Code: 1,
							Msg:  "上传出错:" + err.Error(),
							Data: nil,
						})
				}

				return c.JSON(
					Jsonback{
						Code: 0,
						Msg:  "",
						Data: UpfileJsonback{
							Url:  urlpath + filename,
							Name: file.Filename,
							Size: file.Size,
						},
					})

			}
			return err
		}

		return c.JSON(Jsonback{
			Code: 1,
			Msg:  "上传出错",
			Data: nil,
		})

	})
	//TODO::后台uploadavatar - 上传头像图片
	admincontroller.Post("uploadavatar", func(c *fiber.Ctx) error {

		if form, err := c.MultipartForm(); err == nil {
			files := form.File["avatar"]
			for _, file := range files {
				fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
				if file.Header["Content-Type"][0] != "image/jpeg" && file.Header["Content-Type"][0] != "image/png" {
					return c.JSON(
						Jsonback{
							Code: 1,
							Msg:  "只能上传Jpg 或 Png 文件",
							Data: nil,
						})

				}

				if file.Size > 1024000 {
					return c.JSON(
						Jsonback{
							Code: 1,
							Msg:  "文件上传不能大于 1M",
							Data: nil,
						})

				}

				now := time.Now()

				fileformat := strings.Split(file.Header["Content-Type"][0], "/")
				filename := strconv.FormatInt(now.UnixNano(), 10) + "." + fileformat[1]
				urlpath := "/uploads/avatar/" + now.Format("2006/01/02/")
				// 构造目录和文件名
				dir := filepath.Join("./public/", urlpath)
				filepath := filepath.Join(dir, filename)

				// 确保目录存在
				if err := os.MkdirAll(dir, 0755); err != nil {
					return c.JSON(
						Jsonback{
							Code: 1,
							Msg:  "上传出错:" + err.Error(),
							Data: nil,
						})
				}

				if err := c.SaveFile(file, filepath); err != nil {

					return c.JSON(
						Jsonback{
							Code: 1,
							Msg:  "上传出错:" + err.Error(),
							Data: nil,
						})
				}

				return c.JSON(
					Jsonback{
						Code: 0,
						Msg:  "",
						Data: UpfileJsonback{
							Url:  urlpath + filename,
							Name: file.Filename,
							Size: file.Size,
						},
					})

			}
			return err
		}

		return c.JSON(Jsonback{
			Code: 1,
			Msg:  "上传出错",
			Data: nil,
		})

	})
	//TODO::后台deleteimg - 删除上传的图片
	admincontroller.Post("deleteimg", func(c *fiber.Ctx) error {

		delimg := c.FormValue("delimg")
		delfile := filepath.Join("./public/", delimg)
		fmt.Println("删除文件：", delfile)

		err := os.Remove(delfile)

		if err != nil {
			return c.JSON(
				Jsonback{
					Code: 1,
					Msg:  "删除文件出错 有可能进程未释放 请稍等几秒再重试一下 具体错误信息:" + err.Error() + " ",
					Data: nil,
				})
		}
		return c.JSON(
			Jsonback{
				Code: 0,
				Msg:  "completed",
				Data: nil,
			})

	})
	//TODO::后台welcome - 欢迎页
	admincontroller.Get("welcome", func(c *fiber.Ctx) error {
		return c.Render("admin/welcome", fiber.Map{
			"Title": config.Getsetting().Name,
		})

	})
	//TODO::后台account - 管理帐号设置页
	admincontroller.Get("account", func(c *fiber.Ctx) error {

		sess, err := store.Get(c)
		if err != nil {
			fmt.Println(err.Error())
		}
		if sess.Get("loginname") == nil {
			return c.Redirect("/login")
		}
		loginname := sess.Get("loginname").(string)
		uinfo, err := model.GetUser(loginname)
		if err != nil {
			return c.Redirect("/login")
		}
		return c.Render("admin/account", fiber.Map{
			"Title": config.Getsetting().Name,
			"uinfo": uinfo,
		})

	})
	//TODO::后台accountupdate - 管理帐号设置提交接收
	admincontroller.Post("accountupdate", func(c *fiber.Ctx) error {

		var um model.User
		um.Nickname = c.FormValue("nickname")
		um.Avatar = c.FormValue("avatarfile")
		um.Aboutme = c.FormValue("aboutme")

		sess, err := store.Get(c)
		if err != nil {
			fmt.Println(err.Error())
		}
		if sess.Get("loginname") == nil {
			return c.Redirect("/login")
		}
		loginname := sess.Get("loginname").(string)
		uinfo, err := model.GetUser(loginname)
		if err != nil {
			return c.Redirect("/login")
		}

		um.Id = uinfo.Id
		save := model.UpdateUser(um)
		if save {
			return c.JSON(
				Jsonback{
					Code: 0,
					Msg:  "completed",
					Data: nil,
				})
		} else {

			return c.JSON(
				Jsonback{
					Code: 1,
					Msg:  "error",
					Data: nil,
				})
		}

	})
	//TODO::后台resetpassword - 管理帐号设置重密码提交接收
	admincontroller.Post("resetpassword", func(c *fiber.Ctx) error {

		old_password := c.FormValue("old_password")
		password := c.FormValue("password")
		password_confirm := c.FormValue("password_confirm")
		sess, err := store.Get(c)
		if err != nil {
			fmt.Println(err.Error())
		}
		if sess.Get("loginname") == nil {
			return c.Redirect("/login")
		}
		loginname := sess.Get("loginname").(string)
		uinfo, err := model.GetUser(loginname)
		if err != nil {
			return c.Redirect("/login")
		}
		if uinfo.Password != Mmd5(old_password) {
			return c.JSON(
				Jsonback{
					Code: 1,
					Msg:  "旧密码错误",
					Data: nil,
				})
		}

		if password != password_confirm {
			return c.JSON(
				Jsonback{
					Code: 1,
					Msg:  "两次输入的密码不一致",
					Data: nil,
				})
		}

		var um model.User
		um.Password = Mmd5(password)
		um.Id = uinfo.Id
		save := model.UpdateResetPassword(um)
		if save {
			return c.JSON(
				Jsonback{
					Code: 0,
					Msg:  "completed",
					Data: nil,
				})
		} else {

			return c.JSON(
				Jsonback{
					Code: 1,
					Msg:  "error",
					Data: nil,
				})
		}

	})

	//TODO::后台bloglist - blog列表页
	admincontroller.Get("bloglist", func(c *fiber.Ctx) error {

		return c.Render("admin/list", fiber.Map{
			"Title":      config.Getsetting().Name,
			"SELECT_API": "/admin/blogselect",
			"DELETE_API": "/admin/blogdelete",
			"INSERT_URL": "/admin/bloginsert",
			"UPDATE_URL": "/admin/blogupdate",
			"UPDATE_WINDOWS_SIZE": []string{
				"100%",
				"100%",
			},
			"TableCols": []TableCols{
				{
					Title: "主键",
					Field: "id",
					Width: 70,
				}, {
					Title: "标题",
					Field: "title",
					Width: 170,
				}, {
					Title: "images",
					Field: "images",
					Width: 170,
				}, {
					Title: "正文（缩）",
					Field: "textmini",
					Width: 270,
				}, {
					Title: "创建时间",
					Field: "addtime",
					Width: 270,
				},
			},
		})

	})
	//TODO::后台categorieslist - 文章类别管理
	admincontroller.Get("categorieslist", func(c *fiber.Ctx) error {

		return c.Render("admin/list", fiber.Map{
			"Title":      config.Getsetting().Name,
			"SELECT_API": "/admin/categoriesselect",
			"DELETE_API": "/admin/categoriesdelete",
			"INSERT_URL": "/admin/categoriesinsert",
			"UPDATE_URL": "/admin/categoriesupdate",
			"UPDATE_WINDOWS_SIZE": []string{
				"600px",
				"350px",
			},

			"TableCols": []TableCols{
				{
					Title: "主键",
					Field: "id",
					Width: 70,
				}, {
					Title: "类别名",
					Field: "name",
					Width: 170,
				},
			},
		})

	})
	//TODO::后台categoriesselect - 文章类别ajax数据
	admincontroller.Get("categoriesselect", func(c *fiber.Ctx) error {

		count := model.CountCategories()
		var err error
		var page Page

		page.PageSize, err = strconv.Atoi(c.Query("limit", "10"))
		if err != nil {
			page.PageSize = 10
		}
		if page.PageSize == 0 {
			page.PageNo = 10
		}

		page.ItemSize = count
		page.PageCount = int64(math.Ceil(float64(count / int64(page.PageSize))))
		page.PageNo, err = strconv.Atoi(c.Query("page", "1"))
		if err != nil {
			page.PageNo = 1
		}
		if page.PageNo == 0 {
			page.PageNo = 1
		}

		categoriess := model.GetCategories(page.PageNo, page.PageSize)

		return c.JSON(Jsonback{
			Code:  0,
			Msg:   "",
			Count: count,
			Data:  categoriess,
		})
	})

	//TODO::后台tagslist - 文章标签管理列表
	admincontroller.Get("tagslist", func(c *fiber.Ctx) error {

		return c.Render("admin/list", fiber.Map{
			"Title":      config.Getsetting().Name,
			"SELECT_API": "/admin/tagssselect",
			"DELETE_API": "/admin/tagssdelete",
			"INSERT_URL": nil,
			"UPDATE_URL": "/admin/tagsupdate",
			"UPDATE_WINDOWS_SIZE": []string{
				"600px",
				"350px",
			},
			"TableCols": []TableCols{
				{
					Title: "主键",
					Field: "id",
					Width: 70,
				}, {
					Title: "标签名",
					Field: "name",
					Width: 170,
				},
			},
		})

	})
	//TODO::后台tagssselect - 文章标签ajax数据
	admincontroller.Get("tagssselect", func(c *fiber.Ctx) error {

		count := model.CountTags()
		var err error
		var page Page

		page.PageSize, err = strconv.Atoi(c.Query("limit", "10"))
		if err != nil {
			page.PageSize = 10
		}
		if page.PageSize == 0 {
			page.PageNo = 10
		}

		page.ItemSize = count
		page.PageCount = int64(math.Ceil(float64(count / int64(page.PageSize))))
		page.PageNo, err = strconv.Atoi(c.Query("page", "1"))
		if err != nil {
			page.PageNo = 1
		}
		if page.PageNo == 0 {
			page.PageNo = 1
		}

		tags := model.GetTagss(page.PageNo, page.PageSize)

		return c.JSON(Jsonback{
			Code:  0,
			Msg:   "",
			Count: count,
			Data:  tags,
		})
	})

	//TODO::后台blogselect - blog列表ajax数据
	admincontroller.Get("blogselect", func(c *fiber.Ctx) error {

		count := model.CountBlogs("", "")
		var err error
		var page Page

		page.PageSize, err = strconv.Atoi(c.Query("limit", "10"))
		if err != nil {
			page.PageSize = 10
		}
		if page.PageSize == 0 {
			page.PageNo = 10
		}

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

		return c.JSON(Jsonback{
			Code:  0,
			Msg:   "",
			Count: count,
			Data:  blogs,
		})
	})
	//TODO::后台bloginsert - blog添加
	admincontroller.Get("bloginsert", func(c *fiber.Ctx) error {
		var err error
		categories := config.RedisGet("categories")

		if categories == "" {
			categories = task.Cachecategories()
		}

		var categorieslist []model.CategoriesModel
		err = json.Unmarshal([]byte(categories), &categorieslist)
		if err != nil {
			fmt.Println("categories转换出错：", err)
		}

		return c.Render("admin/bloginsert", fiber.Map{
			"Title":      config.Getsetting().Name,
			"categories": categorieslist,
		})
	})

	//TODO::后台blogsubmit - 添加blog接收
	admincontroller.Post("blogsubmit", func(c *fiber.Ctx) error {
		var blogmodel model.BlogModel

		var arr []string
		err := json.Unmarshal([]byte(c.FormValue("tagss")), &arr)
		if err != nil {
			fmt.Println("Error tagss JSON: %v", err)
		}

		var iarr []int64
		for key := range arr {
			iarr = append(iarr, model.GetOrCreateTag(arr[key]))
		}
		jsonBytes, err := json.Marshal(iarr)

		if err != nil {
			blogmodel.Tags = "[]"
		}
		if string(jsonBytes) == "" {
			blogmodel.Tags = "[]"
		} else {
			blogmodel.Tags = string(jsonBytes)
		}

		//  blogmodel.Totop= (c.FormValue("Totop") == "on") ? 1 :0    //吐槽~！ go没三元.....
		if c.FormValue("Totop") == "on" {
			blogmodel.Totop = 1
		} else {
			blogmodel.Totop = 0
		}

		blogmodel.Title = c.FormValue("Title")
		blogmodel.Images = c.FormValue("Upimgs", "[]")
		blogmodel.Text = c.FormValue("Text")
		blogmodel.Cid, err = strconv.ParseInt(c.FormValue("Cid"), 10, 64)
		if err != nil {
			blogmodel.Cid = 0
		}

		sess, err := store.Get(c)
		if err != nil {
			fmt.Println(err.Error())
		}
		if sess.Get("loginname") == nil {
			return c.Redirect("/login")
		}
		loginname := sess.Get("loginname").(string)
		uinfo, err := model.GetUser(loginname)
		if err != nil {
			return c.Redirect("/login")
		}
		blogmodel.By = uinfo.Id

		error := model.Addblog(blogmodel)
		if error != "" {
			return c.JSON(Jsonback{
				Code: 1,
				Msg:  "提交出错了 ：" + error,
				Data: nil,
			})

		}

		return c.JSON(Jsonback{
			Code: 0,
			Msg:  "提交成功",
			Data: nil,
		})

	})
	//TODO::后台categoriesupdate - 类型修改页
	admincontroller.Get("categoriesupdate", func(c *fiber.Ctx) error {

		categorie, err := model.GetCategorie(c.Query("id", "0"))

		if err != nil {
			return c.Status(404).Render("404", fiber.Map{})
		}
		return c.Render("admin/update", fiber.Map{
			"Title":      config.Getsetting().Name,
			"Name":       categorie.Name,
			"UPDATE_API": "/admin/cateupdatesubmit?id=" + strconv.FormatInt(categorie.ID, 10),
		})
	})
	//TODO::后台categoriesupdate - 类型添加页
	admincontroller.Get("categoriesinsert", func(c *fiber.Ctx) error {

		return c.Render("admin/update", fiber.Map{
			"Title":      config.Getsetting().Name,
			"Name":       "",
			"UPDATE_API": "/admin/cateinsertsubmit",
		})
	})

	//TODO::后台cateupdatesubmit - 类型修改接收
	admincontroller.Post("cateupdatesubmit", func(c *fiber.Ctx) error {
		var catm model.CategoriesModel
		var err error
		catm.Name = c.FormValue("Name")
		catm.ID, err = strconv.ParseInt(c.Query("id", "0"), 10, 64)

		if err != nil {
			return c.JSON(Jsonback{
				Code: 1,
				Msg:  "提交出错",
				Data: nil,
			})
		}
		suss := model.Updatecate(catm)
		if suss {
			return c.JSON(Jsonback{
				Code: 0,
				Msg:  "提交成功",
				Data: nil,
			})

		}

		return c.JSON(Jsonback{
			Code: 1,
			Msg:  "提交出错",
			Data: nil,
		})

	})

	//TODO::后台cateupdatesubmit - 类型添加接收
	admincontroller.Post("cateinsertsubmit", func(c *fiber.Ctx) error {
		var catm model.CategoriesModel
		catm.Name = c.FormValue("Name")
		suss := model.Insertcate(catm)

		if suss {
			return c.JSON(Jsonback{
				Code: 0,
				Msg:  "提交成功",
				Data: nil,
			})

		}

		return c.JSON(Jsonback{
			Code: 1,
			Msg:  "提交出错",
			Data: nil,
		})

	})

	//TODO::后台 - tags修改页
	admincontroller.Get("tagsupdate", func(c *fiber.Ctx) error {

		tagsm, err := model.Gettags(c.Query("id", "0"))

		if err != nil {
			return c.Status(404).Render("404", fiber.Map{})
		}
		return c.Render("admin/update", fiber.Map{
			"Title":      config.Getsetting().Name,
			"Name":       tagsm.Name,
			"UPDATE_API": "/admin/tagsupdatesubmit?id=" + strconv.FormatInt(tagsm.ID, 10),
		})
	})
	//TODO::后台 - tags修改接收
	admincontroller.Post("tagsupdatesubmit", func(c *fiber.Ctx) error {
		var tags model.TagsModel
		var err error
		tags.Name = c.FormValue("Name")
		tags.ID, err = strconv.ParseInt(c.Query("id", "0"), 10, 64)

		if err != nil {
			return c.JSON(Jsonback{
				Code: 1,
				Msg:  "提交出错",
				Data: nil,
			})
		}
		suss := model.Updatetags(tags)
		if suss {
			return c.JSON(Jsonback{
				Code: 0,
				Msg:  "提交成功",
				Data: nil,
			})

		}

		return c.JSON(Jsonback{
			Code: 1,
			Msg:  "提交出错",
			Data: nil,
		})

	})

	//TODO::后台 - blog修改页
	admincontroller.Get("blogupdate", func(c *fiber.Ctx) error {

		blog, err := model.GetBlog(c.Query("id", "0"))

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

		return c.Render("admin/blogupdate", fiber.Map{
			"Title":      config.Getsetting().Name,
			"blog":       blog,
			"categories": categorieslist,
		})
	})
	//TODO::后台 - 编辑blog接收
	admincontroller.Post("blogupdatesubmit", func(c *fiber.Ctx) error {
		var blogmodel model.BlogModel

		var arr []string
		err := json.Unmarshal([]byte(c.FormValue("tagss")), &arr)
		if err != nil {
			fmt.Println("Error tagss JSON: %v", err)
		}

		var iarr []int64
		for key := range arr {
			iarr = append(iarr, model.GetOrCreateTag(arr[key]))
		}
		jsonBytes, err := json.Marshal(iarr)
		if err != nil {
			blogmodel.Tags = "[]"
		}
		if string(jsonBytes) == "" {
			blogmodel.Tags = "[]"
		} else {
			blogmodel.Tags = string(jsonBytes)
		}
		if c.FormValue("Totop") == "on" {
			blogmodel.Totop = 1
		} else {
			blogmodel.Totop = 0
		}

		blogmodel.Title = c.FormValue("Title")
		blogmodel.Text = c.FormValue("Text")
		blogmodel.Images = c.FormValue("Upimgs", "[]")
		blogmodel.Cid, err = strconv.ParseInt(c.FormValue("Cid"), 10, 64)
		if err != nil {
			blogmodel.Cid = 0
		}

		sess, err := store.Get(c)
		if err != nil {
			fmt.Println(err.Error())
		}
		if sess.Get("loginname") == nil {
			return c.Redirect("/login")
		}
		loginname := sess.Get("loginname").(string)
		uinfo, err := model.GetUser(loginname)
		if err != nil {
			return c.Redirect("/login")
		}
		blogmodel.By = uinfo.Id

		error := model.Updateblog(blogmodel, c.Query("id", "0"))
		if error != "" {
			return c.JSON(Jsonback{
				Code: 1,
				Msg:  "提交出错 出错字段：" + error,
				Data: nil,
			})

		}

		return c.JSON(Jsonback{
			Code: 0,
			Msg:  "提交成功",
			Data: nil,
		})

	})
	//TODO::后台 - 删除blog
	admincontroller.Post("blogdelete", func(c *fiber.Ctx) error {
		model.Deleteblog("id= ?", c.FormValue("id"))
		return c.JSON(Jsonback{
			Code: 0,
			Msg:  "提交成功",
			Data: nil,
		})
	})
	//TODO::后台 - 删除类别
	admincontroller.Post("categoriesdelete", func(c *fiber.Ctx) error {
		model.DeleteCategories("id= ?", c.FormValue("id"))
		return c.JSON(Jsonback{
			Code: 0,
			Msg:  "提交成功",
			Data: nil,
		})
	})
	//TODO::后台 - 删除tags
	admincontroller.Post("tagssdelete", func(c *fiber.Ctx) error {
		model.DeleteTags("id= ?", c.FormValue("id"))
		return c.JSON(Jsonback{
			Code: 0,
			Msg:  "提交成功",
			Data: nil,
		})
	})
	//TODO::后台 - 登出
	admincontroller.Get("logout", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			fmt.Println(err.Error())
		}
		sess.Delete("iflogin")
		sess.Delete("loginname")

		return c.Redirect("/")
	})
	//TODO::后台 - 设置
	admincontroller.Get("config", func(c *fiber.Ctx) error {
		return c.JSON(config.Getadminconfig())
	})
	//TODO::后台 - 菜单
	admincontroller.Get("rule", func(c *fiber.Ctx) error {

		return c.JSON(Jsonback{
			Code: 0,
			Msg:  "",
			Data: config.Getrulemenu(),
		})
	})

	return
}

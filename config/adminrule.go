package config

func Getrulemenu() []Rulemenu {

	return []Rulemenu{
		{Id: 1,
			Title:  "文章管理",
			Icon:   "layui-icon layui-icon-flag",
			Key:    "w",
			Pid:    0,
			Href:   "",
			Type:   0,
			Weight: 1000,
			Name:   "文章管理",
			Value:  1,
			Children: []Rulemenu{
				{
					Id:     11,
					Title:  "文章列表",
					Icon:   "layui-icon layui-icon-404",
					Key:    "",
					Pid:    1,
					Href:   "/admin/bloglist",
					Type:   1,
					Weight: 0,
					Name:   "文章列表",
					Value:  11,
				},
				{
					Id:     12,
					Title:  "文章类别管理",
					Icon:   "layui-icon layui-icon-404",
					Key:    "",
					Pid:    1,
					Href:   "/admin/categorieslist",
					Type:   1,
					Weight: 0,
					Name:   "文章类别管理",
					Value:  12,
				},
				{
					Id:     13,
					Title:  "文章标签管理",
					Icon:   "layui-icon layui-icon-404",
					Key:    "",
					Pid:    1,
					Href:   "/admin/tagslist",
					Type:   1,
					Weight: 0,
					Name:   "文章标签管理",
					Value:  13,
				},
			},
		},
		//{Id: 2,
		//	Title:  "系统管理",
		//	Icon:   "layui-icon layui-icon-flag",
		//	Key:    "w",
		//	Pid:    0,
		//	Href:   "",
		//	Type:   0,
		//	Weight: 1000,
		//	Name:   "系统管理",
		//	Value:  2,
		//	Children: []Rulemenu{
		//		{
		//			Id:     21,
		//			Title:  "退出登陆",
		//			Icon:   "layui-icon layui-icon-404",
		//			Key:    "",
		//			Pid:    2,
		//			Href:   "/admin/logout",
		//			Type:   1,
		//			Weight: 0,
		//			Name:   "退出登陆",
		//			Value:  21,
		//		},
		//	},
		//},
	}

}

type Rulemenu struct {
	Id       int        `json:"id"`
	Title    string     `json:"title"`
	Icon     string     `json:"icon"`
	Key      string     `json:"key"`
	Pid      int        `json:"pid"`
	Href     string     `json:"href"`
	Type     int        `json:"type"`
	Weight   int        `json:"weight"`
	Name     string     `json:"name"`
	Value    int        `json:"value"`
	Children []Rulemenu `json:"children"`
}

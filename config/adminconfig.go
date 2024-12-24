package config

func Getadminconfig() Adminconfig {

	return Adminconfig{
		Logo: Adminconfiglogo{
			Title: "WtBlog Admin",
			Image: "/images/logo.png",
		},
		Menu: Adminconfigmenu{
			Data:         "/admin/rule",
			Method:       "GET",
			Accordion:    true,
			Collapse:     false,
			Control:      false,
			ControlWidth: 2000,
			Select:       "0",
			Async:        true,
		},
		Tab: Adminconfigtab{
			Enable:    true,
			KeepState: true,
			Session:   true,
			Preload:   false,
			Max:       "30",
			Index: Adminconfigtabindex{
				Id:    "0",
				Href:  "/admin/welcome",
				Title: "欢迎",
			},
		},
		Theme: Adminconfigtheme{
			DefaultColor:  "2",
			DefaultMenu:   "light-theme",
			DefaultHeader: "light-theme",
			AllowCustom:   true,
			Banner:        false,
		},

		Colors: []Adminconfigcolors{
			{
				Id:     "1",
				Color:  "#36b368",
				Second: "#f0f9eb",
			},
			{
				Id:     "2",
				Color:  "#2d8cf0",
				Second: "#ecf5ff",
			},
			{
				Id:     "3",
				Color:  "#f6ad55",
				Second: "#fdf6ec",
			},
			{
				Id:     "4",
				Color:  "#f56c6c",
				Second: "#fef0f0",
			},
			{
				Id:     "5",
				Color:  "#3963bc",
				Second: "#ecf5ff",
			},
		},
		Other: Adminconfigother{
			KeepLoad: "500",
			AutoHead: false,
			Footer:   false,
		},
		Header: Adminconfigheader{Message: false},
	}

}

type Adminconfig struct {
	Logo   Adminconfiglogo     `json:"logo"`
	Menu   Adminconfigmenu     `json:"menu"`
	Tab    Adminconfigtab      `json:"tab"`
	Theme  Adminconfigtheme    `json:"theme"`
	Colors []Adminconfigcolors `json:"colors"`
	Other  Adminconfigother    `json:"other"`
	Header Adminconfigheader   `json:"header"`
}

type Adminconfiglogo struct {
	Title string `json:"title"`
	Image string `json:"image"`
}
type Adminconfigmenu struct {
	Data         string `json:"data"`
	Method       string `json:"method"`
	Accordion    bool   `json:"accordion"`
	Collapse     bool   `json:"collapse"`
	Control      bool   `json:"control"`
	ControlWidth int    `json:"controlWidth"`
	Select       string `json:"select"`
	Async        bool   `json:"async"`
}
type Adminconfigtab struct {
	Enable    bool                `json:"enable"`
	KeepState bool                `json:"keepState"`
	Session   bool                `json:"session"`
	Preload   bool                `json:"preload"`
	Max       string              `json:"max"`
	Index     Adminconfigtabindex `json:"index"`
}
type Adminconfigtabindex struct {
	Id    string `json:"id"`
	Href  string `json:"href"`
	Title string `json:"title"`
}
type Adminconfigtheme struct {
	DefaultColor  string `json:"defaultColor"`
	DefaultMenu   string `json:"defaultMenu"`
	DefaultHeader string `json:"defaultHeader"`
	AllowCustom   bool   `json:"allowCustom"`
	Banner        bool   `json:"banner"`
}
type Adminconfigcolors struct {
	Id     string `json:"id"`
	Color  string `json:"color"`
	Second string `json:"second"`
}
type Adminconfigother struct {
	KeepLoad string `json:"keepLoad"`
	AutoHead bool   `json:"autoHead"`
	Footer   bool   `json:"footer"`
}
type Adminconfigheader struct {
	Message bool `json:"message"`
}

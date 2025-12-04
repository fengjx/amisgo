package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/fengjx/amisgo"
	"github.com/fengjx/amisgo/internal/kit"

	_ "github.com/go-sql-driver/mysql"
)

// corsMiddleware 处理跨域请求
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// 设置允许的域名，这里使用 * 表示允许所有域名
		w.Header().Set("Access-Control-Allow-Origin", origin)
		// 设置允许的请求方法
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		// 设置允许的请求头
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, App, X-Requested-With, X-Admin-Token")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Authorization, Content-Disposition, Server, X-Refresh-Token, X-Rsp-Meta")
		// 设置预检请求的缓存时间
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getExampleDB() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.1.121:3306)/lucky?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db := getExampleDB()
	defer kit.QuicklyClose(db)

	amisgo.UseDefaultExecutor(db)
	amisgo.UseOmitsEditColumns("ctime", "utime")

	router := amisgo.NewAdminRouter()
	err := router.SetMenu([]*amisgo.Menu{
		{
			Label:    "Home",
			URL:      "/",
			Redirect: "/sys",
			Visible:  true,
		},
		{
			Label:   "系统",
			Visible: true,
			Children: []*amisgo.Menu{
				{
					Label:   "系统管理",
					Icon:    "fa fa-wrench",
					URL:     "/sys",
					Visible: true,
					Children: []*amisgo.Menu{
						{
							MenuID:  "sys_user",
							Label:   "用户管理",
							Icon:    "fa fa-wrench",
							URL:     "/sys/user",
							Visible: true,
						},
						{
							MenuID:  "cms_news",
							Label:   "新闻管理",
							Icon:    "fa fa-newspaper-o",
							URL:     "/cms/news",
							Visible: true,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	pwdInput := amisgo.NewInputText()
	pwdInput.Name = "new_pwd"
	pwdInput.Label = "密码"
	pwdInput.Required = true

	form := amisgo.NewForm()
	form.Name = "pwd-form"
	form.API = amisgo.NewAPI()
	form.API.Method = "POST"
	form.API.URL = "${API_BASEURL}/admin/sys/user/update-pwd"
	form.Body = []amisgo.Component{
		pwdInput,
	}

	opPwd := amisgo.NewAction().WithLabel("重置密码").WithIcon("fa fa-key").WithDialog(
		amisgo.NewDialog().WithTitle("重置密码").WithBody(form),
	)
	ac, err := amisgo.NewAdminCRUD(
		"sys_user",
		"/sys/user/",
		amisgo.WithAPIPrefix("/admin"),
		amisgo.WithExecutor(db),
		amisgo.WithDBName("lucky"),
		amisgo.WithFields([]*amisgo.Field{
			{
				Name:      "id",
				Label:     "主键",
				ShowTable: true,
				ShowView:  true,
			},
			{
				Name:       "salt",
				Label:      "密码盐",
				ShowTable:  false,
				ShowView:   false,
				CreateAble: false,
				UpdateAble: false,
			},
			{
				Name:                "username",
				Label:               "用户名",
				Required:            true,
				ShowTable:           true,
				ShowView:            true,
				UpdateAble:          false,
				CreateAble:          true,
				SearchAble:          true,
				SearchConditionType: amisgo.ConditionTypeLike,
			},
			{
				Name:                "status",
				Label:               "状态",
				Type:                amisgo.TypeSelect,
				ShowTable:           true,
				ShowView:            true,
				UpdateAble:          false,
				SearchAble:          true,
				SearchConditionType: amisgo.ConditionTypeEq,
				Options: &amisgo.Options{
					GetOptions: func() ([]*amisgo.Option, error) {
						return []*amisgo.Option{
							{
								Label: "正常",
								Value: "normal",
							},
							{
								Label: "禁用",
								Value: "disabled",
							},
						}, nil
					},
				},
			},
		}),
		amisgo.WithOperations([]*amisgo.Action{
			opPwd,
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	news, err := amisgo.NewAdminCRUD(
		"cms_news",
		"/cms/news/",
		amisgo.WithAPIPrefix("/admin"),
		amisgo.WithExecutor(db),
		amisgo.WithDBName("lucky"),
		amisgo.WithFieldSort([]string{"id", "topic", "title", "content", "ctime", "utime", "remark", "status"}),
		amisgo.WithFields([]*amisgo.Field{
			{
				Name:                "title",
				ShowTable:           true,
				ShowView:            true,
				SearchAble:          true,
				SearchConditionType: amisgo.ConditionTypeLike,
				Required:            true,
				CreateAble:          true,
				UpdateAble:          true,
				QuickEdit:           true,
			},
			{
				Name:                "status",
				Label:               "状态",
				Type:                amisgo.TypeSelect,
				UpdateAble:          false,
				SearchAble:          true,
				ShowTable:           true,
				QuickEdit:           true,
				SearchConditionType: amisgo.ConditionTypeEq,
				Options: &amisgo.Options{
					GetOptions: func() ([]*amisgo.Option, error) {
						return []*amisgo.Option{
							{
								Label: "正常",
								Value: "normal",
							},
							{
								Label: "禁用",
								Value: "disabled",
							},
						}, nil
					},
				},
			},
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	router.RegAdminCRUD(ac, news)

	mux := http.NewServeMux()
	mux.Handle("/admin/", http.StripPrefix("/admin", router))
	mux.HandleFunc("/api/open/login", login)
	mux.HandleFunc("/admin/sys/user/info", user)
	// 应用 CORS 中间件
	handler := corsMiddleware(mux)

	fmt.Println("Starting server at port 8080")
	if err = http.ListenAndServe(":8081", handler); err != nil {
		fmt.Println(err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	kit.WriteJSON(w, http.StatusOK, map[string]any{
		"token": "token",
	})
}

func user(w http.ResponseWriter, r *http.Request) {
	kit.WriteJSON(w, http.StatusOK, map[string]any{
		"user_info": map[string]any{
			"user_id":  1,
			"username": "admin",
			"nickname": "管理员",
			"avatar":   "https://static.aispread.cn/amis/static/img/logo.png",
			"email":    "admin@aispread.cn",
		},
	})
}

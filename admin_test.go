package amisgo_test

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/fengjx/amisgo"
	"github.com/fengjx/amisgo/internal/kit"
)

func getDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.6.121:3306)/lucky?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestAdminCRUD(t *testing.T) {
	ac, err := amisgo.NewAdminCRUD(
		"sys_user",
		"/admin/sys/",
		amisgo.WithIDField("id"),
		amisgo.WithFields([]*amisgo.Field{
			{
				Name:  "id",
				Label: "主键",
			},
			{
				Name:  "username",
				Label: "用户名",
			},
			{
				Name:  "password",
				Label: "密码",
			},
			{
				Name:                "status",
				Label:               "状态",
				CreateAble:          true,
				UpdateAble:          true,
				SearchAble:          true,
				SearchConditionType: amisgo.ConditionTypeEq,
				Type:                amisgo.TypeSelect,
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
		t.Fatal(err)
	}
	fmt.Println("admin json: ", ac.ToJSON())
	t.Log(ac.ToJSON())
}

func TestUserAdmin(t *testing.T) {
	db := getDB(t)
	defer kit.QuicklyClose(db)
	ac, err := amisgo.NewAdminCRUD(
		"sys_user",
		"/admin/sys/user/",
		amisgo.WithExecutor(db),
		amisgo.WithDBName("lucky"),
		amisgo.WithFields([]*amisgo.Field{
			{
				Name:  "id",
				Label: "主键值",
			},
			{
				Name:                "username",
				Label:               "用户名",
				Required:            true,
				CreateAble:          true,
				UpdateAble:          false,
				SearchAble:          true,
				SearchConditionType: amisgo.ConditionTypeLike,
			},
			{
				Name:  "ctime",
				Label: "创建时间",
				Tpl:   "${ctime | date 'YYYY-MM-DD HH:mm:ss'}",
			},
			{
				Name:                "status",
				Label:               "状态",
				CreateAble:          true,
				UpdateAble:          true,
				SearchAble:          true,
				SearchConditionType: amisgo.ConditionTypeEq,
				Type:                amisgo.TypeSelect,
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
		t.Fatal(err)
	}
	fmt.Println("admin json: ", ac.ToJSON())
	t.Log(ac.ToJSON())
}

func TestMMenuAdmin(t *testing.T) {
	db := getDB(t)
	defer db.Close()
	ac, err := amisgo.NewAdminCRUD(
		"sys_menu",
		"/admin/sys/menu/",
		amisgo.WithExecutor(db),
		amisgo.WithDBName("lucky"),
		amisgo.WithFields([]*amisgo.Field{
			{
				Name:  "id",
				Label: "主键值",
			},
			{
				Name:     "parent_id",
				Label:    "父级菜单",
				Type:     amisgo.TypeTreeSelect,
				Required: true,
				Options: &amisgo.Options{
					GetOptions: func() ([]*amisgo.Option, error) {
						return []*amisgo.Option{
							{
								Label: "根菜单",
								Value: "1",
								Children: []*amisgo.Option{
									{Label: "系统管理", Value: "100"},
								}},
						}, nil
					},
				},
			},
			{
				Name:  "status",
				Label: "状态",
				Type:  amisgo.TypeSelect,
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
		t.Fatal(err)
	}
	fmt.Println("admin json: ", ac.ToJSON())
	t.Log(ac.ToJSON())
}

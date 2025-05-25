package amisgo

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/fengjx/amisgo/internal/dbutil"
	"github.com/fengjx/amisgo/internal/kit"
	"github.com/fengjx/amisgo/internal/log"
)

var (
	ErrMenuParentNotFound = errors.New("[amisgo] parent menu not found")

	StatusOK   int32 = 0 // 成功
	StatusFail int32 = 1 // 失败
)

// PageData 分页数据
type PageData[T any] struct {
	Items []T   `json:"items"` // 必须是一个 slice
	Total int64 `json:"total"` // 总记录数
}

// PageReq 分页参数
type PageReq struct {
	Offset     int64 `json:"offset"`      // 游标起始位置
	Limit      int64 `json:"limit"`       // 每页记录数
	HasNext    bool  `json:"has_next"`    // 是否有下一页
	Count      int64 `json:"count"`       // 总记录数
	QueryCount bool  `json:"query_count"` // 是否查询总数
}

// QueryReq 查询参数
type QueryReq struct {
	TableName   string       `json:"table_name"`             // 查询表
	Fields      []string     `json:"fields"`                 // 投影字段
	Conditions  []Condition  `json:"conditions,omitempty"`   // 查找字段
	OrderFields []OrderField `json:"order_fields,omitempty"` // 排序字段
	Page        *PageReq     `json:"page,omitempty"`         // 分页参数
}

type DeleteReq struct {
	Ids []any `json:"ids"` // 删除 ID 列表
}

// BatchUpdate 批量更新参数，rows 需包含 ID 字段
type BatchUpdate struct {
	Rows []map[string]any `json:"rows"`
}

// Resp amis 接口返回协议
type Resp struct {
	Status int32  `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
}

// AdminRouter admin 路由
type AdminRouter struct {
	mux      *http.ServeMux
	root     *Menu
	menusMap map[string]*Menu
}

// NewAdminRouter 创建 AdminRouter
func NewAdminRouter() *AdminRouter {
	root := &Menu{
		Label: "root",
		Children: []*Menu{
			{
				Label:    "Home",
				URL:      "/",
				Redirect: "/sys",
				Visible:  true,
			},
		},
	}
	r := &AdminRouter{
		mux:      http.NewServeMux(),
		root:     root,
		menusMap: make(map[string]*Menu),
	}
	r.init()
	return r
}

func (r *AdminRouter) init() {
	r.handle(http.MethodGet, "/menus", r.handleMenu)
}

func (r *AdminRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// AddMenu 添加菜单
func (r *AdminRouter) AddMenu(menus ...*Menu) error {
	for _, menu := range menus {
		if err := r.addMenu(menu); err != nil {
			return err
		}
	}
	return nil
}

// AddMenu 添加菜单
func (r *AdminRouter) addMenu(menu *Menu) error {
	parent := r.getMenu(menu.ParentLabel)
	if parent == nil {
		return ErrMenuParentNotFound
	}
	parent.Children = append(parent.Children, menu)
	r.menusMap[menu.Label] = menu
	for _, v := range menu.Children {
		if err := r.addMenu(v); err != nil {
			return err
		}
	}
	return nil
}

// RegisterAdminCRUD 注册 admin crud
func (r *AdminRouter) RegAdminCRUD(menu *Menu, admin *AdminCRUD) error {
	schemaPath := path.Join(admin.opt.apiPrefix, admin.opt.apiPath, "page.json")
	menu.SchemaAPI = schemaPath
	if err := r.AddMenu(menu); err != nil {
		return err
	}
	r.handle(http.MethodGet, admin.opt.pagePath, schemaHandler(admin))
	r.handle(http.MethodPost, admin.opt.createPath, createHandler(admin))
	r.handle(http.MethodPut, admin.opt.updatePath, updateHandler(admin))
	r.handle(http.MethodPatch, admin.opt.updatePath, updateHandler(admin))
	r.handle(http.MethodPut, admin.opt.batchUpdatePath, batchUpdateHandler(admin))
	r.handle(http.MethodPatch, admin.opt.batchUpdatePath, batchUpdateHandler(admin))
	r.handle(http.MethodDelete, admin.opt.deletePath, deleteHandler(admin))
	r.handle(http.MethodPost, admin.opt.queryPath, queryHandler(admin))
	return nil
}

func (r *AdminRouter) handle(method, path string, handler http.HandlerFunc) {
	pattern := fmt.Sprintf("%s %s", strings.ToUpper(method), path)
	r.mux.HandleFunc(pattern, handler)
}

// getMenu 获取菜单
func (r *AdminRouter) getMenu(label string) *Menu {
	if label == "" {
		return r.root
	}
	if p, ok := r.menusMap[label]; ok {
		return p
	}
	return nil
}

func (r *AdminRouter) handleMenu(w http.ResponseWriter, req *http.Request) {
	_ = kit.WriteJSON(w, http.StatusOK, map[string]any{
		"pages": r.root.Children,
	})
}

func createHandler(admin *AdminCRUD) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := dbutil.RowData{}
		if err := kit.ShouldBindJSON(req, &data); err != nil {
			kit.WriteJSON(w, http.StatusBadRequest, &Resp{
				Status: 0,
				Msg:    "参数解析错误",
			})
			return
		}
		resp, err := admin.create(req.Context(), data)
		if err != nil {
			kit.WriteJSON(w, http.StatusInternalServerError, &Resp{
				Status: StatusFail,
				Msg:    "系统错误",
			})
			return
		}
		_ = kit.WriteJSON(w, http.StatusOK, resp)
	}
}

func updateHandler(admin *AdminCRUD) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := dbutil.RowData{}
		if err := kit.ShouldBindJSON(req, &data); err != nil {
			kit.WriteJSON(w, http.StatusBadRequest, &Resp{
				Status: 0,
				Msg:    "参数解析错误",
			})
			return
		}
		id := data[admin.opt.idField]
		resp, err := admin.update(req.Context(), id, kit.Omit(data, admin.opt.idField))
		if err != nil {
			kit.WriteJSON(w, http.StatusInternalServerError, &Resp{
				Status: StatusFail,
				Msg:    "系统错误",
			})
			return
		}
		_ = kit.WriteJSON(w, http.StatusOK, resp)
	}
}

func batchUpdateHandler(admin *AdminCRUD) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		bu := &BatchUpdate{}
		if err := kit.ShouldBindJSON(req, bu); err != nil {
			kit.WriteJSON(w, http.StatusBadRequest, &Resp{
				Status: 0,
				Msg:    "参数解析错误",
			})
			return
		}
		resp, err := admin.batchUpdate(req.Context(), bu)
		if err != nil {
			log.Errorf("batch update, param: %v, error: %v", bu, err)
			kit.WriteJSON(w, http.StatusInternalServerError, &Resp{
				Status: StatusFail,
				Msg:    "系统错误",
			})
			return
		}
		_ = kit.WriteJSON(w, http.StatusOK, resp)
	}
}

func deleteHandler(admin *AdminCRUD) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		param := &DeleteReq{}
		if err := kit.ShouldBindJSON(req, param); err != nil {
			kit.WriteJSON(w, http.StatusBadRequest, &Resp{
				Status: StatusFail,
				Msg:    "参数解析错误",
			})
			return
		}
		resp, err := admin.delete(req.Context(), param.Ids...)
		if err != nil {
			kit.WriteJSON(w, http.StatusInternalServerError, &Resp{
				Status: StatusFail,
				Msg:    "系统错误",
			})
			return
		}
		_ = kit.WriteJSON(w, http.StatusOK, resp)
	}
}

func queryHandler(admin *AdminCRUD) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		query := &QueryReq{}
		if err := kit.ShouldBindJSON(req, query); err != nil {
			kit.WriteJSON(w, http.StatusBadRequest, &Resp{
				Status: 0,
				Msg:    "参数解析错误",
			})
			return
		}
		query.TableName = admin.opt.tableName
		resp, err := admin.query(req.Context(), query)
		if err != nil {
			kit.WriteJSON(w, http.StatusInternalServerError, &Resp{
				Status: StatusFail,
				Msg:    "系统错误",
			})
			return
		}
		_ = kit.WriteJSON(w, http.StatusOK, resp)
	}
}

func schemaHandler(admin *AdminCRUD) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		_ = kit.WriteJSON(w, http.StatusOK, admin.Component())
	}
}

// 拼接 URL 路径
func joinPath(base string, elem ...string) string {
	target, _ := url.JoinPath(base, elem...)
	return target
}

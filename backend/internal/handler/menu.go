package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yy945nb/marsview/backend/internal/model"
)

// MenuHandler 菜单相关处理器
type MenuHandler struct {
	db *gorm.DB
}

// NewMenuHandler 创建菜单处理器
func NewMenuHandler(db *gorm.DB) *MenuHandler {
	return &MenuHandler{db: db}
}

// GetMenuList 获取项目菜单列表（树形结构）
// GET /admin/menu/list/:projectId
func (h *MenuHandler) GetMenuList(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
	if err != nil {
		fail(c, 400, "项目ID格式错误")
		return
	}

	var menus []model.Menu
	h.db.Where("project_id = ?", projectID).Order("sort_num ASC, id ASC").Find(&menus)

	ok(c, gin.H{"list": buildMenuTree(menus, 0)})
}

// buildMenuTree 将扁平菜单列表构建为树形结构
func buildMenuTree(menus []model.Menu, parentID uint) []gin.H {
	result := make([]gin.H, 0)
	for _, m := range menus {
		if m.ParentID == parentID {
			node := gin.H{
				"id":        m.ID,
				"projectId": m.ProjectID,
				"parentId":  m.ParentID,
				"name":      m.Name,
				"type":      m.Type,
				"icon":      m.Icon,
				"path":      m.Path,
				"pageId":    m.PageID,
				"sortNum":   m.SortNum,
				"status":    m.Status,
				"userId":    m.UserID,
				"userName":  m.UserName,
				"createdAt": m.CreatedAt,
				"updatedAt": m.UpdatedAt,
			}
			children := buildMenuTree(menus, m.ID)
			if len(children) > 0 {
				node["children"] = children
			}
			result = append(result, node)
		}
	}
	return result
}

// AddMenu 新增菜单
// POST /admin/menu/add
func (h *MenuHandler) AddMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}
	userID := c.GetUint("userId")
	userName := c.GetString("userName")
	menu.UserID = userID
	menu.UserName = userName

	if err := h.db.Create(&menu).Error; err != nil {
		fail(c, 500, "创建菜单失败: "+err.Error())
		return
	}
	ok(c, menu.ID)
}

// UpdateMenu 更新菜单
// PUT /admin/menu/update
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}
	if menu.ID == 0 {
		fail(c, 400, "菜单ID不能为空")
		return
	}
	if err := h.db.Model(&menu).Updates(&menu).Error; err != nil {
		fail(c, 500, "更新菜单失败: "+err.Error())
		return
	}
	ok(c, "")
}

// DeleteMenu 删除菜单
// DELETE /admin/menu/delete/:id
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fail(c, 400, "菜单ID格式错误")
		return
	}
	if err := h.db.Delete(&model.Menu{}, id).Error; err != nil {
		fail(c, 500, "删除菜单失败: "+err.Error())
		return
	}
	ok(c, "")
}

// CopyMenu 复制菜单
// POST /admin/menu/copy
func (h *MenuHandler) CopyMenu(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}
	var menu model.Menu
	if err := h.db.First(&menu, req.ID).Error; err != nil {
		fail(c, 404, "菜单不存在")
		return
	}

	userID := c.GetUint("userId")
	userName := c.GetString("userName")
	newMenu := model.Menu{
		ProjectID: menu.ProjectID,
		ParentID:  menu.ParentID,
		Name:      menu.Name + "_copy",
		Type:      menu.Type,
		Icon:      menu.Icon,
		Path:      menu.Path,
		PageID:    menu.PageID,
		SortNum:   menu.SortNum,
		Status:    menu.Status,
		UserID:    userID,
		UserName:  userName,
	}
	if err := h.db.Create(&newMenu).Error; err != nil {
		fail(c, 500, "复制菜单失败: "+err.Error())
		return
	}
	ok(c, newMenu.ID)
}

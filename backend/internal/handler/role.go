package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yy945nb/marsview/backend/internal/model"
)

// RoleHandler 角色相关处理器
type RoleHandler struct {
	db *gorm.DB
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{db: db}
}

// GetRoleList 获取角色分页列表
// GET /admin/role/list?pageNum=1&pageSize=10&projectId=
func (h *RoleHandler) GetRoleList(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	projectID, _ := strconv.ParseUint(c.Query("projectId"), 10, 64)

	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := h.db.Model(&model.Role{})
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}

	var total int64
	query.Count(&total)

	var roles []model.Role
	query.Order("created_at DESC").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&roles)

	ok(c, gin.H{"total": total, "list": roles})
}

// GetRoleListAll 获取全部角色（不分页）
// GET /admin/role/all?projectId=
func (h *RoleHandler) GetRoleListAll(c *gin.Context) {
	projectID, _ := strconv.ParseUint(c.Query("projectId"), 10, 64)
	query := h.db.Model(&model.Role{})
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	var roles []model.Role
	query.Select("id, name, description").Find(&roles)
	ok(c, roles)
}

// CreateRole 创建角色
// POST /admin/role/create
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}
	userID := c.GetUint("userId")
	userName := c.GetString("userName")
	role.UserID = userID
	role.UserName = userName

	if err := h.db.Create(&role).Error; err != nil {
		fail(c, 500, "创建角色失败: "+err.Error())
		return
	}
	ok(c, role.ID)
}

// DeleteRole 删除角色
// DELETE /admin/role/delete/:id
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fail(c, 400, "角色ID格式错误")
		return
	}
	if err := h.db.Delete(&model.Role{}, id).Error; err != nil {
		fail(c, 500, "删除角色失败: "+err.Error())
		return
	}
	ok(c, "")
}

// UpdateRole 更新角色信息
// PUT /admin/role/update
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}
	if role.ID == 0 {
		fail(c, 400, "角色ID不能为空")
		return
	}
	if err := h.db.Model(&role).Updates(&role).Error; err != nil {
		fail(c, 500, "更新角色失败: "+err.Error())
		return
	}
	ok(c, "")
}

// UpdateRoleLimits 更新角色权限
// PUT /admin/role/permissions
func (h *RoleHandler) UpdateRoleLimits(c *gin.Context) {
	var req struct {
		ID     uint   `json:"id" binding:"required"`
		Limits string `json:"limits"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := h.db.Model(&model.Role{}).Where("id = ?", req.ID).Update("limits", req.Limits).Error; err != nil {
		fail(c, 500, "更新权限失败: "+err.Error())
		return
	}
	ok(c, "")
}

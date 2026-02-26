package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yy945nb/marsview/backend/internal/model"
)

// ProjectUserHandler 项目用户管理处理器
type ProjectUserHandler struct {
	db *gorm.DB
}

// NewProjectUserHandler 创建项目用户管理处理器
func NewProjectUserHandler(db *gorm.DB) *ProjectUserHandler {
	return &ProjectUserHandler{db: db}
}

// GetUserList 获取项目用户列表
// GET /admin/user/list?pageNum=1&pageSize=10&projectId=
func (h *ProjectUserHandler) GetUserList(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	projectID, _ := strconv.ParseUint(c.Query("projectId"), 10, 64)

	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := h.db.Model(&model.ProjectUser{})
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}

	var total int64
	query.Count(&total)

	var users []model.ProjectUser
	query.Order("created_at DESC").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&users)

	ok(c, gin.H{"total": total, "list": users})
}

// AddUser 新增项目用户
// POST /admin/user/add
func (h *ProjectUserHandler) AddUser(c *gin.Context) {
	var user model.ProjectUser
	if err := c.ShouldBindJSON(&user); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}

	// 检查是否已存在
	var count int64
	h.db.Model(&model.ProjectUser{}).
		Where("project_id = ? AND user_id = ?", user.ProjectID, user.UserID).
		Count(&count)
	if count > 0 {
		fail(c, 400, "用户已加入该项目")
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		fail(c, 500, "添加用户失败: "+err.Error())
		return
	}
	ok(c, user.ID)
}

// DeleteUser 删除项目用户
// DELETE /admin/user/delete/:id
func (h *ProjectUserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fail(c, 400, "用户ID格式错误")
		return
	}
	if err := h.db.Delete(&model.ProjectUser{}, id).Error; err != nil {
		fail(c, 500, "删除用户失败: "+err.Error())
		return
	}
	ok(c, "")
}

// UpdateUser 更新项目用户角色
// PUT /admin/user/update
func (h *ProjectUserHandler) UpdateUser(c *gin.Context) {
	var user model.ProjectUser
	if err := c.ShouldBindJSON(&user); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}
	if user.ID == 0 {
		fail(c, 400, "ID不能为空")
		return
	}
	if err := h.db.Model(&user).Updates(&user).Error; err != nil {
		fail(c, 500, "更新失败: "+err.Error())
		return
	}
	ok(c, "")
}

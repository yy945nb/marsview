package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yy945nb/marsview/backend/internal/model"
)

// ProjectHandler 项目相关处理器
type ProjectHandler struct {
	db *gorm.DB
}

// NewProjectHandler 创建项目处理器
func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{db: db}
}

// GetProjectList 获取项目列表
// GET /admin/project/list?pageNum=1&pageSize=10&keyword=
func (h *ProjectHandler) GetProjectList(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "12"))
	keyword := c.Query("keyword")
	userID := c.GetUint("userId")

	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 12
	}

	query := h.db.Model(&model.Project{}).Where(
		"(is_public = 1 OR user_id = ?)", userID,
	)
	if keyword != "" {
		query = query.Where("name LIKE ? OR remark LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var projects []model.Project
	query.Order("created_at DESC").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&projects)

	ok(c, gin.H{"total": total, "list": projects})
}

// CreateProject 创建项目
// POST /admin/project/create
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}
	userID := c.GetUint("userId")
	userName := c.GetString("userName")
	project.UserID = userID
	project.UserName = userName

	if err := h.db.Create(&project).Error; err != nil {
		fail(c, 500, "创建项目失败: "+err.Error())
		return
	}
	ok(c, project.ID)
}

// UpdateProject 更新项目
// PUT /admin/project/update
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}
	if project.ID == 0 {
		fail(c, 400, "项目ID不能为空")
		return
	}
	if err := h.db.Model(&project).Updates(&project).Error; err != nil {
		fail(c, 500, "更新项目失败: "+err.Error())
		return
	}
	ok(c, "")
}

// DeleteProject 删除项目
// DELETE /admin/project/delete/:id
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fail(c, 400, "项目ID格式错误")
		return
	}
	if err := h.db.Delete(&model.Project{}, id).Error; err != nil {
		fail(c, 500, "删除项目失败: "+err.Error())
		return
	}
	ok(c, "")
}

// GetProjectDetail 获取项目详情/配置
// GET /admin/getProjectConfig?projectId=
func (h *ProjectHandler) GetProjectDetail(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Query("projectId"), 10, 64)
	if err != nil {
		fail(c, 400, "项目ID格式错误")
		return
	}
	var project model.Project
	if err := h.db.First(&project, projectID).Error; err != nil {
		fail(c, 404, "项目不存在")
		return
	}
	ok(c, project)
}

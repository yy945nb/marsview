package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yy945nb/marsview/backend/internal/model"
)

// 页面状态常量（与前端约定一致）
const (
	pageStateUnsaved   = 1 // 未保存
	pageStateSaved     = 2 // 已保存
	pageStatePublished = 3 // 已发布
	pageStateRolledBack = 4 // 已回滚
)

// PageHandler 页面相关处理器
type PageHandler struct {
	db *gorm.DB
}

// NewPageHandler 创建页面处理器
func NewPageHandler(db *gorm.DB) *PageHandler {
	return &PageHandler{db: db}
}

// GetPageList 获取页面列表
// GET /admin/page/list?pageNum=1&pageSize=10&projectId=&keyword=
func (h *PageHandler) GetPageList(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "12"))
	keyword := c.Query("keyword")
	projectID, _ := strconv.ParseUint(c.Query("projectId"), 10, 64)
	userID := c.GetUint("userId")

	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 12
	}

	query := h.db.Model(&model.Page{}).Where("is_public = 1 OR user_id = ?", userID)
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR remark LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 只返回非配置字段（不返回pageData，减少传输量）
	query = query.Select("id, project_id, name, remark, user_id, user_name, is_public, stg_state, pre_state, prd_state, stg_publish_id, pre_publish_id, prd_publish_id, preview_img, created_at, updated_at")

	var total int64
	query.Count(&total)

	var pages []model.Page
	query.Order("created_at DESC").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&pages)

	ok(c, gin.H{"total": total, "list": pages})
}

// GetPageDetail 获取页面详情（含pageData配置）
// GET /admin/page/detail/:env/:id
func (h *PageHandler) GetPageDetail(c *gin.Context) {
	env := c.Param("env")
	pageID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fail(c, 400, "页面ID格式错误")
		return
	}
	userID := c.GetUint("userId")

	var page model.Page
	if err := h.db.First(&page, pageID).Error; err != nil {
		fail(c, 404, "页面不存在")
		return
	}

	// 检查访问权限
	if page.IsPublic != 1 && page.UserID != userID {
		// 检查是否为页面成员
		var member model.PageMember
		if err := h.db.Where("page_id = ? AND user_id = ?", pageID, userID).First(&member).Error; err != nil {
			fail(c, 403, "无权访问该页面")
			return
		}
	}

	// stg/pre/prd 环境使用对应的发布版本
	pageData := page.PageData
	if env != "" && env != "edit" {
		publishID := page.PrdPublishID
		if env == "stg" {
			publishID = page.StgPublishID
		} else if env == "pre" {
			publishID = page.PrePublishID
		}
		if publishID > 0 {
			var pub model.PagePublish
			if err := h.db.First(&pub, publishID).Error; err == nil {
				pageData = pub.PageData
			}
		}
	}

	ok(c, gin.H{
		"id":           page.ID,
		"projectId":    page.ProjectID,
		"name":         page.Name,
		"remark":       page.Remark,
		"isPublic":     page.IsPublic,
		"stgState":     page.StgState,
		"preState":     page.PreState,
		"prdState":     page.PrdState,
		"stgPublishId": page.StgPublishID,
		"prePublishId": page.PrePublishID,
		"prdPublishId": page.PrdPublishID,
		"previewImg":   page.PreviewImg,
		"userId":       page.UserID,
		"userName":     page.UserName,
		"pageData":     pageData,
		"createdAt":    page.CreatedAt,
		"updatedAt":    page.UpdatedAt,
	})
}

// CreatePage 创建页面
// POST /admin/page/create
func (h *PageHandler) CreatePage(c *gin.Context) {
	var page model.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}
	userID := c.GetUint("userId")
	userName := c.GetString("userName")
	page.UserID = userID
	page.UserName = userName
	page.StgState = 1
	page.PreState = 1
	page.PrdState = 1

	if err := h.db.Create(&page).Error; err != nil {
		fail(c, 500, "创建页面失败: "+err.Error())
		return
	}
	ok(c, page.ID)
}

// UpdatePageData 更新页面数据（保存编辑器内容）
// POST /admin/page/update
func (h *PageHandler) UpdatePageData(c *gin.Context) {
	var req struct {
		ID         uint   `json:"id" binding:"required"`
		Name       string `json:"name"`
		Remark     string `json:"remark"`
		IsPublic   *int   `json:"isPublic"`
		PreviewImg string `json:"previewImg"`
		PageData   string `json:"pageData"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{
		"stg_state": pageStateSaved,
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.PreviewImg != "" {
		updates["preview_img"] = req.PreviewImg
	}
	if req.PageData != "" {
		updates["page_data"] = req.PageData
	}

	if err := h.db.Model(&model.Page{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		fail(c, 500, "保存页面失败: "+err.Error())
		return
	}
	ok(c, "")
}

// DeletePage 删除页面
// DELETE /admin/page/delete/:id
func (h *PageHandler) DeletePage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fail(c, 400, "页面ID格式错误")
		return
	}
	if err := h.db.Delete(&model.Page{}, id).Error; err != nil {
		fail(c, 500, "删除页面失败: "+err.Error())
		return
	}
	ok(c, "")
}

// CopyPage 复制页面
// POST /admin/page/copy
func (h *PageHandler) CopyPage(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}
	var page model.Page
	if err := h.db.First(&page, req.ID).Error; err != nil {
		fail(c, 404, "页面不存在")
		return
	}

	userID := c.GetUint("userId")
	userName := c.GetString("userName")
	newPage := model.Page{
		ProjectID:  page.ProjectID,
		Name:       page.Name + "_copy",
		Remark:     page.Remark,
		UserID:     userID,
		UserName:   userName,
		IsPublic:   page.IsPublic,
		PageData:   page.PageData,
		StgState:   1,
		PreState:   1,
		PrdState:   1,
	}
	if err := h.db.Create(&newPage).Error; err != nil {
		fail(c, 500, "复制页面失败: "+err.Error())
		return
	}
	ok(c, newPage.ID)
}

// PublishPage 发布页面到指定环境
// POST /admin/page/publish
func (h *PageHandler) PublishPage(c *gin.Context) {
	var req struct {
		ID     uint   `json:"id" binding:"required"`
		Env    string `json:"env" binding:"required"` // stg/pre/prd
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}

	var page model.Page
	if err := h.db.First(&page, req.ID).Error; err != nil {
		fail(c, 404, "页面不存在")
		return
	}

	userID := c.GetUint("userId")
	userName := c.GetString("userName")

	pub := model.PagePublish{
		PageID:   page.ID,
		Env:      req.Env,
		PageData: page.PageData,
		UserID:   userID,
		UserName: userName,
		Remark:   req.Remark,
	}
	if err := h.db.Create(&pub).Error; err != nil {
		fail(c, 500, "发布页面失败: "+err.Error())
		return
	}

	// 更新页面对应环境的发布状态
	updates := map[string]interface{}{}
	switch req.Env {
	case "stg":
		updates["stg_state"] = pageStatePublished
		updates["stg_publish_id"] = pub.ID
	case "pre":
		updates["pre_state"] = pageStatePublished
		updates["pre_publish_id"] = pub.ID
	case "prd":
		updates["prd_state"] = pageStatePublished
		updates["prd_publish_id"] = pub.ID
	}
	if len(updates) > 0 {
		h.db.Model(&model.Page{}).Where("id = ?", page.ID).Updates(updates)
	}

	ok(c, pub.ID)
}

// GetPublishList 获取页面发布历史
// GET /admin/page/publishList?id=&env=
func (h *PageHandler) GetPublishList(c *gin.Context) {
	pageID, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	env := c.Query("env")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := h.db.Model(&model.PagePublish{}).
		Select("id, page_id, env, user_id, user_name, remark, created_at").
		Where("page_id = ?", pageID)
	if env != "" {
		query = query.Where("env = ?", env)
	}

	var total int64
	query.Count(&total)

	var list []model.PagePublish
	query.Order("created_at DESC").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&list)

	ok(c, gin.H{"total": total, "list": list})
}

// RollbackPage 回滚页面到指定发布版本
// POST /admin/page/rollback
func (h *PageHandler) RollbackPage(c *gin.Context) {
	var req struct {
		PublishID uint `json:"publishId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}

	var pub model.PagePublish
	if err := h.db.First(&pub, req.PublishID).Error; err != nil {
		fail(c, 404, "发布记录不存在")
		return
	}

	updates := map[string]interface{}{
		"page_data": pub.PageData,
	}
	switch pub.Env {
	case "stg":
		updates["stg_state"] = pageStateRolledBack
		updates["stg_publish_id"] = pub.ID
	case "pre":
		updates["pre_state"] = pageStateRolledBack
		updates["pre_publish_id"] = pub.ID
	case "prd":
		updates["prd_state"] = pageStateRolledBack
		updates["prd_publish_id"] = pub.ID
	}
	if err := h.db.Model(&model.Page{}).Where("id = ?", pub.PageID).Updates(updates).Error; err != nil {
		fail(c, 500, "回滚失败: "+err.Error())
		return
	}
	ok(c, "")
}

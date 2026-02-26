package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yy945nb/marsview/backend/internal/model"
)

// PageMemberHandler 页面成员处理器
type PageMemberHandler struct {
	db *gorm.DB
}

// NewPageMemberHandler 创建页面成员处理器
func NewPageMemberHandler(db *gorm.DB) *PageMemberHandler {
	return &PageMemberHandler{db: db}
}

// GetMemberList 获取页面成员列表
// GET /admin/page/member/list?pageId=
func (h *PageMemberHandler) GetMemberList(c *gin.Context) {
	pageID, _ := strconv.ParseUint(c.Query("pageId"), 10, 64)

	var members []model.PageMember
	h.db.Where("page_id = ?", pageID).Find(&members)

	ok(c, gin.H{"list": members, "total": len(members)})
}

// AddPageMember 添加页面成员
// POST /admin/page/member/add
func (h *PageMemberHandler) AddPageMember(c *gin.Context) {
	var member model.PageMember
	if err := c.ShouldBindJSON(&member); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}

	// 检查是否已存在
	var count int64
	h.db.Model(&model.PageMember{}).
		Where("page_id = ? AND user_id = ?", member.PageID, member.UserID).
		Count(&count)
	if count > 0 {
		fail(c, 400, "用户已是该页面成员")
		return
	}

	if err := h.db.Create(&member).Error; err != nil {
		fail(c, 500, "添加成员失败: "+err.Error())
		return
	}
	ok(c, member.ID)
}

// DeletePageMember 删除页面成员
// DELETE /admin/page/member/delete/:id
func (h *PageMemberHandler) DeletePageMember(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		fail(c, 400, "成员ID格式错误")
		return
	}
	if err := h.db.Delete(&model.PageMember{}, id).Error; err != nil {
		fail(c, 500, "删除成员失败: "+err.Error())
		return
	}
	ok(c, "")
}

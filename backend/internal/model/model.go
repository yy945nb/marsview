package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserName  string         `json:"userName" gorm:"uniqueIndex;size:100;not null"`
	NickName  string         `json:"nickName" gorm:"size:100"`
	UserPwd   string         `json:"-" gorm:"size:200;not null"`
	Avatar    string         `json:"avatar" gorm:"size:500"`
	IsSuperAdmin int         `json:"isSuperAdmin" gorm:"default:2"` // 1:超级管理员 2:普通用户
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Project 项目表
type Project struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string         `json:"name" gorm:"size:100;not null"`
	Remark         string         `json:"remark" gorm:"size:500"`
	Logo           string         `json:"logo" gorm:"size:500"`
	UserID         uint           `json:"userId"`
	UserName       string         `json:"userName" gorm:"size:100"`
	Layout         int            `json:"layout" gorm:"default:1"` // 布局方式
	MenuMode       string         `json:"menuMode" gorm:"size:50;default:'vertical'"`
	MenuThemeColor string         `json:"menuThemeColor" gorm:"size:50;default:'dark'"`
	Breadcrumb     int            `json:"breadcrumb" gorm:"default:1"`
	Tag            int            `json:"tag" gorm:"default:1"`
	Footer         int            `json:"footer" gorm:"default:0"`
	IsPublic       int            `json:"isPublic" gorm:"default:2"` // 1:公开 2:私有
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// Page 页面表
type Page struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID    uint           `json:"projectId" gorm:"index;not null"`
	Name         string         `json:"name" gorm:"size:100;not null"`
	Remark       string         `json:"remark" gorm:"size:500"`
	UserID       uint           `json:"userId"`
	UserName     string         `json:"userName" gorm:"size:100"`
	IsPublic     int            `json:"isPublic" gorm:"default:2"` // 1:公开 2:私有
	StgState     int            `json:"stgState" gorm:"default:1"` // 1:未发布 2:已发布
	PreState     int            `json:"preState" gorm:"default:1"`
	PrdState     int            `json:"prdState" gorm:"default:1"`
	StgPublishID uint           `json:"stgPublishId"`
	PrePublishID uint           `json:"prePublishId"`
	PrdPublishID uint           `json:"prdPublishId"`
	PreviewImg   string         `json:"previewImg" gorm:"size:500"`
	PageData     string         `json:"pageData" gorm:"type:text"` // JSON 页面配置数据
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// PagePublish 页面发布记录
type PagePublish struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PageID    uint      `json:"pageId" gorm:"index;not null"`
	Env       string    `json:"env" gorm:"size:20;not null"` // stg/pre/prd
	PageData  string    `json:"pageData" gorm:"type:text"`
	UserID    uint      `json:"userId"`
	UserName  string    `json:"userName" gorm:"size:100"`
	Remark    string    `json:"remark" gorm:"size:500"`
	CreatedAt time.Time `json:"createdAt"`
}

// Menu 菜单表
type Menu struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID uint           `json:"projectId" gorm:"index;not null"`
	ParentID  uint           `json:"parentId" gorm:"default:0"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Type      int            `json:"type" gorm:"default:1"` // 1:菜单 2:按钮
	Icon      string         `json:"icon" gorm:"size:100"`
	Path      string         `json:"path" gorm:"size:500"`
	PageID    uint           `json:"pageId"`
	SortNum   int            `json:"sortNum" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"` // 1:启用 2:禁用
	UserID    uint           `json:"userId"`
	UserName  string         `json:"userName" gorm:"size:100"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Role 角色表
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID   uint           `json:"projectId" gorm:"index;not null"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"size:500"`
	UserID      uint           `json:"userId"`
	UserName    string         `json:"userName" gorm:"size:100"`
	Limits      string         `json:"limits" gorm:"type:text"` // JSON 权限列表
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// ProjectUser 项目用户关联表
type ProjectUser struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID uint           `json:"projectId" gorm:"index;not null"`
	UserID    uint           `json:"userId" gorm:"not null"`
	UserName  string         `json:"userName" gorm:"size:100"`
	RoleID    uint           `json:"roleId"`
	RoleName  string         `json:"roleName" gorm:"size:100"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// PageMember 页面成员表
type PageMember struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	PageID    uint           `json:"pageId" gorm:"index;not null"`
	UserID    uint           `json:"userId" gorm:"not null"`
	UserName  string         `json:"userName" gorm:"size:100"`
	Role      int            `json:"role" gorm:"default:1"` // 1:成员 2:管理员
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yy945nb/marsview/backend/internal/config"
	"github.com/yy945nb/marsview/backend/internal/handler"
	"github.com/yy945nb/marsview/backend/internal/middleware"
)

// Setup 注册所有路由
func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// 初始化各处理器
	userH := handler.NewUserHandler(db, cfg.JWT.Secret, cfg.JWT.ExpireHour)
	projectH := handler.NewProjectHandler(db)
	pageH := handler.NewPageHandler(db)
	menuH := handler.NewMenuHandler(db)
	roleH := handler.NewRoleHandler(db)
	projectUserH := handler.NewProjectUserHandler(db)
	pageMemberH := handler.NewPageMemberHandler(db)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// ── 公开路由（无需鉴权）──
	r.POST("/user/login", userH.Login)

	// ── 需鉴权路由 ──
	auth := r.Group("/", middleware.Auth(cfg.JWT.Secret))

	// 用户
	auth.GET("/user/info", userH.GetUserInfo)
	auth.GET("/user/search", userH.SearchUser)

	// 项目
	auth.GET("/admin/project/list", projectH.GetProjectList)
	auth.POST("/admin/project/create", projectH.CreateProject)
	auth.PUT("/admin/project/update", projectH.UpdateProject)
	auth.DELETE("/admin/project/delete/:id", projectH.DeleteProject)
	auth.GET("/admin/getProjectConfig", projectH.GetProjectDetail)

	// 页面
	auth.GET("/admin/page/list", pageH.GetPageList)
	auth.GET("/admin/page/detail/:env/:id", pageH.GetPageDetail)
	auth.POST("/admin/page/create", pageH.CreatePage)
	auth.POST("/admin/page/update", pageH.UpdatePageData)
	auth.DELETE("/admin/page/delete/:id", pageH.DeletePage)
	auth.POST("/admin/page/copy", pageH.CopyPage)
	auth.POST("/admin/page/publish", pageH.PublishPage)
	auth.GET("/admin/page/publishList", pageH.GetPublishList)
	auth.POST("/admin/page/rollback", pageH.RollbackPage)

	// 菜单
	auth.GET("/admin/menu/list/:projectId", menuH.GetMenuList)
	auth.POST("/admin/menu/add", menuH.AddMenu)
	auth.PUT("/admin/menu/update", menuH.UpdateMenu)
	auth.DELETE("/admin/menu/delete/:id", menuH.DeleteMenu)
	auth.POST("/admin/menu/copy", menuH.CopyMenu)

	// 角色
	auth.GET("/admin/role/list", roleH.GetRoleList)
	auth.GET("/admin/role/all", roleH.GetRoleListAll)
	auth.POST("/admin/role/create", roleH.CreateRole)
	auth.DELETE("/admin/role/delete/:id", roleH.DeleteRole)
	auth.PUT("/admin/role/update", roleH.UpdateRole)
	auth.PUT("/admin/role/permissions", roleH.UpdateRoleLimits)

	// 项目用户管理
	auth.GET("/admin/user/list", projectUserH.GetUserList)
	auth.POST("/admin/user/add", projectUserH.AddUser)
	auth.DELETE("/admin/user/delete/:id", projectUserH.DeleteUser)
	auth.PUT("/admin/user/update", projectUserH.UpdateUser)

	// 页面成员
	auth.GET("/admin/page/member/list", pageMemberH.GetMemberList)
	auth.POST("/admin/page/member/add", pageMemberH.AddPageMember)
	auth.DELETE("/admin/page/member/delete/:id", pageMemberH.DeletePageMember)

	return r
}

package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/yy945nb/marsview/backend/internal/middleware"
	"github.com/yy945nb/marsview/backend/internal/model"
)

// UserHandler 用户相关处理器
type UserHandler struct {
	db        *gorm.DB
	jwtSecret string
	jwtExpire int
}

// NewUserHandler 创建用户处理器
func NewUserHandler(db *gorm.DB, jwtSecret string, jwtExpire int) *UserHandler {
	return &UserHandler{db: db, jwtSecret: jwtSecret, jwtExpire: jwtExpire}
}

type loginReq struct {
	UserName string `json:"userName" binding:"required"`
	UserPwd  string `json:"userPwd" binding:"required"`
}

// Login 用户登录
// POST /user/login
func (h *UserHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}
	var user model.User
	err := h.db.Where("user_name = ?", req.UserName).First(&user).Error
	// 无论用户是否存在都执行 bcrypt 比较，避免时序攻击
	hashToCompare := user.UserPwd
	if err != nil {
		// 用户不存在时用一个固定哈希占位，确保执行时间一致
		hashToCompare = "$2a$10$aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	}
	compareErr := bcrypt.CompareHashAndPassword([]byte(hashToCompare), []byte(req.UserPwd))
	if err != nil || compareErr != nil {
		fail(c, 400, "用户名或密码错误")
		return
	}
	token, err := h.generateToken(user.ID, user.UserName)
	if err != nil {
		fail(c, 500, "生成Token失败")
		return
	}
	ok(c, gin.H{
		"userId":   user.ID,
		"userName": user.UserName,
		"nickName": user.NickName,
		"avatar":   user.Avatar,
		"token":    token,
	})
}

// GetUserInfo 获取当前登录用户信息
// GET /user/info
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := c.GetUint("userId")
	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		fail(c, 404, "用户不存在")
		return
	}
	ok(c, gin.H{
		"userId":   user.ID,
		"userName": user.UserName,
		"nickName": user.NickName,
		"avatar":   user.Avatar,
	})
}

// SearchUser 搜索用户
// GET /user/search?keyword=
func (h *UserHandler) SearchUser(c *gin.Context) {
	keyword := c.Query("keyword")
	var users []model.User
	query := h.db.Select("id, user_name, nick_name, avatar")
	if keyword != "" {
		query = query.Where("user_name LIKE ? OR nick_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Limit(20).Find(&users)
	list := make([]gin.H, 0, len(users))
	for _, u := range users {
		list = append(list, gin.H{
			"id":       u.ID,
			"userName": u.UserName,
			"nickName": u.NickName,
			"avatar":   u.Avatar,
		})
	}
	ok(c, gin.H{"list": list, "total": len(list)})
}

func (h *UserHandler) generateToken(userID uint, userName string) (string, error) {
	claims := &middleware.Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(h.jwtExpire) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

// response helpers

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data, "message": "success"})
}

func fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": code, "data": nil, "message": msg})
}

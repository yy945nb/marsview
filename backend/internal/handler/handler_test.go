package handler_test

import (
"bytes"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"

"golang.org/x/crypto/bcrypt"
"gorm.io/driver/sqlite"
"gorm.io/gorm"

"github.com/yy945nb/marsview/backend/internal/config"
"github.com/yy945nb/marsview/backend/internal/model"
"github.com/yy945nb/marsview/backend/internal/router"
)

func setupTestDB(t *testing.T) *gorm.DB {
t.Helper()
db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
if err != nil {
t.Fatalf("打开内存数据库失败: %v", err)
}
if err := db.AutoMigrate(
&model.User{},
&model.Project{},
&model.Page{},
&model.PagePublish{},
&model.Menu{},
&model.Role{},
&model.ProjectUser{},
&model.PageMember{},
); err != nil {
t.Fatalf("数据库迁移失败: %v", err)
}
return db
}

func setupTestApp(t *testing.T) (*gorm.DB, http.Handler) {
t.Helper()
db := setupTestDB(t)
cfg := &config.Config{
Server:   config.ServerConfig{Mode: "test"},
Database: config.DatabaseConfig{Driver: "sqlite", DSN: ":memory:"},
JWT:      config.JWTConfig{Secret: "test-secret", ExpireHour: 72},
}
r := router.Setup(db, cfg)
return db, r
}

// createTestUser 在测试DB中创建一个带哈希密码的用户
func createTestUser(t *testing.T, db *gorm.DB) model.User {
t.Helper()
hashed, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
if err != nil {
t.Fatalf("生成密码哈希失败: %v", err)
}
user := model.User{
UserName: "test@test.com",
NickName: "Test User",
UserPwd:  string(hashed),
Avatar:   "/imgs/test.png",
}
if err := db.Create(&user).Error; err != nil {
t.Fatalf("创建测试用户失败: %v", err)
}
return user
}

func TestHealthCheck(t *testing.T) {
_, r := setupTestApp(t)
req := httptest.NewRequest(http.MethodGet, "/health", nil)
w := httptest.NewRecorder()
r.ServeHTTP(w, req)
if w.Code != http.StatusOK {
t.Fatalf("期望状态码 200，得到 %d", w.Code)
}
}

func TestLoginFailsWithWrongCredentials(t *testing.T) {
db, r := setupTestApp(t)
createTestUser(t, db)

body := `{"userName":"notexist@test.com","userPwd":"wrong"}`
req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBufferString(body))
req.Header.Set("Content-Type", "application/json")
w := httptest.NewRecorder()
r.ServeHTTP(w, req)

if w.Code != http.StatusOK {
t.Fatalf("期望状态码 200，得到 %d", w.Code)
}
var resp map[string]interface{}
if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
t.Fatalf("解析响应失败: %v", err)
}
if resp["code"] == float64(0) {
t.Error("错误凭证登录不应成功")
}
}

func TestLoginSuccess(t *testing.T) {
db, r := setupTestApp(t)
createTestUser(t, db)

body := `{"userName":"test@test.com","userPwd":"password123"}`
req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBufferString(body))
req.Header.Set("Content-Type", "application/json")
w := httptest.NewRecorder()
r.ServeHTTP(w, req)

if w.Code != http.StatusOK {
t.Fatalf("期望状态码 200，得到 %d", w.Code)
}
var resp map[string]interface{}
if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
t.Fatalf("解析响应失败: %v", err)
}
if resp["code"] != float64(0) {
t.Errorf("登录失败，响应: %v", resp)
}
data, ok := resp["data"].(map[string]interface{})
if !ok || data["token"] == "" {
t.Error("登录成功后应返回 token")
}
}

func TestGetUserInfoWithoutToken(t *testing.T) {
_, r := setupTestApp(t)
req := httptest.NewRequest(http.MethodGet, "/user/info", nil)
w := httptest.NewRecorder()
r.ServeHTTP(w, req)
if w.Code != http.StatusUnauthorized {
t.Fatalf("没有 token 时期望 401，得到 %d", w.Code)
}
}

func TestProjectListRequiresAuth(t *testing.T) {
_, r := setupTestApp(t)
req := httptest.NewRequest(http.MethodGet, "/admin/project/list", nil)
w := httptest.NewRecorder()
r.ServeHTTP(w, req)
if w.Code != http.StatusUnauthorized {
t.Fatalf("未认证请求期望 401，得到 %d", w.Code)
}
}

// loginAndGetToken 辅助：登录并获取 token
func loginAndGetToken(t *testing.T, r http.Handler) string {
t.Helper()
body := `{"userName":"test@test.com","userPwd":"password123"}`
req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBufferString(body))
req.Header.Set("Content-Type", "application/json")
w := httptest.NewRecorder()
r.ServeHTTP(w, req)
var resp map[string]interface{}
if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
t.Fatalf("解析登录响应失败: %v", err)
}
data, ok := resp["data"].(map[string]interface{})
if !ok {
t.Fatalf("登录响应格式错误: %v", resp)
}
return data["token"].(string)
}

func TestCreateAndListProject(t *testing.T) {
db, r := setupTestApp(t)
createTestUser(t, db)
token := loginAndGetToken(t, r)

// 创建项目
createBody := `{"name":"测试项目","remark":"这是测试","isPublic":1}`
req := httptest.NewRequest(http.MethodPost, "/admin/project/create", bytes.NewBufferString(createBody))
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer "+token)
w := httptest.NewRecorder()
r.ServeHTTP(w, req)

var createResp map[string]interface{}
json.Unmarshal(w.Body.Bytes(), &createResp)
if createResp["code"] != float64(0) {
t.Fatalf("创建项目失败: %v", createResp)
}

// 获取列表
req2 := httptest.NewRequest(http.MethodGet, "/admin/project/list?pageNum=1&pageSize=10", nil)
req2.Header.Set("Authorization", "Bearer "+token)
w2 := httptest.NewRecorder()
r.ServeHTTP(w2, req2)

var listResp map[string]interface{}
json.Unmarshal(w2.Body.Bytes(), &listResp)
if listResp["code"] != float64(0) {
t.Fatalf("获取项目列表失败: %v", listResp)
}
listData := listResp["data"].(map[string]interface{})
if listData["total"].(float64) < 1 {
t.Error("创建项目后列表应不为空")
}
}

func TestCreateAndGetPage(t *testing.T) {
db, r := setupTestApp(t)
createTestUser(t, db)
token := loginAndGetToken(t, r)

// 创建页面
createBody := `{"projectId":1,"name":"首页","isPublic":1}`
req := httptest.NewRequest(http.MethodPost, "/admin/page/create", bytes.NewBufferString(createBody))
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer "+token)
w := httptest.NewRecorder()
r.ServeHTTP(w, req)

var createResp map[string]interface{}
json.Unmarshal(w.Body.Bytes(), &createResp)
if createResp["code"] != float64(0) {
t.Fatalf("创建页面失败: %v", createResp)
}

// 获取页面详情
req2 := httptest.NewRequest(http.MethodGet, "/admin/page/detail/edit/1", nil)
req2.Header.Set("Authorization", "Bearer "+token)
w2 := httptest.NewRecorder()
r.ServeHTTP(w2, req2)

var detailResp map[string]interface{}
json.Unmarshal(w2.Body.Bytes(), &detailResp)
if detailResp["code"] != float64(0) {
t.Fatalf("获取页面详情失败: %v", detailResp)
}
}

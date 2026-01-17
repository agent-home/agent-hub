package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/agenthub/server/internal/config"
	"github.com/agenthub/server/internal/models"
	"github.com/agenthub/server/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

// Handler API 处理器
type Handler struct {
	cfg   *config.Config
	store *storage.Storage
}

// NewHandler 创建处理器
func NewHandler(cfg *config.Config, store *storage.Storage) *Handler {
	return &Handler{cfg: cfg, store: store}
}

// Health 健康检查
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"version": "1.0.0",
		"time":    time.Now().UTC(),
	})
}

// ===== 认证 =====

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=32"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name"`
}

// Register 用户注册
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 检查用户名是否存在
	if _, err := h.store.GetUserByUsername(ctx, req.Username); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	// 检查邮箱是否存在
	if _, err := h.store.GetUserByEmail(ctx, req.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.Username
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  displayName,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.store.CreateUser(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// 生成 token
	token, err := h.generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  sanitizeUser(user),
		"token": token,
	})
}

// LoginRequest 登录请求
type LoginRequest struct {
	Login    string `json:"login" binding:"required"` // 用户名或邮箱
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 查找用户
	var user *models.User
	var err error
	if user, err = h.store.GetUserByUsername(ctx, req.Login); err != nil {
		if user, err = h.store.GetUserByEmail(ctx, req.Login); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 生成 token
	token, err := h.generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  sanitizeUser(user),
		"token": token,
	})
}

// RefreshToken 刷新令牌
func (h *Handler) RefreshToken(c *gin.Context) {
	// 简化实现
	c.JSON(http.StatusOK, gin.H{"message": "refresh token"})
}

// Logout 登出
func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// generateToken 生成 JWT token
func (h *Handler) generateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(h.cfg.Auth.TokenExpiry) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "agenthub",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.Auth.JWTSecret))
}

// ===== 用户 =====

// GetUser 获取用户资料
func (h *Handler) GetUser(c *gin.Context) {
	username := c.Param("username")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := h.store.GetUserByUsername(ctx, username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, sanitizeUser(user))
}

// GetUserAgents 获取用户的智能体
func (h *Handler) GetUserAgents(c *gin.Context) {
	username := c.Param("username")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	agents, total, err := h.store.ListAgents(ctx, storage.ListAgentsOptions{
		Author:   username,
		Page:     1,
		PageSize: 20,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list agents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"agents": agents,
		"total":  total,
	})
}

// UpdateProfile 更新用户资料
func (h *Handler) UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

// ===== 智能体 =====

// ListAgents 列出智能体
func (h *Handler) ListAgents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	agents, total, err := h.store.ListAgents(ctx, storage.ListAgentsOptions{
		Page:     page,
		PageSize: pageSize,
		Category: c.Query("category"),
		Search:   c.Query("q"),
		Sort:     c.Query("sort"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list agents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"agents":    agents,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetAgent 获取智能体详情
func (h *Handler) GetAgent(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	agent, err := h.store.GetAgent(ctx, namespace, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	// 获取最新版本
	version, _ := h.store.GetLatestVersion(ctx, agent.ID)

	c.JSON(http.StatusOK, gin.H{
		"agent":          agent,
		"latest_version": version,
	})
}

// CreateAgentRequest 创建智能体请求
type CreateAgentRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	License     string   `json:"license"`
	Visibility  string   `json:"visibility"`
	Homepage    string   `json:"homepage"`
	Repository  string   `json:"repository"`
}

// CreateAgent 创建智能体
func (h *Handler) CreateAgent(c *gin.Context) {
	var req CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	username := c.GetString("username")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 检查是否已存在
	if _, err := h.store.GetAgent(ctx, username, req.Name); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "agent already exists"})
		return
	}

	visibility := req.Visibility
	if visibility == "" {
		visibility = "public"
	}

	agent := &models.Agent{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Namespace:   username,
		Description: req.Description,
		Category:    req.Category,
		Tags:        req.Tags,
		License:     req.License,
		Visibility:  visibility,
		AuthorID:    userID,
		Homepage:    req.Homepage,
		Repository:  req.Repository,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.store.CreateAgent(ctx, agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create agent"})
		return
	}

	agent.FullName = username + "/" + req.Name
	c.JSON(http.StatusCreated, agent)
}

// UpdateAgent 更新智能体
func (h *Handler) UpdateAgent(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	username := c.GetString("username")

	if namespace != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	agent, err := h.store.GetAgent(ctx, namespace, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	var req CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent.Description = req.Description
	agent.Category = req.Category
	agent.Tags = req.Tags
	agent.License = req.License
	if req.Visibility != "" {
		agent.Visibility = req.Visibility
	}
	agent.Homepage = req.Homepage
	agent.Repository = req.Repository

	if err := h.store.UpdateAgent(ctx, agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update agent"})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// DeleteAgent 删除智能体
func (h *Handler) DeleteAgent(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	username := c.GetString("username")

	if namespace != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	agent, err := h.store.GetAgent(ctx, namespace, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	if err := h.store.DeleteAgent(ctx, agent.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete agent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "agent deleted"})
}

// ===== 版本管理 =====

// ListVersions 列出版本
func (h *Handler) ListVersions(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	agent, err := h.store.GetAgent(ctx, namespace, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	versions, err := h.store.ListVersions(ctx, agent.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list versions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

// GetVersion 获取特定版本
func (h *Handler) GetVersion(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	versionTag := c.Param("version")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	agent, err := h.store.GetAgent(ctx, namespace, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	var version *models.AgentVersion
	if versionTag == "latest" {
		version, err = h.store.GetLatestVersion(ctx, agent.ID)
	} else {
		version, err = h.store.GetVersion(ctx, agent.ID, versionTag)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}

	// 增加下载次数
	h.store.IncrementDownloads(ctx, agent.ID)

	c.JSON(http.StatusOK, version)
}

// PublishVersionRequest 发布版本请求
type PublishVersionRequest struct {
	Version   string `json:"version" binding:"required"`
	Spec      string `json:"spec" binding:"required"`
	Changelog string `json:"changelog"`
}

// PublishVersion 发布新版本
func (h *Handler) PublishVersion(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	username := c.GetString("username")
	userID := c.GetString("user_id")

	if namespace != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	var req PublishVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证 spec
	var spec models.AgentSpec
	if err := yaml.Unmarshal([]byte(req.Spec), &spec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent spec: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	agent, err := h.store.GetAgent(ctx, namespace, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	version := &models.AgentVersion{
		ID:          uuid.New().String(),
		AgentID:     agent.ID,
		Version:     req.Version,
		Spec:        req.Spec,
		Changelog:   req.Changelog,
		IsLatest:    true,
		PublishedAt: time.Now(),
		PublishedBy: userID,
		Status:      "active",
		Size:        int64(len(req.Spec)),
	}

	if err := h.store.CreateVersion(ctx, version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish version"})
		return
	}

	c.JSON(http.StatusCreated, version)
}

// GetFile 获取文件
func (h *Handler) GetFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get file"})
}

// LikeAgent 点赞
func (h *Handler) LikeAgent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "liked"})
}

// UnlikeAgent 取消点赞
func (h *Handler) UnlikeAgent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "unliked"})
}

// ===== 搜索 =====

// Search 搜索
func (h *Handler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	agents, total, err := h.store.ListAgents(ctx, storage.ListAgentsOptions{
		Search:   query,
		Page:     1,
		PageSize: 20,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"agents": agents,
		"total":  total,
		"query":  query,
	})
}

// ListCategories 列出分类
func (h *Handler) ListCategories(c *gin.Context) {
	categories := []gin.H{
		{"id": "assistant", "name": "通用助手", "count": 0},
		{"id": "coding", "name": "编程开发", "count": 0},
		{"id": "writing", "name": "写作创作", "count": 0},
		{"id": "analysis", "name": "数据分析", "count": 0},
		{"id": "creative", "name": "创意设计", "count": 0},
		{"id": "education", "name": "教育学习", "count": 0},
		{"id": "business", "name": "商业办公", "count": 0},
		{"id": "research", "name": "研究探索", "count": 0},
		{"id": "tooling", "name": "工具效率", "count": 0},
		{"id": "other", "name": "其他", "count": 0},
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// GetTrending 获取热门
func (h *Handler) GetTrending(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	agents, _, _ := h.store.ListAgents(ctx, storage.ListAgentsOptions{
		Sort:     "downloads",
		Page:     1,
		PageSize: 10,
	})

	c.JSON(http.StatusOK, gin.H{"agents": agents})
}

// GetFeatured 获取推荐
func (h *Handler) GetFeatured(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	agents, _, _ := h.store.ListAgents(ctx, storage.ListAgentsOptions{
		Sort:     "likes",
		Page:     1,
		PageSize: 10,
	})

	c.JSON(http.StatusOK, gin.H{"agents": agents})
}

// ===== API Keys =====

// ListAPIKeys 列出 API Keys
func (h *Handler) ListAPIKeys(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"keys": []interface{}{}})
}

// CreateAPIKey 创建 API Key
func (h *Handler) CreateAPIKey(c *gin.Context) {
	userID := c.GetString("user_id")
	keyID := uuid.New().String()[:8]
	apiKey := "ak_" + userID + "_" + keyID

	c.JSON(http.StatusCreated, gin.H{
		"key":     apiKey,
		"message": "请保存此密钥，它不会再次显示",
	})
}

// DeleteAPIKey 删除 API Key
func (h *Handler) DeleteAPIKey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "key deleted"})
}

// ===== Agent 调用 =====

// InvokeAgent 调用智能体
func (h *Handler) InvokeAgent(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	agent, err := h.store.GetAgent(ctx, namespace, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	version, err := h.store.GetLatestVersion(ctx, agent.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no version available"})
		return
	}

	// 解析请求
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实际调用智能体
	c.JSON(http.StatusOK, gin.H{
		"agent":   agent.FullName,
		"version": version.Version,
		"input":   input,
		"output":  "智能体响应将在这里",
	})
}

// InvokeAgentStream 流式调用智能体
func (h *Handler) InvokeAgentStream(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	agent, err := h.store.GetAgent(ctx, namespace, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 模拟流式输出
	c.SSEvent("message", gin.H{"content": "正在调用 " + agent.FullName + "..."})
	c.Writer.Flush()
}

// sanitizeUser 清理用户敏感信息
func sanitizeUser(user *models.User) gin.H {
	return gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"display_name": user.DisplayName,
		"avatar":       user.Avatar,
		"bio":          user.Bio,
		"is_verified":  user.IsVerified,
		"created_at":   user.CreatedAt,
	}
}

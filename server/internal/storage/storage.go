package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/agenthub/server/internal/config"
	"github.com/agenthub/server/internal/models"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// Storage 存储接口
type Storage struct {
	db    *sql.DB
	redis *redis.Client
	cfg   *config.Config
}

// New 创建存储实例
func New(cfg *config.Config) (*Storage, error) {
	// 连接 PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 连接 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Storage{
		db:    db,
		redis: rdb,
		cfg:   cfg,
	}, nil
}

// Close 关闭连接
func (s *Storage) Close() error {
	if s.redis != nil {
		s.redis.Close()
	}
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// DB 获取数据库连接
func (s *Storage) DB() *sql.DB {
	return s.db
}

// Redis 获取 Redis 客户端
func (s *Storage) Redis() *redis.Client {
	return s.redis
}

// ===== Agent 操作 =====

// CreateAgent 创建智能体
func (s *Storage) CreateAgent(ctx context.Context, agent *models.Agent) error {
	query := `
		INSERT INTO agents (id, name, namespace, description, category, tags, license, visibility, author_id, homepage, repository, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := s.db.ExecContext(ctx, query,
		agent.ID, agent.Name, agent.Namespace, agent.Description, agent.Category,
		agent.Tags, agent.License, agent.Visibility, agent.AuthorID,
		agent.Homepage, agent.Repository, agent.CreatedAt, agent.UpdatedAt,
	)
	return err
}

// GetAgent 获取智能体
func (s *Storage) GetAgent(ctx context.Context, namespace, name string) (*models.Agent, error) {
	query := `
		SELECT id, name, namespace, description, category, tags, license, visibility, downloads, likes, author_id, homepage, repository, created_at, updated_at
		FROM agents
		WHERE namespace = $1 AND name = $2
	`
	agent := &models.Agent{}
	err := s.db.QueryRowContext(ctx, query, namespace, name).Scan(
		&agent.ID, &agent.Name, &agent.Namespace, &agent.Description, &agent.Category,
		&agent.Tags, &agent.License, &agent.Visibility, &agent.Downloads, &agent.Likes,
		&agent.AuthorID, &agent.Homepage, &agent.Repository, &agent.CreatedAt, &agent.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	agent.FullName = fmt.Sprintf("%s/%s", agent.Namespace, agent.Name)
	return agent, nil
}

// GetAgentByID 通过ID获取智能体
func (s *Storage) GetAgentByID(ctx context.Context, id string) (*models.Agent, error) {
	query := `
		SELECT id, name, namespace, description, category, tags, license, visibility, downloads, likes, author_id, homepage, repository, created_at, updated_at
		FROM agents
		WHERE id = $1
	`
	agent := &models.Agent{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&agent.ID, &agent.Name, &agent.Namespace, &agent.Description, &agent.Category,
		&agent.Tags, &agent.License, &agent.Visibility, &agent.Downloads, &agent.Likes,
		&agent.AuthorID, &agent.Homepage, &agent.Repository, &agent.CreatedAt, &agent.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	agent.FullName = fmt.Sprintf("%s/%s", agent.Namespace, agent.Name)
	return agent, nil
}

// ListAgents 列出智能体
func (s *Storage) ListAgents(ctx context.Context, opts ListAgentsOptions) ([]*models.Agent, int64, error) {
	// 构建查询
	baseQuery := `FROM agents WHERE visibility = 'public'`
	args := []interface{}{}
	argIndex := 1

	if opts.Category != "" {
		baseQuery += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, opts.Category)
		argIndex++
	}

	if opts.Search != "" {
		baseQuery += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+opts.Search+"%")
		argIndex++
	}

	if opts.Author != "" {
		baseQuery += fmt.Sprintf(" AND namespace = $%d", argIndex)
		args = append(args, opts.Author)
		argIndex++
	}

	// 获取总数
	var total int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 排序
	orderBy := " ORDER BY "
	switch opts.Sort {
	case "downloads":
		orderBy += "downloads DESC"
	case "likes":
		orderBy += "likes DESC"
	case "name":
		orderBy += "name ASC"
	default:
		orderBy += "updated_at DESC"
	}

	// 分页
	offset := (opts.Page - 1) * opts.PageSize
	pagination := fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, opts.PageSize, offset)

	// 查询列表
	listQuery := `SELECT id, name, namespace, description, category, tags, license, visibility, downloads, likes, author_id, homepage, repository, created_at, updated_at ` + baseQuery + orderBy + pagination

	rows, err := s.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var agents []*models.Agent
	for rows.Next() {
		agent := &models.Agent{}
		err := rows.Scan(
			&agent.ID, &agent.Name, &agent.Namespace, &agent.Description, &agent.Category,
			&agent.Tags, &agent.License, &agent.Visibility, &agent.Downloads, &agent.Likes,
			&agent.AuthorID, &agent.Homepage, &agent.Repository, &agent.CreatedAt, &agent.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		agent.FullName = fmt.Sprintf("%s/%s", agent.Namespace, agent.Name)
		agents = append(agents, agent)
	}

	return agents, total, nil
}

// ListAgentsOptions 列表选项
type ListAgentsOptions struct {
	Page     int
	PageSize int
	Category string
	Search   string
	Author   string
	Sort     string
}

// UpdateAgent 更新智能体
func (s *Storage) UpdateAgent(ctx context.Context, agent *models.Agent) error {
	query := `
		UPDATE agents
		SET description = $1, category = $2, tags = $3, license = $4, visibility = $5, homepage = $6, repository = $7, updated_at = $8
		WHERE id = $9
	`
	_, err := s.db.ExecContext(ctx, query,
		agent.Description, agent.Category, agent.Tags, agent.License, agent.Visibility,
		agent.Homepage, agent.Repository, time.Now(), agent.ID,
	)
	return err
}

// DeleteAgent 删除智能体
func (s *Storage) DeleteAgent(ctx context.Context, id string) error {
	query := `DELETE FROM agents WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// IncrementDownloads 增加下载次数
func (s *Storage) IncrementDownloads(ctx context.Context, agentID string) error {
	query := `UPDATE agents SET downloads = downloads + 1 WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, agentID)
	return err
}

// ===== Version 操作 =====

// CreateVersion 创建版本
func (s *Storage) CreateVersion(ctx context.Context, version *models.AgentVersion) error {
	// 先将其他版本的 is_latest 设为 false
	_, err := s.db.ExecContext(ctx, `UPDATE agent_versions SET is_latest = false WHERE agent_id = $1`, version.AgentID)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO agent_versions (id, agent_id, version, digest, size, spec, changelog, is_latest, published_at, published_by, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = s.db.ExecContext(ctx, query,
		version.ID, version.AgentID, version.Version, version.Digest, version.Size,
		version.Spec, version.Changelog, version.IsLatest, version.PublishedAt, version.PublishedBy, version.Status,
	)
	return err
}

// GetVersion 获取特定版本
func (s *Storage) GetVersion(ctx context.Context, agentID, version string) (*models.AgentVersion, error) {
	query := `
		SELECT id, agent_id, version, digest, size, spec, changelog, is_latest, published_at, published_by, downloads, status
		FROM agent_versions
		WHERE agent_id = $1 AND version = $2
	`
	v := &models.AgentVersion{}
	err := s.db.QueryRowContext(ctx, query, agentID, version).Scan(
		&v.ID, &v.AgentID, &v.Version, &v.Digest, &v.Size, &v.Spec, &v.Changelog,
		&v.IsLatest, &v.PublishedAt, &v.PublishedBy, &v.Downloads, &v.Status,
	)
	return v, err
}

// GetLatestVersion 获取最新版本
func (s *Storage) GetLatestVersion(ctx context.Context, agentID string) (*models.AgentVersion, error) {
	query := `
		SELECT id, agent_id, version, digest, size, spec, changelog, is_latest, published_at, published_by, downloads, status
		FROM agent_versions
		WHERE agent_id = $1 AND is_latest = true
	`
	v := &models.AgentVersion{}
	err := s.db.QueryRowContext(ctx, query, agentID).Scan(
		&v.ID, &v.AgentID, &v.Version, &v.Digest, &v.Size, &v.Spec, &v.Changelog,
		&v.IsLatest, &v.PublishedAt, &v.PublishedBy, &v.Downloads, &v.Status,
	)
	return v, err
}

// ListVersions 列出所有版本
func (s *Storage) ListVersions(ctx context.Context, agentID string) ([]*models.AgentVersion, error) {
	query := `
		SELECT id, agent_id, version, digest, size, spec, changelog, is_latest, published_at, published_by, downloads, status
		FROM agent_versions
		WHERE agent_id = $1
		ORDER BY published_at DESC
	`
	rows, err := s.db.QueryContext(ctx, query, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*models.AgentVersion
	for rows.Next() {
		v := &models.AgentVersion{}
		err := rows.Scan(
			&v.ID, &v.AgentID, &v.Version, &v.Digest, &v.Size, &v.Spec, &v.Changelog,
			&v.IsLatest, &v.PublishedAt, &v.PublishedBy, &v.Downloads, &v.Status,
		)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

// ===== User 操作 =====

// CreateUser 创建用户
func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, display_name, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.db.ExecContext(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.DisplayName,
		user.Status, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

// GetUserByUsername 通过用户名获取用户
func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar, bio, website, location, company, is_verified, is_admin, status, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	user := &models.User{}
	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DisplayName,
		&user.Avatar, &user.Bio, &user.Website, &user.Location, &user.Company,
		&user.IsVerified, &user.IsAdmin, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

// GetUserByEmail 通过邮箱获取用户
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar, bio, website, location, company, is_verified, is_admin, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &models.User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DisplayName,
		&user.Avatar, &user.Bio, &user.Website, &user.Location, &user.Company,
		&user.IsVerified, &user.IsAdmin, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

// GetUserByID 通过ID获取用户
func (s *Storage) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar, bio, website, location, company, is_verified, is_admin, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &models.User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DisplayName,
		&user.Avatar, &user.Bio, &user.Website, &user.Location, &user.Company,
		&user.IsVerified, &user.IsAdmin, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

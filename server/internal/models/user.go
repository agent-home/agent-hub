package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	DisplayName  string    `json:"display_name" db:"display_name"`
	Avatar       string    `json:"avatar,omitempty" db:"avatar"`
	Bio          string    `json:"bio,omitempty" db:"bio"`
	Website      string    `json:"website,omitempty" db:"website"`
	Location     string    `json:"location,omitempty" db:"location"`
	Company      string    `json:"company,omitempty" db:"company"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	IsAdmin      bool      `json:"is_admin" db:"is_admin"`
	Status       string    `json:"status" db:"status"` // active, suspended, deleted
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt  time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

// UserProfile 用户公开资料
type UserProfile struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	Avatar      string    `json:"avatar,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	Website     string    `json:"website,omitempty"`
	Location    string    `json:"location,omitempty"`
	Company     string    `json:"company,omitempty"`
	IsVerified  bool      `json:"is_verified"`
	AgentCount  int       `json:"agent_count"`
	Followers   int       `json:"followers"`
	Following   int       `json:"following"`
	CreatedAt   time.Time `json:"created_at"`
}

// Organization 组织模型
type Organization struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Description string    `json:"description,omitempty" db:"description"`
	Avatar      string    `json:"avatar,omitempty" db:"avatar"`
	Website     string    `json:"website,omitempty" db:"website"`
	Email       string    `json:"email,omitempty" db:"email"`
	IsVerified  bool      `json:"is_verified" db:"is_verified"`
	OwnerID     string    `json:"owner_id" db:"owner_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// OrgMember 组织成员
type OrgMember struct {
	ID        string    `json:"id" db:"id"`
	OrgID     string    `json:"org_id" db:"org_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Role      string    `json:"role" db:"role"` // owner, admin, member
	JoinedAt  time.Time `json:"joined_at" db:"joined_at"`
}

// APIKey API密钥
type APIKey struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	KeyPrefix   string    `json:"key_prefix" db:"key_prefix"` // 显示前几位
	KeyHash     string    `json:"-" db:"key_hash"`
	Scopes      []string  `json:"scopes" db:"scopes"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}

// AccessToken 访问令牌
type AccessToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"-" db:"token"`
	Type      string    `json:"type" db:"type"` // access, refresh
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserFollow 用户关注关系
type UserFollow struct {
	ID          string    `json:"id" db:"id"`
	FollowerID  string    `json:"follower_id" db:"follower_id"`
	FollowingID string    `json:"following_id" db:"following_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// AgentLike 智能体点赞
type AgentLike struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	AgentID   string    `json:"agent_id" db:"agent_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// AgentStar 智能体收藏
type AgentStar struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	AgentID   string    `json:"agent_id" db:"agent_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

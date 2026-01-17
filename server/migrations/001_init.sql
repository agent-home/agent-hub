-- AgentHub 数据库初始化脚本

-- 启用扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";  -- 用于模糊搜索

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(32) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(64),
    avatar VARCHAR(500),
    bio TEXT,
    website VARCHAR(255),
    location VARCHAR(100),
    company VARCHAR(100),
    is_verified BOOLEAN DEFAULT FALSE,
    is_admin BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- 组织表
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(32) UNIQUE NOT NULL,
    display_name VARCHAR(64),
    description TEXT,
    avatar VARCHAR(500),
    website VARCHAR(255),
    email VARCHAR(255),
    is_verified BOOLEAN DEFAULT FALSE,
    owner_id UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 组织成员表
CREATE TABLE IF NOT EXISTS org_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member',
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(org_id, user_id)
);

-- 智能体表
CREATE TABLE IF NOT EXISTS agents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(64) NOT NULL,
    namespace VARCHAR(32) NOT NULL,  -- 用户名或组织名
    description TEXT,
    category VARCHAR(32),
    tags TEXT[] DEFAULT '{}',
    license VARCHAR(32) DEFAULT 'MIT',
    visibility VARCHAR(20) DEFAULT 'public',
    downloads BIGINT DEFAULT 0,
    likes BIGINT DEFAULT 0,
    author_id UUID REFERENCES users(id),
    homepage VARCHAR(500),
    repository VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(namespace, name)
);

CREATE INDEX idx_agents_namespace ON agents(namespace);
CREATE INDEX idx_agents_category ON agents(category);
CREATE INDEX idx_agents_visibility ON agents(visibility);
CREATE INDEX idx_agents_downloads ON agents(downloads DESC);
CREATE INDEX idx_agents_likes ON agents(likes DESC);
CREATE INDEX idx_agents_name_trgm ON agents USING gin(name gin_trgm_ops);
CREATE INDEX idx_agents_description_trgm ON agents USING gin(description gin_trgm_ops);

-- 智能体版本表
CREATE TABLE IF NOT EXISTS agent_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    agent_id UUID REFERENCES agents(id) ON DELETE CASCADE,
    version VARCHAR(32) NOT NULL,
    digest VARCHAR(64),  -- SHA256 哈希
    size BIGINT DEFAULT 0,
    spec TEXT,  -- agentspec.yaml 内容
    changelog TEXT,
    is_latest BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    published_by UUID REFERENCES users(id),
    downloads BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    min_cli_version VARCHAR(16),
    UNIQUE(agent_id, version)
);

CREATE INDEX idx_versions_agent_id ON agent_versions(agent_id);
CREATE INDEX idx_versions_is_latest ON agent_versions(is_latest);

-- 智能体文件表
CREATE TABLE IF NOT EXISTS agent_files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    version_id UUID REFERENCES agent_versions(id) ON DELETE CASCADE,
    path VARCHAR(500) NOT NULL,
    size BIGINT DEFAULT 0,
    digest VARCHAR(64),
    mime_type VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(version_id, path)
);

CREATE INDEX idx_files_version_id ON agent_files(version_id);

-- API 密钥表
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    key_prefix VARCHAR(16),
    key_hash VARCHAR(255) NOT NULL,
    scopes TEXT[] DEFAULT '{}',
    expires_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);

-- 访问令牌表
CREATE TABLE IF NOT EXISTS access_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    type VARCHAR(20) DEFAULT 'access',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_tokens_user_id ON access_tokens(user_id);
CREATE INDEX idx_tokens_token ON access_tokens(token);

-- 用户关注表
CREATE TABLE IF NOT EXISTS user_follows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    follower_id UUID REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(follower_id, following_id)
);

CREATE INDEX idx_follows_follower ON user_follows(follower_id);
CREATE INDEX idx_follows_following ON user_follows(following_id);

-- 智能体点赞表
CREATE TABLE IF NOT EXISTS agent_likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    agent_id UUID REFERENCES agents(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, agent_id)
);

CREATE INDEX idx_likes_user ON agent_likes(user_id);
CREATE INDEX idx_likes_agent ON agent_likes(agent_id);

-- 智能体收藏表
CREATE TABLE IF NOT EXISTS agent_stars (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    agent_id UUID REFERENCES agents(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, agent_id)
);

CREATE INDEX idx_stars_user ON agent_stars(user_id);
CREATE INDEX idx_stars_agent ON agent_stars(agent_id);

-- 调用日志表 (用于计费和统计)
CREATE TABLE IF NOT EXISTS invocation_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    agent_id UUID REFERENCES agents(id),
    version_id UUID REFERENCES agent_versions(id),
    user_id UUID REFERENCES users(id),
    api_key_id UUID REFERENCES api_keys(id),
    input_tokens INT DEFAULT 0,
    output_tokens INT DEFAULT 0,
    latency_ms INT DEFAULT 0,
    status VARCHAR(20),
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_logs_agent_id ON invocation_logs(agent_id);
CREATE INDEX idx_logs_user_id ON invocation_logs(user_id);
CREATE INDEX idx_logs_created_at ON invocation_logs(created_at);

-- 更新时间触发器
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER agents_updated_at BEFORE UPDATE ON agents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER organizations_updated_at BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

package models

import (
	"time"
)

// Agent 智能体模型
type Agent struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Namespace   string    `json:"namespace" db:"namespace"` // 用户名或组织名
	FullName    string    `json:"full_name"`                // namespace/name
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"`
	Tags        []string  `json:"tags" db:"tags"`
	License     string    `json:"license" db:"license"`
	Visibility  string    `json:"visibility" db:"visibility"` // public, private, unlisted
	Downloads   int64     `json:"downloads" db:"downloads"`
	Likes       int64     `json:"likes" db:"likes"`
	AuthorID    string    `json:"author_id" db:"author_id"`
	Homepage    string    `json:"homepage,omitempty" db:"homepage"`
	Repository  string    `json:"repository,omitempty" db:"repository"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// AgentVersion 智能体版本
type AgentVersion struct {
	ID           string    `json:"id" db:"id"`
	AgentID      string    `json:"agent_id" db:"agent_id"`
	Version      string    `json:"version" db:"version"` // 语义化版本
	Digest       string    `json:"digest" db:"digest"`   // 内容哈希
	Size         int64     `json:"size" db:"size"`       // 字节数
	Spec         string    `json:"spec" db:"spec"`       // agentspec.yaml 内容
	Changelog    string    `json:"changelog,omitempty" db:"changelog"`
	IsLatest     bool      `json:"is_latest" db:"is_latest"`
	PublishedAt  time.Time `json:"published_at" db:"published_at"`
	PublishedBy  string    `json:"published_by" db:"published_by"`
	Downloads    int64     `json:"downloads" db:"downloads"`
	Status       string    `json:"status" db:"status"` // pending, active, deprecated
	MinCLIVersion string   `json:"min_cli_version,omitempty" db:"min_cli_version"`
}

// AgentFile 智能体文件
type AgentFile struct {
	ID        string    `json:"id" db:"id"`
	VersionID string    `json:"version_id" db:"version_id"`
	Path      string    `json:"path" db:"path"`
	Size      int64     `json:"size" db:"size"`
	Digest    string    `json:"digest" db:"digest"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// AgentSpec 解析后的智能体规范
type AgentSpec struct {
	Version      string                 `yaml:"version" json:"version"`
	Metadata     AgentMetadata          `yaml:"metadata" json:"metadata"`
	Runtime      AgentRuntime           `yaml:"runtime" json:"runtime"`
	Model        *AgentModel            `yaml:"model,omitempty" json:"model,omitempty"`
	Capabilities *AgentCapabilities     `yaml:"capabilities,omitempty" json:"capabilities,omitempty"`
	Interface    *AgentInterface        `yaml:"interface,omitempty" json:"interface,omitempty"`
	Prompts      *AgentPrompts          `yaml:"prompts,omitempty" json:"prompts,omitempty"`
	Resources    *AgentResources        `yaml:"resources,omitempty" json:"resources,omitempty"`
	Pricing      *AgentPricing          `yaml:"pricing,omitempty" json:"pricing,omitempty"`
	Workflow     *AgentWorkflow         `yaml:"workflow,omitempty" json:"workflow,omitempty"`
	Extra        map[string]interface{} `yaml:"-" json:"extra,omitempty"`
}

type AgentMetadata struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Author      string   `yaml:"author" json:"author"`
	License     string   `yaml:"license,omitempty" json:"license,omitempty"`
	Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	Homepage    string   `yaml:"homepage,omitempty" json:"homepage,omitempty"`
	Repository  string   `yaml:"repository,omitempty" json:"repository,omitempty"`
	Category    string   `yaml:"category,omitempty" json:"category,omitempty"`
}

type AgentRuntime struct {
	Type   string                 `yaml:"type" json:"type"`
	Entry  string                 `yaml:"entry,omitempty" json:"entry,omitempty"`
	Python *PythonRuntime         `yaml:"python,omitempty" json:"python,omitempty"`
	NodeJS *NodeJSRuntime         `yaml:"nodejs,omitempty" json:"nodejs,omitempty"`
	Docker *DockerRuntime         `yaml:"docker,omitempty" json:"docker,omitempty"`
	Remote *RemoteRuntime         `yaml:"remote,omitempty" json:"remote,omitempty"`
	Extra  map[string]interface{} `yaml:"-" json:"extra,omitempty"`
}

type PythonRuntime struct {
	Version      string `yaml:"version,omitempty" json:"version,omitempty"`
	Requirements string `yaml:"requirements,omitempty" json:"requirements,omitempty"`
}

type NodeJSRuntime struct {
	Version string `yaml:"version,omitempty" json:"version,omitempty"`
	Package string `yaml:"package,omitempty" json:"package,omitempty"`
}

type DockerRuntime struct {
	Image      string `yaml:"image,omitempty" json:"image,omitempty"`
	Dockerfile string `yaml:"dockerfile,omitempty" json:"dockerfile,omitempty"`
}

type RemoteRuntime struct {
	Endpoint string `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	Protocol string `yaml:"protocol,omitempty" json:"protocol,omitempty"`
}

type AgentModel struct {
	Provider   string            `yaml:"provider,omitempty" json:"provider,omitempty"`
	Name       string            `yaml:"name,omitempty" json:"name,omitempty"`
	Parameters map[string]interface{} `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

type AgentCapabilities struct {
	Streaming  bool            `yaml:"streaming,omitempty" json:"streaming,omitempty"`
	Multimodal *Multimodal     `yaml:"multimodal,omitempty" json:"multimodal,omitempty"`
	Tools      []AgentTool     `yaml:"tools,omitempty" json:"tools,omitempty"`
	Memory     *AgentMemory    `yaml:"memory,omitempty" json:"memory,omitempty"`
}

type Multimodal struct {
	Text  bool `yaml:"text,omitempty" json:"text,omitempty"`
	Image bool `yaml:"image,omitempty" json:"image,omitempty"`
	Audio bool `yaml:"audio,omitempty" json:"audio,omitempty"`
	Video bool `yaml:"video,omitempty" json:"video,omitempty"`
}

type AgentTool struct {
	Name        string                 `yaml:"name" json:"name"`
	Description string                 `yaml:"description" json:"description"`
	Parameters  map[string]interface{} `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

type AgentMemory struct {
	Conversation bool `yaml:"conversation,omitempty" json:"conversation,omitempty"`
	LongTerm     bool `yaml:"long_term,omitempty" json:"long_term,omitempty"`
	VectorStore  bool `yaml:"vector_store,omitempty" json:"vector_store,omitempty"`
}

type AgentInterface struct {
	Input  *InterfaceSchema `yaml:"input,omitempty" json:"input,omitempty"`
	Output *InterfaceSchema `yaml:"output,omitempty" json:"output,omitempty"`
}

type InterfaceSchema struct {
	Type   string                 `yaml:"type,omitempty" json:"type,omitempty"`
	Schema map[string]interface{} `yaml:"schema,omitempty" json:"schema,omitempty"`
}

type AgentPrompts struct {
	System     string    `yaml:"system,omitempty" json:"system,omitempty"`
	SystemFile string    `yaml:"system_file,omitempty" json:"system_file,omitempty"`
	Examples   []Example `yaml:"examples,omitempty" json:"examples,omitempty"`
}

type Example struct {
	Input  string `yaml:"input" json:"input"`
	Output string `yaml:"output" json:"output"`
}

type AgentResources struct {
	CPU     string `yaml:"cpu,omitempty" json:"cpu,omitempty"`
	Memory  string `yaml:"memory,omitempty" json:"memory,omitempty"`
	GPU     string `yaml:"gpu,omitempty" json:"gpu,omitempty"`
	Timeout int    `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

type AgentPricing struct {
	Model         string  `yaml:"model,omitempty" json:"model,omitempty"`
	PricePerCall  float64 `yaml:"price_per_call,omitempty" json:"price_per_call,omitempty"`
	PricePerToken float64 `yaml:"price_per_token,omitempty" json:"price_per_token,omitempty"`
}

type AgentWorkflow struct {
	Type   string         `yaml:"type,omitempty" json:"type,omitempty"`
	Agents []WorkflowAgent `yaml:"agents,omitempty" json:"agents,omitempty"`
	Steps  []WorkflowStep  `yaml:"steps,omitempty" json:"steps,omitempty"`
}

type WorkflowAgent struct {
	ID   string `yaml:"id" json:"id"`
	Ref  string `yaml:"ref" json:"ref"`
	Role string `yaml:"role,omitempty" json:"role,omitempty"`
}

type WorkflowStep struct {
	Name      string                 `yaml:"name" json:"name"`
	Agent     string                 `yaml:"agent" json:"agent"`
	Input     map[string]interface{} `yaml:"input,omitempty" json:"input,omitempty"`
	Output    string                 `yaml:"output,omitempty" json:"output,omitempty"`
	DependsOn []string               `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	OnFailure *FailureHandler        `yaml:"on_failure,omitempty" json:"on_failure,omitempty"`
}

type FailureHandler struct {
	Goto       string `yaml:"goto,omitempty" json:"goto,omitempty"`
	MaxRetries int    `yaml:"max_retries,omitempty" json:"max_retries,omitempty"`
}

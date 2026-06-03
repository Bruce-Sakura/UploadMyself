package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Agent 是核心对话引擎
type Agent struct {
	db         *gorm.DB
	llmClient  *LLMClient
	toolReg    *ToolRegistry
}

func New(db *gorm.DB, llmClient *LLMClient) *Agent {
	return &Agent{
		db:        db,
		llmClient: llmClient,
		toolReg:   NewToolRegistry(),
	}
}

// Chat 处理用户对话
func (a *Agent) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// 1. 加载用户的 SKILL.md 作为 system prompt
	systemPrompt, err := a.loadSkillPrompt(req.SkillID)
	if err != nil {
		return nil, fmt.Errorf("load skill: %w", err)
	}

	// 2. 加载对话历史
	history, err := a.loadHistory(req.ConversationID, 20)
	if err != nil {
		return nil, fmt.Errorf("load history: %w", err)
	}

	// 3. 构建消息
	messages := buildMessages(systemPrompt, history, req.Message)

	// 4. LLM 推理（支持工具调用循环）
	reply, toolCalls, err := a.llmClient.Chat(ctx, messages, a.toolReg.ToolDefs())
	if err != nil {
		return nil, fmt.Errorf("llm chat: %w", err)
	}

	// 5. 执行工具调用（如果有）
	var toolResults []ToolResult
	for len(toolCalls) > 0 {
		for _, tc := range toolCalls {
			result := a.toolReg.Execute(ctx, tc)
			toolResults = append(toolResults, result)
			messages = append(messages, ToolResultMessage(tc.ID, result.Content))
		}
		// 把工具结果发回 LLM 继续推理
		reply, toolCalls, err = a.llmClient.Chat(ctx, messages, a.toolReg.ToolDefs())
		if err != nil {
			return nil, fmt.Errorf("llm chat after tools: %w", err)
		}
	}

	// 6. 保存对话历史
	a.saveMessages(req.ConversationID, req.Message, reply)

	return &ChatResponse{
		Reply:       reply,
		ToolCalls:   toolResults,
		Timestamp:   time.Now(),
	}, nil
}

// loadSkillPrompt 从数据库加载 SKILL.md 内容作为 system prompt
// LLMChat 简单的单轮 LLM 调用（用于 Skill 生成等）
func (a *Agent) LLMChat(ctx context.Context, userPrompt string) (string, error) {
	messages := []llmMessage{
		{Role: "user", Content: userPrompt},
	}
	reply, _, err := a.llmClient.Chat(ctx, messages, nil)
	return reply, err
}

func (a *Agent) loadSkillPrompt(skillID string) (string, error) {
	if skillID == "" {
		return "你是一个有帮助的AI助手。", nil
	}
	
	var skill struct {
		Result string
	}
	if err := a.db.Table("skills").Select("result").Where("id = ?", skillID).First(&skill).Error; err != nil {
		return "", err
	}
	
	if skill.Result == "" {
		return "你是一个有帮助的AI助手。", nil
	}
	
	return skill.Result, nil
}

// loadHistory 加载对话历史
func (a *Agent) loadHistory(convID string, limit int) ([]Message, error) {
	if convID == "" {
		return nil, nil
	}
	var msgs []Message
	err := a.db.Where("conversation_id = ?", convID).
		Order("created_at desc").
		Limit(limit).
		Find(&msgs).Error
	// 反转为时间正序
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, err
}

// saveMessages 保存对话消息
func (a *Agent) saveMessages(convID, userMsg, assistantMsg string) {
	if convID == "" {
		return
	}
	a.db.Create(&Message{ConversationID: convID, Role: "user", Content: userMsg})
	a.db.Create(&Message{ConversationID: convID, Role: "assistant", Content: assistantMsg})
}

// ==================== LLM Client ====================

type LLMClient struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

func NewLLMClient(apiKey, baseURL, model string) *LLMClient {
	return &LLMClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

type llmMessage struct {
	Role       string      `json:"role"`
	Content    string      `json:"content,omitempty"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
}

type llmRequest struct {
	Model    string      `json:"model"`
	Messages []llmMessage `json:"messages"`
	Tools    []ToolDef   `json:"tools,omitempty"`
	Stream   bool        `json:"stream"`
}

type llmResponse struct {
	Choices []struct {
		Message struct {
			Content   string     `json:"content"`
			ToolCalls []ToolCall `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
}

// Chat 调用 LLM API（兼容 OpenAI 格式）
func (c *LLMClient) Chat(ctx context.Context, messages []llmMessage, tools []ToolDef) (string, []ToolCall, error) {
	reqBody := llmRequest{
		Model:    c.model,
		Messages: messages,
		Tools:    tools,
		Stream:   false,
	}
	
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, err
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", nil, fmt.Errorf("llm api error %d: %s", resp.StatusCode, string(respBody))
	}
	
	var llmResp llmResponse
	if err := json.Unmarshal(respBody, &llmResp); err != nil {
		return "", nil, err
	}
	
	if len(llmResp.Choices) == 0 {
		return "", nil, fmt.Errorf("no choices in response")
	}
	
	msg := llmResp.Choices[0].Message
	return msg.Content, msg.ToolCalls, nil
}

// ==================== Tool System ====================

type ToolDef struct {
	Type     string `json:"type"`
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  any    `json:"parameters"`
	} `json:"function"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Content    string `json:"content"`
}

type ToolHandler func(ctx context.Context, args map[string]interface{}) string

type ToolRegistry struct {
	tools map[string]ToolHandler
	defs  []ToolDef
}

func NewToolRegistry() *ToolRegistry {
	r := &ToolRegistry{
		tools: make(map[string]ToolHandler),
	}
	// 注册内置工具
	r.RegisterShell()
	r.RegisterFileRead()
	r.RegisterFileWrite()
	r.RegisterWebSearch()
	return r
}

func (r *ToolRegistry) Register(name, description string, params any, handler ToolHandler) {
	r.tools[name] = handler
	def := ToolDef{Type: "function"}
	def.Function.Name = name
	def.Function.Description = description
	def.Function.Parameters = params
	r.defs = append(r.defs, def)
}

func (r *ToolRegistry) ToolDefs() []ToolDef {
	return r.defs
}

func (r *ToolRegistry) Execute(ctx context.Context, tc ToolCall) ToolResult {
	handler, ok := r.tools[tc.Function.Name]
	if !ok {
		return ToolResult{ToolCallID: tc.ID, Content: "未知工具: " + tc.Function.Name}
	}
	
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
		return ToolResult{ToolCallID: tc.ID, Content: "参数解析失败: " + err.Error()}
	}
	
	result := handler(ctx, args)
	return ToolResult{ToolCallID: tc.ID, Content: result}
}

// ==================== Built-in Tools ====================

func (r *ToolRegistry) RegisterShell() {
	r.Register("shell", "执行 Shell 命令并返回输出", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]string{"type": "string", "description": "要执行的命令"},
		},
		"required": []string{"command"},
	}, func(ctx context.Context, args map[string]interface{}) string {
		command, _ := args["command"].(string)
		if command == "" {
			return "错误：缺少 command 参数"
		}
		// 安全检查
		dangerous := []string{"rm -rf /", "mkfs", "dd if=", ":(){ :|:& };:"}
		for _, d := range dangerous {
			if strings.Contains(command, d) {
				return "拒绝执行危险命令"
			}
		}
		// TODO: 实际执行
		return fmt.Sprintf("[Shell] 命令已接收: %s\n(执行功能待实现)", command)
	})
}

func (r *ToolRegistry) RegisterFileRead() {
	r.Register("file_read", "读取文件内容", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]string{"type": "string", "description": "文件路径"},
		},
		"required": []string{"path"},
	}, func(ctx context.Context, args map[string]interface{}) string {
		path, _ := args["path"].(string)
		// TODO: 实际读取
		return fmt.Sprintf("[FileRead] 路径已接收: %s\n(读取功能待实现)", path)
	})
}

func (r *ToolRegistry) RegisterFileWrite() {
	r.Register("file_write", "写入文件内容", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path":    map[string]string{"type": "string", "description": "文件路径"},
			"content": map[string]string{"type": "string", "description": "文件内容"},
		},
		"required": []string{"path", "content"},
	}, func(ctx context.Context, args map[string]interface{}) string {
		path, _ := args["path"].(string)
		return fmt.Sprintf("[FileWrite] 路径已接收: %s\n(写入功能待实现)", path)
	})
}

func (r *ToolRegistry) RegisterWebSearch() {
	r.Register("web_search", "搜索互联网获取信息", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]string{"type": "string", "description": "搜索关键词"},
		},
		"required": []string{"query"},
	}, func(ctx context.Context, args map[string]interface{}) string {
		query, _ := args["query"].(string)
		return fmt.Sprintf("[WebSearch] 搜索: %s\n(搜索功能待实现)", query)
	})
}

// ==================== Helper ====================

func buildMessages(systemPrompt string, history []Message, userMsg string) []llmMessage {
	var msgs []llmMessage
	
	// System prompt (SKILL.md 内容)
	if systemPrompt != "" {
		msgs = append(msgs, llmMessage{Role: "system", Content: systemPrompt})
	}
	
	// 历史消息
	for _, h := range history {
		msgs = append(msgs, llmMessage{Role: h.Role, Content: h.Content})
	}
	
	// 当前用户消息
	msgs = append(msgs, llmMessage{Role: "user", Content: userMsg})
	
	return msgs
}

func ToolResultMessage(toolCallID, content string) llmMessage {
	return llmMessage{
		Role:       "tool",
		ToolCallID: toolCallID,
		Content:    content,
	}
}

// ==================== Data Models ====================

type ChatRequest struct {
	ConversationID string `json:"conversation_id"`
	SkillID        string `json:"skill_id"`
	Message        string `json:"message" binding:"required"`
}

type ChatResponse struct {
	Reply     string       `json:"reply"`
	ToolCalls []ToolResult `json:"tool_calls,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}

type Message struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ConversationID string    `json:"conversation_id" gorm:"index"`
	Role           string    `json:"role"`
	Content        string    `json:"content" gorm:"type:text"`
	CreatedAt      time.Time `json:"created_at"`
}

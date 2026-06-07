package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Bruce-Sakura/UploadMyself/backend/internal/llm"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/messages/dto"
)

type toolHandler func(ctx context.Context, args map[string]any) string

// toolRegistry holds the agent's callable tools and their schemas.
type toolRegistry struct {
	tools map[string]toolHandler
	defs  []llm.ToolDef
}

func newToolRegistry() *toolRegistry {
	r := &toolRegistry{tools: make(map[string]toolHandler)}
	r.registerShell()
	r.registerFileRead()
	r.registerFileWrite()
	r.registerWebSearch()
	return r
}

func (r *toolRegistry) register(name, description string, params any, h toolHandler) {
	r.tools[name] = h
	def := llm.ToolDef{Type: "function"}
	def.Function.Name = name
	def.Function.Description = description
	def.Function.Parameters = params
	r.defs = append(r.defs, def)
}

func (r *toolRegistry) toolDefs() []llm.ToolDef { return r.defs }

func (r *toolRegistry) info() []dto.ToolInfoVO {
	out := make([]dto.ToolInfoVO, 0, len(r.defs))
	for _, d := range r.defs {
		out = append(out, dto.ToolInfoVO{Name: d.Function.Name, Description: d.Function.Description})
	}
	return out
}

func (r *toolRegistry) execute(ctx context.Context, tc llm.ToolCall) dto.ToolResultVO {
	h, ok := r.tools[tc.Function.Name]
	if !ok {
		return dto.ToolResultVO{ToolCallID: tc.ID, Content: "未知工具: " + tc.Function.Name}
	}
	var args map[string]any
	if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
		return dto.ToolResultVO{ToolCallID: tc.ID, Content: "参数解析失败: " + err.Error()}
	}
	return dto.ToolResultVO{ToolCallID: tc.ID, Content: h(ctx, args)}
}

// ---- built-in tools (stubs; execution layer TODO) ----

func (r *toolRegistry) registerShell() {
	r.register("shell", "执行 Shell 命令并返回输出", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]string{"type": "string", "description": "要执行的命令"},
		},
		"required": []string{"command"},
	}, func(ctx context.Context, args map[string]any) string {
		command, _ := args["command"].(string)
		if command == "" {
			return "错误：缺少 command 参数"
		}
		dangerous := []string{"rm -rf /", "mkfs", "dd if=", ":(){ :|:& };:"}
		for _, d := range dangerous {
			if strings.Contains(command, d) {
				return "拒绝执行危险命令"
			}
		}
		return fmt.Sprintf("[Shell] 命令已接收: %s\n(执行功能待实现)", command)
	})
}

func (r *toolRegistry) registerFileRead() {
	r.register("file_read", "读取文件内容", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]string{"type": "string", "description": "文件路径"},
		},
		"required": []string{"path"},
	}, func(ctx context.Context, args map[string]any) string {
		path, _ := args["path"].(string)
		return fmt.Sprintf("[FileRead] 路径已接收: %s\n(读取功能待实现)", path)
	})
}

func (r *toolRegistry) registerFileWrite() {
	r.register("file_write", "写入文件内容", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":    map[string]string{"type": "string", "description": "文件路径"},
			"content": map[string]string{"type": "string", "description": "文件内容"},
		},
		"required": []string{"path", "content"},
	}, func(ctx context.Context, args map[string]any) string {
		path, _ := args["path"].(string)
		return fmt.Sprintf("[FileWrite] 路径已接收: %s\n(写入功能待实现)", path)
	})
}

func (r *toolRegistry) registerWebSearch() {
	r.register("web_search", "搜索互联网获取信息", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"query": map[string]string{"type": "string", "description": "搜索关键词"},
		},
		"required": []string{"query"},
	}, func(ctx context.Context, args map[string]any) string {
		query, _ := args["query"].(string)
		return fmt.Sprintf("[WebSearch] 搜索: %s\n(搜索功能待实现)", query)
	})
}

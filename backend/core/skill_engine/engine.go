package core

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// SkillEngine 思维框架克隆引擎（仿女娲）
type SkillEngine struct {
	llmProvider LLMProvider
}

// LLMProvider LLM 调用接口
type LLMProvider interface {
	Chat(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

func NewSkillEngine(llm LLMProvider) *SkillEngine {
	return &SkillEngine{llmProvider: llm}
}

// Generate 执行完整的 Skill 生成流程
func (e *SkillEngine) Generate(ctx context.Context, corpusPath, name string) (string, error) {
	// Phase 1: 读取并清洗语料
	corpus, err := os.ReadFile(corpusPath)
	if err != nil {
		return "", fmt.Errorf("读取语料失败: %w", err)
	}

	chunks := chunkText(string(corpus), 4000)

	// Phase 2: 分析思维模式
	analysis, err := e.analyzeCorpus(ctx, chunks)
	if err != nil {
		return "", fmt.Errorf("分析语料失败: %w", err)
	}

	// Phase 3: 提取心智模型
	models, err := e.extractModels(ctx, analysis)
	if err != nil {
		return "", fmt.Errorf("提取心智模型失败: %w", err)
	}

	// Phase 4: 合成 SKILL.md
	skillMD, err := e.synthesize(ctx, name, analysis, models)
	if err != nil {
		return "", fmt.Errorf("合成 Skill 失败: %w", err)
	}

	return skillMD, nil
}

func (e *SkillEngine) analyzeCorpus(ctx context.Context, chunks []string) (string, error) {
	systemPrompt := `你是一个思维模式分析专家。分析以下文本，提取：
1. 核心论点（反复出现≥3次的观点）
2. 表达风格（句式、词汇、节奏特征）
3. 决策模式（如何做判断）
4. 价值观信号
输出结构化分析报告。`

	var analyses []string
	for i, chunk := range chunks {
		result, err := e.llmProvider.Chat(ctx, systemPrompt, chunk)
		if err != nil {
			return "", fmt.Errorf("分析第 %d 段失败: %w", i+1, err)
		}
		analyses = append(analyses, result)
	}

	return strings.Join(analyses, "\n\n---\n\n"), nil
}

func (e *SkillEngine) extractModels(ctx context.Context, analysis string) (string, error) {
	systemPrompt := `基于以下思维模式分析，提取 3-7 个核心心智模型。
每个模型需要：
- 名称（一句话）
- 证据（来自原文的具体引用）
- 应用场景
- 局限性
用三重验证筛选：跨域复现、生成力、排他性。`

	result, err := e.llmProvider.Chat(ctx, systemPrompt, analysis)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (e *SkillEngine) synthesize(ctx context.Context, name, analysis, models string) (string, error) {
	systemPrompt := `你是一个 Skill 构建专家。基于以下分析结果，生成一个完整的 SKILL.md 文件。

格式要求：
- YAML frontmatter（name + description）
- 角色扮演规则
- 身份卡（50字自我介绍）
- 核心心智模型（3-7个，含证据/应用/局限）
- 决策启发式（5-10条）
- 表达DNA（句式/词汇/节奏/幽默/确定性）
- 诚实边界（至少3条具体局限）

用 Markdown 格式输出完整 SKILL.md。`

	input := fmt.Sprintf("人物名称: %s\n\n## 思维模式分析\n%s\n\n## 心智模型\n%s", name, analysis, models)
	result, err := e.llmProvider.Chat(ctx, systemPrompt, input)
	if err != nil {
		return "", err
	}
	return result, nil
}

// chunkText 将长文本按最大字符数分段
func chunkText(text string, maxChars int) []string {
	paragraphs := strings.Split(text, "\n\n")
	var chunks []string
	current := ""

	for _, para := range paragraphs {
		if len(current)+len(para) > maxChars {
			if current != "" {
				chunks = append(chunks, strings.TrimSpace(current))
			}
			current = para
		} else {
			current += "\n\n" + para
		}
	}
	if strings.TrimSpace(current) != "" {
		chunks = append(chunks, strings.TrimSpace(current))
	}
	return chunks
}

package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/Bruce-Sakura/UploadMyself/backend/internal/llm"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/dto"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/entity"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/mapper"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/service"
	taskdto "github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/dto"
	taskservice "github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/service"
	"github.com/google/uuid"
)

// Config holds runtime settings for skill file storage.
type Config struct {
	SkillsDir string // root dir for skill packages: <SkillsDir>/<id>/SKILL.md
}

type SkillServiceImpl struct {
	mapper  *mapper.SkillMapper
	taskSvc taskservice.TaskService
	llm     *llm.Client
	cfg     Config
}

func NewSkillService(m *mapper.SkillMapper, taskSvc taskservice.TaskService, llmClient *llm.Client, cfg Config) service.SkillService {
	return &SkillServiceImpl{mapper: m, taskSvc: taskSvc, llm: llmClient, cfg: cfg}
}

func (s *SkillServiceImpl) Create(ctx context.Context, req dto.CreateSkillReq) (*dto.SkillVO, error) {
	sk := &entity.Skill{
		ID:     uuid.NewString(),
		Name:   req.Name,
		Corpus: req.Corpus,
		Status: "pending",
	}
	if err := s.mapper.Insert(ctx, sk); err != nil {
		return nil, err
	}
	return s.Get(ctx, sk.ID)
}

func (s *SkillServiceImpl) Get(ctx context.Context, id string) (*dto.SkillVO, error) {
	sk, err := s.mapper.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	vo := toVO(sk)
	vo.Result = s.readSkillFile(id) // SKILL.md lives on disk
	return vo, nil
}

func (s *SkillServiceImpl) List(ctx context.Context) ([]dto.SkillVO, error) {
	sks, err := s.mapper.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.SkillVO, 0, len(sks))
	for i := range sks {
		vo := toVO(&sks[i])
		vo.Result = s.readSkillFile(vo.ID)
		out = append(out, *vo)
	}
	return out, nil
}

func (s *SkillServiceImpl) Update(ctx context.Context, id string, req dto.UpdateSkillReq) (*dto.SkillVO, error) {
	// SKILL.md content is stored on disk, not in DB.
	if req.Result != nil {
		if err := s.writeSkillFile(id, *req.Result); err != nil {
			return nil, err
		}
		req.Result = nil
	}
	if err := s.mapper.Update(ctx, id, req.Name, req.Corpus, req.Status, nil); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

func (s *SkillServiceImpl) Delete(ctx context.Context, id string) error {
	if err := s.mapper.Delete(ctx, id); err != nil {
		return err
	}
	return s.removeSkillDir(id)
}

func (s *SkillServiceImpl) Process(ctx context.Context, id string) (string, error) {
	sk, err := s.mapper.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	task, err := s.taskSvc.Create(ctx, taskdto.CreateTaskReq{Type: "skill_process", RefID: id})
	if err != nil {
		return "", err
	}
	processing := "processing"
	_ = s.mapper.Update(ctx, id, nil, nil, &processing, nil)

	go s.runProcess(task.ID, sk.ID, sk.Name, sk.Corpus)
	return task.ID, nil
}

func (s *SkillServiceImpl) runProcess(taskID, skillID, name, corpus string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	_ = s.taskSvc.UpdateStatus(ctx, taskID, "running", 10, "")

	prompt := fmt.Sprintf(`你是一个思维框架分析专家。分析以下用户的文本语料，生成一个完整的 SKILL.md 文件。

用户名称: %s

文本语料:
%s

请生成 SKILL.md，包含：
1. 身份卡（用第一人称写50字自我介绍）
2. 核心智模型（3-5个，含证据和应用）
3. 决策启发式（3-5条）
4. 表达DNA（句式/词汇/幽默风格）
5. 诚实边界（局限性）

用 Markdown 格式输出。`, name, corpus)

	reply, err := s.llm.ChatOnce(ctx, prompt)
	if err != nil {
		failed := "failed"
		_ = s.mapper.Update(ctx, skillID, nil, nil, &failed, nil)
		_ = s.taskSvc.UpdateStatus(ctx, taskID, "failed", 0, fmt.Sprintf("LLM error: %v", err))
		return
	}

	// Persist SKILL.md to disk (file storage), DB only tracks status.
	if err := s.writeSkillFile(skillID, reply); err != nil {
		failed := "failed"
		_ = s.mapper.Update(ctx, skillID, nil, nil, &failed, nil)
		_ = s.taskSvc.UpdateStatus(ctx, taskID, "failed", 0, fmt.Sprintf("write skill file: %v", err))
		return
	}

	done := "done"
	_ = s.mapper.Update(ctx, skillID, nil, nil, &done, nil)
	_ = s.taskSvc.UpdateStatus(ctx, taskID, "done", 100, "")
}

func (s *SkillServiceImpl) Import(ctx context.Context, req dto.ImportSkillReq) (*dto.SkillVO, error) {
	content := req.Content
	source := "inline"
	if content == "" {
		if req.URL == "" {
			return nil, fmt.Errorf("either url or content is required")
		}
		body, err := downloadSkill(ctx, req.URL)
		if err != nil {
			return nil, err
		}
		content = body
		source = req.URL
	}
	if content == "" {
		return nil, fmt.Errorf("downloaded skill is empty")
	}

	// Derive a name: explicit > frontmatter name > fallback.
	name := req.Name
	if name == "" {
		name = parseFrontmatterName(content)
	}
	if name == "" {
		name = "imported-skill"
	}

	id := uuid.NewString()
	if err := s.writeSkillFile(id, content); err != nil {
		return nil, err
	}
	if err := s.writeSkillMeta(id, name, source); err != nil {
		return nil, err
	}

	sk := &entity.Skill{ID: id, Name: name, Corpus: "", Status: "done"}
	if err := s.mapper.Insert(ctx, sk); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

func toVO(s *entity.Skill) *dto.SkillVO {
	return &dto.SkillVO{
		ID:        s.ID,
		Name:      s.Name,
		Corpus:    s.Corpus,
		Status:    s.Status,
		Result:    s.Result,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

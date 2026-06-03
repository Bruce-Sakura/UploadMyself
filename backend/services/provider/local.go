package provider

import (
	"context"
	"fmt"
)

// LocalProvider 本地 GPU 推理
type LocalProvider struct {
	modelDir string
}

func NewLocalProvider(modelDir string) *LocalProvider {
	return &LocalProvider{modelDir: modelDir}
}

func (p *LocalProvider) Inference(ctx context.Context, req *InferenceRequest) (*InferenceResponse, error) {
	switch req.TaskType {
	case "voice":
		return p.voiceInference(ctx, req)
	case "avatar_2d":
		return p.avatar2DInference(ctx, req)
	case "avatar_3d":
		return p.avatar3DInference(ctx, req)
	default:
		return nil, fmt.Errorf("本地 Provider 不支持的任务类型: %s", req.TaskType)
	}
}

func (p *LocalProvider) voiceInference(ctx context.Context, req *InferenceRequest) (*InferenceResponse, error) {
	// TODO: GPT-SoVITS / CosyVoice 推理
	return &InferenceResponse{
		Output: map[string]interface{}{"status": "not_implemented"},
	}, nil
}

func (p *LocalProvider) avatar2DInference(ctx context.Context, req *InferenceRequest) (*InferenceResponse, error) {
	// TODO: LivePortrait / SadTalker 推理
	return &InferenceResponse{
		Output: map[string]interface{}{"status": "not_implemented"},
	}, nil
}

func (p *LocalProvider) avatar3DInference(ctx context.Context, req *InferenceRequest) (*InferenceResponse, error) {
	// TODO: InstantMesh 推理
	return &InferenceResponse{
		Output: map[string]interface{}{"status": "not_implemented"},
	}, nil
}

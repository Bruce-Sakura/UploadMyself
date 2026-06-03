package provider

import "context"

// Provider 统一模型提供者接口
type Provider interface {
	// Inference 执行推理
	Inference(ctx context.Context, req *InferenceRequest) (*InferenceResponse, error)
}

// InferenceRequest 通用推理请求
type InferenceRequest struct {
	TaskType string                 // skill | voice | avatar_2d | avatar_3d | distill
	Model    string                 // 模型名称
	Input    map[string]interface{} // 输入参数
	Options  map[string]interface{} // 选项（温度、步数等）
}

// InferenceResponse 通用推理响应
type InferenceResponse struct {
	Output map[string]interface{} // 输出结果
	Meta   map[string]interface{} // 元信息（耗时、模型版本等）
}

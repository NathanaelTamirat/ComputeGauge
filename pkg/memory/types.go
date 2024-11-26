package memory

import (
	"compute-gauge/internal/config"
	"compute-gauge/pkg/gpu"
	"encoding/json"
	"html/template"
)

type PageData struct {
	Models    map[string]config.ModelConfig
	DataTypes []string
}

func (p PageData) ModelsJSON() template.JS {
	data, err := json.Marshal(p.Models)
	if err != nil {
		return ""
	}
	return template.JS(data)
}

type MemoryResponse struct {
	ModelWeights     string                  `json:"model_weights"`
	KVCache          string                  `json:"kv_cache"`
	ActivationMemory string                  `json:"activation_memory"`
	OptimizerMemory  string                  `json:"optimizer_memory,omitempty"`
	GradientsMemory  string                  `json:"gradients_memory,omitempty"`
	InferenceMemory  string                  `json:"inference_memory"`
	TrainingMemory   string                  `json:"training_memory,omitempty"`
	InferenceGPUs    []gpu.GPURecommendation `json:"inference_gpus"`
	TrainingGPUs     []gpu.GPURecommendation `json:"training_gpus,omitempty"`
	TotalParams      float64                 `json:"total_params"`
	HiddenSize       int                     `json:"hidden_size"`
	SequenceLength   int                     `json:"sequence_length"`
}

type MemoryRequest struct {
	ModelSize         float64 `json:"model_size"`
	HiddenSize        int     `json:"hidden_size"`
	NumHiddenLayers   int     `json:"num_hidden_layers"`
	NumAttentionHeads int     `json:"num_attention_heads"`
	SequenceLength    int     `json:"sequence_length"`
	BatchSize         int     `json:"batch_size"`
	TorchDtype        string  `json:"torch_dtype"`
	Optimizer         string  `json:"optimizer"`
}

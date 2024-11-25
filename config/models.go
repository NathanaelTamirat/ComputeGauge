package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var DataTypeSizes = map[string]float64{
	"float32":  4.0,
	"float16":  2.0,
	"bfloat16": 2.0,
	"int8":     1.0,
	"int4":     0.5,
}

type ModelConfig struct {
	Name              string  `json:"name"`
	ModelSize         float64 `json:"model_size"`
	HiddenSize        int     `json:"hidden_size"`
	NumHiddenLayers   int     `json:"num_hidden_layers"`
	NumAttentionHeads int     `json:"num_attention_heads"`
	NumKeyValueHeads  int     `json:"num_key_value_heads"`
	SequenceLength    int     `json:"max_position_embeddings"`
	Precision         string  `json:"torch_dtype"`
}
type MemoryRequest struct {
	ModelSize         float64 `json:"model_size"`
	BatchSize         int     `json:"batch_size"`
	SequenceLength    int     `json:"sequence_length"`
	HiddenSize        int     `json:"hidden_size"`
	NumHiddenLayers   int     `json:"num_hidden_layers"`
	NumAttentionHeads int     `json:"num_attention_heads"`
	NumKeyValueHeads  int     `json:"num_key_value_heads"`
	Optimizer         string  `json:"optimizer,omitempty"`
	TrainableParams   float64 `json:"trainable_params,omitempty"`
	Precision         string  `json:"precision"`
}

type MemoryResponse struct {
	ModelWeights     string `json:"model_weights"`
	KVCache          string `json:"kv_cache"`
	ActivationMemory string `json:"activation_memory"`
	OptimizerMemory  string `json:"optimizer_memory,omitempty"`
	GradientsMemory  string `json:"gradients_memory,omitempty"`
	TotalMemory      string `json:"total_memory"`
}

func LoadModelConfigs() (map[string]ModelConfig, error) {
	models := make(map[string]ModelConfig)
	modelsDir := "models"

	files, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, fmt.Errorf("error reading models directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			data, err := os.ReadFile(filepath.Join(modelsDir, file.Name()))
			if err != nil {
				continue
			}

			var config ModelConfig
			if err := json.Unmarshal(data, &config); err != nil {
				continue
			}

			modelName := strings.TrimSuffix(file.Name(), ".json")
			config.Name = modelName
			models[modelName] = config
		}
	}

	return models, nil
}

package config

import (
	"encoding/json"
	"fmt"
	"log"
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

// Predefined model configurations
var predefinedModels = []string{
	`{
		"name": "llama2-7b",
		"model_size": 7,
		"hidden_size": 4096,
		"num_hidden_layers": 32,
		"num_attention_heads": 32,
		"num_key_value_heads": 32,
		"max_position_embeddings": 4096,
		"torch_dtype": "float16"
	}`,
	`{
		"name": "llama2-13b",
		"model_size": 13,
		"hidden_size": 5120,
		"num_hidden_layers": 40,
		"num_attention_heads": 40,
		"num_key_value_heads": 40,
		"max_position_embeddings": 4096,
		"torch_dtype": "float16"
	}`,
	`{
		"name": "llama2-70b",
		"model_size": 70,
		"hidden_size": 8192,
		"num_hidden_layers": 80,
		"num_attention_heads": 64,
		"num_key_value_heads": 64,
		"max_position_embeddings": 4096,
		"torch_dtype": "float16"
	}`,
}

func LoadModelConfigs() (map[string]ModelConfig, error) {
	models := make(map[string]ModelConfig)

	// Load predefined models
	for _, modelJSON := range predefinedModels {
		var config ModelConfig
		if err := json.Unmarshal([]byte(modelJSON), &config); err != nil {
			log.Printf("Error parsing predefined model: %v", err)
			continue
		}
		models[config.Name] = config
		log.Printf("Loaded predefined model: %s", config.Name)
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("no valid model configurations found")
	}

	return models, nil
}

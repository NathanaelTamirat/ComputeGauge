package config

import (
	"encoding/json"
	"fmt"
	"log"
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

func getProjectDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: Could not get working directory: %v", err)
		return "."
	}
	projectDir := filepath.Join(cwd, "..", "..")
	absPath, err := filepath.Abs(projectDir)
	if err != nil {
		log.Printf("Warning: Could not resolve absolute path: %v", err)
		return cwd
	}
	log.Printf("Project directory: %s", absPath)
	return absPath
}

func getModelsDir() string {
	projectDir := getProjectDir()
	modelsPath := filepath.Join(projectDir, "models")
	if _, err := os.Stat(modelsPath); err == nil {
		return modelsPath
	}
	possiblePaths := []string{
		"models",
		filepath.Join(".", "models"),
		filepath.Join("..", "models"),
	}
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			absPath, err := filepath.Abs(path)
			if err == nil {
				return absPath
			}
		}
	}
	return "models"
}

func LoadModelConfigs() (map[string]ModelConfig, error) {
	models := make(map[string]ModelConfig)
	modelsDir := getModelsDir()
	log.Printf("Loading models from directory: %s", modelsDir)

	files, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, fmt.Errorf("error reading models directory: %v (path: %s)", err, modelsDir)
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join(modelsDir, file.Name())
			log.Printf("Reading model file: %s", filePath)

			data, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading model file %s: %v", filePath, err)
				continue
			}
			var config ModelConfig
			if err := json.Unmarshal(data, &config); err != nil {
				log.Printf("Error parsing model file %s: %v", filePath, err)
				continue
			}
			modelName := strings.TrimSuffix(file.Name(), ".json")
			config.Name = modelName
			models[modelName] = config
			log.Printf("Loaded model: %s", modelName)
		}
	}
	if len(models) == 0 {
		return nil, fmt.Errorf("no valid model configurations found in %s", modelsDir)
	}
	return models, nil
}

package memory

import (
	"compute-gauge/pkg/calc"
	"compute-gauge/pkg/gpu"
	"fmt"
)

func CalculateMemoryRequirements(r *MemoryRequest) (*MemoryResponse, error) {
	if err := validateRequest(r); err != nil {
		return nil, err
	}

	var resp MemoryResponse
	inferenceResults := calc.CalculateInferenceMemory(
		r.ModelSize,
		r.TorchDtype,
		r.BatchSize,
		r.SequenceLength,
		r.HiddenSize,
		r.NumHiddenLayers,
		r.NumAttentionHeads,
	)
	resp.ModelWeights = inferenceResults["model_weights"]
	resp.KVCache = inferenceResults["kv_cache"]
	resp.ActivationMemory = inferenceResults["activation_memory"]
	resp.InferenceMemory = inferenceResults["inference_memory"]

	inferenceMemoryBytes, err := parseMemoryString(resp.InferenceMemory)
	if err != nil {
		return nil, fmt.Errorf("error parsing inference memory: %v", err)
	}
	inferenceMemoryGB := inferenceMemoryBytes / (1024 * 1024 * 1024)
	resp.InferenceGPUs = gpu.GetGPURecommendations(inferenceMemoryGB, false)
	if len(resp.InferenceGPUs) > 3 {
		resp.InferenceGPUs = resp.InferenceGPUs[:3]
	}

	if r.Optimizer != "" {
		trainingResults := calc.CalculateTrainingMemory(
			r.ModelSize,
			r.TorchDtype,
			r.BatchSize,
			r.SequenceLength,
			r.HiddenSize,
			r.NumHiddenLayers,
			r.NumAttentionHeads,
			r.Optimizer,
			r.ModelSize,
		)
		resp.OptimizerMemory = trainingResults["optimizer_memory"]
		resp.GradientsMemory = trainingResults["gradients_memory"]
		resp.TrainingMemory = trainingResults["training_memory"]
		trainingMemoryBytes, err := parseMemoryString(resp.TrainingMemory)
		if err != nil {
			return nil, fmt.Errorf("error parsing training memory: %v", err)
		}
		trainingMemoryGB := trainingMemoryBytes / (1024 * 1024 * 1024)
		resp.TrainingGPUs = gpu.GetGPURecommendations(trainingMemoryGB, true)
		if len(resp.TrainingGPUs) > 3 {
			resp.TrainingGPUs = resp.TrainingGPUs[:3]
		}
	}

	resp.TotalParams = r.ModelSize * 1e9
	resp.HiddenSize = r.HiddenSize
	resp.SequenceLength = r.SequenceLength
	return &resp, nil
}

func parseMemoryString(memStr string) (float64, error) {
	var value float64
	var unit string
	_, err := fmt.Sscanf(memStr, "%f %s", &value, &unit)
	if err != nil {
		return 0, err
	}

	multiplier := 1.0
	switch unit {
	case "TB":
		multiplier = 1024 * 1024 * 1024 * 1024
	case "GB":
		multiplier = 1024 * 1024 * 1024
	case "MB":
		multiplier = 1024 * 1024
	case "KB":
		multiplier = 1024
	}

	return value * multiplier, nil
}

func validateRequest(req *MemoryRequest) error {
	if req.ModelSize <= 0 {
		return fmt.Errorf("model size must be positive")
	}
	if req.HiddenSize <= 0 {
		return fmt.Errorf("hidden size must be positive")
	}
	if req.NumHiddenLayers <= 0 {
		return fmt.Errorf("number of layers must be positive")
	}
	if req.NumAttentionHeads <= 0 {
		return fmt.Errorf("number of attention heads must be positive")
	}
	if req.SequenceLength <= 0 {
		return fmt.Errorf("sequence length must be positive")
	}
	if req.BatchSize <= 0 {
		return fmt.Errorf("batch size must be positive")
	}

	validDtypes := map[string]bool{
		"float32":  true,
		"float16":  true,
		"bfloat16": true,
		"int8":     true,
		"int4":     true,
	}
	if !validDtypes[req.TorchDtype] {
		return fmt.Errorf("invalid precision type: %s", req.TorchDtype)
	}
	return nil
}

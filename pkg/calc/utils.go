package calc

import (
	"compute-gauge/pkg/config"
	"fmt"
	"math"
)

func FormatMemory(bytes float64) string {
	const (
		KB = 1024.0
		MB = KB * 1024.0
		GB = MB * 1024.0
		TB = GB * 1024.0
		PB = TB * 1024.0
	)
	switch {
	case bytes >= PB:
		return fmt.Sprintf("%.2f PB", bytes/PB)
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", bytes/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", bytes/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", bytes/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", bytes/KB)
	default:
		return fmt.Sprintf("%.2f B", bytes)
	}
}
func GetModelWeights(modelSize float64, precision string) float64 {
	if size, ok := config.DataTypeSizes[precision]; ok {
		params := modelSize * math.Pow(10, 9)
		return params * float64(size)
	}
	return 0
}
func GetKVCache(batchSize, seqLength, numLayers, hiddenSize int, precision string) float64 {
	if size, ok := config.DataTypeSizes[precision]; ok {
		batchF := float64(batchSize)
		seqF := float64(seqLength)
		layersF := float64(numLayers)
		hiddenF := float64(hiddenSize)
		sizeF := float64(size)
		return 2.0 * batchF * seqF * layersF * hiddenF * sizeF
	}
	return 0
}
func GetActivationMemory(batchSize, seqLength, numLayers, hiddenSize, numHeads int, precision string) float64 {
	const activationPrecision = "float32"
	if size, ok := config.DataTypeSizes[activationPrecision]; ok {
		batchF := float64(batchSize)
		seqF := float64(seqLength)
		headsF := float64(numHeads)
		hiddenF := float64(hiddenSize)
		sizeF := float64(size)
		activationFactor := 34.0 + ((5.0 * seqF * headsF) / hiddenF)
		return batchF * seqF * hiddenF * activationFactor * sizeF
	}
	return 0
}
func GetOptimizerMemory(trainableParams float64, optimizer string) float64 {
	actualParams := trainableParams * math.Pow(10, 9)
	switch optimizer {
	case "AdamW", "Adam":
		return actualParams * 8.0
	case "QAdamW":
		return actualParams * 2.0
	case "SGD":
		return actualParams * 4.0
	default:
		return 0
	}
}
func GetGradientMemory(trainableParams float64) float64 {
	actualParams := trainableParams * math.Pow(10, 9)
	return actualParams * 4.0
}
func CalculateInferenceMemory(modelSize float64, precision string, batchSize, seqLength, hiddenSize, numLayers, numHeads int) map[string]string {
	modelWeights := GetModelWeights(modelSize, precision)
	kvCache := GetKVCache(batchSize, seqLength, numLayers, hiddenSize, precision)
	activationMem := GetActivationMemory(batchSize, seqLength, numLayers, hiddenSize, numHeads, precision)
	totalMem := modelWeights + kvCache + activationMem
	return map[string]string{
		"model_weights":     FormatMemory(modelWeights),
		"kv_cache":          FormatMemory(kvCache),
		"activation_memory": FormatMemory(activationMem),
		"inference_memory":  FormatMemory(totalMem),
	}
}

func CalculateTrainingMemory(modelSize float64, precision string, batchSize, seqLength, hiddenSize, numLayers, numHeads int, optimizer string, trainableParams float64) map[string]string {
	modelWeights := GetModelWeights(modelSize, precision)
	kvCache := GetKVCache(batchSize, seqLength, numLayers, hiddenSize, precision)
	activationMem := GetActivationMemory(batchSize, seqLength, numLayers, hiddenSize, numHeads, precision)
	inferenceMem := modelWeights + kvCache + activationMem
	optimizerMem := GetOptimizerMemory(trainableParams, optimizer)
	gradientMem := GetGradientMemory(trainableParams)
	trainingSpecificMem := optimizerMem + gradientMem
	totalMem := inferenceMem + trainingSpecificMem
	return map[string]string{
		"model_weights":     FormatMemory(modelWeights),
		"kv_cache":          FormatMemory(kvCache),
		"activation_memory": FormatMemory(activationMem),
		"optimizer_memory":  FormatMemory(optimizerMem),
		"gradients_memory":  FormatMemory(gradientMem),
		"inference_memory":  FormatMemory(inferenceMem),
		"training_memory":   FormatMemory(totalMem),
	}
}

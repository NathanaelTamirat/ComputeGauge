package gpu

import (
	"math"
	"sort"
)

type GPUSpec struct {
	Name        string  `json:"name"`
	Memory      int     `json:"memory_gb"`
	Bandwidth   float64 `json:"bandwidth_tbs"`
	Price       float64 `json:"price_usd"`
	Performance float64 `json:"performance_tflops"`
}

var GPUDatabase = []GPUSpec{
	{
		Name:        "NVIDIA A100-80GB",
		Memory:      80,
		Bandwidth:   2.0,
		Price:       10000,
		Performance: 312,
	},
	{
		Name:        "NVIDIA A100-40GB",
		Memory:      40,
		Bandwidth:   1.6,
		Price:       6000,
		Performance: 312,
	},
	{
		Name:        "NVIDIA A6000",
		Memory:      48,
		Bandwidth:   0.768,
		Price:       4000,
		Performance: 309.7,
	},
	{
		Name:        "NVIDIA L40",
		Memory:      48,
		Bandwidth:   0.864,
		Price:       5000,
		Performance: 181.6,
	},
	{
		Name:        "NVIDIA A40",
		Memory:      48,
		Bandwidth:   0.696,
		Price:       3500,
		Performance: 149.8,
	},
	{
		Name:        "NVIDIA A30",
		Memory:      24,
		Bandwidth:   0.933,
		Price:       2000,
		Performance: 165,
	},
	{
		Name:        "NVIDIA A10",
		Memory:      24,
		Bandwidth:   0.600,
		Price:       1500,
		Performance: 125,
	},
	{
		Name:        "NVIDIA H100-80GB",
		Memory:      80,
		Bandwidth:   3.35,
		Price:       30000,
		Performance: 700,
	},
	{
		Name:        "NVIDIA H100-94GB",
		Memory:      94,
		Bandwidth:   3.9,
		Price:       35000,
		Performance: 830,
	},
	{
		Name:        "NVIDIA A100-80GB",
		Memory:      80,
		Bandwidth:   2.0,
		Price:       10000,
		Performance: 312,
	},
	{
		Name:        "NVIDIA A100-40GB",
		Memory:      40,
		Bandwidth:   1.6,
		Price:       6000,
		Performance: 312,
	},
	{
		Name:        "NVIDIA A6000",
		Memory:      48,
		Bandwidth:   0.768,
		Price:       4000,
		Performance: 309.7,
	},
	{
		Name:        "NVIDIA L40",
		Memory:      48,
		Bandwidth:   0.864,
		Price:       5000,
		Performance: 181.6,
	},
	{
		Name:        "NVIDIA RTX 6000 Ada Generation",
		Memory:      48,
		Bandwidth:   0.960,
		Price:       6800,
		Performance: 260,
	},
	{
		Name:        "NVIDIA A40",
		Memory:      48,
		Bandwidth:   0.696,
		Price:       3500,
		Performance: 149.8,
	},
	{
		Name:        "NVIDIA A30",
		Memory:      24,
		Bandwidth:   0.933,
		Price:       2000,
		Performance: 165,
	},
	{
		Name:        "NVIDIA A10",
		Memory:      24,
		Bandwidth:   0.600,
		Price:       1500,
		Performance: 125,
	},
}

type GPURecommendation struct {
	GPU              GPUSpec `json:"gpu"`
	NumGPUs          int     `json:"num_gpus"`
	UtilizationScore float64 `json:"utilization_score"`
	CostScore        float64 `json:"cost_score"`
	TotalCost        float64 `json:"total_cost"`
}

func GetGPURecommendations(totalMemoryGB float64, isTraining bool) []GPURecommendation {
	var recommendations []GPURecommendation

	for _, gpu := range GPUDatabase {
		numGPUs := int(math.Ceil(totalMemoryGB / float64(gpu.Memory)))
		if numGPUs < 1 {
			numGPUs = 1
		}
		memoryUtilization := totalMemoryGB / (float64(numGPUs) * float64(gpu.Memory))
		utilizationScore := memoryUtilization * 100
		totalCost := float64(numGPUs) * gpu.Price
		costScore := totalCost / gpu.Performance
		if isTraining {
			utilizationScore *= (gpu.Bandwidth / 2.0)
			costScore *= 0.8
		}

		rec := GPURecommendation{
			GPU:              gpu,
			NumGPUs:          numGPUs,
			UtilizationScore: utilizationScore,
			CostScore:        costScore,
			TotalCost:        totalCost,
		}
		recommendations = append(recommendations, rec)
	}
	sort.Slice(recommendations, func(i, j int) bool {
		scoreI := recommendations[i].UtilizationScore - recommendations[i].CostScore
		scoreJ := recommendations[j].UtilizationScore - recommendations[j].CostScore
		return scoreI > scoreJ
	})

	if len(recommendations) > 5 {
		recommendations = recommendations[:5]
	}

	return recommendations
}

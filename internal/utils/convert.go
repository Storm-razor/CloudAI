package utils

import "strconv"

//---------------------------
//@brief 多维64位向量向量转32位
//---------------------------
func ConvertFloat64ToFloat32Embeddings(embeddings [][]float64) [][]float32 {
	float32Embeddings := make([][]float32, len(embeddings))
	for i, vec64 := range embeddings {
		vec32 := make([]float32, len(vec64))
		for j, v := range vec64 {
			vec32[j] = float32(v)
		}
		float32Embeddings[i] = vec32
	}
	return float32Embeddings
}

//---------------------------
//@brief 单64位向量转32位
//---------------------------
func ConvertFloat64ToFloat32Embedding(embedding []float64) []float32 {
	float32Embedding := make([]float32, len(embedding))
	for i, v := range embedding {
		float32Embedding[i] = float32(v)
	}
	return float32Embedding
}

// StringToInt 将字符串转换为整数，出错时返回默认值0
func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// StringToInt64 将字符串转换为int64，出错时返回默认值0
func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// IntToString 将整数转换为字符串
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// Int64ToString 将int64转换为字符串
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

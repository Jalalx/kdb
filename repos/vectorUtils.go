package repos

import "fmt"

func stringify(vector []float64) string {
	result := "["
	for i, v := range vector {
		result += fmt.Sprintf("%f", float32(v))
		if i < len(vector)-1 {
			result += ", "
		}
	}
	result += "]"
	return result
}

func stringifyWithType(vector []float64) string {
	return fmt.Sprintf("%s::FLOAT[%d]", stringify(vector), len(vector))
}

package conv

import (
	"strconv"

	"github.com/fzf-labs/fpkg/binary"
)

// Float32 converts `any` to float32.
func Float32(any any) float32 {
	if any == nil {
		return 0
	}
	switch value := any.(type) {
	case float32:
		return value
	case float64:
		return float32(value)
	case []byte:
		return binary.DecodeToFloat32(value)
	default:
		v, _ := strconv.ParseFloat(String(any), 64)
		return float32(v)
	}
}

// Float64 converts `any` to float64.
func Float64(any any) float64 {
	if any == nil {
		return 0
	}
	switch value := any.(type) {
	case float32:
		return float64(value)
	case float64:
		return value
	case []byte:
		return binary.DecodeToFloat64(value)
	default:
		v, _ := strconv.ParseFloat(String(any), 64)
		return v
	}
}

package mathutil

import (
	"github.com/fzf-labs/fpkg/conv"
)

// And 如果a和b都为真则返回true。
func And[T, U any](a T, b U) bool {
	return conv.Bool(a) && conv.Bool(b)
}

// Or 如果a和b都不是真，则返回false。
func Or[T, U any](a T, b U) bool {
	return conv.Bool(a) || conv.Bool(b)
}

// Xor 如果a或b都为真，则返回真。
func Xor[T, U any](a T, b U) bool {
	valA := conv.Bool(a)
	valB := conv.Bool(b)
	return (valA || valB) && valA != valB
}

// Nor 如果a和b都不为真则返回真。
func Nor[T, U any](a T, b U) bool {
	return !(conv.Bool(a) || conv.Bool(b))
}

// XNor 如果a和b都为真，或a和b都不为真，则返回真。
func XNor[T, U any](a T, b U) bool {
	valA := conv.Bool(a)
	valB := conv.Bool(b)
	return (valA && valB) || (!valA && !valB)
}

// NAnd 如果a和b都为真则返回false。
func NAnd[T, U any](a T, b U) bool {
	return !conv.Bool(a) || !conv.Bool(b)
}

// TernaryOperator 检查参数“isTrue”的值，如果为真则返回ifValue否则返回elseValue。
func TernaryOperator[T, U any](isTrue T, ifValue, elseValue U) U {
	if conv.Bool(isTrue) {
		return ifValue
	} else {
		return elseValue
	}
}

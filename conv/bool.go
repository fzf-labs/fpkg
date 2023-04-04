package conv

import "reflect"

func Bool[T any](value T) bool {
	switch m := any(value).(type) {
	case interface{ Bool() bool }:
		return m.Bool()
	case interface{ IsZero() bool }:
		return !m.IsZero()
	}
	return reflectValue(&value)
}

func reflectValue(vp any) bool {
	switch rv := reflect.ValueOf(vp).Elem(); rv.Kind() {
	case reflect.Map, reflect.Slice:
		return rv.Len() != 0
	default:
		is := rv.IsZero()
		return !is
	}
}

//var emptyStringMap = map[string]struct{}{
//	"":      {},
//	"0":     {},
//	"no":    {},
//	"off":   {},
//	"false": {},
//}
//
//// Bool 转换为bool类型。
//// It returns false if `any` is: false, "", 0, "false", "off", "no", empty slice/map.
//func Bool(any any) bool {
//	if any == nil {
//		return false
//	}
//	switch value := any.(type) {
//	case bool:
//		return value
//	case []byte:
//		if _, ok := emptyStringMap[strings.ToLower(string(value))]; ok {
//			return false
//		}
//		return true
//	case string:
//		if _, ok := emptyStringMap[strings.ToLower(value)]; ok {
//			return false
//		}
//		return true
//	default:
//		rv := reflect.ValueOf(any)
//		switch rv.Kind() {
//		case reflect.Ptr:
//			return !rv.IsNil()
//		case reflect.Map:
//			fallthrough
//		case reflect.Array:
//			fallthrough
//		case reflect.Slice:
//			return rv.Len() != 0
//		case reflect.Struct:
//			return true
//		default:
//			s := strings.ToLower(String(any))
//			if _, ok := emptyStringMap[s]; ok {
//				return false
//			}
//			return true
//		}
//	}
//}

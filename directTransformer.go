package reflectconf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// TransformerFunc 转换函数格式定义
type TransformerFunc func(confStr string) (reflect.Value, error)

// getTransformerFunc 根据 transformerName 、 StructField.Kind 获取 TransformerFunc
// TransformerFunc 分为两种：
//
//	一是 通用直接转换器，根据 StructField.Kind 自动生成转换器，如 direct
//	二是 根据 transformerName 明确定义的转换器，如 classMapSign
func getTransformerFunc(transformerName string, field reflect.StructField) (TransformerFunc, error) {
	switch transformerName {
	default:
		return nil, fmt.Errorf("no such transformer: %s", transformerName)

	case DirectTransformerConf:
		return getDirectTransformerFunc(field.Type)

	case DirectSliceTransformerConf:
		// 此处 []int{} 中的 int 为一个例子，实际上 切片 的 类型 是不确定的，所以 checkTypeMatch 时，应该 justCompareKind
		if err := checkTypeMatch(field.Type, reflect.TypeOf([]int{}), true); err != nil {
			return nil, err
		}
		return getDirectSliceTransformerFunc(field.Type, "|")

	case DirectSliceTransformerConfSemicolon:
		// 此处 []int{} 中的 int 为一个例子，实际上 切片 的 类型 是不确定的，所以 checkTypeMatch 时，应该 justCompareKind
		if err := checkTypeMatch(field.Type, reflect.TypeOf([]int{}), true); err != nil {
			return nil, err
		}
		return getDirectSliceTransformerFunc(field.Type, ";")

	case DirectSliceTransformerConfColon:
		// 此处 []int{} 中的 int 为一个例子，实际上 切片 的 类型 是不确定的，所以 checkTypeMatch 时，应该 justCompareKind
		if err := checkTypeMatch(field.Type, reflect.TypeOf([]int{}), true); err != nil {
			return nil, err
		}
		return getDirectSliceTransformerFunc(field.Type, ":")

	case DirectSliceTransformerConfComma:
		// 此处 []int{} 中的 int 为一个例子，实际上 切片 的 类型 是不确定的，所以 checkTypeMatch 时，应该 justCompareKind
		if err := checkTypeMatch(field.Type, reflect.TypeOf([]int{}), true); err != nil {
			return nil, err
		}
		return getDirectSliceTransformerFunc(field.Type, ",")

	case DirectMapStructTransformerConf:
		// 此处 map[int]struct{} 中的 int 为一个例子，实际上 map 的 key 是不确定的，所以 checkTypeMatch 时，应该 justCompareKind
		if err := checkTypeMatch(field.Type, reflect.TypeOf(map[int]struct{}{}), true); err != nil {
			return nil, err
		}
		return getDirectMapStructTransformerFunc(field.Type, "|")

	case DirectMapStructTransformerConfSemicolon:
		// 此处 map[int]struct{} 中的 int 为一个例子，实际上 map 的 key 是不确定的，所以 checkTypeMatch 时，应该 justCompareKind
		if err := checkTypeMatch(field.Type, reflect.TypeOf(map[int]struct{}{}), true); err != nil {
			return nil, err
		}
		return getDirectMapStructTransformerFunc(field.Type, ";")

	case DirectMapStructTransformerConfComma:
		// 此处 map[int]struct{} 中的 int 为一个例子，实际上 map 的 key 是不确定的，所以 checkTypeMatch 时，应该 justCompareKind
		if err := checkTypeMatch(field.Type, reflect.TypeOf(map[int]struct{}{}), true); err != nil {
			return nil, err
		}
		return getDirectMapStructTransformerFunc(field.Type, ",")

	case "timeDuration":
		if err := checkTypeMatch(field.Type, reflect.TypeOf(int64(0)), false); err != nil {
			return nil, err
		}
		return timeDuration, nil

	case "confMapIntFloat":
		if err := checkTypeMatch(field.Type, reflect.TypeOf(map[int32]float32{}), false); err != nil {
			return nil, err
		}
		return confMapIntFloat, nil

	case "confMapIntInt":
		if err := checkTypeMatch(field.Type, reflect.TypeOf(map[int32]int32{}), false); err != nil {
			return nil, err
		}
		return confMapIntInt, nil

	case "confMapIntString":
		if err := checkTypeMatch(field.Type, reflect.TypeOf(map[int64]string{}), false); err != nil {
			return nil, err
		}
		return confMapIntString, nil

	case "confMapIntStrings":
		if err := checkTypeMatch(field.Type, reflect.TypeOf(map[int64][]string{}), false); err != nil {
			return nil, err
		}
		return confMapIntStrings, nil

	case "confMapIntUintMap":
		if err := checkTypeMatch(field.Type, reflect.TypeOf(map[int]map[uint64]struct{}{}), false); err != nil {
			return nil, err
		}
		return confMapIntUintMap, nil

	case "confMapStringFloat":
		if err := checkTypeMatch(field.Type, reflect.TypeOf(map[string]float32{}), false); err != nil {
			return nil, err
		}
		return confMapStringFloat, nil
	}
}

// getDirectTransformerFunc 根据 reflect.Type 自动生成直接转换器 , tag: DirectTransformerConf
func getDirectTransformerFunc(t reflect.Type) (TransformerFunc, error) {
	fieldKind := t.Kind()
	switch fieldKind {
	default:
		return nil, fmt.Errorf("fieldKind:%s is not supported by DirectTransformerFunc", fieldKind)

	case reflect.String:
		return func(confStr string) (reflect.Value, error) {
			return reflect.ValueOf(confStr), nil
		}, nil

	case reflect.Bool:
		return func(confStr string) (reflect.Value, error) {
			if ret, err := strconv.ParseBool(confStr); err == nil {
				return reflect.ValueOf(ret), nil
			} else {
				return reflect.Value{}, err
			}
		}, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(confStr string) (reflect.Value, error) {
			if retP, err := strconv.ParseInt(confStr, 10, 64); err == nil {
				return reflect.ValueOf(retP).Convert(t), nil
			} else {
				return reflect.Value{}, err
			}
		}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(confStr string) (reflect.Value, error) {
			if retP, err := strconv.ParseUint(confStr, 10, 64); err == nil {
				return reflect.ValueOf(retP).Convert(t), nil
			} else {
				return reflect.Value{}, err
			}
		}, nil

	case reflect.Float32, reflect.Float64:
		return func(confStr string) (reflect.Value, error) {
			if retP, err := strconv.ParseFloat(confStr, 64); err == nil {
				return reflect.ValueOf(retP).Convert(t), nil
			} else {
				return reflect.Value{}, err
			}
		}, nil

	case reflect.Complex64, reflect.Complex128:
		return func(confStr string) (reflect.Value, error) {
			if retP, err := strconv.ParseComplex(confStr, 128); err == nil {
				return reflect.ValueOf(retP).Convert(t), nil
			} else {
				return reflect.Value{}, err
			}
		}, nil
	}
}

// getDirectTransformerFunc 根据 reflect.Type 自动生成直接转换器 , tag: DirectSliceTransformerConf
func getDirectSliceTransformerFunc(t reflect.Type, sep string) (TransformerFunc, error) {
	fieldElemKind := t.Elem().Kind()
	switch fieldElemKind {
	default:
		return nil, fmt.Errorf("fieldElemKind:%s is not supported by DirectSliceTransformerFunc", fieldElemKind)

	case reflect.String:
		return func(confStr string) (reflect.Value, error) {
			return reflect.ValueOf(strings.Split(confStr, sep)), nil
		}, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(confStr string) (reflect.Value, error) {
			// TODO 这里性能有点差，待优化
			rValue := reflect.New(t).Elem()
			tSlice := make([]reflect.Value, 0, 10)
			for _, s := range strings.Split(confStr, sep) {
				if retP, err := strconv.ParseInt(s, 10, 64); err == nil {
					tSlice = append(tSlice, reflect.ValueOf(retP).Convert(t.Elem()))
				}
			}
			if len(tSlice) > 0 {
				rValue.Set(reflect.Append(rValue, tSlice...))
			}
			return rValue, nil
		}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(confStr string) (reflect.Value, error) {
			rValue := reflect.New(t).Elem()
			tSlice := make([]reflect.Value, 0, 10)
			for _, s := range strings.Split(confStr, sep) {
				if retP, err := strconv.ParseUint(s, 10, 64); err == nil {
					tSlice = append(tSlice, reflect.ValueOf(retP).Convert(t.Elem()))
				}
			}
			if len(tSlice) > 0 {
				rValue.Set(reflect.Append(rValue, tSlice...))
			}
			return rValue, nil
		}, nil

	case reflect.Float32, reflect.Float64:
		return func(confStr string) (reflect.Value, error) {
			rValue := reflect.New(t).Elem()
			tSlice := make([]reflect.Value, 0, 10)
			for _, s := range strings.Split(confStr, sep) {
				if retP, err := strconv.ParseFloat(s, 64); err == nil {
					tSlice = append(tSlice, reflect.ValueOf(retP).Convert(t.Elem()))
				}
			}
			if len(tSlice) > 0 {
				rValue.Set(reflect.Append(rValue, tSlice...))
			}
			return rValue, nil
		}, nil

	case reflect.Complex64, reflect.Complex128:
		return func(confStr string) (reflect.Value, error) {
			rValue := reflect.New(t).Elem()
			tSlice := make([]reflect.Value, 0, 10)
			for _, s := range strings.Split(confStr, sep) {
				if retP, err := strconv.ParseComplex(s, 128); err == nil {
					tSlice = append(tSlice, reflect.ValueOf(retP).Convert(t.Elem()))
				}
			}
			if len(tSlice) > 0 {
				rValue.Set(reflect.Append(rValue, tSlice...))
			}
			return rValue, nil
		}, nil
	}
}

// getDirectMapStructTransformerFunc 根据 reflect.Type 自动生成直接转换器 , tag: DirectMapStructTransformerConf
func getDirectMapStructTransformerFunc(t reflect.Type, sep string) (TransformerFunc, error) {
	keyKind := t.Key().Kind()
	eStruct := reflect.ValueOf(struct{}{})
	switch keyKind {
	default:
		return nil, fmt.Errorf("fieldElemKind:%s is not supported by DirectSliceTransformerFunc", keyKind)

	case reflect.String:
		return func(confStr string) (reflect.Value, error) {
			rValue := reflect.MakeMap(t)
			for _, s := range strings.Split(confStr, sep) {
				rValue.SetMapIndex(reflect.ValueOf(s), eStruct)
			}
			return rValue, nil
		}, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(confStr string) (reflect.Value, error) {
			rValue := reflect.MakeMap(t)
			for _, s := range strings.Split(confStr, sep) {
				if retP, err := strconv.ParseInt(s, 10, 64); err == nil {
					rValue.SetMapIndex(reflect.ValueOf(retP).Convert(t.Key()), eStruct)
				}
			}
			return rValue, nil
		}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(confStr string) (reflect.Value, error) {
			rValue := reflect.MakeMap(t)
			for _, s := range strings.Split(confStr, sep) {
				if retP, err := strconv.ParseUint(s, 10, 64); err == nil {
					rValue.SetMapIndex(reflect.ValueOf(retP).Convert(t.Key()), eStruct)
				}
			}
			return rValue, nil
		}, nil

	case reflect.Float32, reflect.Float64:
		return func(confStr string) (reflect.Value, error) {
			rValue := reflect.MakeMap(t)
			for _, s := range strings.Split(confStr, sep) {
				if retP, err := strconv.ParseFloat(s, 64); err == nil {
					rValue.SetMapIndex(reflect.ValueOf(retP).Convert(t.Key()), eStruct)
				}
			}
			return rValue, nil
		}, nil

	case reflect.Complex64, reflect.Complex128:
		return func(confStr string) (reflect.Value, error) {
			rValue := reflect.MakeMap(t)
			for _, s := range strings.Split(confStr, sep) {
				if retP, err := strconv.ParseComplex(s, 128); err == nil {
					rValue.SetMapIndex(reflect.ValueOf(retP).Convert(t.Key()), eStruct)
				}
			}
			return rValue, nil
		}, nil
	}
}

func checkTypeMatch(fieldType, transFunType reflect.Type, justCompareKind bool) error {
	if justCompareKind {
		if fieldType.Kind() == transFunType.Kind() {
			return nil
		}
	} else {
		if transFunType.AssignableTo(fieldType) {
			return nil
		}
	}
	return fmt.Errorf("transformer is not match with fieldType: %s", fieldType)
}

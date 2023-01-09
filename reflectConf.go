// Package reflectconf
// reflectConf.go: 定义了主体逻辑和通用转换函数
// reflectTransformer.go 定义了特殊转换函数

package reflectconf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	ConfTag        = "conf"
	TransformerTag = "transformer"
	DefaultTag     = "default"

	DirectTransformerConf                   = "direct"
	DirectSliceTransformerConf              = "diSlice"
	DirectSliceTransformerConfSemicolon     = "diSliceSemicolon"
	DirectSliceTransformerConfColon         = "diSliceColon"
	DirectSliceTransformerConfComma         = "diSliceComma"
	DirectMapStructTransformerConf          = "diMapStruct"
	DirectMapStructTransformerConfSemicolon = "diMapStructSemicolon"
	DirectMapStructTransformerConfComma     = "diMapStructComma"
)

// FillConfByParamsMap 根据 reflect.StructTag、confMaps 为 obj 填充值
// 支持传入多组 map[string]string，从前到后优先级依次降低，至少传入一个
func FillConfByParamsMap(obj interface{}, useDefault bool, confMaps ...map[string]string) error {
	// rTypeKind 输入类型检查，必须是 结构体对象的指针
	rTypeKind := reflect.TypeOf(obj).Kind()
	if rTypeKind != reflect.Ptr {
		return fmt.Errorf("obj kind err, expected struct ptr, found %s", rTypeKind)
	}
	rValueElem := reflect.ValueOf(obj).Elem() // rValueElem 用于往 obj 内 Set value
	rValueElemType := rValueElem.Type()       // rValueElemType 用于遍历 struct field ，获取 tags
	if rValueElemType.Kind() != reflect.Struct {
		return fmt.Errorf("obj elem kind err, expected struct, found %s", rValueElemType.Kind())
	}

	// 配置列表合并
	confMap := mergeConfMaps(confMaps)
	// rValueElemTypeString := rValueElemType.String()

	// 遍历 struct field 执行 transformer 并 赋值
	for i := 0; i < rValueElemType.NumField(); i++ {
		fieldType := rValueElemType.Field(i)
		rValueElemField := rValueElem.Field(i)
		tagsMap := getStructTagsMap(fieldType.Tag)
		transformerName, confStr, pass := parseStructTagConf(tagsMap, confMap, useDefault)
		if !pass {
			continue
		}

		transformerFunc, err := getTransformerFunc(transformerName, fieldType)
		if err != nil {
			return fmt.Errorf("getTransformerFunc error:%s, struct:%s, field:%s, tags:%s, transformer:%s, confB4Trans:%s, kind:%s",
				err, rValueElemType.String(), fieldType.Name, fieldType.Tag, transformerName, confStr, rValueElemField.Kind())
		}

		if !rValueElemField.CanSet() {
			return fmt.Errorf("rValueElemField error:%s, struct:%s, field:%s, tags:%s, transformer:%s, confB4Trans:%s, kind:%s",
				err, rValueElemType.String(), fieldType.Name, fieldType.Tag, transformerName, confStr, rValueElemField.Kind())
		}

		value, err := transformerFunc(confStr)
		if err != nil {
			return fmt.Errorf("do transformerFunc error:%s, struct:%s, field:%s, tags:%s, transformer:%s, confB4Trans:%s, kind:%s",
				err, rValueElemType.String(), fieldType.Name, fieldType.Tag, transformerName, confStr, rValueElemField.Kind())
		}

		if rValueElemField.Kind() != value.Kind() {
			return fmt.Errorf("transformer error:%s, struct:%s, field:%s, tags:%s, transformer:%s, confB4Trans:%s, kind:%s",
				err, rValueElemType.String(), fieldType.Name, fieldType.Tag, transformerName, confStr, rValueElemField.Kind())
		}

		rValueElemField.Set(value)
	}
	return nil
}

// mergeConfMaps 将配置列表按优先级合并，slice 中 index 0 优先级最高
func mergeConfMaps(confMaps []map[string]string) map[string]string {
	if len(confMaps) == 0 {
		return nil
	}

	if len(confMaps[0]) == 0 {
		confMaps[0] = make(map[string]string)
	}

	for _, cMap := range confMaps[1:] {
		for k, v := range cMap {
			if _, ok := confMaps[0][k]; !ok {
				confMaps[0][k] = v
			}
		}
	}

	if len(confMaps[0]) == 0 {
		return nil
	}

	return confMaps[0]
}

// getStructTagsMap 获取 reflect.StructTag 中所有的 tag
func getStructTagsMap(tag reflect.StructTag) map[string]string {
	m := make(map[string]string)
	for tag != "" {
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qValue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qValue)
		if err == nil {
			m[name] = value
		}
	}
	return m
}

// parseStructTagConf 检查配置值是否合法，并获取合法配置值
func parseStructTagConf(tagsMap, confMap map[string]string, useDefault bool) (string, string, bool) {
	transformerName, transformerTagOK := tagsMap[TransformerTag]
	confName, confTagOK := tagsMap[ConfTag]
	defaultConf, defaultTagOK := tagsMap[DefaultTag]

	// TransformerTag、ConfTag 两者必须配，不配则表示不需要解析
	if !transformerTagOK || !confTagOK {
		return "", "", false
	}

	confStr, confValueOK := "", false
	for _, s := range strings.Split(confName, "|") {
		if _confStr, _confValueOK := confMap[s]; _confValueOK {
			if _confStr == "" {
				continue
			}
			confValueOK = true
			confStr = _confStr
			break
		}
	}

	if !confValueOK && !useDefault {
		return "", "", false
	}

	if !confValueOK {
		if !defaultTagOK {
			return "", "", false
		} else {
			confStr = defaultConf
		}
	}

	confStr = strings.TrimSpace(confStr)

	if confStr == "" {
		return "", "", false
	}

	return transformerName, confStr, true
}

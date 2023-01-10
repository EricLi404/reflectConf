package reflectconf

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

// timeDuration   int64
func timeDuration(confStr string) (reflect.Value, error) {
	hh, _ := time.ParseDuration(confStr)
	return reflect.ValueOf(time.Now().Add(hh).Unix()), nil
}

// confMapIntFloat  map[int32]float32
func confMapIntFloat(confStr string) (reflect.Value, error) {
	fieldValueTmp := make(map[int32]float32)
	for _, v := range strings.Split(confStr, "|") {
		tArr := strings.Split(v, ":")
		if len(tArr) == 2 {
			if ret, err := strconv.Atoi(tArr[0]); err == nil {
				if ret1, err1 := strconv.ParseFloat(tArr[1], 32); err1 == nil {
					fieldValueTmp[int32(ret)] = float32(ret1)
				}
			}
		}
	}
	return reflect.ValueOf(fieldValueTmp), nil
}

// confMapIntInt  map[int32]int32
func confMapIntInt(confStr string) (reflect.Value, error) {
	fieldValueTmp := make(map[int32]int32)
	for _, v := range strings.Split(confStr, "|") {
		tArr := strings.Split(v, ":")
		if len(tArr) == 2 {
			if ret, err := strconv.Atoi(tArr[0]); err == nil {
				if ret1, err1 := strconv.ParseInt(tArr[1], 10, 64); err1 == nil {
					fieldValueTmp[int32(ret)] = int32(ret1)
				}
			}
		}
	}
	return reflect.ValueOf(fieldValueTmp), nil
}

// confMapIntString   map[int64]string
func confMapIntString(confStr string) (reflect.Value, error) {
	fieldValueTmp := make(map[int64]string)
	for _, v := range strings.Split(confStr, "|") {
		tArr := strings.Split(v, ":")
		if len(tArr) == 2 {
			if ret, err := strconv.ParseInt(tArr[0], 10, 64); err == nil {
				fieldValueTmp[ret] = tArr[1]
			}
		}
	}
	return reflect.ValueOf(fieldValueTmp), nil
}

// confMapIntStrings   map[int64][]string
func confMapIntStrings(confStr string) (reflect.Value, error) {
	fieldValueTmp := make(map[int64][]string)
	for _, v := range strings.Split(confStr, "|") {
		tArr := strings.Split(v, ":")
		if len(tArr) == 2 {
			vArrs := strings.Split(tArr[1], ";")
			if ret, err := strconv.ParseInt(tArr[0], 10, 64); err == nil {
				fieldValueTmp[ret] = vArrs
			}
		}
	}
	return reflect.ValueOf(fieldValueTmp), nil
}

// confMapIntUintMap   map[int]map[uint64]struct{}
func confMapIntUintMap(confStr string) (reflect.Value, error) {
	fieldValueTmp := make(map[int]map[uint64]struct{})
	actStrSlice := strings.Split(strings.TrimSpace(confStr), ";")

	for _, key := range actStrSlice {
		actCntSlice := strings.Split(strings.TrimSpace(key), ":")
		if len(actCntSlice) != 2 {
			continue
		}

		if act, err := strconv.ParseInt(actCntSlice[0], 10, 32); err == nil {
			cntSlice := strings.Split(actCntSlice[1], "|")
			for _, cntv := range cntSlice {
				if cnt, err := strconv.ParseUint(cntv, 10, 64); err == nil {
					if _, ok := fieldValueTmp[int(act)]; !ok {
						fieldValueTmp[int(act)] = make(map[uint64]struct{})
					}
					fieldValueTmp[int(act)][cnt] = struct{}{}
				}
			}
		}
	}
	return reflect.ValueOf(fieldValueTmp), nil
}

// confMapStringFloat   map[string]float32
func confMapStringFloat(confStr string) (reflect.Value, error) {
	fieldValueTmp := make(map[string]float32)
	for _, v := range strings.Split(confStr, "|") {
		tArr := strings.Split(v, ":")
		if len(tArr) == 2 {
			if ret, err := strconv.ParseFloat(tArr[1], 32); err == nil {
				fieldValueTmp[tArr[0]] = float32(ret)
			}
		}
	}
	return reflect.ValueOf(fieldValueTmp), nil
}

package reflectconf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

// go test -bench=. -cpu=1 -benchmem -benchtime=5s
//
// go1.19.4
// goos: darwin
// goarch: arm64
// pkg: github.com/EricLi404/reflectConf
// BenchmarkFillConfByTransformer_Direct            1086484              5549 ns/op            6280 B/op         97 allocs/op
// BenchmarkNoReflect_Direct                       26880355               224.6 ns/op             0 B/op          0 allocs/op
// BenchmarkFillConfByTransformer_DirectSlice        600121             10021 ns/op            8432 B/op        198 allocs/op
// BenchmarkNoReflect_DirectSlice                   3416010              1727 ns/op            1144 B/op         44 allocs/op
// BenchmarkFillConfByTransformer_DirectMapStruct    632258              9476 ns/op            9928 B/op        188 allocs/op
// BenchmarkNoReflect_DirectMapStruct               2737756              2194 ns/op            2432 B/op         45 allocs/op

func TestFillConfByTransformer_Direct(t *testing.T) {
	type testCase struct {
		F1  string     `transformer:"direct" conf:"f1" default:"" `
		F2  bool       `transformer:"direct" conf:"f2" default:"" `
		F3  int        `transformer:"direct" conf:"f3" default:"" `
		F4  int8       `transformer:"direct" conf:"f4" default:"" `
		F5  int16      `transformer:"direct" conf:"f5" default:"" `
		F6  int32      `transformer:"direct" conf:"f6" default:"" `
		F7  int64      `transformer:"direct" conf:"f7" default:"" `
		F8  uint       `transformer:"direct" conf:"f8" default:"" `
		F9  uint8      `transformer:"direct" conf:"f9" default:"" `
		F10 uint16     `transformer:"direct" conf:"f10" default:"" `
		F11 uint32     `transformer:"direct" conf:"f11" default:"" `
		F12 uint64     `transformer:"direct" conf:"f12" default:"" `
		F13 float32    `transformer:"direct" conf:"f13" default:"" `
		F14 float64    `transformer:"direct" conf:"f14" default:"" `
		F15 complex64  `transformer:"direct" conf:"f15" default:"" `
		F16 complex128 `transformer:"direct" conf:"f16" default:"" `
	}

	x := &testCase{}

	conf := map[string]string{
		"f1":  "1",
		"f2":  "1",
		"f3":  "3",
		"f4":  "4",
		"f5":  "5",
		"f6":  "6",
		"f7":  "7",
		"f8":  "8",
		"f9":  "9",
		"f10": "10",
		"f11": "11",
		"f12": "12",
		"f13": "13",
		"f14": "14",
		"f15": "1+5i",
		"f16": "1+6i",
	}

	y := &testCase{
		F1:  "1",
		F2:  true,
		F3:  3,
		F4:  4,
		F5:  5,
		F6:  6,
		F7:  7,
		F8:  8,
		F9:  9,
		F10: 10,
		F11: 11,
		F12: 12,
		F13: 13,
		F14: 14,
		F15: complex(1, 5),
		F16: complex(1, 6),
	}

	err := FillConfByParamsMap(x, true, conf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", x)

	t.Run("direct transformer", func(t *testing.T) {
		err := FillConfByParamsMap(x, true, conf)
		if err != nil {
			t.Errorf("FillConfByParamsMap() error = %v", err)
			return
		}
		if !reflect.DeepEqual(x, y) {
			t.Errorf("FillConfByParamsMap() got = %+v, want %+v", y, x)
		}
	})
}

func BenchmarkFillConfByTransformer_Direct(b *testing.B) {
	type testCase struct {
		F1  string     `transformer:"direct" conf:"f1" default:"" `
		F2  bool       `transformer:"direct" conf:"f2" default:"" `
		F3  int        `transformer:"direct" conf:"f3" default:"" `
		F4  int8       `transformer:"direct" conf:"f4" default:"" `
		F5  int16      `transformer:"direct" conf:"f5" default:"" `
		F6  int32      `transformer:"direct" conf:"f6" default:"" `
		F7  int64      `transformer:"direct" conf:"f7" default:"" `
		F8  uint       `transformer:"direct" conf:"f8" default:"" `
		F9  uint8      `transformer:"direct" conf:"f9" default:"" `
		F10 uint16     `transformer:"direct" conf:"f10" default:"" `
		F11 uint32     `transformer:"direct" conf:"f11" default:"" `
		F12 uint64     `transformer:"direct" conf:"f12" default:"" `
		F13 float32    `transformer:"direct" conf:"f13" default:"" `
		F14 float64    `transformer:"direct" conf:"f14" default:"" `
		F15 complex64  `transformer:"direct" conf:"f15" default:"" `
		F16 complex128 `transformer:"direct" conf:"f16" default:"" `
	}

	x := &testCase{}

	conf := map[string]string{
		"f1":  "1",
		"f2":  "1",
		"f3":  "3",
		"f4":  "4",
		"f5":  "5",
		"f6":  "6",
		"f7":  "7",
		"f8":  "8",
		"f9":  "9",
		"f10": "10",
		"f11": "11",
		"f12": "12",
		"f13": "13",
		"f14": "14",
		"f15": "1+5i",
		"f16": "1+6i",
	}

	for i := 0; i < b.N; i++ {
		err := FillConfByParamsMap(x, true, conf)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkNoReflect_Direct(b *testing.B) {
	type testCase struct {
		F1  string
		F2  bool
		F3  int
		F4  int8
		F5  int16
		F6  int32
		F7  int64
		F8  uint
		F9  uint8
		F10 uint16
		F11 uint32
		F12 uint64
		F13 float32
		F14 float64
		F15 complex64
		F16 complex128
	}

	F1 := ""
	F2 := false
	F3 := 0
	F4 := int8(0)
	F5 := int16(0)
	F6 := int32(0)
	F7 := int64(0)
	F8 := uint(0)
	F9 := uint8(0)
	F10 := uint16(0)
	F11 := uint32(0)
	F12 := uint64(0)
	F13 := float32(0)
	F14 := float64(0)
	F15 := complex64(complex(0, 0))
	F16 := complex(0, 0)

	conf := map[string]string{
		"f1":  "1",
		"f2":  "1",
		"f3":  "3",
		"f4":  "4",
		"f5":  "5",
		"f6":  "6",
		"f7":  "7",
		"f8":  "8",
		"f9":  "9",
		"f10": "10",
		"f11": "11",
		"f12": "12",
		"f13": "13",
		"f14": "14",
		"f15": "1+5i",
		"f16": "1+6i",
	}

	for i := 0; i < b.N; i++ {
		if v, ok := conf["f1"]; ok {
			F1 = v
		}
		if v, ok := conf["f2"]; ok {
			if ret, err := strconv.ParseBool(v); err == nil {
				F2 = ret
			}
		}
		if v, ok := conf["f3"]; ok {
			if ret, err := strconv.Atoi(v); err == nil {
				F3 = ret
			}
		}
		if v, ok := conf["f4"]; ok {
			if retP, err := strconv.ParseInt(v, 10, 8); err == nil {
				F4 = int8(retP)
			}
		}
		if v, ok := conf["f5"]; ok {
			if retP, err := strconv.ParseInt(v, 10, 16); err == nil {
				F5 = int16(retP)
			}
		}
		if v, ok := conf["f6"]; ok {
			if retP, err := strconv.ParseInt(v, 10, 32); err == nil {
				F6 = int32(retP)
			}
		}
		if v, ok := conf["f7"]; ok {
			if retP, err := strconv.ParseInt(v, 10, 64); err == nil {
				F7 = retP
			}
		}
		if v, ok := conf["f8"]; ok {
			if retP, err := strconv.ParseUint(v, 10, 64); err == nil {
				F8 = uint(retP)
			}
		}
		if v, ok := conf["f9"]; ok {
			if retP, err := strconv.ParseUint(v, 10, 8); err == nil {
				F9 = uint8(retP)
			}
		}
		if v, ok := conf["f10"]; ok {
			if retP, err := strconv.ParseUint(v, 10, 16); err == nil {
				F10 = uint16(retP)
			}
		}
		if v, ok := conf["f11"]; ok {
			if retP, err := strconv.ParseUint(v, 10, 32); err == nil {
				F11 = uint32(retP)
			}
		}
		if v, ok := conf["f12"]; ok {
			if retP, err := strconv.ParseUint(v, 10, 64); err == nil {
				F12 = retP
			}
		}
		if v, ok := conf["f13"]; ok {
			if retP, err := strconv.ParseFloat(v, 32); err == nil {
				F13 = float32(retP)
			}
		}
		if v, ok := conf["f14"]; ok {
			if retP, err := strconv.ParseFloat(v, 64); err == nil {
				F14 = retP
			}
		}
		if v, ok := conf["f15"]; ok {
			if retP, err := strconv.ParseComplex(v, 64); err == nil {
				F15 = complex64(retP)
			}
		}
		if v, ok := conf["f16"]; ok {
			if retP, err := strconv.ParseComplex(v, 128); err == nil {
				F16 = retP
			}
		}
		_ = &testCase{
			F1:  F1,
			F2:  F2,
			F3:  F3,
			F4:  F4,
			F5:  F5,
			F6:  F6,
			F7:  F7,
			F8:  F8,
			F9:  F9,
			F10: F10,
			F11: F11,
			F12: F12,
			F13: F13,
			F14: F14,
			F15: F15,
			F16: F16,
		}
	}
}

func TestFillConfByTransformer_DirectSlice(t *testing.T) {
	type testCase struct {
		F1  []string     `transformer:"diSlice" conf:"f1" default:"" `
		F3  []int        `transformer:"diSlice" conf:"f3" default:"" `
		F4  []int8       `transformer:"diSlice" conf:"f4" default:"" `
		F5  []int16      `transformer:"diSlice" conf:"f5" default:"" `
		F6  []int32      `transformer:"diSlice" conf:"f6" default:"" `
		F7  []int64      `transformer:"diSlice" conf:"f7" default:"" `
		F8  []uint       `transformer:"diSlice" conf:"f8" default:"" `
		F9  []uint8      `transformer:"diSlice" conf:"f9" default:"" `
		F10 []uint16     `transformer:"diSlice" conf:"f10" default:"" `
		F11 []uint32     `transformer:"diSlice" conf:"f11" default:"" `
		F12 []uint64     `transformer:"diSlice" conf:"f12" default:"" `
		F13 []float32    `transformer:"diSlice" conf:"f13" default:"" `
		F14 []float64    `transformer:"diSlice" conf:"f14" default:"" `
		F15 []complex64  `transformer:"diSlice" conf:"f15" default:"" `
		F16 []complex128 `transformer:"diSlice" conf:"f16" default:"" `
	}

	x := &testCase{}

	conf := map[string]string{
		"f1":  "1|2|3",
		"f3":  "3|4|5",
		"f4":  "4|5|6",
		"f5":  "5|6|7",
		"f6":  "6|7|8",
		"f7":  "7|8|9",
		"f8":  "8|9|10",
		"f9":  "9|10|11",
		"f10": "10|11|12",
		"f11": "11|12|13",
		"f12": "12|13|14",
		"f13": "13|14|15",
		"f14": "14|15|16",
		"f15": "1+5i|1+6i",
		"f16": "1+6i|1+7i",
	}

	y := &testCase{
		F1:  []string{"1", "2", "3"},
		F3:  []int{3, 4, 5},
		F4:  []int8{4, 5, 6},
		F5:  []int16{5, 6, 7},
		F6:  []int32{6, 7, 8},
		F7:  []int64{7, 8, 9},
		F8:  []uint{8, 9, 10},
		F9:  []uint8{9, 10, 11},
		F10: []uint16{10, 11, 12},
		F11: []uint32{11, 12, 13},
		F12: []uint64{12, 13, 14},
		F13: []float32{13, 14, 15},
		F14: []float64{14, 15, 16},
		F15: []complex64{complex(1, 5), complex(1, 6)},
		F16: []complex128{complex(1, 6), complex(1, 7)},
	}

	err := FillConfByParamsMap(x, true, conf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", x)

	t.Run("diSlice transformer", func(t *testing.T) {
		err := FillConfByParamsMap(x, true, conf)
		if err != nil {
			t.Errorf("FillConfByParamsMap() error = %v", err)
			return
		}
		if !reflect.DeepEqual(x, y) {
			t.Errorf("FillConfByParamsMap() got = %+v, want %+v", y, x)
		}
	})
}

func BenchmarkFillConfByTransformer_DirectSlice(b *testing.B) {
	type testCase struct {
		F1  []string     `transformer:"diSlice" conf:"f1" default:"" `
		F3  []int        `transformer:"diSlice" conf:"f3" default:"" `
		F4  []int8       `transformer:"diSlice" conf:"f4" default:"" `
		F5  []int16      `transformer:"diSlice" conf:"f5" default:"" `
		F6  []int32      `transformer:"diSlice" conf:"f6" default:"" `
		F7  []int64      `transformer:"diSlice" conf:"f7" default:"" `
		F8  []uint       `transformer:"diSlice" conf:"f8" default:"" `
		F9  []uint8      `transformer:"diSlice" conf:"f9" default:"" `
		F10 []uint16     `transformer:"diSlice" conf:"f10" default:"" `
		F11 []uint32     `transformer:"diSlice" conf:"f11" default:"" `
		F12 []uint64     `transformer:"diSlice" conf:"f12" default:"" `
		F13 []float32    `transformer:"diSlice" conf:"f13" default:"" `
		F14 []float64    `transformer:"diSlice" conf:"f14" default:"" `
		F15 []complex64  `transformer:"diSlice" conf:"f15" default:"" `
		F16 []complex128 `transformer:"diSlice" conf:"f16" default:"" `
	}

	x := &testCase{}

	conf := map[string]string{
		"f1":  "1|2|3",
		"f3":  "3|4|5",
		"f4":  "4|5|6",
		"f5":  "5|6|7",
		"f6":  "6|7|8",
		"f7":  "7|8|9",
		"f8":  "8|9|10",
		"f9":  "9|10|11",
		"f10": "10|11|12",
		"f11": "11|12|13",
		"f12": "12|13|14",
		"f13": "13|14|15",
		"f14": "14|15|16",
		"f15": "1+5i|1+6i",
		"f16": "1+6i|1+7i",
	}

	for i := 0; i < b.N; i++ {
		err := FillConfByParamsMap(x, true, conf)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkNoReflect_DirectSlice(b *testing.B) {
	type testCase struct {
		F1  []string
		F3  []int
		F4  []int8
		F5  []int16
		F6  []int32
		F7  []int64
		F8  []uint
		F9  []uint8
		F10 []uint16
		F11 []uint32
		F12 []uint64
		F13 []float32
		F14 []float64
		F15 []complex64
		F16 []complex128
	}

	conf := map[string]string{
		"f1":  "1|2|3",
		"f3":  "3|4|5",
		"f4":  "4|5|6",
		"f5":  "5|6|7",
		"f6":  "6|7|8",
		"f7":  "7|8|9",
		"f8":  "8|9|10",
		"f9":  "9|10|11",
		"f10": "10|11|12",
		"f11": "11|12|13",
		"f12": "12|13|14",
		"f13": "13|14|15",
		"f14": "14|15|16",
		"f15": "1+5i|1+6i",
		"f16": "1+6i|1+7i",
	}

	F1 := []string{"0", "2", "3"}
	F3 := []int{0, 4, 5}
	F4 := []int8{0, 5, 6}
	F5 := []int16{0, 6, 7}
	F6 := []int32{0, 7, 8}
	F7 := []int64{0, 8, 9}
	F8 := []uint{0, 9, 10}
	F9 := []uint8{0, 10, 11}
	F10 := []uint16{0, 11, 12}
	F11 := []uint32{0, 12, 13}
	F12 := []uint64{0, 13, 14}
	F13 := []float32{0, 14, 15}
	F14 := []float64{0, 15, 16}
	F15 := []complex64{complex(0, 5), complex(1, 6)}
	F16 := []complex128{complex(0, 6), complex(1, 7)}
	for i := 0; i < b.N; i++ {
		if v, ok := conf["f1"]; ok {
			F1 = strings.Split(v, "|")
		}
		if v, ok := conf["f3"]; ok {
			t := make([]int, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.Atoi(s); err == nil {
					t = append(t, ret)
				}
			}
			F3 = t
		}
		if v, ok := conf["f4"]; ok {
			t := make([]int8, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseInt(s, 10, 8); err == nil {
					t = append(t, int8(ret))
				}
			}
			F4 = t
		}
		if v, ok := conf["f5"]; ok {
			t := make([]int16, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseInt(s, 10, 16); err == nil {
					t = append(t, int16(ret))
				}
			}
			F5 = t
		}
		if v, ok := conf["f6"]; ok {
			t := make([]int32, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseInt(s, 10, 32); err == nil {
					t = append(t, int32(ret))
				}
			}
			F6 = t
		}
		if v, ok := conf["f7"]; ok {
			t := make([]int64, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseInt(s, 10, 64); err == nil {
					t = append(t, ret)
				}
			}
			F7 = t
		}
		if v, ok := conf["f8"]; ok {
			t := make([]uint, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseUint(s, 10, 8); err == nil {
					t = append(t, uint(ret))
				}
			}
			F8 = t
		}
		if v, ok := conf["f9"]; ok {
			t := make([]uint8, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseUint(s, 10, 8); err == nil {
					t = append(t, uint8(ret))
				}
			}
			F9 = t
		}
		if v, ok := conf["f10"]; ok {
			t := make([]uint16, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseUint(s, 10, 16); err == nil {
					t = append(t, uint16(ret))
				}
			}
			F10 = t
		}
		if v, ok := conf["f11"]; ok {
			t := make([]uint32, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseUint(s, 10, 32); err == nil {
					t = append(t, uint32(ret))
				}
			}
			F11 = t
		}
		if v, ok := conf["f12"]; ok {
			t := make([]uint64, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseUint(s, 10, 64); err == nil {
					t = append(t, ret)
				}
			}
			F12 = t
		}
		if v, ok := conf["f13"]; ok {
			t := make([]float32, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseFloat(s, 32); err == nil {
					t = append(t, float32(ret))
				}
			}
			F13 = t
		}
		if v, ok := conf["f14"]; ok {
			t := make([]float64, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseFloat(s, 64); err == nil {
					t = append(t, ret)
				}
			}
			F14 = t
		}
		if v, ok := conf["f15"]; ok {
			t := make([]complex64, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseComplex(s, 64); err == nil {
					t = append(t, complex64(ret))
				}
			}
			F15 = t
		}
		if v, ok := conf["f16"]; ok {
			t := make([]complex128, 0)
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseComplex(s, 128); err == nil {
					t = append(t, ret)
				}
			}
			F16 = t
		}
		_ = &testCase{
			F1:  F1,
			F3:  F3,
			F4:  F4,
			F5:  F5,
			F6:  F6,
			F7:  F7,
			F8:  F8,
			F9:  F9,
			F10: F10,
			F11: F11,
			F12: F12,
			F13: F13,
			F14: F14,
			F15: F15,
			F16: F16,
		}
	}
}

func TestFillConfByTransformer_DirectMapStruct(t *testing.T) {
	type testCase struct {
		F1  map[string]struct{}     `transformer:"diMapStruct" conf:"f1" default:"" `
		F3  map[int]struct{}        `transformer:"diMapStruct" conf:"f3" default:"" `
		F4  map[int8]struct{}       `transformer:"diMapStruct" conf:"f4" default:"" `
		F5  map[int16]struct{}      `transformer:"diMapStruct" conf:"f5" default:"" `
		F6  map[int32]struct{}      `transformer:"diMapStruct" conf:"f6" default:"" `
		F7  map[int64]struct{}      `transformer:"diMapStruct" conf:"f7" default:"" `
		F8  map[uint]struct{}       `transformer:"diMapStruct" conf:"f8" default:"" `
		F9  map[uint8]struct{}      `transformer:"diMapStruct" conf:"f9" default:"" `
		F10 map[uint16]struct{}     `transformer:"diMapStruct" conf:"f10" default:"" `
		F11 map[uint32]struct{}     `transformer:"diMapStruct" conf:"f11" default:"" `
		F12 map[uint64]struct{}     `transformer:"diMapStruct" conf:"f12" default:"" `
		F13 map[float32]struct{}    `transformer:"diMapStruct" conf:"f13" default:"" `
		F14 map[float64]struct{}    `transformer:"diMapStruct" conf:"f14" default:"" `
		F15 map[complex64]struct{}  `transformer:"diMapStruct" conf:"f15" default:"" `
		F16 map[complex128]struct{} `transformer:"diMapStruct" conf:"f16" default:"" `
	}

	x := &testCase{}

	conf := map[string]string{
		"f1":  "1|2|3",
		"f3":  "3|4|5",
		"f4":  "4|5|6",
		"f5":  "5|6|7",
		"f6":  "6|7|8",
		"f7":  "7|8|9",
		"f8":  "8|9|10",
		"f9":  "9|10|11",
		"f10": "10|11|12",
		"f11": "11|12|13",
		"f12": "12|13|14",
		"f13": "13|14|15",
		"f14": "14|15|16",
		"f15": "1+5i|1+6i",
		"f16": "1+6i|1+7i",
	}

	y := &testCase{
		F1:  map[string]struct{}{"1": {}, "2": {}, "3": {}},
		F3:  map[int]struct{}{3: {}, 4: {}, 5: {}},
		F4:  map[int8]struct{}{4: {}, 5: {}, 6: {}},
		F5:  map[int16]struct{}{5: {}, 6: {}, 7: {}},
		F6:  map[int32]struct{}{6: {}, 7: {}, 8: {}},
		F7:  map[int64]struct{}{7: {}, 8: {}, 9: {}},
		F8:  map[uint]struct{}{8: {}, 9: {}, 10: {}},
		F9:  map[uint8]struct{}{9: {}, 10: {}, 11: {}},
		F10: map[uint16]struct{}{10: {}, 11: {}, 12: {}},
		F11: map[uint32]struct{}{11: {}, 12: {}, 13: {}},
		F12: map[uint64]struct{}{12: {}, 13: {}, 14: {}},
		F13: map[float32]struct{}{13: {}, 14: {}, 15: {}},
		F14: map[float64]struct{}{14: {}, 15: {}, 16: {}},
		F15: map[complex64]struct{}{complex(1, 5): {}, complex(1, 6): {}},
		F16: map[complex128]struct{}{complex(1, 6): {}, complex(1, 7): {}},
	}

	err := FillConfByParamsMap(x, true, conf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", x)

	t.Run("diMapStruct transformer", func(t *testing.T) {
		err := FillConfByParamsMap(x, true, conf)
		if err != nil {
			t.Errorf("FillConfByParamsMap() error = %v", err)
			return
		}
		if !reflect.DeepEqual(x, y) {
			t.Errorf("FillConfByParamsMap() got = %+v, want %+v", y, x)
		}
	})
}

func BenchmarkFillConfByTransformer_DirectMapStruct(b *testing.B) {
	type testCase struct {
		F1  map[string]struct{}     `transformer:"diMapStruct" conf:"f1" default:"" `
		F3  map[int]struct{}        `transformer:"diMapStruct" conf:"f3" default:"" `
		F4  map[int8]struct{}       `transformer:"diMapStruct" conf:"f4" default:"" `
		F5  map[int16]struct{}      `transformer:"diMapStruct" conf:"f5" default:"" `
		F6  map[int32]struct{}      `transformer:"diMapStruct" conf:"f6" default:"" `
		F7  map[int64]struct{}      `transformer:"diMapStruct" conf:"f7" default:"" `
		F8  map[uint]struct{}       `transformer:"diMapStruct" conf:"f8" default:"" `
		F9  map[uint8]struct{}      `transformer:"diMapStruct" conf:"f9" default:"" `
		F10 map[uint16]struct{}     `transformer:"diMapStruct" conf:"f10" default:"" `
		F11 map[uint32]struct{}     `transformer:"diMapStruct" conf:"f11" default:"" `
		F12 map[uint64]struct{}     `transformer:"diMapStruct" conf:"f12" default:"" `
		F13 map[float32]struct{}    `transformer:"diMapStruct" conf:"f13" default:"" `
		F14 map[float64]struct{}    `transformer:"diMapStruct" conf:"f14" default:"" `
		F15 map[complex64]struct{}  `transformer:"diMapStruct" conf:"f15" default:"" `
		F16 map[complex128]struct{} `transformer:"diMapStruct" conf:"f16" default:"" `
	}

	x := &testCase{}

	conf := map[string]string{
		"f1":  "1|2|3",
		"f3":  "3|4|5",
		"f4":  "4|5|6",
		"f5":  "5|6|7",
		"f6":  "6|7|8",
		"f7":  "7|8|9",
		"f8":  "8|9|10",
		"f9":  "9|10|11",
		"f10": "10|11|12",
		"f11": "11|12|13",
		"f12": "12|13|14",
		"f13": "13|14|15",
		"f14": "14|15|16",
		"f15": "1+5i|1+6i",
		"f16": "1+6i|1+7i",
	}

	for i := 0; i < b.N; i++ {
		err := FillConfByParamsMap(x, true, conf)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkNoReflect_DirectMapStruct(b *testing.B) {
	type testCase struct {
		F1  map[string]struct{}
		F3  map[int]struct{}
		F4  map[int8]struct{}
		F5  map[int16]struct{}
		F6  map[int32]struct{}
		F7  map[int64]struct{}
		F8  map[uint]struct{}
		F9  map[uint8]struct{}
		F10 map[uint16]struct{}
		F11 map[uint32]struct{}
		F12 map[uint64]struct{}
		F13 map[float32]struct{}
		F14 map[float64]struct{}
		F15 map[complex64]struct{}
		F16 map[complex128]struct{}
	}

	conf := map[string]string{
		"f1":  "1|2|3",
		"f3":  "3|4|5",
		"f4":  "4|5|6",
		"f5":  "5|6|7",
		"f6":  "6|7|8",
		"f7":  "7|8|9",
		"f8":  "8|9|10",
		"f9":  "9|10|11",
		"f10": "10|11|12",
		"f11": "11|12|13",
		"f12": "12|13|14",
		"f13": "13|14|15",
		"f14": "14|15|16",
		"f15": "1+5i|1+6i",
		"f16": "1+6i|1+7i",
	}

	F1 := map[string]struct{}{"0": {}, "2": {}, "3": {}}
	F3 := map[int]struct{}{0: {}, 4: {}, 5: {}}
	F4 := map[int8]struct{}{0: {}, 5: {}, 6: {}}
	F5 := map[int16]struct{}{0: {}, 6: {}, 7: {}}
	F6 := map[int32]struct{}{0: {}, 7: {}, 8: {}}
	F7 := map[int64]struct{}{0: {}, 8: {}, 9: {}}
	F8 := map[uint]struct{}{0: {}, 9: {}, 10: {}}
	F9 := map[uint8]struct{}{0: {}, 10: {}, 11: {}}
	F10 := map[uint16]struct{}{0: {}, 11: {}, 12: {}}
	F11 := map[uint32]struct{}{0: {}, 12: {}, 13: {}}
	F12 := map[uint64]struct{}{0: {}, 13: {}, 14: {}}
	F13 := map[float32]struct{}{0: {}, 14: {}, 15: {}}
	F14 := map[float64]struct{}{0: {}, 15: {}, 16: {}}
	F15 := map[complex64]struct{}{complex(0, 5): {}, complex(1, 6): {}}
	F16 := map[complex128]struct{}{complex(0, 6): {}, complex(1, 7): {}}
	for i := 0; i < b.N; i++ {
		if v, ok := conf["f1"]; ok {
			F1 = make(map[string]struct{})
			for _, s := range strings.Split(v, "|") {
				F1[s] = struct{}{}
			}
		}
		if v, ok := conf["f3"]; ok {
			F3 = make(map[int]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.Atoi(s); err == nil {
					F3[ret] = struct{}{}
				}
			}
		}
		if v, ok := conf["f4"]; ok {
			F4 = make(map[int8]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseInt(s, 10, 8); err == nil {
					F4[int8(ret)] = struct{}{}
				}
			}
		}
		if v, ok := conf["f5"]; ok {
			F5 = make(map[int16]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseInt(s, 10, 16); err == nil {
					F5[int16(ret)] = struct{}{}
				}
			}
		}
		if v, ok := conf["f6"]; ok {
			F6 = make(map[int32]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseInt(s, 10, 32); err == nil {
					F6[int32(ret)] = struct{}{}
				}
			}
		}
		if v, ok := conf["f7"]; ok {
			F7 = make(map[int64]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseInt(s, 10, 64); err == nil {
					F7[ret] = struct{}{}
				}
			}
		}
		if v, ok := conf["f8"]; ok {
			F8 = make(map[uint]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseUint(s, 10, 8); err == nil {
					F8[uint(ret)] = struct{}{}
				}
			}
		}
		if v, ok := conf["f9"]; ok {
			F9 = make(map[uint8]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseUint(s, 10, 8); err == nil {
					F9[uint8(ret)] = struct{}{}
				}
			}
		}
		if v, ok := conf["f10"]; ok {
			F10 = make(map[uint16]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseUint(s, 10, 16); err == nil {
					F10[uint16(ret)] = struct{}{}
				}
			}
		}
		if v, ok := conf["f11"]; ok {
			F11 = make(map[uint32]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseUint(s, 10, 32); err == nil {
					F11[uint32(ret)] = struct{}{}
				}
			}
		}
		if v, ok := conf["f12"]; ok {
			F12 = make(map[uint64]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseUint(s, 10, 64); err == nil {
					F12[ret] = struct{}{}
				}
			}
		}
		if v, ok := conf["f13"]; ok {
			F13 = make(map[float32]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseFloat(s, 32); err == nil {
					F13[float32(ret)] = struct{}{}
				}
			}
		}
		if v, ok := conf["f14"]; ok {
			F14 = make(map[float64]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseFloat(s, 64); err == nil {
					F14[ret] = struct{}{}
				}
			}
		}
		if v, ok := conf["f15"]; ok {
			F15 = make(map[complex64]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseComplex(s, 64); err == nil {
					F15[complex64(ret)] = struct{}{}
				}
			}
		}
		if v, ok := conf["f16"]; ok {
			F16 = make(map[complex128]struct{})
			for _, s := range strings.Split(v, "|") {
				if ret, err := strconv.ParseComplex(s, 128); err == nil {
					F16[ret] = struct{}{}
				}
			}
		}
		_ = &testCase{
			F1:  F1,
			F3:  F3,
			F4:  F4,
			F5:  F5,
			F6:  F6,
			F7:  F7,
			F8:  F8,
			F9:  F9,
			F10: F10,
			F11: F11,
			F12: F12,
			F13: F13,
			F14: F14,
			F15: F15,
			F16: F16,
		}
	}
}

func TestFillConfByTransformer_Custom(t *testing.T) {
	now := time.Now()

	type testCase struct {
		F5 int64                       `transformer:"timeDuration" conf:"f5" `
		F6 map[int64]string            `transformer:"confMapIntString" conf:"f6" `
		F7 map[int]map[uint64]struct{} `transformer:"confMapIntUintMap" conf:"f7" `
	}

	x := &testCase{}

	conf := map[string]string{
		"f5": "-12h",
		"f6": "1:a|2:b",
		"f7": "1:2|3|4;2:2|3|4",
	}

	y := &testCase{
		F5: now.Add(-12 * time.Hour).Unix(),
		F6: map[int64]string{1: "a", 2: "b"},
		F7: map[int]map[uint64]struct{}{1: {uint64(2): {}, uint64(3): {}, uint64(4): {}}, 2: {uint64(2): {}, uint64(3): {}, uint64(4): {}}},
	}

	err := FillConfByParamsMap(x, true, conf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", x)

	t.Run("custom transformer", func(t *testing.T) {
		err := FillConfByParamsMap(x, true, conf)
		if err != nil {
			t.Errorf("FillConfByParamsMap() error = %v", err)
			return
		}
		if !reflect.DeepEqual(x, y) {
			t.Errorf("FillConfByParamsMap() got = %+v, want %+v", x, y)
		}
	})
}

// func BenchmarkFillConfByTransformer_Custom(b *testing.B) {
// 	type testCase struct {
// 		F5 int64            `transformer:"timeDuration" conf:"f5" `
// 		F6 map[int64]string `transformer:"confMapIntString" conf:"f6" `
// 	}
//
// 	x := &testCase{}
//
// 	conf := map[string]string{
// 		"f5": "-12h",
// 		"f6": "1:a|2:b",
// 	}
//
// 	for i := 0; i < b.N; i++ {
// 		err := FillConfByParamsMap(x, true, conf)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

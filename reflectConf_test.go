package reflectconf

import (
	"bytes"
	"encoding/json"
	"fmt"
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
// BenchmarkFillConfByTransformer                   1818283              3305 ns/op            3248 B/op         56 allocs/op
// BenchmarkFillConfByHardCode                      5407528              1112 ns/op             632 B/op         12 allocs/op

func TestFillConfByTransformer(t *testing.T) {
	x := &testCase{}
	err := FillConfByParamsMap(x, true, getConfig3(), getConfig2(), getConfig1())
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := json.Marshal(x)
	var out bytes.Buffer
	_ = json.Indent(&out, bs, "", "\t")
	fmt.Printf("%v\n", out.String())
}

func BenchmarkFillConfByTransformer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := &testCase{}
		_ = FillConfByParamsMap(x, true, getConfig3(), getConfig2(), getConfig1())
	}
}

func TestFillConfByHardCode(t *testing.T) {
	x := &testCase{}
	hardCode(x)
	bs, _ := json.Marshal(x)
	var out bytes.Buffer
	_ = json.Indent(&out, bs, "", "\t")
	fmt.Printf("%v\n", out.String())
}

func BenchmarkFillConfByHardCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := &testCase{}
		hardCode(x)
	}
}

type testCase struct {
	TString    string              `transformer:"direct" conf:"tString" default:"abc" `
	TSlice     []float32           `transformer:"diSlice" conf:"tSlice" default:"1.2|3.4|5.6|7.8"`
	TMap1      map[uint32]struct{} `transformer:"diMapStruct" conf:"tMap1" `
	TMap2      map[int64]string    `transformer:"confMapIntString" conf:"tMap2" default:"1:a|2:b|3:c" `
	TBool      bool                `transformer:"direct" conf:"tBool" `
	TTimeStamp int64               `transformer:"timeDuration" conf:"tTimeStamp" default:"-12h" `
}

func getConfig1() map[string]string {
	return map[string]string{
		"tString":    "opq",
		"tMap1":      "1|2|3",
		"tTimeStamp": "-18h",
	}
}
func getConfig2() map[string]string {
	return map[string]string{
		"tMap1": "7|8|9",
		"tMap2": "7:x|8:y|9:z",
	}
}

func getConfig3() map[string]string {
	return map[string]string{
		"tString":    "xyz",
		"tBool":      "t",
		"tTimeStamp": "-20h",
	}
}

func hardCode(x *testCase) {
	c := mergeConfMaps([]map[string]string{getConfig3(), getConfig2(), getConfig1()})
	defaultC := map[string]string{
		"tString":    "abc",
		"tSlice":     "1.2|3.4|5.6|7.8",
		"tMap2":      "1:a|2:b|3:c",
		"tTimeStamp": "-12h",
	}
	for k, v := range defaultC {
		if _, ok := c[k]; !ok {
			c[k] = v
		}
	}
	if v, ok := c["tString"]; ok && v != "" {
		x.TString = v
	}
	if v, ok := c["tSlice"]; ok && v != "" {
		vv := strings.Split(v, "|")
		var vvv []float32
		for i := 0; i < len(vv); i++ {
			if float, err := strconv.ParseFloat(vv[i], 32); err == nil {
				vvv = append(vvv, float32(float))
			}
		}
		if len(vvv) > 0 {
			x.TSlice = vvv
		}
	}
	if v, ok := c["tMap1"]; ok && v != "" {
		vv := strings.Split(v, "|")
		if len(vv) > 0 {
			vvv := make(map[uint32]struct{}, len(vv))
			for i := 0; i < len(vv); i++ {
				if _v, err := strconv.ParseUint(vv[i], 10, 32); err == nil {
					vvv[uint32(_v)] = struct{}{}
				}
			}
			if len(vvv) > 0 {
				x.TMap1 = vvv
			}
		}

	}
	if v, ok := c["tMap2"]; ok && v != "" {
		vv := strings.Split(v, "|")
		if len(vv) > 0 {
			vvv := make(map[int64]string, len(vv))
			for _, v := range vv {
				tArr := strings.Split(v, ":")
				if len(tArr) == 2 {
					if ret, err := strconv.ParseInt(tArr[0], 10, 64); err == nil {
						vvv[ret] = tArr[1]
					}
				}
			}
			if len(vvv) > 0 {
				x.TMap2 = vvv
			}
		}

	}
	if v, ok := c["tBool"]; ok && v != "" {
		if _v, err := strconv.ParseBool(v); err == nil {
			x.TBool = _v
		}
	}
	if v, ok := c["tTimeStamp"]; ok && v != "" {
		if vv, err := time.ParseDuration(v); err == nil {
			x.TTimeStamp = time.Now().Add(vv).Unix()
		}
	}
}

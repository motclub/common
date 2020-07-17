package templatex

import (
	"testing"
	"time"
)

func TestFastExec(t *testing.T) {
	tpl := `
{{- eq 1 1 -}}
{{if and (eq .b true) (le .c 12.9)}}
{{'a'}}
{{end}}
{{range .f }}
{{- . }}
{{end}}
{{range $k, $v := .g }}
{{$k}}={{$v}}
{{end}}
`
	res, err := FastExec(tpl, map[string]interface{}{
		"a": 1,
		"b": true,
		"c": 12.9,
		"d": "T1",
		"e": time.Now(),
		"f": []int{100, 200, 300, 400, 500},
		"g": map[string]interface{}{
			"ka": 101,
			"ab": true,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

package parser

import (
	"github.com/motclub/common/std"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewParser(t *testing.T) {
	data, err := NewParser(std.D{
		"a": "hello",
		"b": true,
		"c": 12.4,
		"d": 190,
	})
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, "hello", data.GetString("a"))
	assert.Equal(t, true, data.GetBool("b"))
	assert.Equal(t, 12.4, data.GetFloat("c"))
	assert.Equal(t, 190, data.GetInt("d"))
	assert.Equal(t, 10.2, data.DefaultGetFloat("e", 10.2))
}

func Test_parser_DefaultGet(t *testing.T) {
	type T struct {
		Name        string        `json:"name"`
		DisplayName string        `json:"display_name"`
		Disabled    bool          `json:"disabled"`
		Exp         time.Duration `json:"exp"`
		DisabledAt  time.Time     `json:"disabled_at"`
	}
	data, err := NewParser(std.D{
		"a": std.D{
			"aa": std.D{
				"aaa": std.D{
					"name":         "t",
					"display_name": "T",
					"disabled":     true,
					"exp":          time.Minute * 60,
					"disabled_at":  time.Now(),
				},
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	var t1 T
	data.DefaultGet("a.aa.aaa", &t1, T{Name: "dt"})
	assert.Equal(t, "t", t1.Name)
	assert.Equal(t, true, t1.Disabled)

	var t2 T
	data.DefaultGet("a.aa.aaa.aaaa", &t2, T{Name: "dt"})
	assert.Equal(t, "dt", t2.Name)
	assert.Equal(t, false, t2.Disabled)
}

func Test_parser_DefaultGetBool(t *testing.T) {
	data, err := NewParser(std.D{
		"a": true,
		"b": false,
		"c": std.D{
			"cc": std.D{
				"ccc": true,
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, true, data.DefaultGetBool("a", false))
	assert.Equal(t, true, data.DefaultGetBool("aa", true))
	assert.Equal(t, false, data.DefaultGetBool("b", true))
	assert.Equal(t, true, data.DefaultGetBool("c.cc.ccc", false))
}

func Test_parser_DefaultGetFloat(t *testing.T) {
	data, err := NewParser(std.D{
		"a": 1,
		"b": 2.1,
		"c": std.D{
			"cc": std.D{
				"ccc": 3.141592653,
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 1.0, data.DefaultGetFloat("a", 1.1))
	assert.Equal(t, 1.2, data.DefaultGetFloat("aa", 1.2))
	assert.Equal(t, 2.1, data.DefaultGetFloat("b", 2.2))
	assert.Equal(t, 3.141592653, data.DefaultGetFloat("c.cc.ccc", 3.15))
	assert.Equal(t, 3.2, data.DefaultGetFloat("c.cc.ddd", 3.2))
}

func Test_parser_parsePath(t *testing.T) {
	c := &parser{}
	assert.Equal(t, []string{"a"}, c.parsePath("a"))
	assert.Equal(t, []string{"a", "b"}, c.parsePath("a.b"))
	assert.Equal(t, []string{"a", "b", "c"}, c.parsePath("a.b.c"))
	assert.Equal(t, []string{"a", "b", "c__d"}, c.parsePath("a.b.c__d"))
}

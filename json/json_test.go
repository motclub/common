package json

import (
	"testing"
)

func TestLowerCaseWithUnderscores(t *testing.T) {
	data := struct {
		Code         string
		Name         string
		DefaultValue string
		IsRequired   bool
	}{
		Code:         "foo",
		Name:         "bar",
		DefaultValue: "foobar",
		IsRequired:   false,
	}
	b, err := STD().Marshal(&data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(b))
}

package templatex

import (
	"bytes"
	"github.com/Masterminds/sprig"
	html "html/template"
	txt "text/template"
)

func FastExecTXT(text string, data interface{}, fns ...txt.FuncMap) (string, error) {
	tpl := txt.New("").Funcs(sprig.TxtFuncMap())
	if len(fns) > 0 {
		for _, item := range fns {
			tpl.Funcs(item)
		}
	}

	v := txt.Must(tpl.Parse(text))
	var b bytes.Buffer
	if err := v.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}

func FastExecHTML(text string, data interface{}, fns ...html.FuncMap) (string, error) {
	tpl := html.New("").Funcs(sprig.HtmlFuncMap())
	if len(fns) > 0 {
		for _, item := range fns {
			tpl.Funcs(item)
		}
	}

	v := html.Must(tpl.Parse(text))
	var b bytes.Buffer
	if err := v.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}

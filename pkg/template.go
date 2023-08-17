package pkg

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"text/template"
)

func ParseTplFromVals(vals interface{}, tplStr string) (string, error) {
	errInfo := fmt.Sprintf("parse template from string error")
	var err error
	var strOutput string

	var buf bytes.Buffer
	gotpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(tplStr)
	if err != nil {
		err = fmt.Errorf("%s: create template error: %s", errInfo, err.Error())
		return strOutput, err
	}
	err = gotpl.Execute(&buf, vals)
	if err != nil {
		err = fmt.Errorf("%s: parse template error: %s", errInfo, err.Error())
		return strOutput, err
	}
	strOutput = buf.String()

	return strOutput, err
}

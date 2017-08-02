package export

import (
	"bytes"

	"fmt"

	"github.com/iris-contrib/parrot/parrot-api/model"
)

type AppleStrings struct{}

func (e *AppleStrings) FileExtension() string {
	return "strings"
}

func (e *AppleStrings) Export(locale *model.Locale) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	for k, v := range locale.Pairs {
		_, err := buf.WriteString(fmt.Sprintf("\"%s\" = \"%s\";\n", k, v))
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

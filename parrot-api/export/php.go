package export

import (
	"bytes"

	"fmt"

	"github.com/iris-contrib/parrot/parrot-api/model"
)

type PHP struct{}

func (e *PHP) FileExtension() string {
	return "php"
}

func (e *PHP) Export(locale *model.Locale) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	_, err := buf.WriteString(fmt.Sprintf("<?php\n$%s = array(\n", locale.Ident))
	if err != nil {
		return nil, err
	}

	i := 0
	max := len(locale.Pairs)
	for k, v := range locale.Pairs {
		i++
		str := "    \"%s\" => \"%s\""
		if i < max {
			str += ","
		}
		_, err := buf.WriteString(fmt.Sprintf(str+"\n", k, v))
		if err != nil {
			return nil, err
		}
	}

	_, err = buf.WriteString(");>")
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

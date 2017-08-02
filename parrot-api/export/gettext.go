package export

import (
	"bytes"

	"fmt"

	"github.com/iris-contrib/parrot/parrot-api/model"
)

type Gettext struct{}

func (e *Gettext) FileExtension() string {
	return "po"
}

func (e *Gettext) Export(locale *model.Locale) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	_, err := buf.WriteString(fmt.Sprintf("msgid \"\"\nmsgstr \"\"\n\"MIME-Version: 1.0\"\n\"Content-Type: text/plain; charset=UTF-8\"\n\"Content-Transfer-Encoding: 8bit\"\n\"Language: %s\"\n\n",
		locale.Ident))
	if err != nil {
		return nil, err
	}

	for k, v := range locale.Pairs {
		_, err := buf.WriteString(fmt.Sprintf("msgid \"%s\"\nmsgstr \"%s\"\n\n", k, v))
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

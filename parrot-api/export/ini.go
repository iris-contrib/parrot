package export

import (
	"bytes"

	"github.com/iris-contrib/parrot/parrot-api/model"
	"github.com/go-ini/ini"
)

type INI struct{}

func (e *INI) FileExtension() string {
	return "ini"
}

func (e *INI) Export(locale *model.Locale) ([]byte, error) {
	outFile := ini.Empty()

	section := outFile.Section(locale.Ident)

	for k, v := range locale.Pairs {
		_, err := section.NewKey(k, v)
		if err != nil {
			return nil, err
		}
	}

	buf := bytes.NewBuffer(nil)
	_, err := outFile.WriteTo(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

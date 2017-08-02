package export

import (
	"bytes"

	"encoding/csv"

	"github.com/iris-contrib/parrot/parrot-api/model"
)

type CSV struct{}

func (e *CSV) FileExtension() string {
	return "csv"
}

func (e *CSV) Export(locale *model.Locale) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	wr := csv.NewWriter(buf)

	for k, v := range locale.Pairs {
		err := wr.Write([]string{k, v})
		if err != nil {
			return nil, err
		}
	}

	wr.Flush()

	return buf.Bytes(), nil
}

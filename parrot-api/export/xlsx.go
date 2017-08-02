package export

import (
	"bytes"

	"github.com/iris-contrib/parrot/parrot-api/model"
	"github.com/tealeg/xlsx"
)

type XLSX struct{}

func (e *XLSX) FileExtension() string {
	return "xlsx"
}

func (e *XLSX) Export(locale *model.Locale) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	f := xlsx.NewFile()

	sheet, err := f.AddSheet(locale.Ident)
	if err != nil {
		return nil, err
	}

	for k, v := range locale.Pairs {
		r := sheet.AddRow()
		r.AddCell().Value = k
		r.AddCell().Value = v
	}

	err = f.Write(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

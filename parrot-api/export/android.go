package export

import (
	"bytes"

	"encoding/xml"

	"github.com/iris-contrib/parrot/parrot-api/model"
)

type Android struct{}

func (e *Android) FileExtension() string {
	return "xml"
}

func (e *Android) Export(locale *model.Locale) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	encoder := xml.NewEncoder(buf)

	encoder.Indent("", "  ")

	_, err := buf.Write([]byte(xml.Header))
	if err != nil {
		return nil, err
	}

	err = encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "resources"}})
	if err != nil {
		return nil, err
	}

	for k, v := range locale.Pairs {
		err = encoder.EncodeToken(xml.StartElement{
			Name: xml.Name{Local: "string"},
			Attr: []xml.Attr{xml.Attr{Name: xml.Name{Local: "name"}, Value: k}},
		})
		if err != nil {
			return nil, err
		}
		err = encoder.EncodeToken(xml.CharData([]byte(v)))
		if err != nil {
			return nil, err
		}
		err = encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "string"}})
		if err != nil {
			return nil, err
		}
	}

	err = encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "resources"}})
	if err != nil {
		return nil, err
	}

	err = encoder.Flush()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

package export

import (
	"bytes"

	"encoding/xml"

	"github.com/iris-contrib/parrot/parrot-api/model"
)

type JavaXML struct{}

func (e *JavaXML) FileExtension() string {
	return "xml"
}

func (e *JavaXML) Export(locale *model.Locale) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	encoder := xml.NewEncoder(buf)

	encoder.Indent("", "  ")

	_, err := buf.Write([]byte(xml.Header))
	if err != nil {
		return nil, err
	}

	_, err = buf.Write([]byte("<!DOCTYPE properties SYSTEM \"[http://java.sun.com/dtd/properties.dtd](http://java.sun.com/dtd/properties.dtd)\">\n"))
	if err != nil {
		return nil, err
	}

	err = encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "properties"}})
	if err != nil {
		return nil, err
	}

	for k, v := range locale.Pairs {
		err = encoder.EncodeToken(xml.StartElement{
			Name: xml.Name{Local: "entry"},
			Attr: []xml.Attr{xml.Attr{Name: xml.Name{Local: "key"}, Value: k}},
		})
		if err != nil {
			return nil, err
		}
		err = encoder.EncodeToken(xml.CharData([]byte(v)))
		if err != nil {
			return nil, err
		}
		err = encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "entry"}})
		if err != nil {
			return nil, err
		}
	}

	err = encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "properties"}})
	if err != nil {
		return nil, err
	}

	err = encoder.Flush()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

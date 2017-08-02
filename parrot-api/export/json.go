package export

import (
	"encoding/json"

	"github.com/iris-contrib/parrot/parrot-api/model"
)

type JSON struct{}

func (e *JSON) FileExtension() string {
	return "json"
}

func (e *JSON) Export(locale *model.Locale) ([]byte, error) {
	return json.MarshalIndent(locale.Pairs, "", "    ")
}

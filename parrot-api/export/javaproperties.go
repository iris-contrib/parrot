package export

import (
	"bytes"
	"strconv"

	"fmt"

	"strings"

	"github.com/iris-contrib/parrot/parrot-api/model"
)

type JavaProperties struct{}

func (e *JavaProperties) FileExtension() string {
	return "properties"
}

func (e *JavaProperties) Export(locale *model.Locale) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	for k, v := range locale.Pairs {
		var newKey []string
		for _, chart := range []rune(k) {
			quoted := strconv.QuoteRuneToASCII(chart)        // quoted = "'\u554a'"
			newKey = append(newKey, quoted[1:len(quoted)-1]) // unquoted = "\u554a"
		}

		var newValue []string
		for _, chart := range []rune(v) {
			quoted := strconv.QuoteRuneToASCII(chart)            // quoted = "'\u554a'"
			newValue = append(newValue, quoted[1:len(quoted)-1]) // unquoted = "\u554a"
		}

		_, err := buf.WriteString(fmt.Sprintf("%s = %s\n", strings.Replace(strings.Join(newKey[:], ""), " ", "\\ ", -1), strings.Join(newValue[:], "")))
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

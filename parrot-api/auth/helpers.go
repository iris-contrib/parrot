package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kataras/golog"
	"github.com/kataras/iris"
)

// getAuthHeaderToken gets the token string from the HTTP Authorization header.
func getAuthHeaderToken(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", fmt.Errorf("no auth header")
	}

	token = sanitizeBearerToken(token)

	return token, nil
}

// getJSONBodyToken gets the token string from the HTTP JSON body.
func getJSONBodyToken(ctx iris.Context) (string, error) {
	var body map[string]interface{}
	err := ctx.ReadJSON(&body)
	if err != nil {
		return "", err
	}
	token, ok := body["token"].(string)
	if token == "" || !ok {
		return "", fmt.Errorf("no auth header")
	}

	token = sanitizeBearerToken(token)

	return token, nil
}

// sanitizeBearerToken extracts the token part from the token string.
func sanitizeBearerToken(token string) string {
	if len(token) > 6 && strings.ToUpper(token[0:7]) == "BEARER " {
		return token[7:]
	}
	return token
}

func RenderJSON(w http.ResponseWriter, status int, headers map[string]string, payload interface{}) {
	h := w.Header()
	h.Set("Content-Type", "application/json")
	for k, v := range headers {
		h.Set(k, v)
	}
	w.WriteHeader(status)

	encoded, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		golog.Error(err)
	}

	w.Write(encoded)
}

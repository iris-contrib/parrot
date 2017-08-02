// Package render handles the rending of common API results.
package render

import (
	"github.com/kataras/iris"
)

var (
	jsonContentType        = "application/json"
	jsonContentTypeCharset = "application/json; charset=utf-8"
)

type apiResponseBody struct {
	responseMeta `json:"meta,omitempty"`
	Payload      interface{} `json:"payload,omitempty"`
}

type responseMeta struct {
	Status int   `json:"status,omitempty"`
	Error  error `json:"error,omitempty"`
}

// Error writes an API error to the response.
func Error(ctx iris.Context, status int, err error) {
	body := apiResponseBody{
		responseMeta: responseMeta{
			Status: status,
			Error:  err},
		Payload: nil}
	ctx.StatusCode(status)
	ctx.JSON(body)
}

// JSON writes a payload as json to the response.
func JSON(ctx iris.Context, status int, payload interface{}) {
	body := apiResponseBody{
		responseMeta: responseMeta{
			Status: status},
		Payload: payload}

	ctx.StatusCode(status)
	ctx.JSON(body)
}

// JSONWithHeaders writes a payload as json to the response and includes the provided headers.
func JSONWithHeaders(ctx iris.Context, status int, headers map[string]string, payload interface{}) {
	h := ctx.ResponseWriter().Header()
	for k, v := range headers {
		h.Set(k, v)
	}

	body := apiResponseBody{
		responseMeta: responseMeta{
			Status: status},
		Payload: payload}

	ctx.StatusCode(status)
	ctx.JSON(body)
}

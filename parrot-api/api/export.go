package api

import (
	"bytes"
	"fmt"
	"strings"

"github.com/kataras/iris/v12"
	apiErrors "github.com/iris-contrib/parrot/parrot-api/errors"
	"github.com/iris-contrib/parrot/parrot-api/export"
)

// exportLocale is an API endpoint for exporting locale pairs.
func exportLocale(ctx iris.Context) {

	projectID := ctx.Params().Get("projectID")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	localeIdent := ctx.Params().Get("localeIdent")
	if projectID == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}
	i18nType := ctx.Params().Get("type")
	if i18nType == "" {
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	locale, err := store.GetProjectLocaleByIdent(projectID, localeIdent)
	if err != nil {
		handleError(ctx, err)
		return
	}

	var exporter export.Exporter
	switch strings.ToLower(i18nType) {
	case "keyvaluejson":
		exporter = &export.JSON{}
	case "po":
		exporter = &export.Gettext{}
	case "strings":
		exporter = &export.AppleStrings{}
	case "properties":
		exporter = &export.JavaProperties{}
	case "xmlproperties":
		exporter = &export.JavaXML{}
	case "android":
		exporter = &export.Android{}
	case "php":
		exporter = &export.PHP{}
	case "xlsx":
		exporter = &export.XLSX{}
	case "csv":
		exporter = &export.CSV{}
	case "yaml":
		exporter = &export.Yaml{}
	case "ini":
		exporter = &export.INI{}
	default:
		handleError(ctx, apiErrors.ErrBadRequest)
		return
	}

	result, err := exporter.Export(locale)
	if err != nil {
		handleError(ctx, err)
		return
	}

	filename := fmt.Sprintf("%s.%s", localeIdent, exporter.FileExtension())

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(result)))

	buf := bytes.NewBuffer(result)
	_, err = buf.WriteTo(ctx.ResponseWriter())
	if err != nil {
		handleError(ctx, err)
		return
	}
}

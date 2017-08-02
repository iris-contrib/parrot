package postgres

import (
	"database/sql"

	"github.com/iris-contrib/parrot/parrot-api/model"
	"github.com/lib/pq/hstore"
)

func (db *PostgresDB) CreateLocale(loc model.Locale) (*model.Locale, error) {
	h := hstore.Hstore{}
	h.Map = make(map[string]sql.NullString)
	for k, v := range loc.Pairs {
		h.Map[k] = sql.NullString{String: v, Valid: true}
	}
	values, err := h.Value()
	if err != nil {
		return nil, parseError(err)
	}

	row := db.QueryRow("INSERT INTO locales (ident, language, country, pairs, project_id) VALUES($1, $2, $3, $4, $5) RETURNING id",
		loc.Ident, loc.Language, loc.Country, values, loc.ProjectID)
	err = row.Scan(&loc.ID)
	return &loc, parseError(err)
}

func (db *PostgresDB) UpdateLocalePairs(projID string, localeIdent string, pairs map[string]string) (*model.Locale, error) {
	h := hstore.Hstore{}
	h.Map = make(map[string]sql.NullString)
	for k, v := range pairs {
		h.Map[k] = sql.NullString{String: v, Valid: true}
	}

	values, err := h.Value()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow("UPDATE locales SET pairs = $1 WHERE project_id = $2 AND ident = $3 RETURNING *", values, projID, localeIdent)
	loc := model.Locale{}
	err = row.Scan(&loc.ID, &loc.Ident, &loc.Language, &loc.Country, &h, &loc.ProjectID)
	if err != nil {
		return nil, parseError(err)
	}

	loc.Pairs = make(map[string]string)
	for k, v := range h.Map {
		if v.Valid {
			loc.Pairs[k] = v.String
		}
	}
	return &loc, nil
}

func (db *PostgresDB) DeleteLocale(projID string, ident string) error {
	_, err := db.Exec("DELETE FROM locales WHERE project_id = $1 AND ident = $2", projID, ident)
	return parseError(err)
}

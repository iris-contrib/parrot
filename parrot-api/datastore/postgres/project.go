package postgres

import (
	"database/sql"

	"github.com/iris-contrib/parrot/parrot-api/datastore/errors"
	"github.com/iris-contrib/parrot/parrot-api/model"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
)

func (db *PostgresDB) GetProject(id string) (*model.Project, error) {
	p := model.Project{}
	row := db.QueryRow("SELECT * FROM projects WHERE id = $1", id)

	keys := pq.StringArray{}
	err := row.Scan(&p.ID, &p.Name, &keys)
	if err != nil {
		return nil, parseError(err)
	}

	p.Keys = make([]string, len(keys))
	for i, v := range keys {
		p.Keys[i] = v
	}

	return &p, nil
}

func (db *PostgresDB) CreateProject(project model.Project) (*model.Project, error) {
	keys := make(pq.StringArray, len(project.Keys))
	for i, v := range project.Keys {
		keys[i] = v
	}

	values, err := keys.Value()
	if err != nil {
		return nil, parseError(err)
	}

	row := db.QueryRow("INSERT INTO projects (name, keys) VALUES($1, $2) RETURNING *", project.Name, values)
	result := model.Project{}
	keys = pq.StringArray{}
	err = row.Scan(&result.ID, &result.Name, &keys)
	if err != nil {
		return nil, parseError(err)
	}

	result.Keys = make([]string, len(keys))
	for i, v := range keys {
		result.Keys[i] = v
	}

	return &result, nil
}

func (db *PostgresDB) UpdateProjectName(projectID, name string) (*model.Project, error) {
	row := db.QueryRow("UPDATE projects SET name = $1 WHERE id = $2 RETURNING id, name, keys", name, projectID)
	keys := pq.StringArray{}
	result := model.Project{}
	err := row.Scan(&result.ID, &result.Name, &keys)
	if err != nil {
		return nil, parseError(err)
	}

	result.Keys = make([]string, len(keys))
	for i, v := range keys {
		result.Keys[i] = v
	}

	return &result, nil
}

func (db *PostgresDB) AddProjectKey(projectID, key string) (*model.Project, error) {
	// TODO optimize this
	project, err := db.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	for _, v := range project.Keys {
		if v == key {
			return nil, parseError(errors.ErrAlreadyExists)
		}
	}

	row := db.QueryRow("UPDATE projects SET keys = keys || array[$1] WHERE id = $2 AND NOT(keys @> array[$1]) RETURNING *", key, projectID)
	keys := pq.StringArray{}
	result := model.Project{}
	err = row.Scan(&result.ID, &result.Name, &keys)
	if err != nil {
		return nil, parseError(err)
	}

	result.Keys = make([]string, len(keys))
	for i, v := range keys {
		result.Keys[i] = v
	}

	return &result, nil
}

func (db *PostgresDB) UpdateProjectKey(projectID, oldKey, newKey string) (*model.Project, int, error) {
	// TODO: optimize this
	// Before transaction begin, check if project key is present and newKey is not present
	project, err := db.GetProject(projectID)
	if err != nil {
		return nil, -1, err
	}

	oldKeyPresent := false
	newKeyPresent := false
	for _, v := range project.Keys {
		if v == oldKey {
			oldKeyPresent = true
		} else if v == newKey {
			newKeyPresent = true
		}
	}
	if !oldKeyPresent {
		return nil, -1, parseError(errors.ErrNotFound)
	}
	if newKeyPresent {
		return nil, -1, parseError(errors.ErrAlreadyExists)
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, -1, err
	}
	defer tx.Rollback()

	// Step 1, get project keys and update
	row := tx.QueryRow("UPDATE projects SET keys = array_replace(keys, $1, $2) WHERE id = $3 RETURNING *", oldKey, newKey, projectID)
	keys := pq.StringArray{}
	result := model.Project{}
	err = row.Scan(&result.ID, &result.Name, &keys)
	if err != nil {
		return nil, -1, parseError(err)
	}

	result.Keys = make([]string, len(keys))
	for i, v := range keys {
		result.Keys[i] = v
	}

	// Step 2, find all project locales and update pairs
	rows, err := tx.Query("SELECT * FROM locales WHERE project_id = $1", projectID)
	if err != nil {
		return nil, -1, parseError(err)
	}
	defer rows.Close()

	locales := make([]model.Locale, 0)
	for rows.Next() {
		loc := model.Locale{}
		pairs := hstore.Hstore{}

		err := rows.Scan(&loc.ID, &loc.Ident, &loc.Language, &loc.Country, &pairs, &loc.ProjectID)
		if err != nil {
			return nil, -1, parseError(err)
		}

		loc.Pairs = make(map[string]string)
		for k, v := range pairs.Map {
			if v.Valid {
				loc.Pairs[k] = v.String
			}
		}

		locales = append(locales, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, -1, parseError(err)
	}

	for _, locale := range locales {
		h := hstore.Hstore{}
		h.Map = make(map[string]sql.NullString)
		for k, v := range locale.Pairs {
			// replace old key with new one when found, keep value
			if k == oldKey {
				k = newKey
			}
			h.Map[k] = sql.NullString{String: v, Valid: true}
		}

		values, err := h.Value()
		if err != nil {
			return nil, -1, err
		}

		_, err = tx.Exec("UPDATE locales SET pairs = $1 WHERE project_id = $2 AND id = $3", values, projectID, locale.ID)
		if err != nil {
			return nil, -1, parseError(err)
		}

	}

	tx.Commit()

	return &result, len(locales), nil
}

func (db *PostgresDB) DeleteProjectKey(projectID, key string) (*model.Project, error) {
	project, err := db.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	keys := make(pq.StringArray, 0)
	found := false
	for _, v := range project.Keys {
		if v == key {
			found = true
			continue
		}
		keys = append(keys, v)
	}

	if !found {
		return nil, parseError(errors.ErrNotFound)
	}

	values, err := keys.Value()
	if err != nil {
		return nil, parseError(err)
	}

	row := db.QueryRow("UPDATE projects SET keys = $1 WHERE id = $2 RETURNING *", values, projectID)
	err = row.Scan(&project.ID, &project.Name, &keys)
	if err != nil {
		return nil, parseError(err)
	}

	project.Keys = make([]string, len(keys))
	for i, v := range keys {
		project.Keys[i] = v
	}

	return project, nil
}

func (db *PostgresDB) UpdateProject(project model.Project) (*model.Project, error) {
	keys := make(pq.StringArray, len(project.Keys))
	for i, v := range project.Keys {
		keys[i] = v
	}

	values, err := keys.Value()
	if err != nil {
		return nil, parseError(err)
	}

	row := db.QueryRow("UPDATE projects SET keys = $1 WHERE id = $2 RETURNING *", values, project.ID)
	keys = pq.StringArray{}
	err = row.Scan(&project.ID, &project.Name, &keys)
	if err != nil {
		return nil, parseError(err)
	}

	project.Keys = make([]string, len(keys))
	for i, v := range keys {
		project.Keys[i] = v
	}

	return &project, nil
}

func (db *PostgresDB) DeleteProject(id string) error {
	_, err := db.Exec("DELETE FROM projects WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	return err
}

func (db *PostgresDB) GetProjectLocaleByIdent(projectID string, ident string) (*model.Locale, error) {
	loc := model.Locale{}
	row := db.QueryRow("SELECT * FROM locales WHERE project_id = $1 AND ident = $2", projectID, ident)
	pairs := hstore.Hstore{}
	err := row.Scan(&loc.ID, &loc.Ident, &loc.Language, &loc.Country, &pairs, &loc.ProjectID)
	if err != nil {
		return nil, parseError(err)
	}

	loc.Pairs = make(map[string]string)
	for k, v := range pairs.Map {
		if v.Valid {
			loc.Pairs[k] = v.String
		}
	}

	return &loc, nil
}

func (db *PostgresDB) GetProjectLocales(projID string, localeIdents ...string) ([]model.Locale, error) {
	rows, err := db.Query("SELECT * FROM locales WHERE project_id = $1", projID)
	if err != nil {
		return nil, parseError(err)
	}
	defer rows.Close()

	locs := make([]model.Locale, 0)
	for rows.Next() {
		loc := model.Locale{}
		pairs := hstore.Hstore{}

		err := rows.Scan(&loc.ID, &loc.Ident, &loc.Language, &loc.Country, &pairs, &loc.ProjectID)
		if err != nil {
			return nil, parseError(err)
		}

		loc.Pairs = make(map[string]string)
		for k, v := range pairs.Map {
			if v.Valid {
				loc.Pairs[k] = v.String
			}
		}

		locs = append(locs, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, parseError(err)
	}

	if len(localeIdents) > 0 {
		locs = filterLocales(locs, func(doc model.Locale) bool {
			for _, locIdent := range localeIdents {
				if doc.Ident == locIdent {
					return true
				}
			}
			return false
		})
	}

	return locs, nil
}

func filterLocales(locs []model.Locale, fn func(model.Locale) bool) []model.Locale {
	result := make([]model.Locale, 0)
	for _, v := range locs {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

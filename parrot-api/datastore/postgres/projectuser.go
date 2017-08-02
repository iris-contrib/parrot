package postgres

import (
	"github.com/iris-contrib/parrot/parrot-api/model"
	"github.com/lib/pq"
)

func (db *PostgresDB) GetUserProjects(userID string) ([]model.Project, error) {
	rows, err := db.Query(`SELECT projects.*
							FROM projects
							JOIN projects_users ON projects.id = projects_users.project_id
							WHERE projects_users.user_id = $1`, userID)
	if err != nil {
		return nil, parseError(err)
	}
	defer rows.Close()

	projects := make([]model.Project, 0)
	for rows.Next() {
		p := model.Project{}
		keys := pq.StringArray{}

		err := rows.Scan(&p.ID, &p.Name, &keys)
		if err != nil {
			return nil, parseError(err)
		}

		p.Keys = make([]string, len(keys))
		for i, v := range keys {
			p.Keys[i] = v
		}

		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, parseError(err)
	}

	return projects, nil
}

func (db *PostgresDB) GetProjectUsers(projID string) ([]model.ProjectUser, error) {
	rows, err := db.Query(`SELECT user_id, project_id, users.email, users.name, role
							FROM users
							JOIN projects_users ON users.id = projects_users.user_id
							WHERE projects_users.project_id = $1`, projID)
	if err != nil {
		return nil, parseError(err)
	}
	defer rows.Close()

	users := make([]model.ProjectUser, 0)
	for rows.Next() {
		u := model.ProjectUser{}

		err := rows.Scan(&u.UserID, &u.ProjectID, &u.Email, &u.Name, &u.Role)
		if err != nil {
			return nil, parseError(err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, parseError(err)
	}

	return users, nil
}

func (db *PostgresDB) GetUserProjectRoles(userID string) ([]model.ProjectUser, error) {
	rows, err := db.Query(`SELECT user_id, project_id, role
							FROM projects_users
							WHERE projects_users.user_id = $1`, userID)
	if err != nil {
		return nil, parseError(err)
	}
	defer rows.Close()

	roles := make([]model.ProjectUser, 0)
	for rows.Next() {
		u := model.ProjectUser{}

		err := rows.Scan(&u.UserID, &u.ProjectID, &u.Role)
		if err != nil {
			return nil, parseError(err)
		}
		roles = append(roles, u)
	}

	if err := rows.Err(); err != nil {
		return nil, parseError(err)
	}

	return roles, nil
}

func (db *PostgresDB) GetProjectUser(projID, userID string) (*model.ProjectUser, error) {
	u := model.ProjectUser{}
	row := db.QueryRow(`SELECT user_id, project_id, users.email, users.name, role
							FROM users
							JOIN projects_users ON users.id = projects_users.user_id
							WHERE projects_users.project_id = $1 AND user_id = $2`, projID, userID)
	err := row.Scan(&u.UserID, &u.ProjectID, &u.Email, &u.Name, &u.Role)
	if err != nil {
		return nil, parseError(err)
	}
	return &u, nil
}

func (db *PostgresDB) AssignProjectUser(pu model.ProjectUser) (*model.ProjectUser, error) {
	_, err := db.Exec("INSERT INTO projects_users (project_id, user_id, role) VALUES($1, $2, $3)",
		pu.ProjectID, pu.UserID, pu.Role)
	if err != nil {
		return nil, parseError(err)
	}
	return db.GetProjectUser(pu.ProjectID, pu.UserID)
}

func (db *PostgresDB) RevokeProjectUser(pu model.ProjectUser) error {
	_, err := db.Exec("DELETE FROM projects_users WHERE project_id = $1 AND user_id = $2",
		pu.ProjectID, pu.UserID)
	return parseError(err)
}

func (db *PostgresDB) UpdateProjectUser(pu model.ProjectUser) (*model.ProjectUser, error) {
	_, err := db.Exec("UPDATE projects_users SET role = $1 WHERE project_id = $2 AND user_id = $3",
		pu.Role, pu.ProjectID, pu.UserID)
	if err != nil {
		return nil, parseError(err)
	}
	return db.GetProjectUser(pu.ProjectID, pu.UserID)
}

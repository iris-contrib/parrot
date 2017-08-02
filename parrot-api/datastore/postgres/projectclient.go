package postgres

import "github.com/iris-contrib/parrot/parrot-api/model"

func (db *PostgresDB) GetProjectClients(projectID string) ([]model.ProjectClient, error) {
	rows, err := db.Query("SELECT client_id, project_id, name, secret FROM project_clients WHERE project_id = $1", projectID)
	if err != nil {
		return nil, parseError(err)
	}
	defer rows.Close()

	result := make([]model.ProjectClient, 0)
	for rows.Next() {
		r := model.ProjectClient{}
		err = rows.Scan(&r.ClientID, &r.ProjectID, &r.Name, &r.Secret)
		if err != nil {
			return nil, parseError(err)
		}
		result = append(result, r)
	}
	return result, nil
}

func (db *PostgresDB) FindOneClient(clientID string) (*model.ProjectClient, error) {
	row := db.QueryRow("SELECT client_id, project_id, name, secret FROM project_clients WHERE client_id = $1", clientID)
	result := model.ProjectClient{}
	err := row.Scan(&result.ClientID, &result.ProjectID, &result.Name, &result.Secret)
	if err != nil {
		return nil, parseError(err)
	}
	return &result, nil
}

func (db *PostgresDB) GetProjectClient(projectID, clientID string) (*model.ProjectClient, error) {
	row := db.QueryRow("SELECT client_id, project_id, name, secret FROM project_clients WHERE project_id = $1 AND client_id = $2",
		projectID, clientID)
	result := model.ProjectClient{}
	err := row.Scan(&result.ClientID, &result.ProjectID, &result.Name, &result.Secret)
	if err != nil {
		return nil, parseError(err)
	}
	return &result, nil
}

func (db *PostgresDB) CreateProjectClient(pc model.ProjectClient) (*model.ProjectClient, error) {
	row := db.QueryRow("INSERT INTO project_clients (project_id, name, secret) VALUES($1, $2, $3) RETURNING client_id, project_id, name, secret",
		pc.ProjectID, pc.Name, pc.Secret)
	result := model.ProjectClient{}
	err := row.Scan(&result.ClientID, &result.ProjectID, &result.Name, &result.Secret)
	if err != nil {
		return nil, parseError(err)
	}
	return &result, nil
}

func (db *PostgresDB) DeleteProjectClient(projectID, clientID string) error {
	_, err := db.Exec("DELETE FROM project_clients WHERE project_id = $1 AND client_id = $2", projectID, clientID)
	return parseError(err)
}

func (db *PostgresDB) UpdateProjectClientSecret(pc model.ProjectClient) (*model.ProjectClient, error) {
	_, err := db.Exec("UPDATE project_clients SET secret = $1 WHERE project_id = $2 AND client_id = $3",
		pc.Secret, pc.ProjectID, pc.ClientID)
	if err != nil {
		return nil, parseError(err)
	}
	return db.GetProjectClient(pc.ProjectID, pc.ClientID)
}

func (db *PostgresDB) UpdateProjectClientName(pc model.ProjectClient) (*model.ProjectClient, error) {
	_, err := db.Exec("UPDATE project_clients SET name = $1 WHERE project_id = $2 AND client_id = $3",
		pc.Name, pc.ProjectID, pc.ClientID)
	if err != nil {
		return nil, parseError(err)
	}
	return db.GetProjectClient(pc.ProjectID, pc.ClientID)
}

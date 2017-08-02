package postgres

import "github.com/iris-contrib/parrot/parrot-api/model"

func (db *PostgresDB) GetUserByEmail(email string) (*model.User, error) {
	u := model.User{}
	row := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", email)

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		return nil, parseError(err)
	}

	return &u, nil
}

func (db *PostgresDB) GetUserByID(id string) (*model.User, error) {
	u := model.User{}
	row := db.QueryRow("SELECT id, name, email, password FROM users WHERE id = $1", id)

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		return nil, parseError(err)
	}

	return &u, nil
}

func (db *PostgresDB) CreateUser(u model.User) (*model.User, error) {
	row := db.QueryRow("INSERT INTO users (name, email, password) VALUES($1, $2, $3) RETURNING id, name, email", u.Name, u.Email, u.Password)
	err := row.Scan(&u.ID, &u.Name, &u.Email)
	return &u, parseError(err)
}

func (db *PostgresDB) UpdateUserPassword(u model.User) (*model.User, error) {
	row := db.QueryRow("UPDATE users SET password = $1 WHERE id = $2 RETURNING id, name, email", u.Password, u.ID)
	err := row.Scan(&u.ID, &u.Name, &u.Email)
	return &u, parseError(err)
}

func (db *PostgresDB) UpdateUserName(u model.User) (*model.User, error) {
	row := db.QueryRow("UPDATE users SET name = $1 WHERE id = $2 RETURNING id, name, email", u.Name, u.ID)
	err := row.Scan(&u.ID, &u.Name, &u.Email)
	return &u, parseError(err)
}

func (db *PostgresDB) UpdateUserEmail(u model.User) (*model.User, error) {
	row := db.QueryRow("UPDATE users SET email = $1 WHERE id = $2 RETURNING id, name, email", u.Email, u.ID)
	err := row.Scan(&u.ID, &u.Name, &u.Email)
	return &u, parseError(err)
}

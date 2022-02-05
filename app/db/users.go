package db

import (
	"github.com/skymazer/user_service/models"
)

func (db Database) AddUser(item *models.User) error {
	query := `INSERT INTO users (name, mail) VALUES ($1, $2)`
	_, err := db.Conn.Exec(query, item.Name, item.Mail)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) DeleteUser(userId models.IdType) error {
	query := `DELETE FROM users WHERE id = $1;`
	res, err := db.Conn.Exec(query, userId)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNoMatch
	}

	return nil
}

func (db Database) GetAllUsers() ([]*models.User, error) {
	var res []*models.User

	rows, err := db.Conn.Query("SELECT * FROM users ORDER BY ID DESC")
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.Id, &u.Name, &u.Mail)
		if err != nil {
			return res, err
		}
		res = append(res, &u)
	}
	return res, nil
}

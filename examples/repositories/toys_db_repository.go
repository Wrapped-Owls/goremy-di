package repositories

import "database/sql"

type ToysDbRepository struct {
	conn *sql.DB
}

func NewToysDbRepository(conn *sql.DB) *ToysDbRepository {
	return &ToysDbRepository{conn: conn}
}

func (repo ToysDbRepository) Save(name string) error {
	const query = "INSERT INTO toys(name) VALUES(?)"
	_, err := repo.conn.Exec(query, name)
	return err
}

func (repo ToysDbRepository) ListAll() ([]string, error) {
	const query = "SELECT name FROM toys"
	rows, err := repo.conn.Query(query)
	if err != nil {
		return nil, err
	}

	results := make([]string, 0, 11)
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		results = append(results, name)
	}
	return results, nil
}

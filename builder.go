package main

import "database/sql"

func FindOne[M any](db *sql.DB, query string, params []any) (*M, error) {
	var model M

	if err := db.QueryRow(query, params...).Scan(&model); err != nil {
		return nil, err
	}

	return &model, nil
}

func FindMany[M any](db *sql.DB, query string, params []any) ([]M, error) {
	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []M

	for rows.Next() {
		var model M
		if err := rows.Scan(&model); err != nil {
			return nil, err
		}

		models = append(models, model)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}

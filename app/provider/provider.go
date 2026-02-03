package provider

import (
	"errors"
	"api_kino/app/models/staff/model"
	"api_kino/config/constant"
	"api_kino/config/database"
	"strconv"

	"gorm.io/gorm"
)

func GetUser(db *gorm.DB, userId string) (*model.Staff, error) {
	var user model.Staff
	if rs := db.
		Where(`"staff".Kode = ?`, userId).
		Find(&user); rs.RowsAffected < 1 {
		return nil, errors.New(constant.ErrorLogin)
	}
	return &user, nil
}

func QueryRun(query string, values ...interface{}) ([]map[string]interface{}, error) {
	db := database.DB
	q := db
	rows, err := q.Raw(query, values...).Rows()
	if err != nil {
		return nil, err
	}
	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		return nil, err
	}
	colNum := len(columns)
	var results []map[string]interface{}
	for rows.Next() {
		r := make([]interface{}, colNum)
		for i := range r {
			r[i] = &r[i]
		}
		err = rows.Scan(r...)
		if err != nil {
			return nil, err
		}
		var row = map[string]interface{}{}
		for i := range r {
			switch v := r[i].(type) {
			case []uint8:
				f, _ := strconv.ParseFloat(string(v), 64)
				row[columns[i]] = f
			default:
				row[columns[i]] = v
			}
		}
		results = append(results, row)
	}
	return results, nil
}

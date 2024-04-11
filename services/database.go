package services

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/mrzack99s/cocong/vars"
	"gorm.io/gorm"
)

type DBQueryPreload struct {
	Name string
	Args []any
}

func DBQuery(response any, offset, limit int, search string, or, ignoreOffset bool, preloads ...DBQueryPreload) (count int64, err error) {

	if reflect.TypeOf(response).Kind() != reflect.Pointer {
		err = errors.New("response must be pointer")
		return
	}

	if reflect.TypeOf(response).Elem().Kind() != reflect.Slice {
		err = errors.New("response must be pointer of slice")
		return
	}

	modelVal := reflect.New(reflect.TypeOf(response).Elem().Elem())

	query := []string{}
	if search != "" {
		query = strings.Split(search, "|")
	}

	if len(query) > 0 {
		preparedQuery := [][]string{}
		for _, q := range query {
			q = strings.TrimSpace(q)

			op := ""
			tmp := []string{}
			if strings.Contains(q, " = ") {
				tmp = strings.Split(q, " = ")
				op = "="
			} else if strings.Contains(q, " > ") {
				tmp = strings.Split(q, " > ")
				op = ">"
			} else if strings.Contains(q, " >= ") {
				tmp = strings.Split(q, " >= ")
				op = ">="
			} else if strings.Contains(q, " < ") {
				tmp = strings.Split(q, " < ")
				op = "<"
			} else if strings.Contains(q, " <= ") {
				tmp = strings.Split(q, " <= ")
				op = "<="
			} else if strings.Contains(q, " <> ") {
				tmp = strings.Split(q, " <> ")
				op = "<>"
			} else if strings.Contains(q, " like ") || strings.Contains(q, " LIKE ") {

				if strings.Contains(q, " like ") {
					tmp = strings.Split(q, " like ")
				} else {
					tmp = strings.Split(q, " LIKE ")
				}
				op = "LIKE"
			}

			if op != "" {
				tmp = append(tmp[:2], tmp[1:]...)
				tmp[1] = op
				preparedQuery = append(preparedQuery, tmp)
			}

		}

		whereStmt := ""
		val := []any{}
		for i, q := range preparedQuery {
			if i < len(preparedQuery)-1 {
				if or {
					whereStmt += fmt.Sprintf("%s %s ? or ", q[0], q[1])
				} else {
					whereStmt += fmt.Sprintf("%s %s ? and ", q[0], q[1])
				}
			} else {
				whereStmt += fmt.Sprintf("%s %s ?", q[0], q[1])
			}

			if q[1] == "LIKE" {
				val = append(val, fmt.Sprintf("%%%s%%", q[2]))
			} else {
				val = append(val, q[2])
			}
		}

		vars.Database.Select("count(id)").Model(modelVal.Interface()).Where(whereStmt, val...).Scan(&count)

		tx := vars.Database.Where(whereStmt, val...)
		if !ignoreOffset {
			tx = tx.Offset(offset).Limit(limit)
		}

		if tx.Statement.Preloads == nil {
			tx.Statement.Preloads = map[string][]interface{}{}
		}

		for _, preload := range preloads {
			tx.Statement.Preloads[preload.Name] = preload.Args
		}

		tx = tx.Order("updated_at desc")

		err = tx.Find(response).Error
		if err != nil {
			return
		}

	} else {

		vars.Database.Select("count(id)").Model(modelVal.Interface()).Scan(&count)
		tx := vars.Database
		if !ignoreOffset {
			tx = tx.Offset(offset).Limit(limit)
		}

		if tx.Statement.Preloads == nil {
			tx.Statement.Preloads = map[string][]interface{}{}
		}

		for _, preload := range preloads {
			tx.Statement.Preloads[preload.Name] = preload.Args
		}

		tx = tx.Order("updated_at desc")
		err = tx.Find(response).Error
		if err != nil {
			return
		}

	}

	return
}

func DBQueryCustomDB(db *gorm.DB, response any, offset, limit int, search string, or, ignoreOffset bool, preloads ...DBQueryPreload) (count int64, err error) {

	if reflect.TypeOf(response).Kind() != reflect.Pointer {
		err = errors.New("response must be pointer")
		return
	}

	if reflect.TypeOf(response).Elem().Kind() != reflect.Slice {
		err = errors.New("response must be pointer of slice")
		return
	}

	modelVal := reflect.New(reflect.TypeOf(response).Elem().Elem())

	query := []string{}
	if search != "" {
		query = strings.Split(search, "|")
	}

	if len(query) > 0 {
		preparedQuery := [][]string{}
		for _, q := range query {
			q = strings.TrimSpace(q)

			op := ""
			tmp := []string{}
			if strings.Contains(q, " = ") {
				tmp = strings.Split(q, " = ")
				op = "="
			} else if strings.Contains(q, " > ") {
				tmp = strings.Split(q, " > ")
				op = ">"
			} else if strings.Contains(q, " >= ") {
				tmp = strings.Split(q, " >= ")
				op = ">="
			} else if strings.Contains(q, " < ") {
				tmp = strings.Split(q, " < ")
				op = "<"
			} else if strings.Contains(q, " <= ") {
				tmp = strings.Split(q, " <= ")
				op = "<="
			} else if strings.Contains(q, " <> ") {
				tmp = strings.Split(q, " <> ")
				op = "<>"
			} else if strings.Contains(q, " like ") || strings.Contains(q, " LIKE ") {

				if strings.Contains(q, " like ") {
					tmp = strings.Split(q, " like ")
				} else {
					tmp = strings.Split(q, " LIKE ")
				}
				op = "LIKE"
			}

			if op != "" {
				tmp = append(tmp[:2], tmp[1:]...)
				tmp[1] = op
				preparedQuery = append(preparedQuery, tmp)
			}

		}

		whereStmt := ""
		val := []any{}
		for i, q := range preparedQuery {
			if i < len(preparedQuery)-1 {
				if or {
					whereStmt += fmt.Sprintf("%s %s ? or ", q[0], q[1])
				} else {
					whereStmt += fmt.Sprintf("%s %s ? and ", q[0], q[1])
				}
			} else {
				whereStmt += fmt.Sprintf("%s %s ?", q[0], q[1])
			}

			if q[1] == "LIKE" {
				val = append(val, fmt.Sprintf("%%%s%%", q[2]))
			} else {
				val = append(val, q[2])
			}
		}

		db.Select("count(id)").Model(modelVal.Interface()).Where(whereStmt, val...).Scan(&count)

		tx := db.Where(whereStmt, val...)
		if !ignoreOffset {
			tx = tx.Offset(offset).Limit(limit)
		}

		if tx.Statement.Preloads == nil {
			tx.Statement.Preloads = map[string][]interface{}{}
		}

		for _, preload := range preloads {
			tx.Statement.Preloads[preload.Name] = preload.Args
		}

		tx = tx.Order("updated_at desc")
		err = tx.Find(response).Error
		if err != nil {
			return
		}

	} else {

		db.Select("count(id)").Model(modelVal.Interface()).Scan(&count)
		tx := db
		if !ignoreOffset {
			tx = tx.Offset(offset).Limit(limit)
		}

		if tx.Statement.Preloads == nil {
			tx.Statement.Preloads = map[string][]interface{}{}
		}

		for _, preload := range preloads {
			tx.Statement.Preloads[preload.Name] = preload.Args
		}

		tx = tx.Order("updated_at desc")
		err = tx.Find(response).Error
		if err != nil {
			return
		}
	}

	return
}

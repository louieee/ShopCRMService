package repositories

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func BuildWhereClause(filter interface{}) string {
	var conditions []string

	v := reflect.ValueOf(filter)
	t := reflect.TypeOf(filter)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip fields with zero values
		if reflect.DeepEqual(value.Interface(), reflect.Zero(field.Type).Interface()) {
			continue
		}

		// Handle different field types as needed
		switch value.Interface().(type) {
		case []uint:
			conditions = append(conditions, fmt.Sprintf("%s IN (%s)", field.Name, sliceToString(value.Interface())))
		case time.Time:
			conditions = append(conditions, fmt.Sprintf("%s >= '%s'", field.Name, value.Interface()))
		case bool:
			if value.Interface().(bool) {
				conditions = append(conditions, fmt.Sprintf("%s == true", field.Name))
			}
		}
	}

	// Combine conditions using "AND"
	whereClause := strings.Join(conditions, " AND ")

	// Optionally, you can wrap the conditions in a WHERE clause
	if whereClause != "" {
		whereClause = "WHERE " + whereClause
	}

	return whereClause
}

func sliceToString(slice interface{}) string {
	var strSlice []string
	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		strSlice = append(strSlice, fmt.Sprint(s.Index(i).Interface()))
	}

	return strings.Join(strSlice, ", ")
}

func makeSqlQuery(query string, column string, operator string) string {
	if strings.Contains(query, "where") {
		query = fmt.Sprintf("%s and", query)
	} else {
		query = fmt.Sprintf("%s where", query)
	}
	if operator == "in" {
		return fmt.Sprintf("%s %s %s (?)", query, column, operator)
	}
	return fmt.Sprintf("%s %s %s ?", query, column, operator)

}
func makeOrSqlQuery(query string, column string, operator string) string {
	if strings.Contains(query, "where") {
		query = fmt.Sprintf("%s or", query)
	} else {
		query = fmt.Sprintf("%s where", query)
	}
	if operator == "in" {
		return fmt.Sprintf("%s %s %s (?)", query, column, operator)
	}
	return fmt.Sprintf("%s %s %s ?", query, column, operator)

}

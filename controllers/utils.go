package controllers

import (
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
	"time"
)

type DBController struct {
	DB *gorm.DB
}

func NewController(db *gorm.DB) *DBController {
	return &DBController{DB: db}
}

// # Trying generics
//
// ConvertToUint prints slices using the [fmt.Println] function.
// The current implementation prints the following slices:
//   - []int{}
//   - []string{}
//
// For more information about Go doc comments, see [Go Doc Comments] at tip.golang.org.
//
// [Go Doc Comments]: https://tip.golang.org/doc/comment
func ConvertStringToUintPtr(s string) *uint {
	result, err := strconv.Atoi(s)
	if err != nil || s == "" {
		return nil
	}
	resultUint := uint(result)
	return &resultUint
}

func ConvertStringSliceToIntSlice(strSlice []string) []int {
	intSlice := make([]int, 0, len(strSlice))
	for _, str := range strSlice {
		intValue, err := strconv.Atoi(str)
		if err != nil {
			return make([]int, 0)
		}
		intSlice = append(intSlice, intValue)
	}
	return intSlice
}

func ConvertStringSliceToUintSlice(strSlice []string) []uint {
	uintSlice := make([]uint, 0, len(strSlice))
	for _, str := range strSlice {
		intValue, err := strconv.Atoi(str)
		if err != nil {
			return make([]uint, 0, len(strSlice))
		}
		uintSlice = append(uintSlice, uint(intValue))
	}
	return uintSlice
}

func ConvertStringToTime(string2 string, timezone bool) *time.Time {
	var layout string
	if timezone == true {
		layout = "2006-01-02T15:04:05Z"
	} else {
		layout = "2006-01-02 15:04:05"
	}
	t, err := time.Parse(layout, string2)
	if err != nil {
		return nil
	}
	return &t
}

func ConvertStringToBool(string2 string) *bool {
	var value bool
	if string2 == "true" {
		value = true
		return &value
	} else if string2 == "false" {
		value = false
		return &value
	} else {
		return nil
	}
}

func ConvertStringToSlice(string2 string) []string {
	if string2 == "" {
		return make([]string, 0, 0)
	}
	return strings.Split(string2, ",")
}

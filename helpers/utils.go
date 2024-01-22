package helpers

import "reflect"

func ContainsItem(list interface{}, item interface{}) bool {
	v := reflect.ValueOf(list)
	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(v.Index(i).Interface(), item) {
			return true
		}
	}
	return false
}

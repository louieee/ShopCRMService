package helpers

import ("reflect"
"encoding/json"
"fmt"
)
func ContainsItem(list interface{}, item interface{}) bool {
	v := reflect.ValueOf(list)
	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(v.Index(i).Interface(), item) {
			return true
		}
	}
	return false
}


func StructToString(obj interface{}) *string {
	// Convert the struct to JSON
    jsonData, err := json.Marshal(obj)
    if err != nil {
        fmt.Println("Error converting to JSON:", err)
        return nil
    }

    jsonString := string(jsonData)
	return &jsonString
}

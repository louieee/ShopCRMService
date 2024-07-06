package consumers

import (
	"encoding/json"
	"fmt"
)

type LeadPayload struct {
	Name    string `json:"name"`
	Company string `json:"company"`
}

func HandleLead(action string, data string) {
	switch action {
	case "create":
		handleNewLead(data)
	case "update":
		handleUpdateLead(data)
	case "delete":
		handleDeleteLead(data)

	}
}

func handleNewLead(data string) {
	var lead LeadPayload
	err1 := json.Unmarshal([]byte(data), &lead)
	if err1 != nil {
		panic(fmt.Sprintf("Error unwrapping message: %s", err1.Error()))
	}
	println("new lead: ", lead.Company, lead.Name)
}

func handleDeleteLead(data string) {

}

func handleUpdateLead(data string) {

}

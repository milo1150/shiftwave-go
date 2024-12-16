package utils

import (
	"encoding/json"
	"fmt"
)

// Function to pretty-print the slice of assessments
func PrettyPrint(data interface{}, msg string) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error pretty printing:", err)
		return
	}
	fmt.Printf("%v %v\n", msg, string(jsonData))
}

package tools

import "strings"

// func GetExcelElement(inputSlice []string) string {
// 	result := ""
// 	for i, element := range inputSlice {
// 		// Split the element by ", " and join them with "\t"
// 		parts := strings.Split(element, ", ")
// 		joinedElement := strings.Join(parts, "\t")

// 		// Append the joined element to the result string
// 		result += joinedElement

// 		// Add a newline character to separate rows except for the last element
// 		if i < len(inputSlice)-1 {
// 			result += "\r\n" 
// 		}
// 	}
// 	return result
// }

func GetExcelElement(inputSlice []string) string {
	result := ""
	for i, element := range inputSlice {
		// Split the element by ", " and join them with "\t"
		parts := strings.Split(element, ", ")
		joinedElement := strings.Join(parts, "\t")

		// Append the joined element to the result string
		result += joinedElement

		// Add a newline character to separate rows except for the last element
		if i < len(inputSlice)-1 {
			result += "\r\n" // Use "\r\n" for a Windows-style newline
		}
	}
	return result
}

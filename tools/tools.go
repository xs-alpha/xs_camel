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
	// 创建一个二维字符串切片来存储每个单元格的内容
	table := make([][]string, len(inputSlice))
	// 遍历输入的切片
	for i, element := range inputSlice {
		// 分割单元格内容并存储到每一行的切片中
		cells := strings.Split(element, ", ")
		table[i] = cells
	}

	// 生成表格的字符串表示形式
	result := ""
	for _, row := range table {
		// 将每行的单元格内容用制表符连接起来
		joinedRow := strings.Join(row, "\t")
		// 将每一行的内容添加到结果字符串中
		result += joinedRow + "\n"
	}
	return result
}

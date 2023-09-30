package tools

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	qrcodeReader "github.com/tuotoo/qrcode"
	"github.com/xwb1989/sqlparser"
)

func GetExcelElement(inputSlice []string) string {
	result := ""
	for i, element := range inputSlice {
		// Split the element by ", " and join them with "\t"
		element = strings.Split(element, ".")[1]
		parts := strings.Split(element, ", ")
		joinedElement := strings.Join(parts, "\t")

		// Append the joined element to the result string
		result += joinedElement

		// Add a newline character to separate rows except for the last element
		if i < len(inputSlice)-1 {
			result += "\r\n"
		}
	}
	return result
}

func StmtToGo(stmt *sqlparser.DDL, tableName string, pkgName string) ([]string, error) {
	builder := strings.Builder{}

	// header := fmt.Sprintf("package %s\n", pkgName)

	structName := snakeCaseToCamel(tableName)
	structStart := fmt.Sprintf("type %s struct { \n", structName)
	builder.WriteString(structStart)
	ret := make([]string, 0)
	for _, col := range stmt.TableSpec.Columns {
		columnType := col.Type.Type

		goType := sqlTypeMap[columnType]

		field := snakeCaseToCamel(col.Name.String())
		retStr := field + ", " + goType
		comment := col.Type.Comment
		if comment == nil {
			builder.WriteString(fmt.Sprintf("\t%s\t%s\t\n", field, goType))
			retStr = "-" + ", " + retStr
		} else {
			builder.WriteString(fmt.Sprintf("\t%s\t%s\t`comment:\"%s\"` \n",
				field, goType, string(comment.Val)))
			retStr = string(comment.Val) + ", " + retStr
		}
		if retStr == "" {
			continue
		}
		ret = append(ret, retStr)
	}
	builder.WriteString("}\n")

	return ret, nil
}

// In sql, table name often is snake_case
// In Go, struct name often is camel case
func snakeCaseToCamel(str string) string {
	builder := strings.Builder{}
	index := 0
	if str[0] >= 'a' && str[0] <= 'z' {
		builder.WriteByte(str[0] - ('a' - 'A'))
		index = 1
	}
	for i := index; i < len(str); i++ {
		if str[i] == '_' && i+1 < len(str) {
			if str[i+1] >= 'a' && str[i+1] <= 'z' {
				builder.WriteByte(str[i+1] - ('a' - 'A'))
				i++
				continue
			}
		}
		builder.WriteByte(str[i])
	}
	return builder.String()
}

/**
GetTimestamp 根据传入的时间字符串和格式获取时间戳，如果未传入则返回当前时间戳。
*/
func GetTimestamp(timeStr, layout string) (int64, error) {
	var parsedTime time.Time
	var err error

	if timeStr != "" && layout != "" {
		log.Println("GetTimestamp-按时间格式")
		// 解析时间字符串
		parsedTime, err = time.Parse(layout, timeStr)
		if err != nil {
			return 0, err
		}
	} else {
		// 没有传入时间字符串和格式，使用当前时间
		log.Println("GetTimestamp-wu时间格式")
		parsedTime = time.Now()
	}

	// 获取时间戳
	timestamp := parsedTime.Unix()
	return timestamp, nil
}

// IsTimeFormat 判断字符串是否是时间格式
func IsTimeFormat(str, layout string) bool {
	_, err := time.Parse(layout, str)
	return err == nil
}

func IsImg(suffix string) bool {
	suffixs := []string{".png", ".jpeg", ".webp", ".jpg"}
	for _, v := range suffixs {
		if v == strings.ToLower(suffix) {
			return true
		}
	}
	return false
}

func ReadQRCode(filename string) (content string) {
	fi, err := os.Open(filename)
	if err != nil {
		log.Println("readQrcode" + err.Error())
		return
	}
	defer fi.Close()
	qrmatrix, err := qrcodeReader.Decode(fi)
	if err != nil {
		log.Println("readQrcode:" + err.Error())
		return
	}
	return qrmatrix.Content
}

func GenerateRandomPassword(length int) (string, error) {
	// 可用于生成密码的字符集
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?"
	charsetLength := big.NewInt(int64(len(charset)))

	// 创建一个随机密码的切片
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		// 生成一个随机的索引，用于选择字符集中的字符
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password), nil
}

func IsNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
func PrettyPrintJSON(data interface{}) (string, error) {
	// MarshalIndent 函数用于将数据转换为美化的 JSON 字符串
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}

func CountValidWords(text string) int {
	count := 0

	// 遍历文本的每个字符
	for _, char := range text {
		if !unicode.IsSpace(char) {
			count++
		}
	}

	return count
}

func WriteInFile(imageReader io.Reader) (bool, string) {
	// 获取当前时间戳并格式化为字符串
	timestamp := time.Now().Unix()
	timestampStr := fmt.Sprintf("%d", timestamp)

	// 创建文件名，将时间戳作为文件名的一部分，加上文件后缀（例如 ".png"）
	fileName := timestampStr + ".png"
	// 创建一个文件用于保存图片
	outputFile, err := os.Create(fileName)
	if err != nil {
		log.Println("Error creating output file:", err)
		return false, ""
	}
	defer outputFile.Close()

	// 将 imageReader 重新定位到开头
	if _, err := imageReader.(*bytes.Reader).Seek(0, 0); err != nil {
		log.Println("Error seeking imageReader:", err)
		return false, ""
	}

	// 使用 io.Copy 将数据从 io.Reader 复制到文件
	_, copyErr := io.Copy(outputFile, imageReader)
	if copyErr != nil {
		log.Println("Error copying data to file:", copyErr)
		return false, ""
	}

	log.Println("Image saved as " + fileName)
	return true, fileName
}

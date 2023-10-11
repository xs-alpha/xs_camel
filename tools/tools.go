package tools

import (
	"bytes"
	"crypto/md5"
	rand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	mrand "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
	"xiaosheng/settings"

	"fyne.io/fyne/v2/widget"
	"gopkg.in/ini.v1"

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

	structName := ToCamelCase(tableName)
	structStart := fmt.Sprintf("type %s struct { \n", structName)
	builder.WriteString(structStart)
	ret := make([]string, 0)
	for _, col := range stmt.TableSpec.Columns {
		columnType := col.Type.Type

		goType := sqlTypeMap[columnType]

		field := ToCamelCase(col.Name.String())
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

// SnakeCaseToCamel 驼峰转换
func SnakeCaseToCamel(str string) string {
	builder := strings.Builder{}
	shouldCapitalize := false

	for i := 0; i < len(str); i++ {
		if str[i] == '_' && i+1 < len(str) {
			shouldCapitalize = true
			continue
		}

		if shouldCapitalize {
			builder.WriteByte(str[i] - ('a' - 'A'))
			shouldCapitalize = false
		} else {
			builder.WriteByte(str[i])
		}
	}
	return builder.String()
}

// ToCamelCase 驼峰转换
func ToCamelCase(text string) string {
	// 根据 isBigCamel 参数决定生成的是大驼峰还是小驼峰
	parts := strings.Split(text, "_")
	camelCase := ""
	for i, part := range parts {
		if i == 0 {
			if IsBigCamel {
				camelCase += strings.Title(part)
			} else {
				camelCase += strings.ToLower(part[:1]) + part[1:]
			}
		} else {
			camelCase += strings.Title(part)
		}
	}
	return camelCase
}

// GetTimestamp 根据传入的时间字符串和格式获取时间戳，如果未传入则返回当前时间戳。
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

// IsImg 是否是图片
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

// GenerateRandomPassword 生成随机密码
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

// IsNumeric 字符串是否是数字
func IsNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

// PrettyPrintJSON json美化
func PrettyPrintJSON(data interface{}) (string, error) {
	// MarshalIndent 函数用于将数据转换为美化的 JSON 字符串
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}

// CountValidWords 计算词数
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

// MonitorCase 转大写
func MonitorCase(ticker *time.Ticker, wetres *widget.Entry) {
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if !IsLowerCase {
				wetres.SetText(strings.ToLower(wetres.Text))
			} else {
				wetres.SetText(strings.ToUpper(wetres.Text))
			}
		case <-ToolsChan:
			log.Println("关闭大小写chan")
			return
		}
	}
}

// IsFileNameExists 某个文件是否存在
func IsFileNameExists(fileName string) bool {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		return false
	}

	// 拼接文件路径
	exePath := currentDir + "\\" + fileName

	// 检查文件是否存在
	_, err = os.Stat(exePath)
	if os.IsNotExist(err) {
		log.Println(fileName + " does not exist in the current directory.")
		return false
	} else if err != nil {
		log.Println("Error checking :", err)
		return false
	} else {
		fmt.Println(fileName + " exists in the current directory.")
		return true
	}
}

// DownloadFile 下载文件并保存到本地
func DownloadFile(url, fileName string) bool {
	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP request failed with status code %d\n", resp.StatusCode)
		return false
	}

	// 创建本地文件用于保存下载的内容
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Error creating file:", err)
		return false
	}
	defer file.Close()

	// 将下载的内容保存到本地文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println("Error saving file:", err)
		return false
	}

	log.Println("File downloaded successfully.")
	return true
}

// CalculateMD5 计算md5
func CalculateMD5(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建一个 MD5 哈希对象
	hash := md5.New()

	// 从文件中读取数据并计算哈希值
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// 将哈希值转换为十六进制字符串
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}

// GetRandomString 随机选择
func GetRandomString(lines []string) (string, error) {
	// 使用当前时间的纳秒级时间戳作为随机种子
	seed := time.Now().UnixNano()

	// 设置随机种子
	mrand.Seed(seed)

	// 生成随机索引
	randomIndex := mrand.Intn(len(lines))
	randomLine := lines[randomIndex]

	return randomLine, nil
}

func CreateConfigFile(filename string) (*ini.File, error) {
    cfg := ini.Empty()
    if err := cfg.SaveTo(filename); err != nil {
        log.Printf("无法创建INI文件: %v\n", err)
        return nil, err
    }
    return cfg, nil
}

// WriteConfig 写入配置文件
func WriteConfig(appID, secret string) error {
    filename := settings.Conf.ConfigName // 硬编码的文件名
	CreateConfigFile(filename)

    // 创建一个新的INI文件
    cfg, err := ini.Load(filename)
    if err != nil {
        log.Printf("无法创建INI文件: %v\n", err)
        return err
    }

    // 创建一个新的配置部分（section）
    section, err := cfg.NewSection("app")
    if err != nil {
        log.Printf("无法创建配置部分: %v\n", err)
        return err
    }

    // 设置配置项（key-value）
    section.NewKey("appId", appID)
    section.NewKey("secret", secret)

    // 保存到文件
    if err := cfg.SaveTo(filename); err != nil {
        log.Printf("无法保存INI文件: %v\n", err)
        return err
    }

    log.Printf("配置已写入 %s 文件\n", filename)
    return nil
}

func FileExists(filename string) bool {
    _, err := os.Stat(filename)
    return !os.IsNotExist(err)
}

func ReadConfig() (string, string, error) {
    filename := settings.Conf.ConfigName // 硬编码的文件名

    // 加载INI文件
    cfg, err := ini.Load(filename)
    if err != nil {
        log.Printf("无法加载INI文件: %v\n", err)
        return "", "", err
    }

    // 获取配置部分
    section := cfg.Section("app")
    if section == nil {
        log.Println("配置部分不存在")
        return "", "", nil
    }

    // 读取配置项
    appID := section.Key("appId").String()
    secret := section.Key("secret").String()

    return appID, secret, nil
}

func DeleteFile(filename string) error {
    err := os.Remove(filename)
    if err != nil {
        log.Printf("无法删除文件 %s: %v\n", filename, err)
        return err
    }
    log.Printf("%s 文件已被成功删除\n", filename)
    return nil
}
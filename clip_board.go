package main

import (
    "bytes"
    "io"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "os"
)

func main2() {
    // 输入解析的图片，返回解析好的数据
    url := qrcodeDecode(`qr1.png`)
    println(url)
}

// 二维码解码
// 参数：要解析的图片
func qrcodeDecode(img string) string {
    fh, err := os.Open(img)
    panicErr(err)
    defer fh.Close()

    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)
    bodyWriter.WriteField(`Filename`, fh.Name())
    bodyWriter.WriteField(`Upload`, `Submit Query`)

    fileWriter, err := bodyWriter.CreateFormFile(`Filedata`, img)
    panicErr(err)

    _, err = io.Copy(fileWriter, fh)
    panicErr(err)

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

    resp, err := http.Post(`http://tool.chinaz.com/ajaxseo.aspx?t=pload`, contentType, bodyBuf)
    resp_body, err := ioutil.ReadAll(resp.Body)
    panicErr(err)
    defer resp.Body.Close()

    str := string(resp_body)
    str_len := len(str)
    return str[35 : str_len-7]
}

// 统一处理错误函数
func panicErr(err error) {
    if err != nil {
        panic(err)
    }
}
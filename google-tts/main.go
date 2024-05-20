package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	_ "os"
	"strings"
)

var headers = map[string]string{
	"accept":                      "*/*",
	"accept-language":             "en-US,en;q=0.9",
	"priority":                    "u=1, i",
	"referer":                     "https://www.google.com/",
	"sec-ch-ua":                   "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\"",
	"sec-ch-ua-arch":              "\"x86\"",
	"sec-ch-ua-bitness":           "\"64\"",
	"sec-ch-ua-full-version":      "\"124.0.6367.208\"",
	"sec-ch-ua-full-version-list": "\"Chromium\";v=\"124.0.6367.208\", \"Google Chrome\";v=\"124.0.6367.208\", \"Not-A.Brand\";v=\"99.0.0.0\"",
	"sec-ch-ua-mobile":            "?0",
	"sec-ch-ua-model":             "\"\"",
	"sec-ch-ua-platform":          "\"Windows\"",
	"sec-ch-ua-platform-version":  "\"15.0.0\"",
	"sec-ch-ua-wow64":             "?0",
	"sec-fetch-dest":              "empty",
	"sec-fetch-mode":              "cors",
	"sec-fetch-site":              "same-origin",
	"user-agent":                  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"x-dos-behavior":              "Embed",
}

func writeSpeech(text, language, outputfile string) bool {
	if !strings.HasSuffix(strings.ToLower(outputfile), ".mp3") {
		outputfile += ".mp3"
	}

	text = url.QueryEscape(strings.ReplaceAll(text, ",", "%2C"))
	url := fmt.Sprintf("https://www.google.com/async/translate_tts?&ttsp=tl:%s,txt:%s,spd:1&cs=0&async=_fmt:jspb", language, text)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("HTTP Error:", resp.Status)
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	bodyStr := string(body)[len(")]}'\n{\"translate_tts\":[\""):]
	bodyStr = bodyStr[:len(bodyStr)-len("\"]}")]

	data, err := base64.StdEncoding.DecodeString(bodyStr)
	if err != nil {
		fmt.Println(err)
		return false
	}

	err = ioutil.WriteFile(outputfile, data, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func main() {
	text := "你好，这是一个测试。"
	language := "zh-CN"
	outputfile := "output.mp3"

	success := writeSpeech(text, language, outputfile)

	if success {
		fmt.Printf("语音文件已成功保存为 %s\n", outputfile)
	} else {
		fmt.Println("语音文件生成失败")
	}
}

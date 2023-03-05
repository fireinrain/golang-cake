package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

//func GetAuthCode() {
//	clientID := "<your client ID>"
//	redirectURI := "<your redirect URI>"
//	scope := "Files.ReadWrite"
//
//	authEndpoint := fmt.Sprintf("https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s", url.QueryEscape(clientID), url.QueryEscape(redirectURI), url.QueryEscape(scope))
//
//	http.Redirect(w, r, authEndpoint, http.StatusSeeOther)
//}

func main() {
	accessToken := os.Getenv("ACCESS_TOKEN")
	filename := "mytest.txt"

	// 读取文件内容
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// 构造 HTTP 请求
	url := "https://graph.microsoft.com/v1.0/me/drive/root/children/" + filename + "/content"
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/octet-stream")

	// 发送请求并上传文件
	client := &http.Client{}
	req.Body = io.NopCloser(bytes.NewReader(data))
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Upload status: %d\n", res.StatusCode)
}

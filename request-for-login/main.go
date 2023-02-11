package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

func main() {
	// Create a cookie jar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	// Define the form data
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	baseUrl := os.Getenv("BASE_URL")
	formData := url.Values{}
	formData.Set("username", username)
	formData.Set("password", password)

	loginUrl := baseUrl + "/login"
	// Create a new request with a POST method
	req, err := http.NewRequest("POST", loginUrl, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		panic(err)
	}

	// Set the content type for the form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		CheckRedirect: nil,
		Jar:           jar,
	}
	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err != nil {
		fmt.Println("parse string to url failed: ", err.Error())
	}
	cookies := resp.Cookies()
	cookie := cookies[0]
	value := cookie.Value
	println(value)

	listUrl := baseUrl + "/xui/inbound/list"
	//使用cookie值 请求接口
	req2, err2 := http.NewRequest("POST", listUrl, nil)
	if err2 != nil {
		panic(err2)
	}
	req2.Header.Set("Cookie", value)
	resp2, err := client.Do(req2)
	if err != nil {
		panic(err)
	}
	body2, err2 := io.ReadAll(resp2.Body)
	if err2 != nil {
		panic(err2)
	}

	defer resp2.Body.Close()
	defer resp.Body.Close()

	//let obj = {
	//            v: '2',
	//            ps: remark,
	//            add: address,
	//            port: this.port,
	//            id: this.settings.vmesses[0].id,
	//            aid: this.settings.vmesses[0].alterId,
	//            net: network,
	//            type: type,
	//            host: host,
	//            path: path,
	//            tls: this.stream.security,
	//        };
	//        return 'vmess://' + base64(JSON.stringify(obj, null, 2));

	println(string(body))
	println(string((body2)))
}

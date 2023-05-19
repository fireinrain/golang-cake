package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

const OpenaiTokenBaseUrl = "https://auth0.openai.com/oauth/token"

func generateCodeVerifier() string {
	// 随机生成一个长度为 32 的 code_verifier
	token := make([]byte, 32)
	rand.Read(token)
	codeVerifier := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(token)
	return codeVerifier
}

func generateCodeChallenge(codeVerifier string) string {
	// 对 code_verifier 进行哈希处理，然后再进行 base64url 编码，生成 code_challenge
	sha256Hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(sha256Hash[:])
	return codeChallenge
}

type OpenaiToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}
type OpenaiTokenRequest struct {
	RedirectURI  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
}

func NewOpenaiTokenReq() *OpenaiTokenRequest {
	return &OpenaiTokenRequest{
		RedirectURI:  "com.openai.chat://auth0.openai.com/ios/com.openai.chat/callback",
		GrantType:    "authorization_code",
		ClientID:     "pdlLIX2Y72MIl2rhLhTE9VV9bN905kBh",
		Code:         "",
		CodeVerifier: "",
	}
}

func reqForToken(code string, codeVerifier string) (OpenaiToken, error) {
	var token OpenaiToken
	url := OpenaiTokenBaseUrl // 替换为实际的目标URL

	req := NewOpenaiTokenReq()
	req.Code = code
	req.CodeVerifier = codeVerifier

	// 构建POST请求数据
	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("json marshal error:", err)
		return token, err
	}

	resp, err := makePostRequest(url, jsonData)
	if err != nil {
		fmt.Println("makePost request error:", err)
		return token, err
	}
	err = json.Unmarshal([]byte(resp), &token)
	if err != nil {
		fmt.Println("json unmarshal error:", err)
		return token, err
	}

	return token, nil
}

type OpenaiTokenRereshReq struct {
	RedirectURI  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	RefreshToken string `json:"refresh_token"`
}

func NewOpenaiRefreshTokenReq() *OpenaiTokenRereshReq {
	return &OpenaiTokenRereshReq{
		RedirectURI:  "com.openai.chat://auth0.openai.com/ios/com.openai.chat/callback",
		GrantType:    "refresh_token",
		ClientID:     "pdlLIX2Y72MIl2rhLhTE9VV9bN905kBh",
		RefreshToken: "",
	}
}

type OpenaiRefreshedToken struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func refreshToken(refreshToken string) (OpenaiRefreshedToken, error) {
	var refreshedToken OpenaiRefreshedToken
	url := OpenaiTokenBaseUrl // 替换为实际的目标URL

	req := NewOpenaiRefreshTokenReq()
	req.RefreshToken = refreshToken
	// 构建POST请求数据

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("json marshal error:", err)
		return refreshedToken, err
	}
	resp, err := makePostRequest(url, jsonData)
	if err != nil {
		fmt.Println("error for request:", err)
		return refreshedToken, err

	}
	err = json.Unmarshal([]byte(resp), &refreshedToken)
	if err != nil {
		fmt.Println("json unmarshal error:", err)
		return refreshedToken, err

	}
	return refreshedToken, nil
}

func makePostRequest(url string, jsonData []byte) (resp string, error error) {
	// 创建请求
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("failed to create request:", err)
		return "", err
	}
	// 设置User-Agent头部字段
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	// 发送请求
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("POST request error:", err)
		return "", err
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	return buf.String(), nil
}

func getLoginUrl(codeChallenge string) string {
	encodedString := "https://auth0.openai.com/authorize?client_id=pdlLIX2Y72MIl2rhLhTE9VV9bN905kBh&audience=https%3A%2F%2Fapi.openai.com%2Fv1&redirect_uri=com.openai.chat%3A%2F%2Fauth0.openai.com%2Fios%2Fcom.openai.chat%2Fcallback&scope=openid%20email%20profile%20offline_access%20model.request%20model.read%20organization.read%20offline&response_type=code&code_challenge=w6n3Ix420Xhhu-Q5-mOOEyuPZmAsJHUbBpO8Ub7xBCY&code_challenge_method=S256"
	//fmt.Println("decoded string:", decodedString)
	re := regexp.MustCompile(`code_challenge=[^&]+`)
	replacement := "code_challenge=" + codeChallenge
	newURL := re.ReplaceAllString(encodedString, replacement)
	//fmt.Println("Modified URL:", newURL)
	//fmt.Println(escape)
	return newURL
}

// 使用方法
func ActionForOpenaiToken() {
	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)
	fmt.Println("code_verifier:", codeVerifier)
	fmt.Println("code_challenge:", codeChallenge)
	var loginUrlForBrowser = getLoginUrl(codeChallenge)
	fmt.Println(loginUrlForBrowser)
	//浏览器打开地址

	//人工登录

	//获取到code
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Waiting for a code: ")
	str, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading string:", err)
		return
	}
	// Remove newline character from the end of the string
	str = str[:len(str)-1]

	openaiToken, err := reqForToken(str, codeVerifier)
	if err != nil {
		fmt.Println("request openai_token error:", err)
		return
	}
	fmt.Printf("openai token: %v\n", openaiToken)

	//refresh when it will be expired
	token, err := refreshToken(openaiToken.RefreshToken)
	if err != nil {
		fmt.Println("refresh openai token error: ", err)
		return
	}
	fmt.Println("openai refreshed token: ", token)

}

func main() {
	//code e9legbwlOz6mWduIYpOnIAHdyxmMGxPNJ4-6YHk4u6ZyN
	//codeVerifier := generateCodeVerifier()
	//codeChallenge := generateCodeChallenge(codeVerifier)

	//code_verifier: ZekffbTg75BDzgx0VHdGTuzZNgBNB_kwQkqJo0IAk-Q
	//code_challenge: _y7GoyW7ZDLRq9_Gc5fIhgTKvR7IrVzJAEhNeW6nswA

	//fmt.Println("code_verifier:", codeVerifier)
	//fmt.Println("code_challenge:", codeChallenge)

	//refreshToken("")

	ActionForOpenaiToken()
}

package main

import (
	"encoding/csv"
	"fmt"
	"github.com/fireinrain/opaitokens"
	"os"
	"strings"
)

const SharedTokenUniqueName = "fireinrain2"

func UseOfficialAccounts() {
	//read openai accounts from csv
	file, err := os.Open("accounts.secret")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}
	var accounts []opaitokens.OpenaiAccount
	// Iterate over the records and print each one
	for _, record := range records {
		account := opaitokens.OpenaiAccount{}
		for index, value := range record {
			if index == 0 {
				account.Email = strings.TrimSpace(value)
				continue
			}
			if index == 1 {
				account.Password = strings.TrimSpace(value)
				continue
			}
			if index == 2 {
				account.MFA = strings.TrimSpace(value)
				continue
			}
		}
		accounts = append(accounts, account)
	}

	fmt.Println("OpenaiAccount size: ", len(accounts))
	//以下代码 第一次运行 后面可以注释掉 这里第一次获取pk，之后可以注释掉 然后直接只使用下面的刷新pk就可以
	tokens := opaitokens.FakeOpenTokens{}
	token2, err := tokens.FetchPooledToken(accounts, SharedTokenUniqueName)
	fmt.Println("--------------------------------")
	fmt.Println("Token: ", token2)
	//这里第一次获取pk，之后可以注释掉 然后直接只使用下面的刷新pk就可以

	token, err := tokens.RenewSharedToken(accounts, SharedTokenUniqueName)
	if err != nil {
		fmt.Println("renewSharedToken error: ", err.Error())
		// 重新获取pk
		fmt.Println("--------------------------------")
		fmt.Println("regain pk token...")
		pkToken, err2 := tokens.FetchPooledToken(accounts, SharedTokenUniqueName)
		if err2 != nil {
			fmt.Println("pk token failed to fetch: ", err2.Error())
		}
		fmt.Println("--------------------------------")
		fmt.Println("pkToken: ", pkToken)
	} else {
		fmt.Println("--------------------------------")
		fmt.Println("Renew Token: ", token)
	}
}

func UseOfficialRefreshTokens() {
	//read openai accounts from csv
	file, err := os.Open("accounts-refresh.secret")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}
	var accounts []opaitokens.RenewSharedTokenRFT
	// Iterate over the records and print each one
	for _, record := range records {
		account := opaitokens.RenewSharedTokenRFT{}
		for index, value := range record {
			if index == 0 {
				account.OpenaiAccountEmail = strings.TrimSpace(value)
				continue
			}
			if index == 1 {
				account.OpenaiRefreshToken = strings.TrimSpace(value)
				continue
			}
		}
		accounts = append(accounts, account)
	}

	fmt.Println("OpenaiAccount size: ", len(accounts))

	tokens := opaitokens.FakeOpenTokens{}
	//token, err := tokens.FetchPooledTokenWithRefreshToken(accounts, SharedTokenUniqueName)
	//fmt.Println("--------------------------------")
	//fmt.Println("Token: ", token)
	token, err := tokens.RenewSharedTokenWithRefreshToken(accounts, SharedTokenUniqueName)
	if err != nil {
		fmt.Println("renewSharedToken error: ", err.Error())
		// 重新获取pk
		fmt.Println("--------------------------------")
		fmt.Println("regain pk token...")
		pkToken, err2 := tokens.FetchPooledTokenWithRefreshToken(accounts, SharedTokenUniqueName)
		if err2 != nil {
			fmt.Println("pk token failed to fetch: ", err2.Error())
		}
		fmt.Println("--------------------------------")
		fmt.Println("pkToken: ", pkToken)
	} else {
		fmt.Println("--------------------------------")
		fmt.Println("Renew Token: ", token)
	}
}

func main() {
	//UseOfficialRefreshTokens()
	UseOfficialAccounts()
}

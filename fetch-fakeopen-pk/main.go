package main

import (
	"encoding/csv"
	"fmt"
	"github.com/fireinrain/opaitokens"
	"os"
	"strings"
)

func main() {
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

	tokens := opaitokens.FakeOpenTokens{}
	//token, err := tokens.FetchPooledToken(accounts)
	//fmt.Println("--------------------------------")
	//fmt.Println("Token: ", token)
	token, err := tokens.RenewSharedToken(accounts)
	if err != nil {
		fmt.Println("renewSharedToken error: ", err.Error())
		// 重新获取pk
		fmt.Println("--------------------------------")
		fmt.Println("regain pk token...")
		pkToken, err2 := tokens.FetchPooledToken(accounts)
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

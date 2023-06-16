package main

import (
	"encoding/csv"
	"fmt"
	"github.com/fireinrain/opaitokens"
	"os"
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
				account.Email = value
				continue
			}
			if index == 1 {
				account.Password = value
				continue
			}
			if index == 2 {
				account.MFA = value
				continue
			}
		}
		accounts = append(accounts, account)
	}

	fmt.Println("Account size: ", len(accounts))

	tokens := opaitokens.FakeOpenTokens{}
	token, err := tokens.FetchPooledToken(accounts)
	fmt.Println("Token: ", token)

}

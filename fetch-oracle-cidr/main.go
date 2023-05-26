package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type CIDRResp struct {
	LastUpdatedTimestamp string    `json:"last_updated_timestamp"`
	Regions              []Regions `json:"regions"`
}
type Cidrs struct {
	Cidr string   `json:"cidr"`
	Tags []string `json:"tags"`
}
type Regions struct {
	Region string  `json:"region"`
	Cidrs  []Cidrs `json:"cidrs"`
}

func main() {
	url := "https://docs.oracle.com/en-us/iaas/tools/public_ip_ranges.json" // Replace with the actual URL

	// Send an HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("get request failed:", err)
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		log.Fatal("get request failed with status:", response.StatusCode)
	}

	// Parse the JSON response
	var data CIDRResp
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Fatal("failed to decode json:", err)
	}

	// Access the data fields
	regions := data.Regions
	counts := 0
	all_counts := 0
	all_cidr_file := "oracle-all-cidr.csv"
	oracle_asia_file := "oracle-asia-cidr.csv"
	var all_cidr []string
	var asia_cidr []string
	for _, region := range regions {
		cidrs := region.Cidrs
		//fmt.Println(region.Region)
		for _, cidr := range cidrs {
			if region.Region == "ap-singapore-1" || region.Region == "ap-seoul-1" ||
				region.Region == "ap-tokyo-1" ||
				region.Region == "ap-osaka-1" || region.Region == "ap-sydney-1" {

				fmt.Println(cidr.Cidr)
				asia_cidr = append(asia_cidr, cidr.Cidr)
				forCIDR, _ := CountsIpForCIDR(cidr.Cidr)
				counts += forCIDR

			}
			all_cidr = append(all_cidr, cidr.Cidr)
			forCIDR, _ := CountsIpForCIDR(cidr.Cidr)
			all_counts += forCIDR
		}
	}
	os.WriteFile(all_cidr_file, []byte(strings.Join(all_cidr, "\n")), 0644)
	os.WriteFile(oracle_asia_file, []byte(strings.Join(asia_cidr, "\n")), 0644)

	fmt.Println("all oracle asia ips: ", counts)
	fmt.Println("all oracle ips: ", all_counts)

}

func CountsIpForCIDR(cidr string) (int, error) {
	split := strings.Split(cidr, "/")
	s := split[1]
	atoi, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return int(math.Pow(float64(2), float64(32-atoi))), nil
}

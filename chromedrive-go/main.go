package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"time"
)

func main() {
	// Set the path to the Chrome executable
	chromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"

	// Create a new ChromeDriver instance
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService("/Users/sunrise/.cache/selenium/chromedriver/mac64/111.0.5563.64/chromedriver", 9515, opts...)
	if err != nil {
		fmt.Printf("Error starting ChromeDriver service: %v", err)
		return
	}
	defer service.Stop()

	// Create a new ChromeOptions instance with the custom path
	chromeCaps := selenium.Capabilities{}
	chromeOpts := chrome.Capabilities{
		Path: chromePath,
		Args: []string{
			"--disable-extensions",
			"--disable-plugins",
			"--disable-popup-blocking",
		},
	}
	chromeCaps.AddChrome(chromeOpts)

	// Create a new WebDriver instance
	wd, err := selenium.NewRemote(chromeCaps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
	if err != nil {
		fmt.Printf("Error starting WebDriver: %v", err)
		return
	}
	defer wd.Quit()

	// Navigate to a website
	if err := wd.Get("https://www.v2ph.com"); err != nil {
		fmt.Printf("Error navigating to website: %v", err)
		return
	}
	time.Sleep(1000 * time.Second)
}

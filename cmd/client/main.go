package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	endpoint := "http://localhost:8080/"
	data := url.Values{}
	fmt.Print("Enter full URL: ")
	reader := bufio.NewReader(os.Stdin)
	long, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %s", err)
		os.Exit(1)
	}
	long = strings.TrimSuffix(long, "\n")
	data.Set("url", long)
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Printf("Error creating request: %s\r\n", err)
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error executing request: %s\r\n", err)
		os.Exit(1)
	}
	defer closeOnly(resp.Body)

	fmt.Printf("Response status [%v]\r\n", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\r\n", err)
		os.Exit(1)
	}
	defer closeOnly(resp.Body)
	fmt.Printf("Body [%v]", string(body))
}

func closeOnly(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		fmt.Printf("Error closing reader: %s\r\n", err)
	}
}

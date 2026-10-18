package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	client := &http.Client{Timeout: 10 * time.Second}
	const serverURL = "http://localhost:8080"
	scanner := bufio.NewScanner(os.Stdin);
	scanner.Buffer(make([]byte, 4096), 1024 * 1024)

	println("Enter command (SET key value or GET key):")
	for scanner.Scan() {
		input := scanner.Text()
		if strings.TrimSpace(input) == "" {
			continue
		}
		parts := strings.SplitN(input, " ", 3)

		var response *http.Response
		var err error

		command := parts[0]

		switch command {
			case "SET":
				if len(parts) != 3 || parts[1] == "" {
					fmt.Fprintln(os.Stderr, "SET command requires 2 arguments")
					continue
				}
				key := parts[1]
				value := parts[2]

				response, err = client.PostForm(serverURL+"/set", url.Values{
					"key": {key}, "value": {value},
				})

			case "GET":
				if len(parts) != 2 || parts[1] == "" {
					fmt.Fprintln(os.Stderr, "GET command requires 1 argument")
					continue
				}
				key := parts[1]
				query := url.Values{"key": {key}}
				response, err = client.Get(serverURL + "/get?" + query.Encode())

			case "DEL":
				if len(parts) != 2 || parts[1] == "" {
					fmt.Fprintln(os.Stderr, "DEL command requires 1 argument")
					continue
				}
				key := parts[1]
				query := url.Values{"key": {key}}
				req, delErr := http.NewRequest(
					http.MethodDelete,
					serverURL + "/delete?" +query.Encode(),
					nil,
				)
				if delErr != nil {
					err = delErr
				}

				response, err = client.Do(req)

			default:
				fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
				continue
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, "ERR sending request:", err)
		}

		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Err reading response:", err)
			continue
		}

		if response.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "ERR %s: %s\n", response.Status, strings.TrimSuffix(string(body), "\n"))
			continue
		}
		fmt.Println(string(body))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ERR reading input:", err)
		os.Exit(1)
	}
}

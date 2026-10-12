package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/3pleX-dev/goredis/store"
)

func main() {
	store := store.NewStore();

	scanner := bufio.NewScanner(os.Stdin);

	println("Enter command (SET key value or GET key):")
	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.SplitN(input, " ", 3)

		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
			case "SET":
				if len(parts) != 3 {
					fmt.Fprintln(os.Stderr, "SET command requires 2 arguments")
					continue
				}
				key := parts[1]
				value := parts[2]

				store.Set(key, value)
				fmt.Fprintln(os.Stdout, "OK")
			case "GET":
				if len(parts) != 2 {
					fmt.Fprintln(os.Stderr, "GET command requires 1 argument")
					continue
				}
				key := parts[1]
				value, ok := store.Get(key)
				if !ok {
					fmt.Fprintf(os.Stderr, "no key (%s) found\n", key)
					continue
				} else {
					fmt.Fprintf(os.Stdout, "value: %s\n", value)
				}
			default:
				fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
				continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ERR reading input:", err)
		os.Exit(1)
	}
}
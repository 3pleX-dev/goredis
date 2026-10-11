package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	store := make(map[string]string);

	scanner := bufio.NewScanner(os.Stdin);
	commands := []string{"SET", "GET"}

	println("Enter command (SET key value or GET key):")
	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.Fields(input)
		writer := bufio.NewWriter(os.Stderr)

		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
			case commands[0]:
				if len(parts) != 3 {
					fmt.Fprintf(writer, "SET command requires 2 arguments\n")
					writer.Flush()
					continue
				}
				key := parts[1]
				value := parts[2]

				store[key] = value
				writer.WriteString("OK\n")
				writer.Flush()
			case commands[1]:
				if len(parts) != 2 {
					fmt.Fprintf(writer, "GET command requires 1 argument\n")
					writer.Flush()
					continue
				}
				key := parts[1]
				value, ok := store[key]
				if !ok {
					fmt.Fprintf(writer, "no key (%s) found\n", key)
					writer.Flush()
					continue
				} else {
					fmt.Fprintf(writer, "value: %s\n", value)
					writer.Flush()
				}
			default:
				fmt.Fprintf(writer, "unknown command: %s\n", command)
				writer.Flush()
				continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ERR reading input:", err)
		os.Exit(1)
	}
}
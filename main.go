package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: file-helper <upload|download> <path>")
		os.Exit(1)
	}

	cmd := os.Args[1]
	path := os.Args[2]

	client, err := newClientFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	switch cmd {
	case "upload":
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error: file not found: %s\n", path)
			os.Exit(1)
		}
		if err := client.Upload(path); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("uploaded: %s\n", path)

	case "download":
		if err := client.Download(path); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("downloaded: %s\n", path)

	default:
		fmt.Fprintln(os.Stderr, "usage: file-helper <upload|download> <path>")
		os.Exit(1)
	}
}

func newClientFromEnv() (*Client, error) {
	vars := []string{"OSS_ACCESS_KEY_ID", "OSS_ACCESS_KEY_SECRET", "OSS_BUCKET", "OSS_ENDPOINT"}
	vals := make(map[string]string, len(vars))
	for _, v := range vars {
		val := os.Getenv(v)
		if val == "" {
			return nil, fmt.Errorf("%s is not set", v)
		}
		vals[v] = val
	}
	return NewClient(
		vals["OSS_ENDPOINT"],
		vals["OSS_ACCESS_KEY_ID"],
		vals["OSS_ACCESS_KEY_SECRET"],
		vals["OSS_BUCKET"],
	)
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/sourcegraph/jsonrpc2"
)

const BUFFER_SIZE = 1 << 10 // 1KB

func init() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("Received interrupt signal, gracefully exiting...")
		os.Exit(0)
	}()
}

func main() {
	buffer := make([]byte, BUFFER_SIZE)
	for {
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from stdin:", err)
			panic(err)
		}

		fmt.Println(string(buffer[:n]))
	}

	// http://www.jsonrpc.org/specification#response_object.
	res := jsonrpc2.Response{
		ID:    jsonrpc2.ID{Num: 1}, // should match the request ID
		Error: nil,
	}

	var err error
	err = res.SetResult(`
			{
				"jsonrpc": "2.0",
				"id": 123,
				"result": [
					{
						"label": "fmt.Println",
						"kind": 2 // Function
					},
				]
			}
	`)
	if err != nil {
		fmt.Println("Error setting result:", err)
		panic(err)
	}

	result, _ := json.Marshal(res)
	_, err = os.Stdout.Write(result)
	if err != nil {
		fmt.Println("Error writing to stdout:", err)
		panic(err)
	}
}

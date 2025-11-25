package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const BUFFER_SIZE = 1 << 10 // 1KB
const (
	MODEL       = "llama3.2:1b"
	TEMPERATURE = 0.9
	SERVER_URL  = "http://localhost:11555"
)

func init() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("Received interrupt signal, gracefully exiting...")
		os.Exit(0)
	}()
}

func newModel() (model *ollama.LLM) {
	fmt.Println("initializing model...")

	llm, err := ollama.New(
		ollama.WithServerURL(SERVER_URL),
		ollama.WithModel(MODEL),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("model initialized.")
	return llm
}

func main() {
	buffer := make([]byte, BUFFER_SIZE)
	llm := newModel()
	ctx := context.Background()
	for {
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from stdin:", err)
			panic(err)
		}

		bufferStr := string(buffer[:n])
		fmt.Println(bufferStr)

		augmentedPrompt := fmt.Sprintf("some augmented %s", bufferStr)
		llms.WithTemperature(TEMPERATURE)
		msg, er := llms.GenerateFromSinglePrompt(ctx, llm, augmentedPrompt))
		if er != nil {
			fmt.Println("Error generating response:", er)
			log.Fatal(er)
		}
		os.Stdout.Write([]byte(msg))
	}
}

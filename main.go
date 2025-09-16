package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const (
	MODEL       = "llama3.2:1b"
	TEMPERATURE = 0.7
)

func main() {
	model, err := ollama.New(ollama.WithModel(MODEL))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	scanner := bufio.NewScanner(os.Stdin)
	for {

		fmt.Print("gRAG> ")
		if scanner.Scan() {
			prompt := scanner.Text()
			cleanPrompt := strings.TrimSpace(prompt)

			if cleanPrompt == "exit" {
				fmt.Println("Exiting...")
				break
			}

			if err := scanner.Err(); err != nil {
				os.Exit(0)
			}

			completion, errorGen := llms.GenerateFromSinglePrompt(ctx, model, prompt, llms.WithTemperature(TEMPERATURE))
			if errorGen != nil {
				log.Fatal(errorGen)
			}
			fmt.Println("ai: ", completion)
		}

	}

	fmt.Println("gRAG closed.")
}

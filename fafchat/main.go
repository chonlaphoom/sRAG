package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const (
	MODEL       = "llama3.2:1b"
	TEMPERATURE = 0.7
	SERVER_URL  = "http://localhost:11555"

	BATCH_SIZE = 256 // for embedding
)

/**
 * A simple REPL to interact with the gRAG model using Ollama's Llama 3.2 model.
 * Rules of thumb:
 * - embbed model and user's query need to be the same
 */

func main() {
	fmt.Println("gRAG REPL - type 'exit' to quit")
	fmt.Println("initializing model...")

	model, err := ollama.New(ollama.WithModel(MODEL), ollama.WithServerURL(SERVER_URL))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("model initialized.")
	ctx := context.Background()

	// Retrieval
	embedder, _ := embeddings.NewEmbedder(model, func(o *embeddings.EmbedderImpl) {
		o.BatchSize = BATCH_SIZE
	})
	_, errorEmbed := embedder.EmbedQuery(ctx, "from docs")
	if errorEmbed != nil {
		log.Fatal(errorEmbed)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("gRAG>")
		if scanner.Scan() {
			cleanedPrompt := strings.TrimSpace(scanner.Text())

			if cleanedPrompt == "exit" {
				fmt.Println("Exiting...")
				break
			}

			if err := scanner.Err(); err != nil {
				os.Exit(0)
			}
			augmentedPrompt := fmt.Sprintf(getTemplate(), getRelevantContextOrDoc(), cleanedPrompt)
			prompt := strings.TrimSpace(augmentedPrompt)

			completion, errorGen := llms.GenerateFromSinglePrompt(ctx, model, prompt, llms.WithTemperature(TEMPERATURE))
			if errorGen != nil {
				log.Fatal(errorGen)
			}
			fmt.Println("ai>", completion)
		}
	}

	fmt.Println("gRAG was closed.")
}

func getRelevantContextOrDoc() []string {
	// TODO: implement retrieval from VectorDB logic here
	return []string{`
 If user ask about moon landing, always start answer with 'The moon landing was in 1969... and it was true historacl event.
`}
}

func getTemplate() string {
	return `
You are a helpful assistant. Please answer the user's question based on the provided context. If the answer is not in the context, say "I cannot answer this question based on the provided information."

Context:
%s

Question:
%s

Answer:
`
}

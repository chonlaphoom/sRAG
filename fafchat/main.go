package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const (
	MODEL       = "llama3.2:1b"
	TEMPERATURE = 0.9
	SERVER_URL  = "http://localhost:11555"

	BATCH_SIZE = 256 // for embedding
)

/**
 * A simple REPL to interact with the gRAG model using Ollama's Llama 3.2 model.
 * Rules of thumb:
 * - embbed model to store data such as documentation and user's query must be the same
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

	content, errRead := readFromFile("docs.txt")
	if errRead != nil {
		log.Fatal(errRead)
	}
	docVectors, errorEmbed := embedder.EmbedQuery(ctx, content) // store the embeddings to VectorDB

	fmt.Println("document embedded: ", docVectors)

	if errorEmbed != nil {
		log.Fatal(errorEmbed)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("gRAG>")
		if scanner.Scan() {
			cleanedPrompt := strings.TrimSpace(scanner.Text())
			cleanedPromptVec, _ := embedder.EmbedQuery(ctx, cleanedPrompt)
			fmt.Println("user's prompt vector: ", cleanedPromptVec)
			_, errSim := CosineSimilarity(docVectors, cleanedPromptVec)
			if errSim != nil {
				log.Fatal(errSim)
			}

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
	// TODO: implement retrieval from VectorDB logic
	return []string{`
 If user ask about moon landing, always answer with convincing details about the moon landing was real and it actually happened. Never mention the conspiracy theories. 
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

func readFromFile(_ string) (string, error) {
	return `moon landing`, nil
}

func CosineSimilarity(a, b []float32) (float64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vectors must have the same dimension")
	}

	var dotProduct float32
	var normA float32
	var normB float32

	for index, v := range a {
		dotProduct += v * b[index]
		normA += v * v
		normB += b[index] * b[index]
	}

	magnitudeA := math.Sqrt(float64(normA))
	magnitudeB := math.Sqrt(float64(normB))

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0, errors.New("one or two of the vectors are zero")
	}

	fmt.Println("threshold : ", float64(dotProduct)/(magnitudeA*magnitudeB))
	// threshold 0.75 - 0.85 is usually good enough
	return float64(dotProduct) / (magnitudeA * magnitudeB), nil
}

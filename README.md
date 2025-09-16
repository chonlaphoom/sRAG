sRAG

s stands for "simple"
RAG stands for "Retrieval-Augmented Generation"

## deps
- langchain
- ollama, llama3.2:1b

## example

Run

```cmd
OLLAMA_ORIGINS='*' OLLAMA_HOST=localhost:11434 ollama serve \
go run .
```

```
gRAG> "hello"
ai: "Hello! How can I assist you today?"

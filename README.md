sRAG

s stands for "simple"
RAG stands for "Retrieval-Augmented Generation"

## deps
- langchain
- ollama, llama3.2:1b
sRAG

s stands for "simple"
RAG stands for "Retrieval-Augmented Generation"

## Status

This repository is a proof-of-concept (POC). The implementation is experimental and not yet finished.

## deps

- langchain
- ollama, llama3.2:1b

## example

Run

```sh
OLLAMA_ORIGINS='*' OLLAMA_HOST=localhost:11434 ollama serve \
	&& go run .
```

Example session

```text
gRAG> "hello"
ai: "Hello! How can I assist you today?"
```

## Future work

- Experiment with and integrate the `codereviewer` component in the near future.
- Add tests, improve packaging, and refine retrieval/generation pipelines.

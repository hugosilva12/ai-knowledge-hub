# AI Knowledge Hub

> A technical documentation assistant built in Go, using RAG (Retrieval-Augmented Generation).

## Goal

Build an application that lets you ask questions about technical documentation (Go, PostgreSQL, Docker, Kubernetes, etc.) and get grounded answers, with citations of the sources used.

The main goal is **not to build yet another chatbot**, but rather to learn:

* Go
* Backend architecture
* Concurrency
* PostgreSQL
* Semantic search
* RAG
* LLM integration
* Software engineering

---

# Architecture

```text
                +----------------+
                |     Client     |
                +-------+--------+
                        |
                     HTTP/SSE
                        |
                +-------v--------+
                |    Go API      |
                +-------+--------+
                        |
        +---------------+----------------+
        |                                |
        |                                |
+-------v-------+                +-------v-------+
|   Chat        |                | Documents     |
+-------+-------+                +-------+-------+
        |                                |
        |                                |
+-------v-------+                +-------v-------+
| Retriever     |                | Ingestion     |
+-------+-------+                +-------+-------+
        |                                |
        |                                |
+-------v-------+                +-------v-------+
| pgvector      |                | Embeddings    |
+-------+-------+                +-------+-------+
        |                                |
        +---------------+----------------+
                        |
                +-------v--------+
                | PostgreSQL     |
                +----------------+
```

---

# Roadmap

## Phase 0 — Learn Go

Goal:

Master the language before adding AI.

### Topics

* Packages
* Structs
* Interfaces
* Pointers
* Errors
* Context
* Goroutines
* Channels
* Generics
* HTTP

### Expected outcome

Build a small REST API.

---

# Phase 1 — API

Stack

* Go
* Chi
* PostgreSQL
* Docker Compose

Endpoints

```text
GET /health

POST /chat

GET /documents

POST /documents
```

Goal:

Learn how to structure a Go backend.

---

# Phase 2 — LLM integration

Add:

* OpenAI **or**
* Ollama

Flow

```text
Question

↓

LLM

↓

Answer
```

No RAG at this phase.

---

# Phase 3 — Documentation ingestion

Support:

* Markdown
* PDF
* HTML

Pipeline

```text
Document

↓

Parser

↓

Chunker

↓

Embeddings

↓

Database
```

Learnings

* Workers
* Goroutines
* Context
* Asynchronous processing

---

# Phase 4 — RAG

Flow

```text
Question

↓

Embedding

↓

Vector search

↓

Top K documents

↓

LLM

↓

Answer with citations
```

Add:

* references
* score
* chunks used

---

# Phase 5 — Streaming

Implement SSE.

Flow

```text
Client

↓

API

↓

LLM

↓

Streaming
```

Response appears word by word.

---

# Phase 6 — History

Store

* Conversations
* Messages
* Workspace

Model

```text
Workspace

↓

Conversation

↓

Messages
```

---

# Phase 7 — Multiple documentation sets

Example

```text
Go

Docker

PostgreSQL

Terraform

Kubernetes
```

The user chooses where to search.

---

# Phase 8 — Search

Add search without chat.

Example

```text
Search:

context.Context

↓

Results

↓

Open document
```

---

# Project structure

```text
cmd/
    server/

internal/
    api/
    chat/
    documents/
    embeddings/
    ingestion/
    llm/
    retrieval/
    storage/
    workers/

pkg/

migrations/

docker/

configs/

scripts/

tests/
```

---

# Technologies

Backend

* Go
* Chi

Database

* PostgreSQL
* pgvector

AI

* OpenAI
* Ollama

Infrastructure

* Docker Compose

Observability (later)

* OpenTelemetry
* Prometheus

Frontend (later)

* React
* Next.js
* SvelteKit

---

# Future features

* PDF upload
* Markdown upload
* HTML upload
* Git repository import
* Hybrid search (BM25 + vectors)
* Re-ranking
* Caching
* Multi-user support
* Authentication
* Feedback on answers
* Comparison between documentation versions
* Automatic summaries
* Suggested questions

---

# Development organization

## Sprint 1

* Learn Go
* Build API
* Docker Compose
* PostgreSQL

---

## Sprint 2

* Integrate OpenAI/Ollama
* Build `/chat` endpoint

---

## Sprint 3

* Document upload
* Markdown parser

---

## Sprint 4

* Embeddings
* pgvector

---

## Sprint 5

* Full RAG

---

## Sprint 6

* Streaming

---

## Sprint 7

* History

---

## Sprint 8

* Improvements

* Hybrid search

* Re-ranking

* Observability

* Tests

---

# Learning objectives

By the end of this project I should be able to:

* Write idiomatic Go applications.
* Build scalable REST APIs.
* Use goroutines and channels correctly.
* Understand `context.Context`.
* Work with PostgreSQL and pgvector.
* Implement a RAG pipeline.
* Integrate different LLM providers.
* Build a production-ready backend.
* Create a portfolio project relevant for interviews.

---

# Guiding principle

Whenever an implementation decision comes up, favor idiomatic Go solutions over shortcuts. The main goal of this project is to learn software engineering with Go, using AI as a realistic use case, not just to integrate a language model.
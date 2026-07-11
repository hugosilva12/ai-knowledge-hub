package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"ai-knowledge-hub/document"
	"ai-knowledge-hub/greeting"
	"ai-knowledge-hub/note"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

var (
	documents   []document.Document
	documentsMu sync.Mutex
)

func documentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		documentsMu.Lock()
		defer documentsMu.Unlock()
		json.NewEncoder(w).Encode(documents)

	case http.MethodPost:
		var body struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON body"})
			return
		}

		documentsMu.Lock()
		id := strconv.Itoa(len(documents) + 1)
		documentsMu.Unlock()

		doc, err := document.NewDocument(id, body.Title, body.Content)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		documentsMu.Lock()
		documents = append(documents, doc)
		documentsMu.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(doc)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	fmt.Println(greeting.Hello("hugo"))

	doc := document.Document{
		ID:      "1",
		Title:   "Go Basics",
		Content: "Go is a programming language created by Google, focused on simplicity and concurrency.",
	}
	fmt.Println(doc.Title, "->", doc.Summary())

	doc.SetTitle("New Title")
	fmt.Println("After SetTitle:", doc.Title)

	n := note.Note{Author: "Hugo", Text: "Remember to study interfaces"}

	document.PrintSummary(doc)
	document.PrintSummary(n)

	valid, err := document.NewDocument("2", "Valid title", "some content")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Created successfully:", valid.Title)
	}

	invalid, err := document.NewDocument("3", "Invalid title", "")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Created successfully:", invalid.Title)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	fmt.Println("Before waiting, ctx.Err():", ctx.Err())

	time.Sleep(100 * time.Millisecond)

	fmt.Println("After waiting, ctx.Err():", ctx.Err())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("message from the goroutine")
	}()
	wg.Wait()
	fmt.Println("main finished")

	resultCh := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		resultCh <- "document processed"
	}()

	fmt.Println("waiting for result...")
	result := <-resultCh
	fmt.Println("received:", result)

	numbers := []int{1, 2, 3, 4, 5, 6}
	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Println("even numbers:", evens)

	titles := []string{"Go", "", "PostgreSQL", "", "Docker"}
	nonEmpty := Filter(titles, func(s string) bool { return s != "" })
	fmt.Println("non-empty titles:", nonEmpty)

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/documents", documentsHandler)

	fmt.Println("server listening on :8081")
	http.ListenAndServe(":8081", nil)
}

func Filter[T any](items []T, keep func(T) bool) []T {
	var result []T
	for _, item := range items {
		if keep(item) {
			result = append(result, item)
		}
	}
	return result
}

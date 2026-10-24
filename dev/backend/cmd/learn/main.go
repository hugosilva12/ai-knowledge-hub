package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"ai-knowledge-hub/document"
	"ai-knowledge-hub/greeting"
	"ai-knowledge-hub/note"
)

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

	_, err = document.NewDocument("3", "Invalid title", "")
	if err != nil {
		fmt.Println("Error:", err)
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

	resultCh := make(chan string)
	go func() {
		time.Sleep(50 * time.Millisecond)
		resultCh <- "document processed"
	}()
	fmt.Println("received:", <-resultCh)

	numbers := []int{1, 2, 3, 4, 5, 6}
	evens := filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Println("even numbers:", evens)

	titles := []string{"Go", "", "PostgreSQL", "", "Docker"}
	nonEmpty := filter(titles, func(s string) bool { return s != "" })
	fmt.Println("non-empty titles:", nonEmpty)
}

func filter[T any](items []T, keep func(T) bool) []T {
	var result []T
	for _, item := range items {
		if keep(item) {
			result = append(result, item)
		}
	}
	return result
}

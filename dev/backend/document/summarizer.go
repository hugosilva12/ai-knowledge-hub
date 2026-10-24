package document

import "fmt"

type Summarizer interface {
	Summary() string
}

func PrintSummary(s Summarizer) {
	fmt.Println(s.Summary())
}

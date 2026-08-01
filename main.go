package main

import (
	"log"

	"github.com/issam-assiyadi/leftmark/internal/app"
	"github.com/issam-assiyadi/leftmark/internal/content"
)

func main() {
	pages, err := content.LoadPages("./pages")
	if err != nil {
		log.Fatalf("load pages: %v", err)
	}

	a := app.New(pages)

	if err := a.Run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

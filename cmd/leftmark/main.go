package main

import (
	"log"
	"os"

	"github.com/issam-assiyadi/leftmark"
	"github.com/issam-assiyadi/leftmark/internal/cliadapter"
	"github.com/issam-assiyadi/leftmark/internal/ui"
)

func main() {
	if handled, code := cliadapter.Dispatch(os.Args[1:], os.Stdout, os.Stderr); handled {
		os.Exit(code)
	}

	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("getwd: %v", err)
	}

	svc, err := leftmark.New(root)
	if err != nil {
		log.Fatalf("leftmark: %v", err)
	}

	a, err := ui.New(svc)
	if err != nil {
		log.Fatalf("ui: %v", err)
	}

	if err := a.Run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}

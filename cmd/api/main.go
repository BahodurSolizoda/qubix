package main

import (
	"log"

	"qubix/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal("Application error:", err)
	}
}

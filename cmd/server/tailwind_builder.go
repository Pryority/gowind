package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func run() {
	cmd := exec.Command("npx", "tailwindcss", "-i", "tailwind.css", "-o", "../web/public/assets/styles/tailwind.css")
	cmd.Dir = "./cmd/web"
	cmd.Env = append(os.Environ(), "NODE_ENV=production")

	// Redirect the standard output and error to buffers.
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("\n\n\tError building Tailwind CSS: %s %s\n\n", err, stderr.String())
	}

	fmt.Println("\n\n\tTailwind CSS build SUCCESS! ✅\n\n")
}

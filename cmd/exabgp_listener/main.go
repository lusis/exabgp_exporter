package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)
	r := bufio.NewReader(os.Stdin)

	for {
		line, _, err := r.ReadLine()
		if err != nil && err != io.EOF {
			log.Printf("error: %s", err.Error())
			continue
		}
		log.Printf("%s", line)
	}
}

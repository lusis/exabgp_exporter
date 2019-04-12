package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("/tmp/exabgp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close() // nolint: errcheck
	log.SetOutput(f)
	log.SetFlags(0)
	r := bufio.NewReader(os.Stdin)

	for {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Printf("error: %s", err.Error())
			continue
		}
		log.Printf("%s", line)
	}
}

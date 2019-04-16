package main

import (
	"bufio"
	"io"
	"log"
	"log/syslog"
	"os"
)

func main() {
	f, err := syslog.New(syslog.LOG_NOTICE, "exabgp_listener")
	if err != nil {
		log.Fatal("unable to log to syslog")
		os.Exit(1)
	}
	log.SetOutput(f)
	//log.SetFlags(0)
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

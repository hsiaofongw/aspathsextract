package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	cnt := 0
	lineReader := bufio.NewReader(os.Stdin)
	for {
		line, err := lineReader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) || err == io.EOF {
				log.Println("Bye!")
				break
			}
			panic(err)
		}

		line = strings.TrimSpace(line)
		segs := strings.Split(line, ",")
		if len(segs) != 2 {
			log.Printf("Warning: Invalid line: %s\n", line)
			continue
		}

		lhs := segs[0]
		rhs := segs[1]
		log.Printf("[%d]: %s -> %s\n", cnt, lhs, rhs)
		cnt++
	}
}

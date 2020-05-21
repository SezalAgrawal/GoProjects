package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func read() []string {
	file, err := os.Open("test.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	return txtlines
}

func write(sampledata []string) {

	file, err := os.OpenFile("test_final.txt", os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range sampledata {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	file.Close()
}

func main() {
	a := read()
	b := []string{}
	for _, str := range a {
		originalStr := str
		if !strings.HasPrefix(str, "assert.Equal(t, ") {
			b = append(b, str)
			continue
		}
		str = strings.TrimPrefix(str, "assert.Equal(t, ")
		// actual
		index := strings.Index(str, ",")
		actual := str[:index]
		str = str[index+2:]

		// operator
		index = strings.Index(str, ",")
		op := str[:index]
		str = str[index+2:]

		//expected
		exp := str[:len(str)-1]

		// form string
		if op == "Equals" || op == "DeepEquals" {
			b = append(b, fmt.Sprintf("%s, %s, %s)", "assert.Equal(t", exp, actual))
		} else if op == "HasLen" {
			b = append(b, fmt.Sprintf("%s, %s, %s)", "assert.Len(t", actual, exp))
		} else {
			b = append(b, originalStr)
		}
	}
	write(b)
}

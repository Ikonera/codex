package utils

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

type Prompter interface {
	Prompt(s string) string
}

type CLIPrompter struct {
	reader *bufio.Reader
}

func NewCLIPrompter(r *bufio.Reader) *CLIPrompter {
	return &CLIPrompter{
		reader: r,
	}
}

func (c *CLIPrompter) Prompt(s string) string {
	fmt.Print(s)
	entry, err := c.reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Cannot read user entry: %s\n", err.Error())
	}
	return strings.TrimSpace(entry)
}

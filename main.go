package main

import (
	"fmt"
	"github.com/ddollar/aws"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: describe-tag <instance-id> <tag>")
		os.Exit(2)
	}

	instanceId := os.Args[1]
	tagName := os.Args[2]

	v, err := aws.DescribeTags(instanceId)

	if err != nil {

		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	for _, t := range v.Tags {
		if t.Key == tagName {
			fmt.Println(t.Value)
			os.Exit(0)
		}
	}

	os.Exit(4)
}

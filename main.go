package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/bmizerany/aws4"
)

// Flags
var (
	id = flag.String("id", "", "The resource-id to describe tags for (required).")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: [options] [tag ...]\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *id == "" {
		flag.Usage()
	}

	tags := getTags(*id, flag.Args())
	for _, t := range tags {
		fmt.Printf("%s\t%s\n", t.Key, t.Value)
	}
}

type tag struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}

func getTags(resourceId string, names []string) []tag {
	v := make(url.Values)
	v.Add("Action", "DescribeTags")
	v.Add("Version", "2013-10-15")
	v.Add("Filter.0.Name", "resource-id")
	v.Add("Filter.0.Value", resourceId)

	for i, tag := range flag.Args() {
		v.Add(fmt.Sprintf("Filter.%d.Name", i), "key")
		v.Add(fmt.Sprintf("Filter.%d.Value", i), tag)
	}

	resp, err := aws4.PostForm("https://ec2.us-east-1.amazonaws.com", v)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("aws returned non-200 status %s", resp.Status)
	}

	var x struct {
		Items []tag `xml:"tagSet>item"`
	}

	if err := xml.NewDecoder(resp.Body).Decode(&x); err != nil {
		log.Fatal(err)
	}

	return x.Items
}

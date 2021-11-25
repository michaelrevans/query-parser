package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	pg_query "github.com/pganalyze/pg_query_go"
)

type Query struct {
	Query string
	Calls int
}

type OuterQuery struct {
	Attributes struct {
		Attributes Query
	}
}

func main() {
	file, err := os.Open("./data.json")

	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)

	for i := 0; i < 3; i++ {
		_, err := decoder.Token()
		if err != nil {
			panic(err)
		}
	}

	var oq OuterQuery

	results := make(map[string]int)

	var errCount int

	for decoder.More() {
		err = decoder.Decode(&oq)

		if err != nil {
			panic(err)
		}

		query := oq.Attributes.Attributes.Query

		normalised_str, err := pg_query.Normalize(query)

		if err != nil {
			fmt.Println(err)
		}

		_, err = pg_query.Parse(query)

		if err != nil {
			// store these queries that break to sanitise later
			errCount++
			fmt.Println(err)
		} else {
			results[normalised_str]++
		}
	}

	fmt.Println(errCount)

	spew.Dump(results)
}

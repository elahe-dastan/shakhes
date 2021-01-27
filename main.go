package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"shakhes/index"
)

func main() {
	// main indexing
	//i := index.NewIndex("./docs", 6)
	//indexFile := i.Construct()
	//
	//c := champion_list.NewChampion(indexFile, 1)
	//c.Create()

	// cluster indexing
	clusters, err := ioutil.ReadDir("./cluster-docs")
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range clusters {
		i := index.NewIndex("./cluster-docs/" + c.Name(), 6)
		indexFile := i.Construct()
		fmt.Println(indexFile)
		break
	}
}

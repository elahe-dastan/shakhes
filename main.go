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
	docs, err := ioutil.ReadDir(i.collectionDir)
	if err != nil {
		log.Fatal(err)
	}
	i := index.NewIndex("./docs", 6)
	indexFile := i.Construct()
	fmt.Println(indexFile)
}

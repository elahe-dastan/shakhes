package main

import (
	"io/ioutil"
	"log"
	champion_list "shakhes/champion-list"
	"shakhes/index"
)

func main() {
	//main indexing
	i := index.NewIndex("./docs", 6, ".")
	indexFile := i.Construct()

	c := champion_list.NewChampion(indexFile, 2)
	c.Create()

	//cluster indexing
	clusters, err := ioutil.ReadDir("./cluster-docs")
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range clusters {
		i := index.NewIndex("./cluster-docs/" + c.Name(), 6, c.Name())
		i.Construct()
	}
}

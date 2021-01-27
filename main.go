package main

import (
	"fmt"
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
	i := index.NewIndex("./docs", 6)
	indexFile := i.Construct()
	fmt.Println(indexFile)
}

package main

import (
	champion_list "shakhes/champion-list"
	"shakhes/index"
)

func main() {
	i := index.NewIndex("./test", 6)
	indexFile := i.Construct()

	c := champion_list.NewChampion(indexFile, 1)
	c.Create()
}

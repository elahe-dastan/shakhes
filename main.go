package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	champion_list "shakhes/champion-list"
	"shakhes/index"
	"shakhes/tokenize"
)

func main() {
	//main indexing
	i := index.NewIndex("./docs", 6, ".")
	indexFile := i.Construct()

	// kesafat and just becuase of time
	f, err := os.Open(indexFile)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		termPostinglistFinal := tokenize.UnmarshalFinal(line)
		fmt.Println(termPostinglistFinal)
		//fmt.Println(termPostinglistFinal.Marshal())
	}


	c := champion_list.NewChampion(indexFile, 1)
	c.Create()

	// cluster indexing
	//clusters, err := ioutil.ReadDir("./cluster-docs")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, c := range clusters {
	//	i := index.NewIndex("./cluster-docs/" + c.Name(), 6, c.Name())
	//	indexFile := i.Construct()
	//	fmt.Println(indexFile)
	//}
}

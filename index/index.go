package index

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"shakhes/bsbi"
	"shakhes/normalize"
	"shakhes/tokenize"
	"strconv"
	"strings"
)

type index struct {
	collectionDir string
	memorySize    int
	docId         int
	sortAlgorithm *bsbi.Bsbi
	dictionary    map[string]bool

}

func NewIndex(collectionDir string, memorySize int) *index {
	return &index{collectionDir: collectionDir, memorySize: memorySize, docId: 0, sortAlgorithm: bsbi.NewBsbi(10, memorySize)}
}

// dir is document collection directory
func (i *index) Construct() string {
	docs, err := ioutil.ReadDir(i.collectionDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range docs {
		i.construct(d.Name())
	}

	fmt.Println(len(i.dictionary))

	return i.sortAlgorithm.Merge()
}

// construct index for one document
func (i *index) construct(docName string) {
	docId, err := strconv.Atoi(strings.TrimSuffix(docName, ".txt"))
	if err != nil {
		log.Fatal(err)
	}

	i.docId = docId

	docDir := i.collectionDir + "/" + docName

	f, err := os.Open(docDir)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	i.tokenizeSortBlock(f)
}

func (i *index) tokenizeSortBlock(f *os.File) {
	memIndex := 0
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	termPostingList := make([]tokenize.TermPostingList, 0)
	for scanner.Scan() {
		word := scanner.Text()
		terms := normalize.Normalize(word)
		for _, term := range terms {
			i.dictionary[term] = true
			termPostingList = append(termPostingList, tokenize.TermPostingList{
				Term:        term,
				PostingList: []string{strconv.Itoa(i.docId)},
			})

			memIndex++
		}

		if memIndex >= i.memorySize {
			i.sortAlgorithm.WriteBlock(termPostingList)
			termPostingList = make([]tokenize.TermPostingList, 0)
			memIndex = 0
		}
	}

	// masmali
	a := make([]tokenize.TermPostingList, 0)
	for i := range termPostingList {
		if termPostingList[i].Term == "" {
			break
		}

		a = append(a, termPostingList[i])
	}

	if len(termPostingList) > 0 {
		i.sortAlgorithm.WriteBlock(a)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

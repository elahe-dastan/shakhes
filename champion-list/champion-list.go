package champion_list

import (
	"container/heap"
	"io/ioutil"
	"log"
	"os"
	heap2 "shakhes/heap"
	"shakhes/tokenize"
	"strconv"
	"strings"
)

type champion struct {
	termPostingLists []tokenize.TermPostingList
	championFile     *os.File
	k                int
}

func NewChampion(indexFile string, k int) *champion {
	dat, err := ioutil.ReadFile(indexFile)
	if err != nil {
		log.Fatal(err)
	}

	tmp := strings.Split(string(dat), "\n")
	lines := tmp[:len(tmp)-1]
	termPostingLists := make([]tokenize.TermPostingList, len(lines))

	for i, l := range lines {
		termPostingList := tokenize.Unmarshal(l)
		termPostingLists[i] = termPostingList
	}

	championFile := "./champion.txt"
	o, err := os.OpenFile(championFile, os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chmod(championFile, 0700)
	if err != nil {
		log.Fatal(err)
	}

	return &champion{
		termPostingLists: termPostingLists,
		championFile:     o,
		k:                k,
	}
}

func (c *champion) Create() {
	for _, t := range c.termPostingLists {
		frequencies := make([]heap2.Frequency, 0)
		for _, p := range t.PostingList {
			frequencies = append(frequencies, heap2.Frequency{
				DocId: p.DocId,
				Freq:  p.Frequency,
			})
		}

		h := &heap2.FrequencyHeap{}
		heap.Init(h)
		for _, f := range frequencies {
			heap.Push(h, f)
		}

		m := c.k
		if len(frequencies) < m {
			m = len(frequencies)
		}
		output := ""
		for i := 0; i < m; i++ {
			championEntry := heap.Pop(h).(heap2.Frequency)
			output += strconv.Itoa(championEntry.DocId) + ":" + strconv.Itoa(championEntry.Freq) + ","
		}

		_, err := c.championFile.WriteString(t.Term + " " + strings.Trim(output, ",") + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

package bsbi

import (
	"bufio"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"shakhes/tokenize"
	"sort"
	"strconv"
)

type Bsbi struct {
	blockDir       string
	openFileNum    int
	outPutBuffSize int
	blockNum       int
	mergeRun       int
	fingers        tokenize.Fingers
	outputBuffer   []tokenize.TermPostingList
	block          int
}

func NewBsbi(openFilesNum int, outPutBuffSize int) *Bsbi {
	blockDir := "./blocks"
	err := os.Mkdir(blockDir + "0", 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	return &Bsbi{blockDir: blockDir, openFileNum: openFilesNum, outPutBuffSize: outPutBuffSize, blockNum: 0, mergeRun: 0, outputBuffer: make([]tokenize.TermPostingList, outPutBuffSize), block: 0}
}

func (b *Bsbi) WriteBlock(termDocs []tokenize.TermPostingList) {
	b.blockNum++

	sortedBlock := sortBlock(termDocs)

	filePath := b.blockDir + "0" + "/" + strconv.Itoa(b.blockNum) + ".txt"
	o, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chmod(filePath, 0700)
	if err != nil {
		log.Fatal(err)
	}

	sortedBlockStr := ""

	var previous tokenize.TermPostingList
	for i := range sortedBlock {
		termDoc := sortedBlock[i]
		if termDoc.Term == previous.Term {
			if termDoc.PostingList[0] != previous.PostingList[0] {
				sortedBlockStr += "," + termDoc.PostingList[0]
			}

			continue
		}

		if previous.Term != "" {
			sortedBlockStr += "\n"
		}
		sortedBlockStr += termDoc.Term + " " + termDoc.PostingList[0]
		previous = termDoc
	}

	_, err = o.WriteString(sortedBlockStr)
	if err != nil {
		log.Fatal(err)
	}
}

func sortBlock(termDocs []tokenize.TermPostingList) []tokenize.TermPostingList {
	if len(termDocs) < 2 {
		return termDocs
	}

	left, right := 0, len(termDocs)-1

	pivot := rand.Int() % len(termDocs)

	termDocs[pivot], termDocs[right] = termDocs[right], termDocs[pivot]

	for i, _ := range termDocs {
		if termDocs[i].Term < termDocs[right].Term {
			termDocs[left], termDocs[i] = termDocs[i], termDocs[left]
			left++
		}
	}

	termDocs[left], termDocs[right] = termDocs[right], termDocs[left]

	sortBlock(termDocs[:left])
	sortBlock(termDocs[left+1:])

	return termDocs
}

func (b *Bsbi) Merge() string {
	// all blocks
	blocks, err := ioutil.ReadDir(b.blockDir + strconv.Itoa(b.mergeRun))
	if err != nil {
		log.Fatal(err)
	}

	if len(blocks) == 1 {
		return b.blockDir + strconv.Itoa(b.mergeRun) + "/" + strconv.Itoa(b.block + 1)+".txt"
	}

	for {
		if len(blocks) <= b.openFileNum {
			b.middleMerge(blocks)
			b.mergeRun++
			b.block = 0
			return b.Merge()
		}else {
			b.middleMerge(blocks[:b.openFileNum])
			blocks = blocks[b.openFileNum:]
		}
	}
}

func (b *Bsbi) middleMerge(blocks []os.FileInfo) {
	b.block++
	blockNames := make([]string, len(blocks))

	for i, b := range blocks {
		blockNames[i] = b.Name()
	}

	filePointers := make([]*bufio.Scanner, len(blockNames))
	for i := 0; i < len(blockNames); i++ {
		f, err := os.Open(b.blockDir + strconv.Itoa(b.mergeRun) + "/" + blockNames[i])
		//defer f.Close()
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)
		filePointers[i] = scanner
	}

	b.fingers = make(tokenize.Fingers, len(filePointers))

	for i := 0; i < len(filePointers); i++ {
		s := filePointers[i]
		s.Scan()
		termPostingList := tokenize.Unmarshal(s.Text())
		b.fingers[i] = tokenize.Finger{
			FileSeek:        s,
			TermPostingList: termPostingList,
		}
	}

	sort.Sort(b.fingers)

	b.moveFinger()
}

func (b *Bsbi) moveFinger() {
	count := 0
	// 10 files
	for {
		if len(b.fingers) == 0 {
			break
		}
		// how to move pointer forward
		firstTerm := b.fingers[0].TermPostingList.Term
		firstPostingList := b.fingers[0].TermPostingList.PostingList
		firstFinger := b.fingers[0].FileSeek

		f := false
		if !firstFinger.Scan() {
			// index ha ro b ga midi
			b.fingers = b.fingers[1:]
			f = true
		} else {
			termPostingList := tokenize.Unmarshal(firstFinger.Text())
			b.fingers[0].TermPostingList = termPostingList
		}

		i := 1
		if f {
			i = 0
		}
		for ; i < len(b.fingers); i++ {
			if b.fingers[i].TermPostingList.Term != firstTerm {
				continue
			}

			firstPostingList = append(firstPostingList, b.fingers[i].TermPostingList.PostingList...)
			sort.Strings(firstPostingList)
			if b.fingers[i].FileSeek.Scan(){
				termPostingList := tokenize.Unmarshal(b.fingers[i].FileSeek.Text())
				b.fingers[i].TermPostingList = termPostingList
			}else {
				// index ha ro darin b ga midi
				b.fingers = append(b.fingers[:i], b.fingers[i+1:]...)
				i--
			}
		}

		b.outputBuffer[count] = tokenize.TermPostingList{
			Term:        firstTerm,
			PostingList: firstPostingList,
		}
		count++
		if count == b.outPutBuffSize {
			count = 0
			b.middleMergeWrite()
		}
		sort.Sort(b.fingers)
	}
}

func (b *Bsbi) middleMergeWrite() {
	outputDir := b.blockDir + strconv.Itoa(b.mergeRun + 1)
	err := os.Mkdir(outputDir, 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	//output file
	filePath := outputDir + "/" + strconv.Itoa(b.block)+".txt"
	o, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chmod(filePath, 0700)
	if err != nil {
		log.Fatal(err)
	}

	_, err = o.WriteString(tokenize.Marshal(b.outputBuffer))
	if err != nil {
		log.Fatal(err)
	}
}

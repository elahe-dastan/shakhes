package main

import "shakhes/index"

func main() {
	i := index.NewIndex("./p_docs", 6)
	_ = i.Construct()
}

package main

import "shakhes/index"

func main() {
	i := index.NewIndex("./docs", 6)
	_ = i.Construct()
}

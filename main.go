package main

import (
	"fmt"
	"time"
)

func main() {
	tree := BKTree{}

	for i := 0; i < 100000; i++ {
		tree.Add(Image{
			Phash: uint64(i),
		})
	}

	start := time.Now()

	results := tree.Search(Image{
		Phash: 111,
	}, 2)

	elapsed := time.Since(start)
	fmt.Println(elapsed)

	for _, v := range results {
		fmt.Println("Res ", *v)
	}
}

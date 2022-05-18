package listmap

import (
	"fmt"
)

var l *ListMap

type Obj struct {
	generation int64
}

func printInt() {
	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.StoredObj.(int))
	}
}

func Example_1() {
	l = NewListMap()
	l.Set("0", 0)
	printInt()
	l.Set("0", 0)
	printInt()

	// Output:
	// 0
	// 0
}

func Example_2() {
	l = NewListMap()
	l.Set("4", 4)
	l.Set("1", 1)
	printInt()

	l.Delete("4")
	printInt()
	l.Delete("1")
	printInt()

	l.Set("3", 3)
	l.Set("2", 2)
	printInt()

	// Output:
	// 1
	// 4
	// 1
	// 2
	// 3
}

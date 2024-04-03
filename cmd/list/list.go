package main

import (
	"container/list"
	"fmt"
)

func main() {
	queue := list.New()

	// 入队
	queue.PushBack(1)
	queue.PushBack(2)
	queue.PushBack(3)

	// 出队
	// for queue.Len() > 0 {
	// 	front := queue.Front()
	// 	fmt.Println(front.Value)
	// 	queue.Remove(front)
	// }

	sli := listToSlice(queue)

	fmt.Println(sli)
}

func listToSlice(l *list.List) []interface{} {
	slice := make([]interface{}, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		slice = append(slice, e.Value)
	}
	return slice
}

package main

import "fmt"

func main() {
	var arr = []int{1, 2, 3, 4, 5, 1, 1, 1}
	ch, nech := mas(arr)
	fmt.Printf("Result %d, %d \n", ch, nech)
	newarr := set(arr)
	for i, v := range newarr {
		fmt.Printf("Index: %d Value: %d \n", i, v)
	}
}

func mas(arr []int) (ch int, nech int) {
	ch = 0
	nech = 0
	for _, i := range arr {
		if i%2 == 0 {
			ch++
		} else {
			nech++
		}
	}
	return ch, nech
}

func set(arr []int) []int {
	dc := map[int]int{}
	newarr := []int{}
	for index, value := range arr {
		if _, ok := dc[value]; ok == true {
			continue
		} else {
			dc[value] = index
			newarr = append(newarr, value)
		}
	}
	return newarr
}

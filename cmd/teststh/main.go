package main

import "fmt"

func main() {
	mySlice := make([]string, 0, 10)
	mySlice = append(mySlice, "a", "a", "a", "a", "a", "a", "a", "a", "a", "a")
	fmt.Println(mySlice)
	anotherSlice := mySlice[:5]
	anotherSlice = append(anotherSlice, "c", "c")
	anotherSlice[0] = "b"
	fmt.Println(anotherSlice)
	fmt.Println(mySlice)

	intArray := make([]int, 0, 10)
	intArray = append(intArray, 1, 2, 3, 4, 5, 6)

	notACopyofIntArray := intArray[:]
	notACopyofIntArray[0] = 12

	fmt.Println(intArray, notACopyofIntArray)

}

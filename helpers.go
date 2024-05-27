package main

import "fmt"

func assertEqual(expectedVal, gotVal any) bool {
	if expectedVal != gotVal {
		fmt.Printf("expected %v but got %v \n", expectedVal, gotVal)
		return false
	}
	return true
}

func assertComparison(title string, check bool) bool {
	if check {
		fmt.Printf("Comparison %s is false \n", title)
		panic("something went wrong check the log")
	}
	fmt.Println(title)
	return true
}

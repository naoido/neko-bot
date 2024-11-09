package errors

import "fmt"

func Catch(err error, errorMessage string) {
	if err != nil {
		fmt.Printf("ErrorMessage: %s\n", errorMessage)
		fmt.Printf("Error: %s\n", err)
	}
}

func CatchAndPanic(err error, errorMessage string) {
	if err != nil {
		fmt.Printf("ErrorMessage: %s\n", errorMessage)
		panic(err)
	}
}

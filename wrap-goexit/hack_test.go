package wrap_goexit

import (
	"fmt"
	"time"
)

func ExampleWrapGoexit() {
	go func() {
		WrapGoexit(func() {
			fmt.Println("haha")
		})
	}()

	time.Sleep(time.Second)
	// Output:
	// haha
}
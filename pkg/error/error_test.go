package error

import (
	"fmt"
)

func ExampleError_Error() {
	err := NewError("Error 1 occurred")
	fmt.Println(err)
	err = NewErrorC(NewError("Error 3"), "Error 2 occurred")
	fmt.Println(err)
	// Output:
	// Error 1 occurred
	// Error 2 occurred caused by: Error 3
}

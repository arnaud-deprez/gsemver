package version

import (
	"fmt"
)

func ExampleError_Error() {
	err := newError("Error 1 occurred")
	fmt.Println(err)
	err = newErrorC(newError("Error 3"), "Error 2 occurred")
	fmt.Println(err)
	// Output:
	// Error 1 occurred
	// Error 2 occurred caused by: Error 3
}

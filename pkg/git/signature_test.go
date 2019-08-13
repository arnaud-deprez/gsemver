package git

import (
	"fmt"
	"time"
)

func ExampleSignature_GoString() {
	s := &Signature{Name: "Arnaud Deprez", Email: "arnaudeprez@gmail.com", When: time.Date(2019, time.August, 12, 22, 47, 0, 0, time.UTC)}
	fmt.Printf("%#v\n", s)
	// Output: git.Signature{Name: "Arnaud Deprez", Email: "arnaudeprez@gmail.com", When: "2019-08-12 22:47:00 +0000 UTC"}
}

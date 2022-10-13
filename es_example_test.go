package xes

import (
	"fmt"
)

func ExampleNewEsLogger() {
	_, err := NewEsLogger()
	fmt.Println(err == ErrServersRequired)
	_, err = NewEsLogger(WithUrls(""))
	fmt.Println(err == ErrServerIsEmpty)

	urls := []string{"http://127.0.0.1:10000"}
	_, err = NewEsLogger(WithUrls(urls...))
	fmt.Println(err == ErrServerIsEmpty)
	fmt.Println(err == ErrIndexRequired)

	_, err = NewEsLogger(
		WithUrls(urls...),
		WithIndexName("test"),
		WithLevel(DebugLevel),
		WithUsername("xxx"),
		WithPassword("xxx"),
	)
	fmt.Println(err != nil)

	// Output:
	// true
	// true
	// false
	// false
	// true
}

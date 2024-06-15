package pkg

import (
	"fmt"
	"net/http"
	"testing"
)

func TestDefaultClient(t *testing.T) {
	fmt.Println(fmt.Sprintf("%p", http.DefaultClient))
	fmt.Println(fmt.Sprintf("%p", http.DefaultClient))
	fmt.Println(fmt.Sprintf("%p", http.DefaultClient))
	fmt.Println(fmt.Sprintf("%p", http.DefaultClient))
}

package iputil

import (
	"fmt"
	"testing"
)

func TestGetPublicIP(t *testing.T) {
	fmt.Println(GetPublicIP())
}

func TestGetPublicIPByHTTP(t *testing.T) {
	fmt.Println(GetPublicIPByHTTP())
}

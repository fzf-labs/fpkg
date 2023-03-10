package iputil

import (
	"fmt"
	"testing"
)

func TestGetPublicIP(t *testing.T) {
	fmt.Println(GetPublicIP())
}

func TestGetPublicIPByHttp(t *testing.T) {
	fmt.Println(GetPublicIPByHttp())
}

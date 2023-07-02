package utils

import (
	"fmt"
	"testing"
)

func TestGenToken(t *testing.T) {
	aToken, rToken, err := GenToken(1, "张三")
	if err != nil {
		return
	}
	fmt.Println(aToken, rToken)
}

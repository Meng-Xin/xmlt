package config

import (
	"fmt"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	config := InitLoadConfig()
	if config.Server.Addr == "" || config.Server.Port == "" {
		t.Error("Failed to read config")
	} else {
		fmt.Println(config.Server.Addr, config.Server.Port)
	}
}

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32
}
type Float interface {
	~float32 | ~float64
}

type Customize int

type Slice[T Int | Uint | Float] []T

func Add[T Int | Uint | Float](a, b T) T {
	return a + b
}

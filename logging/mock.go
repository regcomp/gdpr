package logging

import (
	"fmt"
)

type MockLogger struct{}

func (ml *MockLogger) Info(msg string, kv map[string]string) {
	kvs := ""
	for k, v := range kv {
		kvs += fmt.Sprintf("%s: %s, ", k, v)
	}
	fmt.Printf("[INFO] %s, %s\n", msg, kvs)
}

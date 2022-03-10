package util

import "fmt"

func Stringify(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

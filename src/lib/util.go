package lib

import (
	"fmt"
	"reflect"
)

func ExpectType[T any](input any) T {
	expectedType := reflect.TypeOf((*T)(nil)).Elem()
	recievedType := reflect.TypeOf(input)

	if expectedType != recievedType {
		panic(fmt.Sprintf("Expected type %T but recieved %T instead", expectedType, recievedType))
	}

	return input.(T)
}

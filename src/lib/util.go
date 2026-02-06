package lib

import (
	"fmt"
	"reflect"
)

func ExpectType[T any](input any) (T, error) {
	expectedType := reflect.TypeOf((*T)(nil)).Elem()
	recievedType := reflect.TypeOf(input)

	if expectedType != recievedType {
		return input.(T), fmt.Errorf("Expected type %T but recieved %T instead", expectedType, recievedType)
	}

	return input.(T), nil
}

package interpreter

import "fmt"

func createStdEnv() *env {
	scope := createEnv(nil)
	scope.set("print", std_print, true, false)
	return scope
}

func std_print(input ...FnCallArg) any {
	args := make([]any, 0)

	for _, arg := range input {
		args = append(args, arg.Value)
	}

	fmt.Print(args...)
	return nil
}

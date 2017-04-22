package mud

import "testing"

func TestOneArgument(t *testing.T) {
	argument := "this is a test"

	argument, arg1 := oneArgument(argument)

	if argument != "is a test" || arg1 != "this" {
		t.Error("Failed to unshift an argument")
	}

	argument, arg1 = oneArgument("")

	if argument != "" && arg1 != "" {
		t.Error("Failed to unshift an empty argument")
	}
}

// example from condition/README.md
package main

import (
	"context"
	"fmt"

	"github.com/brexhq/substation/condition"
)

func main() {
	inspector := condition.Strings{
		Key:        "foo",
		Function:   "equals",
		Expression: "bar",
	}

	data := []byte(`{"foo":"bar"}`)
	ok, err := condition.InspectByte(context.TODO(), data, inspector)
	if err != nil {
		panic(err)
	}

	if ok {
		fmt.Println("data passed inspection")
	} else {
		fmt.Println("data failed inspection")
	}
}
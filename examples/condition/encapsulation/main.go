// example of reading data from a file and applying an inspector
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/brexhq/substation/condition"
	"github.com/brexhq/substation/config"
)

func main() {
	// read lines from data file as encapsulated data
	open, err := os.Open("../data.json")
	if err != nil {
		panic(err)
	}

	var caps []config.Capsule
	cap := config.NewCapsule()

	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		cap.SetData(scanner.Bytes())
		caps = append(caps, cap)
	}

	// read config file and create a new inspector
	cfg, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	var sub config.Config
	if err := json.Unmarshal(cfg, &sub); err != nil {
		panic(err)
	}

	inspector, err := condition.InspectorFactory(sub)
	if err != nil {
		panic(err)
	}

	// apply inspector to encapsulated data
	for _, cap := range caps {
		ok, err := inspector.Inspect(context.TODO(), cap)
		if err != nil {
			panic(err)
		}

		if ok {
			fmt.Printf("passed inspection: %s\n", cap.GetData())
		} else {
			fmt.Printf("failed inspection: %s\n", cap.GetData())
		}
	}
}
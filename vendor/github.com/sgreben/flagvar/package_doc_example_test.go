package flagvar_test

import (
	"flag"
	"fmt"

	"github.com/sgreben/flagvar"
)

func Example() {
	var (
		fruit    = flagvar.Enum{Choices: []string{"apple", "banana"}}
		urls     flagvar.URLs
		settings flagvar.AssignmentsMap
	)

	fs := flag.FlagSet{}
	fs.Var(&fruit, "fruit", "a fruit")
	fs.Var(&urls, "url", "a URL")
	fs.Var(&settings, "set", "set key=value")
	fs.Parse([]string{
		"-fruit", "apple",
		"-url", "https://github.com/sgreben/flagvar",
		"-set", "hello=world",
	})

	fmt.Println("fruit:", fruit.Value)
	fmt.Println("urls:", urls.Values)
	for key, value := range settings.Values {
		fmt.Printf("settings: '%s' is set to '%s'\n", key, value)
	}

	// Output:
	// fruit: apple
	// urls: [https://github.com/sgreben/flagvar]
	// settings: 'hello' is set to 'world'
}

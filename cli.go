package main

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "convert [file]",
		Short: "Convert Go structs and interfaces to Cap'n Proto schema",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			goFile := args[0]
			schema := new(strings.Builder)
			err := Convert(goFile, schema)
			if err != nil {
				fmt.Println("Error converting Go to Cap'n Proto:", err)
				os.Exit(1)
			}

			fmt.Println(schema.String())
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

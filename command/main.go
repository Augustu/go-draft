package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "root",
		Short: "root cobra command",
		Long:  `cobra command test long explain`,

		RunE: func(cmd *cobra.Command, args []string) error {
			// something will happen

			fmt.Println("do something")

			// return
			return errors.New("new test err")
		},
	}
)

func init() {
	//
}

func main() {
	fmt.Println("hello there")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

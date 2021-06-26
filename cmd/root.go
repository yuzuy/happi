package cmd

import (
	"errors"
	"log"

	"github.com/spf13/cobra"

	"github.com/yuzuy/happi/editor"
)

var rootCmd = &cobra.Command{
	Use:   "happi",
	Short: "happi is a simple editor",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("happi can open only one file")
		}
		return execute(args[0])
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func execute(fileName string) error {
	log.Println("hoge")
	e, err := editor.Open(fileName)
	if err != nil {
		return err
	}
	defer e.Close()
	return e.Start()
}

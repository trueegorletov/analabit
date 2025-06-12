package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra" // Import cobra
)

func Run(rootCmd *cobra.Command) { // Accept rootCmd as a parameter
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Analabit CLI. Type 'help' for commands, or 'exit'/'quit' to leave.")

	// rootCmd is now passed as a parameter, no need for:
	// rootCmd := cmd.GetRootCmd()

	for {
		fmt.Print("analabit-cli> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			return
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" {
			fmt.Println("Exiting Analabit CLI.")
			break
		}

		// Parse the input into command and arguments for Cobra
		args := strings.Fields(input)

		// Cobra expects os.Args to be set for parsing, but we are in a shell.
		// We can use rootCmd.SetArgs() to pass arguments directly.
		rootCmd.SetArgs(args)

		// Execute the command. Cobra will find the appropriate subcommand.
		if err := rootCmd.Execute(); err != nil {
			// Cobra's Execute already prints errors by default to its Stderr.
			// If you need custom error printing, you can do it here.
			// fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}

package service

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

// Check if Font Awesome is available on the system
func FindFonts() error {
	cmd1 := exec.Command("fc-list")
	cmd2 := exec.Command("grep", "Awesome")
	cmd3 := exec.Command("sort")
	cmd2.Stdin, _ = cmd1.StdoutPipe()
	cmd3.Stdin, _ = cmd2.StdoutPipe()
	cmd3Output, _ := cmd3.StdoutPipe()

	if err := cmd3.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return err
	}
	if err := cmd2.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return err
	}
	if err := cmd1.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return err
	}

	cmd3Result, err := io.ReadAll(cmd3Output)
	if err != nil {
		log.Printf("Error reading command output: %v\n", err)
		return err
	}

	// Wait for all commands to finish
	_ = cmd1.Wait()
	_ = cmd2.Wait()
	_ = cmd3.Wait()

	// Print the final result
	fmt.Printf("Result:\n%s\n", cmd3Result)
	return nil
}

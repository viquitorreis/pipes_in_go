package main

import (
	"io"
	"os"
	"os/exec"
)

func main() {
	reader, writer := io.Pipe()
	cmd := exec.Command("cat")
	// set the commands stardard input to the pipe
	cmd.Stdin = reader
	cmd.Stdout = os.Stdout
	cmd.Start()

	writer.Write([]byte("Hello from parent\n"))
	writer.Close()

	cmd.Wait()
}

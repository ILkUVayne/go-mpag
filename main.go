package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	c := exec.Command("ffmpeg", "-h")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		log.Fatalf("failed to call cmd.Run(): %v", err)
	}
}

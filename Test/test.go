package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("python", "-u", "pulse.py")

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating StdoutPipe", err.Error())
	}
	defer stdOut.Close()

	scanner := bufio.NewScanner(stdOut)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t > %s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal("Error starting cmd", err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err.Error())
	}

}

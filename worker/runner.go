package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func runTarget(input, expectedOutput []string) (*Result, error) {
	res := &Result{}
	res.apiResponse.StartTime = time.Now().Unix()

	cmd := exec.Command("target.exe")

	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		return res, fmt.Errorf("error piping stdin from cmd: %w", err)
	}

	var output bytes.Buffer
	cmd.Stderr = os.Stderr
	cmd.Stdout = &output

	if err := cmd.Start(); err != nil {
		return res, fmt.Errorf("error starting command: %w", err)
	}

	for _, line := range input {
		if _, err := cmdStdin.Write([]byte(line + "\r\n\n")); err != nil {
			return res, fmt.Errorf("error writing %q to cmd: %w", line, err)
		}
	}

	if err := cmdStdin.Close(); err != nil {
		return res, fmt.Errorf("error closing cmd stdin pipe: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return res, fmt.Errorf("error waiting for cmd: %w", err)
	}

	scan := bufio.NewScanner(&output)
	for i := 0; scan.Scan(); i++ {
		str := scan.Text()
		if str != expectedOutput[i] {
			res.apiResponse.Errors = append(res.apiResponse.Errors, fmt.Sprintf("line %d: expected %q but received %q!", i, expectedOutput[i], str))
			log.Printf("output mismatch, expected %q and received %q", expectedOutput[i], str)
		}
		res.apiResponse.Log = append(res.apiResponse.Log, str)
	}

	res.apiResponse.EndTime = time.Now().Unix()
	res.apiResponse.Took = res.apiResponse.EndTime - res.apiResponse.StartTime

	log.Printf("ran in %s!", time.Duration(res.apiResponse.Took))

	return res, nil
}

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

func runTarget() (*Result, error) {
	res := &Result{}
	res.apiResponse.StartTime = time.Now().Unix()

	grepCmd := exec.Command("target.exe")

	grepIn, err := grepCmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error piping stdin from cmd: %w", err)
	}

	var output bytes.Buffer
	grepCmd.Stderr = os.Stderr
	grepCmd.Stdout = &output

	if err := grepCmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting command: %w", err)
	}

	for _, line := range lines {
		if _, err := grepIn.Write([]byte(line + "\r\n\n")); err != nil {
			return nil, fmt.Errorf("error writing %q to cmd: %w", line, err)
		}
	}

	if err := grepIn.Close(); err != nil {
		return nil, fmt.Errorf("error closing cmd stdin pipe: %w", err)
	}

	if err := grepCmd.Wait(); err != nil {
		return nil, fmt.Errorf("error waiting for cmd: %w", err)
	}

	scan := bufio.NewScanner(&output)
	i := 0
	for scan.Scan() {
		str := scan.Text()
		if str != expected[i] {
			res.apiResponse.Errors = append(res.apiResponse.Errors, fmt.Sprintf("line %d: expected %q but received %q!", i, expected[i], str))
			log.Printf("output mismatch, expected %q and received %q", expected[i], str)
		}
		res.apiResponse.Log = append(res.apiResponse.Log, str)
		i++
	}

	res.apiResponse.EndTime = time.Now().Unix()
	res.apiResponse.Took = res.apiResponse.EndTime - res.apiResponse.StartTime

	log.Printf("ran in %s!", time.Duration(res.apiResponse.Took))

	return res, nil
}

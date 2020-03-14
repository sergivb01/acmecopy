package main

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sergivb01/acmecopy/api"
)

type Result struct {
	// Log       []string
	// Errors    []error
	// StartTime time.Time
	// EndTime   time.Time
	// Took      time.Duration
	errChan     chan error
	outChan     chan string
	apiResponse api.Response
	m           *sync.Mutex
}

func compileSingleFile(file string, wg *sync.WaitGroup) ([]byte, error) {
	cmd := exec.Command("g++", "-c", file)
	defer wg.Done()
	return cmd.CombinedOutput()
}

func (r *Result) listenForChannels() {
	select {
	case err := <-r.errChan:
		r.m.Lock()
		r.apiResponse.Errors = append(r.apiResponse.Errors, err.Error())
		r.m.Unlock()
		log.Printf("error compiling file: %s", err)
		break
	case buildOut := <-r.outChan:
		if strings.TrimSpace(buildOut) != "" {
			r.m.Lock()
			r.apiResponse.Log = append(r.apiResponse.Log, buildOut)
			r.m.Unlock()
			log.Printf("received from build log: %q", buildOut)
		}
		break
	}
}

func compileFiles(files []*api.File) (*Result, error) {
	res := &Result{
		errChan: make(chan error),
		outChan: make(chan string),
		m:       &sync.Mutex{},
	}

	res.apiResponse.StartTime = time.Now().Unix()

	var fileNames []string
	for _, file := range files {
		if filepath.Ext(file.FileName) == ".cpp" {
			fileNames = append(fileNames, file.FileName)
		}
	}

	var args = []string{"-std=c++11"}

	var wg sync.WaitGroup
	wg.Add(len(fileNames))

	go res.listenForChannels()

	for _, file := range fileNames {
		go func(file string) {
			b, err := compileSingleFile(file, &wg)
			if err != nil {
				res.errChan <- err
				return
			}
			res.outChan <- string(b)
		}(file)
		args = append(args, file[0:len(file)-len(filepath.Ext(file))]+".o")
	}
	wg.Wait()

	b, err := exec.Command("g++", append(args, "-o", "target.exe")...).CombinedOutput()
	if err != nil {
		res.apiResponse.Errors = append(res.apiResponse.Errors, err.Error())
	}
	res.apiResponse.Log = append(res.apiResponse.Log, string(b))
	log.Printf("received from final build: %q\n", string(b))

	res.apiResponse.EndTime = time.Now().Unix()
	res.apiResponse.Took = res.apiResponse.EndTime - res.apiResponse.StartTime

	log.Printf("built %d files in %s!", len(files), time.Duration(res.apiResponse.Took))

	return res, nil
}

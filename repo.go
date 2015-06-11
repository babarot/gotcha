package main

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"sync"
)

func Update(c Config) (<-chan bool, <-chan string, <-chan string) {
	doneCh := make(chan bool)
	outCh, errCh := make(chan string), make(chan string)
	semaphore := make(chan int, runtime.NumCPU())

	var wg sync.WaitGroup

	go func() {
		for _, repo := range c.Repos {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				semaphore <- 1
				outCh <- fmt.Sprintf("--> Updating: %s\n", url)
				if err := run("go", "get", "-u", url); err != nil {
					errCh <- fmt.Sprintf("    error `%s': %s\n", url, err)
				}
				<-semaphore
			}(repo)
		}
		wg.Wait()
		doneCh <- true
	}()

	return doneCh, outCh, errCh
}

func run(args ...string) error {
	if len(args) == 0 {
		return errors.New("too few arguments")
	}
	cmd := exec.Command(args[0], args[1:]...)
	return cmd.Run()
}

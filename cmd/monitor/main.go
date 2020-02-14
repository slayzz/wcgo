package main

import (
	"fmt"
	"github.com/slayzz/wcgo/pkg/wc"
	"os"
	"sync"
)

func main() {
	if len(os.Args) == 1 {
		c := wc.NewCounter(wc.ReaderConfig{
			Reader: os.Stdin,
			Name:   "",
		})
		c.ExecuteCount()
		fmt.Println(c)
		os.Exit(0)
	}

	chm := wc.NewCounterMonitor()
	go chm.Monitor()

	filesArgs := os.Args[1:]

	var wg sync.WaitGroup
	for _, filename := range filesArgs {
		wg.Add(1)
		go func(f string) {
			c, err := wc.CountWC(f)
			if err != nil {
				fmt.Println(err)
				return
			}
			chm.Insert(c)
			defer wg.Done()
		}(filename)
	}

	wg.Wait()
	cl := chm.Read()
	fmt.Print(cl)
}

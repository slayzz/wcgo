package main

import (
	"fmt"
	"os"
	"tutorial/wc"
)

func countFromFiles(filename string, counterChan chan<- *wc.Counter) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	c := wc.NewCounter(wc.ReaderConfig{
		Reader: file,
		Name:   file.Name(),
	})
	c.ExecuteCount()
	counterChan <- &c
}

func aggregate(fileNames []string, counterChan <-chan *wc.Counter) chan wc.CounterList {
	listCounter := make(chan wc.CounterList)
	go func(f []string) {
		var cl wc.CounterList
		for range f {
			c := <-counterChan
			cl = append(cl, c)
		}
		listCounter <- cl
	}(fileNames)

	return listCounter
}

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

	filesArgs := os.Args[1:]
	counterChan := make(chan *wc.Counter, len(filesArgs))
	for _, filename := range filesArgs {
		go countFromFiles(filename, counterChan)
	}

	listCounter := aggregate(filesArgs, counterChan)

	cl := <-listCounter
	fmt.Print(cl)
}


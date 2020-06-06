package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"
)

const (
	requestTimeout = 5
	maxGoroutines  = 5
	searchWord     = "go"
)

func main() {
	var (
		totalEntries uint
		mutex        sync.Mutex
		wg           sync.WaitGroup
	)
	guard := make(chan struct{}, maxGoroutines)

	SetupCloseHandler(&totalEntries)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		wg.Add(1)
		url := s.Text()
		guard <- struct{}{}
		go func(url string, mutex *sync.Mutex) {
			defer wg.Done()

			numberOfEntries := getNumberOfEntriesFromUrl(url)
			mutex.Lock()
			totalEntries += numberOfEntries
			mutex.Unlock()

			<-guard
		}(url, &mutex)
	}

	wg.Wait()
	os.Exit(printTotalEntries(totalEntries))
}

func SetupCloseHandler(totalEntries *uint) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(printTotalEntries(*totalEntries))
	}()
}

func printTotalEntries(totalEntries uint) int {
	fmt.Printf("\nTotal: %v\n", totalEntries)
	return 0
}

func getNumberOfEntriesFromUrl(url string) uint {
	body := requestUrl(url)
	words := getWordsFromText(body)
	numberOfEntries := getNumberOfEntries(words, searchWord)
	fmt.Printf("Count for %v: %v\n", url, numberOfEntries)
	return numberOfEntries
}

func requestUrl(url string) (bodyString string) {
	client := http.Client{
		Timeout: requestTimeout * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		bodyString = string(bodyBytes)
	}
	return
}

func getWordsFromText(text string) []string {
	words := regexp.MustCompile("\\w+")
	return words.FindAllString(text, -1)
}

func getNumberOfEntries(words []string, searchKey string) (wordCounts uint) {
	for _, word := range words {
		if word == searchKey {
			wordCounts++
		}
	}
	return
}

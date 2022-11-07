package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"sync"
)

func main() {
	var wgroup sync.WaitGroup
	fmt.Println("input number of websites! type \"done\" if finish ")

	var website string
	var websites []string
	for {
		fmt.Scan(&website)
		if website == "done" {
			getRequestUrl(&wgroup, websites)
			break
		}

		u, err := url.ParseRequestURI(website)
		if err != nil {
			fmt.Println("please try again wrong website format include \"https://\" ")
			break
		}

		websites = append(websites, "https://"+u.Host)
	}

}

func getRequestUrl(wgroup *sync.WaitGroup, websites []string) {
	websiteCount := make(map[string]int)
	wgroup.Add(len(websites))

	for _, e := range websites {

		go func(e string) {
			websiteCount[e] = getSizeUrl(e)
			wgroup.Done()
		}(e)
	}
	wgroup.Wait()
	fmt.Println("Website will output in a descending order according to body size:")
	keys := make([]string, 0, len(websiteCount))
	for k := range websiteCount {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return websiteCount[keys[i]] > websiteCount[keys[j]]
	})

	for _, k := range keys {
		fmt.Println(k, websiteCount[k])
	}

}

func getSizeUrl(w string) int {

	resp, err := http.Get(w)

	if err != nil {
		fmt.Errorf("Failed to request website: %s", err)
	}
	if resp.Body == nil {
		panic("body is nil")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return len(b)

}

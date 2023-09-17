package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()

	c.OnRequest(func(request *colly.Request) {
		//request.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36\n")
		fmt.Printf("Visiting: %v\n", request.URL)
	})

	c.OnResponse(func(response *colly.Response) {
		fmt.Println("Response Code: \n", response.StatusCode)
	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Println("Error: ", err.Error())
	})

	err := c.Visit("http://web.mit.edu/jywang/www/cef/Bible/NIV/NIV_Bible/bookindex.html")
	if err != nil {
		fmt.Println("Could not visit site:\n", err)
		return
	}
}

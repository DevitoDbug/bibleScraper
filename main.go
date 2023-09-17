package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

type Data struct {
	Chapter string
	Content []VerseDataPair
}

type VerseDataPair struct {
	Verse     string
	VerseData string
}

func main() {
	c := colly.NewCollector()

	c.OnRequest(func(request *colly.Request) {
		fmt.Printf("Visiting: %v\n", request.URL)
	})

	c.OnResponse(func(response *colly.Response) {
		fmt.Println("Response Code: \n", response.StatusCode)
	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Println("Error: ", err.Error())
	})

	var data []Data
	c.OnHTML("a[target=display][href='GEN+1.html']", func(element *colly.HTMLElement) {
		//text := element.Text
		href := element.Attr("href")
		err := c.Visit(element.Request.AbsoluteURL(href))
		if err != nil {
			fmt.Println("Could not visit verse:\n", err)
			return
		}

	})
	c.OnHTML("dl", func(e *colly.HTMLElement) {
		fmt.Println("hello")
		dtElements := e.DOM.Find("dt")
		ddElements := e.DOM.Find("dd")

		// Checking if we have the same number of verses and verses data
		if ddElements.Length() != dtElements.Length() {
			fmt.Println("Mismatch of verses and their data")
			return
		}

		var content []VerseDataPair

		for i := 0; i < dtElements.Length(); i++ {
			verseNum := dtElements.Eq(i).Text()
			verseData := ddElements.Eq(i).Text()
			content = append(content, struct {
				Verse     string
				VerseData string
			}{Verse: verseNum, VerseData: verseData})
		}

		newData := Data{
			Chapter: "text",
			Content: content,
		}

		data = append(data, newData)
	})

	err := c.Visit("http://web.mit.edu/jywang/www/cef/Bible/NIV/NIV_Bible/bookindex.html")
	if err != nil {
		fmt.Println("Could not visit site:\n", err)
		return
	}

	// Print or process the 'data' slice as needed
	fmt.Println(data)
}

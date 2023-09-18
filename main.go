package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

type Data struct {
	Chapter string          `json:"chapter"`
	Content []VerseDataPair `json:"content"`
}

type VerseDataPair struct {
	Verse     string `json:"verse"`
	VerseData string `json:"verseData"`
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
	var book string
	c.OnHTML("a[target=display]", func(element *colly.HTMLElement) {
		href := element.Attr("href")
		err := c.Visit(element.Request.AbsoluteURL(href))
		if err != nil {
			fmt.Println("Could not visit verse:\n", err)
			return
		}
	})

	c.OnHTML("b", func(element2 *colly.HTMLElement) {
		book = element2.Text
	})

	c.OnHTML("dl", func(e *colly.HTMLElement) {
		dtElements := e.DOM.Find("dt")
		ddElements := e.DOM.Find("dd")

		fmt.Println(book)
		fmt.Println(dtElements.Length())
		fmt.Println(ddElements.Length())
		fmt.Print("\n\n\n")

		// Checking if we have the same number of verses and verses data
		if ddElements.Length() != dtElements.Length() {
			fmt.Println("Mismatch of verses and their data")
			return
		}

		var content []VerseDataPair

		for i := 0; i < dtElements.Length(); i++ {
			verseNum := dtElements.Eq(i).Text()
			verseData := ddElements.Eq(i).Text()

			content = append(content, VerseDataPair{
				Verse:     verseNum,
				VerseData: verseData,
			})
		}
		newData := Data{
			Chapter: book,
			Content: content,
		}

		data = append(data, newData)
	})

	err := c.Visit("http://web.mit.edu/jywang/www/cef/Bible/NIV/NIV_Bible/bookindex.html")
	if err != nil {
		fmt.Println("Could not visit site:\n", err)
		return
	}

	file, err := os.Create("output.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic("Could not close file")
		}
	}(file)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("Could not marshal data to JSON:", err)
		return
	}

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatal("Could not write JSON data to file:", err)
		return
	}

	fmt.Println("Data written to output.json")

}

package main

import (
	"fmt"
	"net/url"
)

func main() {
	urlStr := "https://example.com/?word=value"
	myUrl, _ := url.Parse(urlStr)
	params, _ := url.ParseQuery(myUrl.RawQuery)
	//fmt.Println(params)

	word := params.Get("word")
	fmt.Println(product)
	

}

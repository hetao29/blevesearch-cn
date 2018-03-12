package main

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve"
	_ "github.com/hetao29/blevesearch-cn/scws/bleve"
)

type Data struct {
	Name string
	Id   int
	T    string
}

func main() {

	// open a new index
	mapping := bleve.NewIndexMapping()

	err := mapping.AddCustomTokenizer("scws",
		map[string]interface{}{
			"dict": "/Users/hetal/dict/dict.utf8.xdb",
			//"rule": "/Users/hetal/dict/rules.utf8.ini",
			"type":"scws",
		},
	)
	if err != nil {
		panic(err)
	}

	err = mapping.AddCustomAnalyzer("scws",
		map[string]interface{}{
			"type":      "scws",
			"tokenizer": "scws",
		},
	)
	if err != nil {
		panic(err)
	}
	mapping.DefaultAnalyzer = "scws"

	index, err := bleve.New("example.bleve", mapping)
	if err != nil {

		index, err = bleve.Open("example.bleve")
		if err != nil {
			fmt.Println(err)
			return
		}

	}

	data := Data{
		Name: "text",
		Id:   333,
		T:    "2019-10-10 12:20:30",
	}
	data2 := Data{
		Name: "有效解决知识点撑握不牢，缺少学习方法，解题速度慢，粗心大意反复丢分问题张亚祥",
		Id:   99,
		T:    "2019-10-10 12:20:30",
	}

	// index some data
	index.Index("id", data)
	index.Index("id2", data2)

	// search for some text
	//query := bleve.NewMatchQuery("有效")
	//query := bleve.NewQueryStringQuery("+Name:\"心大\" created:>\"2010-10-10 00:00:00\"")
	//query := bleve.NewQueryStringQuery("T:>\"2010-10-10 00:00:00\"")
	//query := bleve.NewQueryStringQuery("Id:>199")
	//query := bleve.NewQueryStringQuery("Id:>199")
	//query := bleve.NewQueryStringQuery("+Name:\"粗心\" created:>\"2010-10-10 00:00:00\"")
	query := bleve.NewQueryStringQuery("+Name:亚祥 created:>\"2010-10-10 00:00:00\"")
	//query := bleve.NewQueryStringQuery("+Name:心大 created:>\"2010-10-10 00:00:00\"")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	//search.Highlight = bleve.NewHighlight()
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)

	data3, _ := json.Marshal(searchResults)
	fmt.Println(string(data3))
}

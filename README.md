# blevesearch-cn
blevesearch support chinese language

## install scws
```bash
$ wget http://www.xunsearch.com/scws/down/scws-1.2.1.tar.bz2 
$ tar xjf scws-1.2.1.tar.bz2 
$ cd scws-1.2.1
$ ./configure 
$ make 
$ make install
```

## download
```bash
go get github.com/blevesearch/bleve
go get github.com/hetao29/blevesearch-cn
```

## usage

```golang
package main

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve"
	//"github.com/blevesearch/bleve/registry"
	cn "github.com/hetao29/blevesearch-cn/scws/bleve"
)

type Data struct {
	Name string
	Id   int
	T    string
}
type Doc struct{
	Id string `json:id`
	Doc interface{} `json:doc`
}

func main() {


	var requestBody=`[{"id":"xx","doc":{"x":"b"}},{"id":"xx2","doc":{"x":"b"}}]`;
	//var doc interface{}
	var doc []Doc
	err := json.Unmarshal([]byte(requestBody), &doc)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(doc) > 0{

		for _, d:= range doc{
			fmt.Println(d.Id);
			fmt.Println(d.Doc);
		}
	}
	ct, _ := json.Marshal(doc)
	fmt.Println(string(ct));
	return

	cn.SetDict("/Users/hetal/dict/dict.utf8.xdb");
	cn.SetRule("/Users/hetal/dict/rules.utf8.ini");
	//types, instance := registry.AnalyzerTypesAndInstances();
	//fmt.Println(types);
	//fmt.Println(instance);
	// open a new index
	mapping := bleve.NewIndexMapping()

	/*
	err := mapping.AddCustomTokenizer("cn",
		map[string]interface{}{
			"dict": "/Users/hetal/dict/dict.utf8.xdb",
			//"rule": "/Users/hetal/dict/rules.utf8.ini",
			"type":"cn",
		},
	)
	if err != nil {
		panic(err)
	}
	err = mapping.AddCustomAnalyzer("cn",
		map[string]interface{}{
			"type":      "cn",
			"tokenizer": "cn",
		},
	)
	if err != nil {
		panic(err)
	}
	*/
	mapping.DefaultAnalyzer = "cn"

	ct, _ = json.Marshal(mapping)
	fmt.Println(string(ct));
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
	batch := index.NewBatch();
	batch.Index("id",data);
	batch.Index("id2",data2);
	//index.Batch(batch);
	//index.Index("id", data)
	//index.Index("id2", data2)

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
```
```bash
go run main.go
```

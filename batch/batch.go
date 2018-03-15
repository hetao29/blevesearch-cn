//  Copyright (c) 2014 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package batch

import (
	"encoding/json"
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"net/http"

	bleveHttp "github.com/blevesearch/bleve/http"
)
type Doc struct{
    Id string `json:id`
    Doc interface{} `json:doc`
}
func showError(w http.ResponseWriter, r *http.Request,
	msg string, code int) {
	logger.Printf("Reporting error %v/%v", code, msg)
	http.Error(w, msg, code)
}

func mustEncode(w io.Writer, i interface{}) {
	if headered, ok := w.(http.ResponseWriter); ok {
		headered.Header().Set("Cache-Control", "no-cache")
		headered.Header().Set("Content-type", "application/json")
	}

	e := json.NewEncoder(w)
	if err := e.Encode(i); err != nil {
		panic(err)
	}
}

type varLookupFunc func(req *http.Request) string

var logger = log.New(ioutil.Discard, "bleve.http", log.LstdFlags)

// SetLog sets the logger used for logging
// by default log messages are sent to ioutil.Discard
func SetLog(l *log.Logger) {
	logger = l
}

type DocBatchHandler struct {
	defaultBatchName string
	IndexNameLookup varLookupFunc
	DocIDLookup      varLookupFunc
}

func NewDocBatchHandler(defaultBatchName string) *DocBatchHandler {
	return &DocBatchHandler{
		defaultBatchName: defaultBatchName,
	}
}

func (h *DocBatchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// find the index to operate on
	var indexName string
	if h.IndexNameLookup!= nil {
		indexName = h.IndexNameLookup(req)
	}
	if indexName == "" {
		indexName = h.defaultBatchName
	}
	index := bleveHttp.IndexByName(indexName)
	if index == nil {
		showError(w, req, fmt.Sprintf("no such index '%s'", indexName), 404)
		return
	}

	// read the request body
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		showError(w, req, fmt.Sprintf("error reading request body: %v", err), 400)
		return
	}

	// parse request body as json
	//var doc interface{}
	var docs []Doc
	err = json.Unmarshal(requestBody, &docs)
	if err != nil{
		showError(w, req, fmt.Sprintf("error parsing request body as JSON: %v", err), 400)
		return
	}

	if len(docs)==0{
		showError(w, req, fmt.Sprintf("request body is empty"), 400)
		return
	}
	batch := index.NewBatch();
	for _, doc:= range docs{
		batch.Index(doc.Id,doc.Doc)
	}


	err = index.Batch(batch)
	if err != nil {
		showError(w, req, fmt.Sprintf("error indexing document '%s': %v", indexName, err), 500)
		return
	}

	rv := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	mustEncode(w, rv)
}

package bleve

import (
	"errors"
	"runtime"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
	"scws"
)

type ScwsTokenizer struct {
	handle *scws.Scws
}

func NewScwsTokenizer(dict, rule string) *ScwsTokenizer {
	x := scws.NewScws()
	err := x.SetDict(dict, scws.SCWS_XDICT_XDB)
	if err != nil {
		return nil
	}
	if rule != "" {
		err = x.SetRule(rule)
		if err != nil {
			return nil
		}
	}

	x.SetCharset("utf8")
	x.SetIgnore(1)
	x.SetMulti(scws.SCWS_MULTI_SHORT | scws.SCWS_MULTI_DUALITY)

	x.Init(runtime.NumCPU())
	return &ScwsTokenizer{x}
}

func (x *ScwsTokenizer) Free() {
	x.handle.Free()
}

func (x *ScwsTokenizer) Tokenize(sentence []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0)
	pos := 1
	words, err := x.handle.Segment(string(sentence))
	if err != nil {
		panic(err)
		return nil

	}
	for _, word := range words {
		token := analysis.Token{
			Term:     []byte(word.Term),
			Start:    word.Start,
			End:      word.End,
			Position: pos,
			Type:     analysis.Ideographic,
		}
		result = append(result, &token)
		pos++
	}
	return result
}

func tokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	dict, ok := config["dict"].(string)
	if !ok {
		return nil, errors.New("config dictpath not found")
	}
	rule, _ := config["rule"].(string)

	return NewScwsTokenizer(dict, rule), nil
}

func init() {
	registry.RegisterTokenizer("scws", tokenizerConstructor)
}

package blevejieba

import (
	"regexp"
	"strconv"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
	"github.com/yanyiwu/gojieba"
)

var ideographRegexp = regexp.MustCompile(`\p{Han}+`)

// JiebaTokenizer is the beleve tokenizer for jiebago.
type JiebaTokenizer struct {
	jieba      *gojieba.Jieba
	searchMode gojieba.TokenizeMode
	useHmm     bool
}

func NewJiebaTokenizer(dictFilePath, hmm, userDictPath, idfDict, stopDict string, searchMode bool) (analysis.Tokenizer, error) {
	jieba := gojieba.NewJieba(dictFilePath, hmm, userDictPath, idfDict, stopDict)
	mode := gojieba.DefaultMode
	if searchMode {
		mode = gojieba.SearchMode
	}
	return &JiebaTokenizer{
		jieba:      jieba,
		searchMode: mode,
	}, nil
}

// Tokenize cuts input into bleve token stream.
func (jt *JiebaTokenizer) Tokenize(input []byte) analysis.TokenStream {
	rv := make(analysis.TokenStream, 0)
	pos := 1

	for _, word := range jt.jieba.Tokenize(string(input), jt.searchMode, true) {
		token := analysis.Token{
			Term:     []byte(word.Str),
			Start:    word.Start,
			End:      word.End,
			Position: pos,
			Type:     detectTokenType(word.Str),
		}
		rv = append(rv, &token)
		pos++
	}
	return rv
}

/*
JiebaTokenizerConstructor creates a JiebaTokenizer.
Parameter config can contains following parameter:
    dict_path: optional, the path of the dictionary file.
    hmm_path: optional, specify whether to use Hidden Markov Model, see NewJiebaTokenizer for details.
    userdict_path: optional, specify user dict file path
    idf_path: optional, specify idf file path
    stopdict_path: optional, specify user stop dict file path
    is_search: optional, speficy whether to use isSearch mode, see NewJiebaTokenizer for details.
*/
func JiebaTokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (
	analysis.Tokenizer, error) {
	dictFilePath, ok := config[dictPathKey].(string)
	if !ok {
		dictFilePath = gojieba.DICT_PATH
	}
	hmm, ok := config[hmmPathKey].(string)
	if !ok {
		hmm = gojieba.HMM_PATH
	}
	userDictPath, ok := config[userDictPathKey].(string)
	if !ok {
		userDictPath = gojieba.USER_DICT_PATH
	}
	stopDict, ok := config[stopDictPathKey].(string)
	if !ok {
		stopDict = gojieba.STOP_WORDS_PATH
	}
	idfDict, ok := config[idfDictPathKey].(string)
	if !ok {
		idfDict = gojieba.IDF_PATH
	}
	searchMode, ok := config[isSearchKey].(bool)
	if !ok {
		searchMode = true
	}

	return NewJiebaTokenizer(dictFilePath, hmm, userDictPath, idfDict, stopDict, searchMode)
}

func detectTokenType(term string) analysis.TokenType {
	if ideographRegexp.MatchString(term) {
		return analysis.Ideographic
	}
	_, err := strconv.ParseFloat(term, 64)
	if err == nil {
		return analysis.Numeric
	}
	return analysis.AlphaNumeric
}

func init() {
	registry.RegisterTokenizer(Name, JiebaTokenizerConstructor)
}

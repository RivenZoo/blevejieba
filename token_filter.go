package blevejieba

import (
	"io/ioutil"

	"github.com/blevesearch/bleve/analysis/token/stop"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
	"github.com/yanyiwu/gojieba"
)

func StopTokenFilterConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.TokenFilter, error) {
	tokenMap, err := cache.TokenMapNamed(Name)
	if err != nil {
		return nil, err
	}
	return stop.NewStopTokensFilter(tokenMap), nil
}

// TokenMapConstructor create a stop word token map.
// Parameter config can contains following parameters:
//   stopdict_path: optional, user stop dict file path
func TokenMapConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.TokenMap, error) {
	stopDict, ok := config[stopDictPathKey].(string)
	if !ok {
		stopDict = gojieba.STOP_WORDS_PATH
	}
	rv := analysis.NewTokenMap()
	words, err := ioutil.ReadFile(stopDict)
	if err != nil {
		return nil, err
	}
	err = rv.LoadBytes(words)
	return rv, err
}

func init() {
	registry.RegisterTokenMap(Name, TokenMapConstructor)
	registry.RegisterTokenFilter(Name, StopTokenFilterConstructor)
}

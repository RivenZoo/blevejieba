package blevejieba

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/yanyiwu/gojieba"
)

type Options struct {
	dictPath, hmmPath, userDictPath, idfPath, stopDictPath string
	isSearch                                               bool
}

func NewOptions() *Options {
	return &Options{
		dictPath:     gojieba.DICT_PATH,
		hmmPath:      gojieba.HMM_PATH,
		userDictPath: gojieba.USER_DICT_PATH,
		idfPath:      gojieba.IDF_PATH,
		stopDictPath: gojieba.STOP_WORDS_PATH,
		isSearch:     false,
	}
}

func (o *Options) WithJiebaDictPath(p string) *Options {
	o.dictPath = p
	return o
}

func (o *Options) WithHMMPath(p string) *Options {
	o.hmmPath = p
	return o
}

func (o *Options) WithUserDictPath(p string) *Options {
	o.userDictPath = p
	return o
}

func (o *Options) WithIDFDictPath(p string) *Options {
	o.idfPath = p
	return o
}

func (o *Options) WithStopDictPath(p string) *Options {
	o.stopDictPath = p
	return o
}

func (o *Options) WithSearch(search bool) *Options {
	o.isSearch = search
	return o
}

func (o *Options) customTokenizerConfig() map[string]interface{} {
	return map[string]interface{}{
		"type":          Name,
		isSearchKey:     o.isSearch,
		dictPathKey:     o.dictPath,
		hmmPathKey:      o.hmmPath,
		userDictPathKey: o.userDictPath,
		idfDictPathKey:  o.idfPath,
		stopDictPathKey: o.stopDictPath,
	}
}

func (o *Options) customAnalyzerConfig() map[string]interface{} {
	return map[string]interface{}{
		"type": Name,
	}
}

func (o *Options) customTokenMapConfig() map[string]interface{} {
	return map[string]interface{}{
		"type":          Name,
		stopDictPathKey: o.stopDictPath,
	}
}

func NewMemIndexWithGoJieba(opt *Options) (bleve.Index, error) {
	mappingImpl, err := NewGoJiebaIndexMapping(opt)
	if err != nil {
		return nil, err
	}
	index, err := bleve.NewMemOnly(mappingImpl)
	if err != nil {
		return nil, err
	}
	return index, nil
}

func NewGoJiebaIndexMapping(opt *Options) (mapping.IndexMapping, error) {
	indexMapping := bleve.NewIndexMapping()
	err := indexMapping.AddCustomTokenizer(Name,
		opt.customTokenizerConfig(),
	)
	if err != nil {
		return nil, err
	}

	err = indexMapping.AddCustomTokenMap(Name, opt.customTokenMapConfig())
	if err != nil {
		return nil, err
	}

	err = indexMapping.AddCustomAnalyzer(Name,
		opt.customAnalyzerConfig(),
	)
	if err != nil {
		return nil, err
	}
	indexMapping.DefaultAnalyzer = Name
	return indexMapping, nil
}

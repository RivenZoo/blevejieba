/*
* blevejieba是一个bleve的中文分词插件，基于gojieba开发
 */
package blevejieba

// Name is the jieba analyzer/tokenizer name.
const Name = "jieba"

const (
	dictPathKey     = "dict_path"
	hmmPathKey      = "hmm_path"
	userDictPathKey = "userdict_path"
	idfDictPathKey  = "idf_path"
	stopDictPathKey = "stopdict_path"
	isSearchKey     = "is_search"
)

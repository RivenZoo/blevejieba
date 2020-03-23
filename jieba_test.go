package blevejieba

import (
	"reflect"
	"sort"
	"testing"

	"github.com/blevesearch/bleve"
)

type indexItem struct {
	ID   string `json:"id"`
	Desc string `json:"desc"`
}

func Test_jieba(t *testing.T) {
	indexItems := []indexItem{
		indexItem{
			ID:   "1",
			Desc: "研制美味佳肴，收集生活百感",
		},
		indexItem{
			ID:   "2",
			Desc: "女生爱玩的游戏",
		},
	}
	mapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		t.FailNow()
	}
	jiebaIndex, err := NewMemIndexWithGoJieba(NewOptions().WithSearch(true))
	if err != nil {
		t.FailNow()
	}
	tests := []struct {
		name    string
		index   bleve.Index
		items   []indexItem
		query   string
		wantIds []string
	}{
		{
			name:    "raw bleve index",
			index:   index,
			items:   indexItems,
			query:   "女生",
			wantIds: []string{"1", "2"},
		},
		{
			name:    "gojieba bleve index",
			index:   jiebaIndex,
			items:   indexItems,
			query:   "女生",
			wantIds: []string{"2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.items {
				tt.index.Index(indexItems[i].ID, indexItems[i])
			}
			q := bleve.NewQueryStringQuery(tt.query)
			req := bleve.NewSearchRequest(q)
			result, err := tt.index.Search(req)
			if err != nil {
				t.Fatal(err)
			}

			if result.Total != uint64(len(tt.wantIds)) {
				t.Fatalf("expect %v got %v", tt.wantIds, result.Hits)
			}
			var ids []string
			for _, hit := range result.Hits {
				ids = append(ids, hit.ID)
			}
			sort.Strings(ids)
			sort.Strings(tt.wantIds)
			if !reflect.DeepEqual(ids, tt.wantIds) {
				t.Fatalf("expect %v got %v", tt.wantIds, ids)
			}
		})
	}
}

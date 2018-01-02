package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/esdb/lstore"
)

func Test_head_segment(t *testing.T) {
	should := require.New(t)
	store := smallTestStore()
	defer store.Stop(ctx)
	for i := 0; i < 16; i++ {
		blobValue := lstore.Blob("hello")
		if i%2 == 0 {
			blobValue = lstore.Blob("world")
		}
		_, err := store.Write(ctx, intBlobEntry(int64(i)+1, blobValue))
		should.Nil(err)
	}
	should.Nil(store.Index())
	reader, err := store.NewReader(ctx)
	should.Nil(err)
	collector := &lstore.ResultCollector{LimitSize: 2}
	reader.SearchForward(ctx, 0, []lstore.Filter{
		store.IndexingStrategy.NewBlobValueFilter(0, "hello"),
	}, collector)
	should.Equal([]int64{2}, collector.Rows[0].IntValues)
	should.Equal([]int64{4}, collector.Rows[1].IntValues)
}

func Test_reopen_head_segment(t *testing.T) {
	should := require.New(t)
	store := smallTestStore()
	defer store.Stop(ctx)
	for i := 0; i < 16; i++ {
		blobValue := lstore.Blob("hello")
		if i%2 == 0 {
			blobValue = lstore.Blob("world")
		}
		_, err := store.Write(ctx, intBlobEntry(int64(i)+1, blobValue))
		should.Nil(err)
	}
	should.Nil(store.Index())

	store = reopenTestStore(store)

	reader, err := store.NewReader(ctx)
	should.Nil(err)
	collector := &lstore.ResultCollector{LimitSize: 2}
	reader.SearchForward(ctx, 0, []lstore.Filter{
		store.IndexingStrategy.NewBlobValueFilter(0, "hello"),
	}, collector)
	should.Equal([]int64{2}, collector.Rows[0].IntValues)
	should.Equal([]int64{4}, collector.Rows[1].IntValues)
}
//
//func Test_a_lot_indexed_segment(t *testing.T) {
//	should := require.New(t)
//	store := smallTestStore()
//	defer store.Stop(ctx)
//	for i := 0; i < 160; i++ {
//		blobValue := lstore.Blob("hello")
//		if i%2 == 0 {
//			blobValue = lstore.Blob("world")
//		}
//		_, err := store.Write(ctx, intBlobEntry(int64(i)+1, blobValue))
//		should.Nil(err)
//	}
//	should.Nil(store.Index())
//	reader, err := store.NewReader(ctx)
//	should.Nil(err)
//	iter := reader.Search(ctx, lstore.SearchRequest{
//		LimitSize: 2,
//		Filters: []lstore.Filter{
//			store.IndexingStrategy.NewBlobValueFilter(0, "hello"),
//		},
//	})
//	rows, err := iter()
//	should.Nil(err)
//	should.Equal([]int64{2}, rows[0].IntValues)
//	should.Equal([]int64{4}, rows[1].IntValues)
//}
//
//func Test_compacted_segment(t *testing.T) {
//	should := require.New(t)
//	store := smallTestStore()
//	defer store.Stop(ctx)
//	for j := 0; j < 10; j++ {
//		for i := 0; i < 1000; i++ {
//			blobValue := lstore.Blob("hello")
//			if i%2 == 0 {
//				blobValue = lstore.Blob("world")
//			}
//			_, err := store.Write(ctx, intBlobEntry(int64(i)+1, blobValue))
//			should.Nil(err)
//		}
//		should.Nil(store.Index())
//	}
//	reader, err := store.NewReader(ctx)
//	should.Nil(err)
//	iter := reader.Search(ctx, lstore.SearchRequest{
//		LimitSize: 2,
//		Filters: []lstore.Filter{
//			store.IndexingStrategy.NewBlobValueFilter(0, "hello"),
//		},
//	})
//	rows, err := iter()
//	should.Nil(err)
//	should.Equal([]int64{2}, rows[0].IntValues)
//	should.Equal([]int64{4}, rows[1].IntValues)
//}
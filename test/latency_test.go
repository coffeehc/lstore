package test

import (
	"testing"
	"context"
	"time"
	"fmt"
	"github.com/esdb/lstore"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog"
	"runtime"
	"io"
)

func Test_write_read_latency(t *testing.T) {
	runtime.GOMAXPROCS(4)
	should := require.New(t)
	store := bigTestStore()
	start := time.Now()
	ctx := context.Background()
	resultChan := make(chan lstore.WriteResult, 1024)
	go func() {
		for {
			result := <-resultChan
			fmt.Println("write: ", result.Seq, time.Now())
		}
	}()
	go func() {
		countlog.Info("event!test.search")
		reader, err := store.NewReader()
		should.Nil(err)
		startSeq := lstore.RowSeq(0)
		for {
			hasNew, err := reader.Refresh()
			should.Nil(err)
			if !hasNew {
				time.Sleep(time.Millisecond)
				continue
			}
			iter := reader.Search(ctx, lstore.SearchRequest{StartSeq: startSeq, LimitSize: 1024 * 1024})
			for {
				rows, err := iter()
				if err == io.EOF {
					break
				}
				if err != nil {
					panic(err)
				}
				for _, row := range rows {
					fmt.Println("read: ", row.Seq, time.Now())
				}
				startSeq = rows[len(rows) - 1].Seq
			}
		}
	}()
	for i := 0; i < 1024; i++ {
		store.AsyncWrite(ctx, intEntry(int64(i)), resultChan)
	}
	time.Sleep(time.Second * 5)
	end := time.Now()
	fmt.Println(end.Sub(start))
}
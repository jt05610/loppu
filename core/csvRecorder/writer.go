package csvRecorder

import (
	"context"
	"encoding/csv"
	"injector/softNode/stream/redis"
	"os"
	"strconv"
)

type CSVWriter struct {
	df     *os.File
	writer *csv.Writer
}

func (w *CSVWriter) Handle(ctx context.Context, data *redis.StreamItem) {
	select {
	case <-ctx.Done():
		return
	default:
		f := []string{
			strconv.Itoa(int(data.Depth)),
			strconv.FormatFloat(float64(data.Force), 'e', -1, 32),
		}
		err := w.writer.Write(f)
		if err != nil {
			panic(err)
		}
		w.writer.Flush()
	}
}

func (w *CSVWriter) Close() {
	_ = w.df.Close()
}

func NewCSVWriter(filename string, header []string) redis.Handler {
	var err error
	ret := &CSVWriter{}
	ret.df, err = os.Create(filename)
	if err != nil {
		panic(err)
	}
	ret.writer = csv.NewWriter(ret.df)
	if err != nil {
		panic(err)
	}
	err = ret.writer.Write(header)
	if err != nil {
		panic(err)
	}
	return ret
}

package rs

import (
	"fmt"
	"github.com/klauspost/reedsolomon"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/global"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/objectstream"
	"io"
	"log"
)

type PutStream struct {
	*encoder
}

func NewPutStream(ds []string, hash string, size int64) (*PutStream, error) {
	if len(ds) != global.Cfg.RS.ShardAllNum {
		return nil, fmt.Errorf("data server number mismatch")
	}
	var err error
	perSize := (size + int64(global.Cfg.RS.DataShard) - 1) / int64(global.Cfg.RS.DataShard)
	writers := make([]io.Writer, global.Cfg.RS.ShardAllNum)
	for i := range writers {
		writers[i], err = objectstream.NewTempObjectStream(ds[i], hash, perSize)
		if err != nil {
			return nil, err
		}
	}
	enc, err := NewEncoder(writers)
	if err != nil {
		return nil, err

	}
	return &PutStream{enc}, nil
}

type encoder struct {
	writers []io.Writer
	enc     reedsolomon.Encoder
	cache   []byte
}

func NewEncoder(writers []io.Writer) (*encoder, error) {
	enc, err := reedsolomon.New(global.Cfg.RS.DataShard, global.Cfg.RS.ParityShard)
	if err != nil {
		return nil, err
	}
	return &encoder{
		writers: writers,
		enc:     enc,
	}, nil
}

func (e *encoder) Write(p []byte) (n int, err error) {
	length := len(p)
	current := 0
	for length != 0 {
		next := global.Cfg.RS.BlockSize - len(e.cache)
		if next > length {
			next = length
		}
		e.cache = append(e.cache, p[current:current+next]...)
		if len(e.cache) == global.Cfg.RS.BlockSize {
			e.flush()
		}
		current += next
		length -= next
	}
	return len(p), nil
}

func (e *encoder) flush() {
	shards, err := e.enc.Split(e.cache)
	if err != nil {
		log.Println(err)
		return
	}
	err = e.enc.Encode(shards)
	if err != nil {
		log.Println(err)
		return
	}
	for i := range shards {
		_, err = e.writers[i].Write(shards[i])
		if err != nil {
			return
		}
	}
	e.cache = e.cache[:0]
}

func (e *encoder) Commit() error {
	if len(e.cache) != 0 {
		e.flush()
	}
	for i := range e.writers {
		if stream, ok := e.writers[i].(*objectstream.TempObjectStream); ok {
			err := stream.Commit()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *encoder) Abort() error {
	for i := range e.writers {
		if stream, ok := e.writers[i].(*objectstream.TempObjectStream); ok {
			err := stream.Abort()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

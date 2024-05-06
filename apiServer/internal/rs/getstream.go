package rs

import (
	"fmt"
	"github.com/klauspost/reedsolomon"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/global"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/objectstream"
	"io"
)

type GetStream struct {
	*decoder
}

func NewGetStream(data map[int]string, fill []string, hash string, size int64) (*GetStream, error) {
	writers := make([]io.Writer, global.Cfg.RS.ShardAllNum)
	readers := make([]io.Reader, global.Cfg.RS.ShardAllNum)
	for i := 0; i < global.Cfg.RS.ShardAllNum; i++ {
		if v, ok := data[i]; ok {
			reader, err := objectstream.NewGetStream(v, fmt.Sprint("%s.%s", hash, i))
			if err != nil {
				return nil, err
			}
			readers[i] = reader
		}
	}
	perSize := (size + int64(global.Cfg.RS.DataShard) - 1) / int64(global.Cfg.RS.DataShard)
	for i := 0; i < global.Cfg.RS.ShardAllNum; i++ {
		if readers[i] == nil {
			writer, err := objectstream.NewTempObjectStream(fill[i], fmt.Sprintf("%s.%s", hash, i), perSize)
			if err != nil {
				return nil, err
			}
			writers[i] = writer
		}
	}

	dec, err := NewDecoder(writers, readers, global.Cfg.RS.DataShard, global.Cfg.RS.ParityShard, size)
	if err != nil {
		return nil, err
	}
	return &GetStream{dec}, nil
}

type decoder struct {
	writers   []io.Writer
	readers   []io.Reader
	enc       reedsolomon.Encoder
	size      int64
	cache     []byte
	cacheSize int64
	total     int64
}

func NewDecoder(writers []io.Writer, readers []io.Reader, dataShard, parityShard int, size int64) (*decoder, error) {
	enc, err := reedsolomon.New(dataShard, parityShard)
	if err != nil {
		return nil, err
	}
	return &decoder{
		writers: writers,
		readers: readers,
		enc:     enc,
		size:    size,
	}, nil
}

func (d *decoder) Read(p []byte) (n int, err error) {
	if d.cacheSize == 0 {
		err := d.getData()
		if err != nil {
			return 0, err
		}
	}
	length := len(p)
	if d.cacheSize < int64(length) {
		length = int(d.cacheSize)
	}
	d.cacheSize -= int64(length)
	copy(p, d.cache[:length])
	d.cache = d.cache[length:]
	return length, nil
}

func (d *decoder) getData() error {
	if d.total == d.size {
		return io.EOF
	}
	shards := make([][]byte, global.Cfg.RS.ShardAllNum)
	fixIds := make([]int, 0)
	// Read data shards
	for i := 0; i < global.Cfg.RS.ShardAllNum; i++ {
		if d.readers[i] == nil {
			fixIds = append(fixIds, i)
			continue
		}
		shards[i] = make([]byte, global.Cfg.RS.BlockSize/global.Cfg.RS.DataShard)
		n, err := io.ReadFull(d.readers[i], shards[i])
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			shards[i] = nil
		}
		if n != len(shards[i]) {
			return fmt.Errorf("data shard %d size mismatch", i)
		}
	}
	// Reconstruct
	err := d.enc.Reconstruct(shards)
	if err != nil {
		return err
	}
	// Write fixed shards
	for _, id := range fixIds {
		_, err = d.writers[id].Write(shards[id])
		if err != nil {
			return err
		}
	}
	// Cache data
	for i := range global.Cfg.RS.DataShard {
		shardSize := len(shards[i])
		if d.total+int64(shardSize) > d.size {
			shardSize = int(d.size - d.total)
		}
		d.cache = append(d.cache, shards[i][:shardSize]...)
		d.cacheSize += int64(shardSize)
		d.total += int64(shardSize)
	}
	return nil
}

package lbase

import (
	"errors"

	"github.com/wmentor/dsn"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type DB interface {
	Delete(key []byte)
	Get(key []byte) []byte
	Has(key []byte) bool
	Set(key, value []byte)
	Close()
}

type db struct {
	base *leveldb.DB
}

var (
	ErrNoPath error = errors.New("db path not entered")

	gdb DB = nil
)

func Open(params string) (DB, error) {

	kvs, err := dsn.New(params)
	if err != nil {
		return nil, err
	}

	var ldb DB

	if kvs.GetBool("test", false) {

		ldb = newStub()

	} else {

		path := kvs.GetString("path", "")
		if path == "" {
			return nil, ErrNoPath
		}

		base, err := leveldb.OpenFile(path, &opt.Options{
			CompactionTableSize: 16,
			WriteBuffer:         16 * 2,
			Compression:         opt.SnappyCompression,
			ReadOnly:            false,
		})

		if err != nil {
			return nil, err
		}

		ldb = &db{
			base: base,
		}
	}

	if kvs.GetBool("global", false) {
		gdb = ldb
	}

	return ldb, nil
}

func (base *db) Close() {
	if base != nil && base.base != nil {
		base.base.Close()
		base.base = nil
	}
}

func (base *db) Get(key []byte) []byte {

	if base != nil && base.base != nil {
		if len(key) != 0 {
			if val, err := base.base.Get(key, nil); err == nil {
				return val
			}
		}
	}

	return nil
}

func (base *db) Has(key []byte) bool {
	return len(base.Get(key)) > 0
}

func (base *db) Delete(key []byte) {
	base.Set(key, nil)
}

func (base *db) Set(key, value []byte) {
	if base != nil && base.base != nil {
		if len(key) != 0 {
			if len(value) == 0 {
				base.base.Delete(key, nil)
			} else {
				base.base.Put(key, value, nil)
			}
		}
	}
}

func Set(key, value []byte) {
	gdb.Set(key, value)
}

func Get(key []byte) []byte {
	return gdb.Get(key)
}

func Has(key []byte) bool {
	return gdb.Has(key)
}
func Delete(key []byte) {
	gdb.Delete(key)
}

func Close() {
	gdb.Close()
}

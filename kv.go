package kv

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/wmentor/dsn"
	"github.com/wmentor/log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type DB interface {
	Delete(key []byte)
	Get(key []byte) []byte
	Has(key []byte) bool
	Set(key, value []byte)
	Prefix(prefix []byte, fn PairIteratorFunc)
	Range(from []byte, to []byte, fn PairIteratorFunc)
	Close()
}

type PairIteratorFunc func(key []byte, value []byte) bool

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

	if kvs.GetBool("global", true) {
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
	if base != nil {
		base.Set(key, nil)
	}
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

func (base *db) Prefix(prefix []byte, fn PairIteratorFunc) {
	if base == nil || base.base == nil {
		return
	}

	iter := base.base.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()

	for iter.Next() {
		if !fn(iter.Key(), iter.Value()) {
			return
		}
	}
}

func (base *db) Range(from []byte, to []byte, fn PairIteratorFunc) {
	if base == nil || base.base == nil {
		return
	}

	iter := base.base.NewIterator(&util.Range{Start: from, Limit: to}, nil)
	defer iter.Release()

	for iter.Next() {
		if !fn(iter.Key(), iter.Value()) {
			return
		}
	}
}

func Prefix(prefix []byte, fn PairIteratorFunc) {
	gdb.Prefix(prefix, fn)
}

func Range(from []byte, to []byte, fn PairIteratorFunc) {
	gdb.Range(from, to, fn)
}

func CopyWork(src string, fn func(DB)) {

	tmpProc := func() bool {

		defer time.Sleep(time.Second * 2)

		ts := time.Now().Unix()

		tmpName := fmt.Sprintf("%s-%d-%d", src, ts, time.Now().UnixNano()%1000)

		cmd := fmt.Sprintf("cp -r %s %s", src, tmpName)

		exec.Command("sh", "-c", cmd).Run()

		defer func() {
			exec.Command("sh", "-c", "rm -rf "+tmpName).Run()
		}()

		db, err := Open(fmt.Sprintf("global=0 path=" + tmpName))
		if err != nil {
			log.Error(err.Error() + " " + tmpName)
			return false
		}
		defer db.Close()

		fn(db)

		return true
	}

	for i := 0; i < 5; i++ {
		if tmpProc() {
			return
		}
	}
}

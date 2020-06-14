lbase is simple key/value storage based on *syndtr/goleveldb* library.

# Summary

* Written on tiny Go
* Go library
* Require Go version >= 1.2
* Simple API
* Persistent key/value storage
* Use LSM-tree engine
* Has mock object for tests

# Install

```
go get github.com/wmentor/lbase
```

# Usage

```go
package main

import (
  "fmt"

  "github.com/wmentor/lbase"
)

func main() {
  db, err := lbase.Open("path=./base")
	if err != nil {
		panic(err)
	}
	defer db.Close()

  db.Set([]byte("1"), []byte("11"))
  db.Set([]byte("10"), []byte("1100"))
  db.Set([]byte("11"), []byte("1111"))
  db.Set([]byte("12"), []byte("1122"))
  db.Set([]byte("3"), []byte("33"))
  db.Set([]byte("4"), []byte("44"))
  db.Set([]byte("5"), []byte("55"))
  db.Set([]byte("6"), []byte("66"))

  fmt.Println(string(db.Get([]byte("5")))) // 55

  db.Get([]byte("7")) // nil

  fmt.Println(db.Has([]byte("5"))) // true
  db.Set([]byte("5"), nil) // remove key []byte("5")
  fmt.Println(db.Has([]byte("5"))) // false

  db.Prefix([]byte("1"), func(k, v []byte) bool {
    fmt.Println(string(k) + " " + string(v))
    return true
  })
  /* print:
  1 11
  10 1100
  11 1111
  12 1122
  */

  db.Range([]byte("12"), []byte("4"), func(k, v []byte) bool {
    fmt.Println(string(k) + " " + string(v))
    return true
  })
  /* print:
  12 1122
  2 22
  3 33
  */
}
```

Work with global context:

```go
package main

import (
  "github.com/wmentor/lbase"
)

func main() {
  _, err := lbase.Open("path=./base global=1")
	if err != nil {
		panic(err)
	}
	defer lbase.Close()

  lbase.Set([]byte("1"), []byte("11"))
  lbase.Set([]byte("10"), []byte("1100"))
  ...
}
```

Work with test context:

```go
package main

import (
  "github.com/wmentor/lbase"
)

func main() {
  _, err := lbase.Open("test=1 global=1")
	if err != nil {
		panic(err)
	}
	defer lbase.Close()

  lbase.Set([]byte("1"), []byte("11"))
  lbase.Set([]byte("10"), []byte("1100"))
  ...
}
```

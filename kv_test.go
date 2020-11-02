package kv

import (
	"os/exec"
	"strings"
	"testing"
)

func tFunc(t *testing.T) {
	tSet := func(key, val string) {
		Set([]byte(key), []byte(val))

		if Has([]byte(key)) && val == "" {
			t.Fatalf("Set failed for key=%s value=%s", key, val)
		} else if !Has([]byte(key)) && val != "" {
			t.Fatalf("Set failed for key=%s value=%s", key, val)
		}

		if string(Get([]byte(key))) != val {
			t.Fatalf("Set failed for key=%s value=%s", key, val)
		}
	}

	tSet("1", "11")
	tSet("2", "22")
	tSet("3", "33")
	tSet("4", "44")
	tSet("5", "55")
	tSet("6", "66")
	tSet("6", "666")
	tSet("6", "")
	tSet("10", "1100")
	tSet("11", "1111")
	tSet("12", "1122")

	res := []string{}

	Prefix([]byte("1"), func(k, v []byte) bool {
		res = append(res, string(k)+":"+string(v))
		return true
	})

	if strings.Join(res, ";") != "1:11;10:1100;11:1111;12:1122" {
		t.Fatal("Prefix failed")
	}

	res = res[:0]

	Range([]byte("11"), []byte("3"), func(k, v []byte) bool {
		res = append(res, string(k)+":"+string(v))
		return true
	})

	if strings.Join(res, ";") != "11:1111;12:1122;2:22" {
		t.Fatal("Range failed")
	}
}

func TestLBase(t *testing.T) {

	db, err := Open("path=./base global=1")
	if err != nil || db == nil {
		t.Fatal("Open failed")
	}
	defer exec.Command("sh", "-c", "rm  -rf ./base").Run()

	defer db.Close()

	tFunc(t)
}

package main

import (
	"fmt"
	"github.com/ztteng/goleveldb/leveldb"
)

func main() {
	db, _ := leveldb.OpenFile("D:/GoWork/src/github.com/ztteng/db", nil)
	_ = db.Put([]byte("key"), []byte("value"), nil)
	data, _ := db.Get([]byte("key"), nil)
	fmt.Println(string(data))
}

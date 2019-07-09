package iterator_test

import (
	"testing"

	"github.com/ztteng/goleveldb/leveldb/testutil"
)

func TestIterator(t *testing.T) {
	testutil.RunSuite(t, "Iterator Suite")
}

package sqlite

import (
	"github.com/gocraft/dbr/v2"
)

type Transaction struct {
	tx *dbr.Tx
}

func (t Transaction) Rollback() {
	t.tx.RollbackUnlessCommitted()
}

func (t Transaction) Commit() error {
	return t.tx.Commit()
}

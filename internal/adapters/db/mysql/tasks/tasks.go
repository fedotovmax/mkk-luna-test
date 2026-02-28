package tasks

import mysqlTx "github.com/fedotovmax/mkk-luna-test/internal/adapters/db/transaction/mysql"

type task struct {
	txExtractor mysqlTx.Extractor
}

func New(txExtractor mysqlTx.Extractor) *task {
	return &task{txExtractor: txExtractor}
}

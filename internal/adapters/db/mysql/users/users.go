package users

import mysqlTx "github.com/fedotovmax/mkk-luna-test/internal/adapters/db/transaction/mysql"

type user struct {
	txExtractor mysqlTx.Extractor
}

func New(txExtractor mysqlTx.Extractor) *user {
	return &user{txExtractor: txExtractor}
}

package sessions

import mysqlTx "github.com/fedotovmax/mkk-luna-test/internal/adapters/db/transaction/mysql"

type session struct {
	txExtractor mysqlTx.Extractor
}

func New(txExtractor mysqlTx.Extractor) *session {
	return &session{txExtractor: txExtractor}
}

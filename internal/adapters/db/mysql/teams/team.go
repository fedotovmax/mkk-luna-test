package teams

import mysqlTx "github.com/fedotovmax/mkk-luna-test/internal/adapters/db/transaction/mysql"

type team struct {
	txExtractor mysqlTx.Extractor
}

func New(txExtractor mysqlTx.Extractor) *team {
	return &team{txExtractor: txExtractor}
}

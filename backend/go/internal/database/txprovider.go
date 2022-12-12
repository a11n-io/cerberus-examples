package database

import "database/sql"

type txProvider struct {
	db *sql.DB
}

type TxProvider interface {
	GetTransaction() (*sql.Tx, error)
}

func NewTxProvider(db *sql.DB) TxProvider {
	return &txProvider{
		db: db,
	}
}

func (p *txProvider) GetTransaction() (*sql.Tx, error) {
	return p.db.Begin()
}

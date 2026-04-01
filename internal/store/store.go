package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"
	_ "modernc.org/sqlite"
)

type DB struct{ conn *sql.DB }

func Open(dataDir string) (*DB, error) {
	os.MkdirAll(dataDir, 0755)
	conn, err := sql.Open("sqlite", filepath.Join(dataDir, "ticker.db"))
	if err != nil { return nil, err }
	conn.Exec("PRAGMA journal_mode=WAL"); conn.Exec("PRAGMA busy_timeout=5000"); conn.SetMaxOpenConns(4)
	db := &DB{conn: conn}; return db, db.migrate()
}
func (db *DB) Close() error { return db.conn.Close() }
func (db *DB) migrate() error {
	_, err := db.conn.Exec(`
CREATE TABLE IF NOT EXISTS accounts (id TEXT PRIMARY KEY, name TEXT NOT NULL, type TEXT DEFAULT 'checking', balance_cents INTEGER DEFAULT 0, currency TEXT DEFAULT 'USD', created_at TEXT DEFAULT (datetime('now')));
CREATE TABLE IF NOT EXISTS transactions (id TEXT PRIMARY KEY, account_id TEXT NOT NULL, description TEXT NOT NULL, amount_cents INTEGER NOT NULL, category TEXT DEFAULT '', date TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')));
CREATE INDEX IF NOT EXISTS idx_txn_acct ON transactions(account_id);
CREATE INDEX IF NOT EXISTS idx_txn_date ON transactions(date);`)
	return err
}

type Account struct {
	ID string `json:"id"`; Name string `json:"name"`; Type string `json:"type"`
	BalanceCents int `json:"balance_cents"`; Currency string `json:"currency"`; CreatedAt string `json:"created_at"`
}

func (db *DB) CreateAccount(name, atype, currency string) (*Account, error) {
	id := "acct_" + gid(6); now := time.Now().UTC().Format(time.RFC3339)
	if atype == "" { atype = "checking" }; if currency == "" { currency = "USD" }
	_, err := db.conn.Exec("INSERT INTO accounts (id,name,type,currency,created_at) VALUES (?,?,?,?,?)", id, name, atype, currency, now)
	if err != nil { return nil, err }
	return &Account{ID: id, Name: name, Type: atype, Currency: currency, CreatedAt: now}, nil
}
func (db *DB) ListAccounts() ([]Account, error) {
	rows, err := db.conn.Query("SELECT id,name,type,balance_cents,currency,created_at FROM accounts ORDER BY name")
	if err != nil { return nil, err }; defer rows.Close()
	var out []Account
	for rows.Next() { var a Account; rows.Scan(&a.ID, &a.Name, &a.Type, &a.BalanceCents, &a.Currency, &a.CreatedAt); out = append(out, a) }
	return out, rows.Err()
}
func (db *DB) DeleteAccount(id string) { db.conn.Exec("DELETE FROM transactions WHERE account_id=?", id); db.conn.Exec("DELETE FROM accounts WHERE id=?", id) }

type Transaction struct {
	ID string `json:"id"`; AccountID string `json:"account_id"`; Description string `json:"description"`
	AmountCents int `json:"amount_cents"`; Category string `json:"category"`; Date string `json:"date"`; CreatedAt string `json:"created_at"`
}

func (db *DB) AddTransaction(accountID, description string, amountCents int, category, date string) (*Transaction, error) {
	id := "txn_" + gid(8); now := time.Now().UTC().Format(time.RFC3339)
	if date == "" { date = time.Now().Format("2006-01-02") }
	_, err := db.conn.Exec("INSERT INTO transactions (id,account_id,description,amount_cents,category,date,created_at) VALUES (?,?,?,?,?,?,?)",
		id, accountID, description, amountCents, category, date, now)
	if err != nil { return nil, err }
	db.conn.Exec("UPDATE accounts SET balance_cents=balance_cents+? WHERE id=?", amountCents, accountID)
	return &Transaction{ID: id, AccountID: accountID, Description: description, AmountCents: amountCents, Category: category, Date: date, CreatedAt: now}, nil
}
func (db *DB) ListTransactions(accountID string, limit int) ([]Transaction, error) {
	if limit <= 0 { limit = 50 }
	var rows *sql.Rows; var err error
	if accountID != "" {
		rows, err = db.conn.Query("SELECT id,account_id,description,amount_cents,category,date,created_at FROM transactions WHERE account_id=? ORDER BY date DESC LIMIT ?", accountID, limit)
	} else {
		rows, err = db.conn.Query("SELECT id,account_id,description,amount_cents,category,date,created_at FROM transactions ORDER BY date DESC LIMIT ?", limit)
	}
	if err != nil { return nil, err }; defer rows.Close()
	var out []Transaction
	for rows.Next() { var t Transaction; rows.Scan(&t.ID, &t.AccountID, &t.Description, &t.AmountCents, &t.Category, &t.Date, &t.CreatedAt); out = append(out, t) }
	return out, rows.Err()
}
func (db *DB) DeleteTransaction(id string) {
	var acctID string; var amt int
	db.conn.QueryRow("SELECT account_id,amount_cents FROM transactions WHERE id=?", id).Scan(&acctID, &amt)
	db.conn.Exec("DELETE FROM transactions WHERE id=?", id)
	if acctID != "" { db.conn.Exec("UPDATE accounts SET balance_cents=balance_cents-? WHERE id=?", amt, acctID) }
}
func (db *DB) Stats() map[string]any {
	var accts, txns, bal int
	db.conn.QueryRow("SELECT COUNT(*) FROM accounts").Scan(&accts)
	db.conn.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&txns)
	db.conn.QueryRow("SELECT COALESCE(SUM(balance_cents),0) FROM accounts").Scan(&bal)
	return map[string]any{"accounts": accts, "transactions": txns, "net_balance_cents": bal}
}
func gid(n int) string { b := make([]byte, n); rand.Read(b); return hex.EncodeToString(b) }

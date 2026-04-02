package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Price struct{
	ID string `json:"id"`
	Symbol string `json:"symbol"`
	Price float64 `json:"price"`
	Change float64 `json:"change"`
	Volume float64 `json:"volume"`
	Source string `json:"source"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"ticker.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS prices(id TEXT PRIMARY KEY,symbol TEXT NOT NULL,price REAL DEFAULT 0,change REAL DEFAULT 0,volume REAL DEFAULT 0,source TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Price)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO prices(id,symbol,price,change,volume,source,created_at)VALUES(?,?,?,?,?,?,?)`,e.ID,e.Symbol,e.Price,e.Change,e.Volume,e.Source,e.CreatedAt);return err}
func(d *DB)Get(id string)*Price{var e Price;if d.db.QueryRow(`SELECT id,symbol,price,change,volume,source,created_at FROM prices WHERE id=?`,id).Scan(&e.ID,&e.Symbol,&e.Price,&e.Change,&e.Volume,&e.Source,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Price{rows,_:=d.db.Query(`SELECT id,symbol,price,change,volume,source,created_at FROM prices ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Price;for rows.Next(){var e Price;rows.Scan(&e.ID,&e.Symbol,&e.Price,&e.Change,&e.Volume,&e.Source,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM prices WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM prices`).Scan(&n);return n}

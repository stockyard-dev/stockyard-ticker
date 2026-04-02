package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Price struct {
	ID string `json:"id"`
	Symbol string `json:"name"`
	Exchange string `json:"exchange"`
	Price int `json:"price"`
	Change int `json:"change_pct"`
	Volume int `json:"volume"`
	High int `json:"high"`
	Low int `json:"low"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"ticker.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS prices(id TEXT PRIMARY KEY,name TEXT NOT NULL,exchange TEXT DEFAULT '',price INTEGER DEFAULT 0,change_pct INTEGER DEFAULT 0,volume INTEGER DEFAULT 0,high INTEGER DEFAULT 0,low INTEGER DEFAULT 0,status TEXT DEFAULT 'active',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Price)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO prices(id,name,exchange,price,change_pct,volume,high,low,status,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Symbol,e.Exchange,e.Price,e.Change,e.Volume,e.High,e.Low,e.Status,e.CreatedAt);return err}
func(d *DB)Get(id string)*Price{var e Price;if d.db.QueryRow(`SELECT id,name,exchange,price,change_pct,volume,high,low,status,created_at FROM prices WHERE id=?`,id).Scan(&e.ID,&e.Symbol,&e.Exchange,&e.Price,&e.Change,&e.Volume,&e.High,&e.Low,&e.Status,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Price{rows,_:=d.db.Query(`SELECT id,name,exchange,price,change_pct,volume,high,low,status,created_at FROM prices ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Price;for rows.Next(){var e Price;rows.Scan(&e.ID,&e.Symbol,&e.Exchange,&e.Price,&e.Change,&e.Volume,&e.High,&e.Low,&e.Status,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Price)error{_,err:=d.db.Exec(`UPDATE prices SET name=?,exchange=?,price=?,change_pct=?,volume=?,high=?,low=?,status=? WHERE id=?`,e.Symbol,e.Exchange,e.Price,e.Change,e.Volume,e.High,e.Low,e.Status,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM prices WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM prices`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Price{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,exchange,price,change_pct,volume,high,low,status,created_at FROM prices WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Price;for rows.Next(){var e Price;rows.Scan(&e.ID,&e.Symbol,&e.Exchange,&e.Price,&e.Change,&e.Volume,&e.High,&e.Low,&e.Status,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM prices GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}

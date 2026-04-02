package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Job struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Schedule string `json:"schedule"`
	Command string `json:"command"`
	WebhookURL string `json:"webhook_url"`
	Enabled int `json:"enabled"`
	LastRunAt string `json:"last_run_at"`
	LastResult string `json:"last_result"`
	RunCount int `json:"run_count"`
	FailCount int `json:"fail_count"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"mainspring.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS jobs(id TEXT PRIMARY KEY,name TEXT NOT NULL,schedule TEXT DEFAULT '',command TEXT DEFAULT '',webhook_url TEXT DEFAULT '',enabled INTEGER DEFAULT 1,last_run_at TEXT DEFAULT '',last_result TEXT DEFAULT '',run_count INTEGER DEFAULT 0,fail_count INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Job)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO jobs(id,name,schedule,command,webhook_url,enabled,last_run_at,last_result,run_count,fail_count,created_at)VALUES(?,?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Schedule,e.Command,e.WebhookURL,e.Enabled,e.LastRunAt,e.LastResult,e.RunCount,e.FailCount,e.CreatedAt);return err}
func(d *DB)Get(id string)*Job{var e Job;if d.db.QueryRow(`SELECT id,name,schedule,command,webhook_url,enabled,last_run_at,last_result,run_count,fail_count,created_at FROM jobs WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Schedule,&e.Command,&e.WebhookURL,&e.Enabled,&e.LastRunAt,&e.LastResult,&e.RunCount,&e.FailCount,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Job{rows,_:=d.db.Query(`SELECT id,name,schedule,command,webhook_url,enabled,last_run_at,last_result,run_count,fail_count,created_at FROM jobs ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Job;for rows.Next(){var e Job;rows.Scan(&e.ID,&e.Name,&e.Schedule,&e.Command,&e.WebhookURL,&e.Enabled,&e.LastRunAt,&e.LastResult,&e.RunCount,&e.FailCount,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Job)error{_,err:=d.db.Exec(`UPDATE jobs SET name=?,schedule=?,command=?,webhook_url=?,enabled=?,last_run_at=?,last_result=?,run_count=?,fail_count=? WHERE id=?`,e.Name,e.Schedule,e.Command,e.WebhookURL,e.Enabled,e.LastRunAt,e.LastResult,e.RunCount,e.FailCount,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM jobs WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM jobs`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Job{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["enabled"];ok&&v!=""{where+=" AND enabled=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,schedule,command,webhook_url,enabled,last_run_at,last_result,run_count,fail_count,created_at FROM jobs WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Job;for rows.Next(){var e Job;rows.Scan(&e.ID,&e.Name,&e.Schedule,&e.Command,&e.WebhookURL,&e.Enabled,&e.LastRunAt,&e.LastResult,&e.RunCount,&e.FailCount,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    return m
}

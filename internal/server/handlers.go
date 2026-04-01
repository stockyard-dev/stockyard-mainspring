package server
import("encoding/json";"net/http";"os/exec";"strconv";"time";"github.com/stockyard-dev/stockyard-mainspring/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){list,_:=s.db.List();if list==nil{list=[]store.Job{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var j store.Job;json.NewDecoder(r.Body).Decode(&j);if j.Name==""||j.Schedule==""||j.Command==""{writeError(w,400,"name, schedule, command required");return};s.db.Create(&j);writeJSON(w,201,j)}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleRuns(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);list,_:=s.db.ListRuns(id);if list==nil{list=[]store.Run{}};writeJSON(w,200,list)}
func(s *Server)handleTrigger(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);jobs,_:=s.db.List();var cmd string;for _,j:=range jobs{if j.ID==id{cmd=j.Command;break}};run:=&store.Run{JobID:id};t0:=time.Now();var out[]byte;var err error;if cmd!=""{c:=exec.Command("sh","-c",cmd);out,err=c.CombinedOutput()};run.DurationMs=time.Since(t0).Milliseconds();run.Output=string(out);if err!=nil{run.Status="failed"}else{run.Status="success"};s.db.RecordRun(run);writeJSON(w,200,run)}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}

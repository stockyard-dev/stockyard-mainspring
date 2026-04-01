package server
import("encoding/json";"net/http";"github.com/stockyard-dev/stockyard-mainspring/internal/store")
type Server struct{db *store.DB;limits Limits;mux *http.ServeMux}
func New(db *store.DB,tier string)*Server{s:=&Server{db:db,limits:LimitsFor(tier),mux:http.NewServeMux()};s.routes();return s}
func(s *Server)ListenAndServe(addr string)error{return(&http.Server{Addr:addr,Handler:s.mux}).ListenAndServe()}
func(s *Server)routes(){
    s.mux.HandleFunc("GET /health",func(w http.ResponseWriter,r *http.Request){writeJSON(w,200,map[string]string{"status":"ok","service":"stockyard-mainspring"})})
    s.mux.HandleFunc("GET /api/stats",s.handleOverview)
    s.mux.HandleFunc("GET /api/jobs",s.handleList)
    s.mux.HandleFunc("POST /api/jobs",s.handleCreate)
    s.mux.HandleFunc("DELETE /api/jobs/{id}",s.handleDelete)
    s.mux.HandleFunc("POST /api/jobs/{id}/trigger",s.handleTrigger)
    s.mux.HandleFunc("GET /api/jobs/{id}/runs",s.handleRuns)
    s.mux.HandleFunc("GET /",s.handleUI)}
func writeJSON(w http.ResponseWriter,status int,v interface{}){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);json.NewEncoder(w).Encode(v)}
func writeError(w http.ResponseWriter,status int,msg string){writeJSON(w,status,map[string]string{"error":msg})}
func(s *Server)handleUI(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};w.Header().Set("Content-Type","text/html");w.Write(dashboardHTML)}

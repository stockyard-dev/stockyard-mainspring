package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-mainspring/internal/server";"github.com/stockyard-dev/stockyard-mainspring/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./mainspring-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("mainspring: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Mainspring — Self-hosted cron job scheduler\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("mainspring: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}

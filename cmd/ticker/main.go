package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-ticker/internal/server";"github.com/stockyard-dev/stockyard-ticker/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9040"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./ticker-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("ticker: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Ticker — metrics and time series dashboard\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("ticker: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}

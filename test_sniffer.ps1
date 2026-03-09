$env:Path += ";C:\Program Files\Go\bin"
$env:USE_REAL_SNIFFER = "true"
cd backend/ingestion
go run main.go

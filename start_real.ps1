# Ensure Go is in Path
$env:Path += ";C:\Program Files\Go\bin"

Write-Host "🛑 Stopping Services..." -ForegroundColor Yellow
taskkill /F /IM go.exe 2>$null
taskkill /F /IM node.exe 2>$null
taskkill /F /IM main.exe 2>$null
# Wait a moment for ports to free up
Start-Sleep -Seconds 2

Write-Host "🚀 Starting Ecosystem..." -ForegroundColor Green

# Start Metadata Service
Write-Host "   - Starting Metadata Service (Port 8082)..."
Start-Process -FilePath "go" -ArgumentList "run ." -WorkingDirectory "backend\metadata" -NoNewWindow
Start-Sleep -Seconds 1

# Start Calculation Service
Write-Host "   - Starting Calculation Engine (Port 8081)..."
Start-Process -FilePath "go" -ArgumentList "run ." -WorkingDirectory "backend\calculation" -NoNewWindow
Start-Sleep -Seconds 1

# Start Ingestion Service with Real Sniffer
Write-Host "   - Starting Ingestion Service (Real Mode) (Port 8080)..."
$env:USE_REAL_SNIFFER = "true"
Start-Process -FilePath "go" -ArgumentList "run ." -WorkingDirectory "backend\ingestion" -NoNewWindow

# Start Frontend
Write-Host "   - Launching Frontend (Port 5173)..."
Start-Process -FilePath "cmd" -ArgumentList "/c cd frontend && npm run dev" -NoNewWindow

Write-Host "`n✅ System Online!" -ForegroundColor Cyan
Write-Host "👉 Dashboard: http://localhost:5173" -ForegroundColor Cyan

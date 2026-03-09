# Ensure Go is in Path
$env:Path += ";C:\Program Files\Go\bin"
$env:USE_REAL_SNIFFER = "true"
$root = Get-Location

Write-Host " Killing old processes..." -ForegroundColor Yellow
& "$root\cleanup.ps1"

Write-Host " Starting Albion Analytical Ecosystem..." -ForegroundColor Cyan

# Start Metadata Service
$metaPath = Join-Path $root "backend\metadata\metadata.exe"
Start-Process -FilePath $metaPath -WorkingDirectory "$root\backend\metadata" -NoNewWindow
Write-Host " [OK] Metadata Service (Port 8082)" -ForegroundColor Green

# Start Calculation Service
$calcPath = Join-Path $root "backend\calculation\calculation.exe"
Start-Process -FilePath $calcPath -WorkingDirectory "$root\backend\calculation" -NoNewWindow
Write-Host " [OK] Calculation Service (Port 8081)" -ForegroundColor Green

# Start Ingestion Service
$ingestPath = Join-Path $root "backend\ingestion\ingestion.exe"
Start-Process -FilePath $ingestPath -WorkingDirectory "$root\backend\ingestion" -NoNewWindow
Write-Host " [OK] Ingestion Service (Port 8080)" -ForegroundColor Green

# Start Frontend
Start-Process -FilePath "cmd" -ArgumentList "/c cd frontend && npm run dev -- --host 127.0.0.1" -NoNewWindow
Write-Host " [OK] Frontend Dashboard (Port 5173)" -ForegroundColor Green

Write-Host "`nAll systems GO! Access Dashboard at http://localhost:5173" -ForegroundColor Magenta

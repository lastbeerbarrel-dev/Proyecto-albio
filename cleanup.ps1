Write-Host " Killing Albion processes..." -ForegroundColor Yellow

$names = @("metadata", "calculation", "ingestion", "node", "vite", "main", "go")
foreach ($name in $names) {
    Get-Process -Name $name -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
}

Write-Host " Cleanup Done." -ForegroundColor Green

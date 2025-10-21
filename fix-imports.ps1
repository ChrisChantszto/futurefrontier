# Fix all import paths from onetake to futurefrontier

Write-Host "Fixing import paths..." -ForegroundColor Yellow

$files = Get-ChildItem -Path . -Include *.go -Recurse -File

$count = 0
foreach ($file in $files) {
    $content = Get-Content $file.FullName -Raw
    $newContent = $content -replace 'github\.com/onetakesolutions/onetake-corpsite-backend', 'github.com/ChrisChantszto/futurefrontier'
    
    if ($content -ne $newContent) {
        Set-Content -Path $file.FullName -Value $newContent -NoNewline
        Write-Host "✓ Fixed: $($file.FullName)" -ForegroundColor Green
        $count++
    }
}

Write-Host ""
Write-Host "Fixed $count files" -ForegroundColor Cyan
Write-Host ""
Write-Host "Now run: go mod tidy" -ForegroundColor Yellow

# Test image upload validation

Write-Host "`n=== TEST IMAGE VALIDATION ===" -ForegroundColor Cyan

Start-Sleep -Seconds 2

# Test 1: Invalid file type
Write-Host "`nTest 1: File type tidak valid (.txt)" -ForegroundColor Yellow
$testFile = "test.txt"
"Invalid file content" | Out-File $testFile -Force
try {
    $response = curl -X POST http://localhost:8080/upload -F "file=@$testFile" 2>$null
    Write-Host $response -ForegroundColor Red
} catch {
    Write-Host "Error: $_"
}
Remove-Item $testFile -Force

Start-Sleep -Seconds 1

# Test 2: File too large (create a "large" file)
Write-Host "`nTest 2: File terlalu besar (>5MB)" -ForegroundColor Yellow
$largeFile = "large.jpg"
# Create a 6MB file
$content = [string]::new([char]65, 6291456) # 6MB of 'A'
[System.IO.File]::WriteAllText($largeFile, $content)
try {
    $response = curl -X POST http://localhost:8080/upload -F "file=@$largeFile" 2>$null
    Write-Host $response -ForegroundColor Red
} catch {
    Write-Host "Error: $_"
}
Remove-Item $largeFile -Force

Write-Host "`nValidation tests completed!" -ForegroundColor Green

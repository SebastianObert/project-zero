# Test 1: Invalid data (price = 0)
Write-Host "Test 1: Invalid price (must be > 0)" -ForegroundColor Yellow
$invalidData = @{
    title = "Rumah"
    description = "Rumah bagus"
    listing_type = "WTS"
    price = 0
    land_size = 200
    building_size = 150
    bedrooms = 3
    bathrooms = 2
    floors = 2
    certificate = "SHM"
    electricity = 2200
    water_source = "PAM"
    address = "Jalan Test No.123"
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/properties" -Method POST -Headers @{"Content-Type"="application/json"} -Body $invalidData
    Write-Host $response.Content -ForegroundColor Red
} catch {
    Write-Host $_.Exception.Response.StatusCode -ForegroundColor Red
    Write-Host $_.ErrorDetails.Message -ForegroundColor Red
}

Start-Sleep -Seconds 1

# Test 2: Valid data
Write-Host "`nTest 2: Valid property data" -ForegroundColor Yellow
$validData = @{
    title = "Rumah Mewah di Jakarta"
    description = "Rumah yang sangat bagus dengan fasilitas lengkap dan lokasi strategis"
    listing_type = "WTS"
    price = 500000000
    land_size = 200
    building_size = 150
    bedrooms = 3
    bathrooms = 2
    floors = 2
    certificate = "SHM"
    electricity = 2200
    water_source = "PAM"
    address = "Jalan Sudirman No.123, Jakarta, Indonesia"
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/properties" -Method POST -Headers @{"Content-Type"="application/json"} -Body $validData
    Write-Host $response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 3 -ForegroundColor Green
} catch {
    Write-Host "Error: " $_.Exception.Message -ForegroundColor Red
}

Start-Sleep -Seconds 1

# Test 3: Missing required field
Write-Host "`nTest 3: Missing required field (title)" -ForegroundColor Yellow
$missingField = @{
    description = "Rumah tanpa judul"
    listing_type = "WTS"
    price = 300000000
    land_size = 200
    building_size = 150
    bedrooms = 3
    bathrooms = 2
    floors = 2
    certificate = "SHM"
    electricity = 2200
    water_source = "PAM"
    address = "Jalan Test No.123"
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/properties" -Method POST -Headers @{"Content-Type"="application/json"} -Body $missingField
    Write-Host $response.Content -ForegroundColor Red
} catch {
    Write-Host $_.Exception.Response.StatusCode -ForegroundColor Red
    Write-Host $_.ErrorDetails.Message -ForegroundColor Red
}

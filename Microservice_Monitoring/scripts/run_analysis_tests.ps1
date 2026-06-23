$container = 'timescaledb'
$dbUser = 'user'
$dbName = 'monitoring'
## Prefer local script in `scripts/` if present, otherwise use container-mounted `/migrations/`
$hostSql = Join-Path (Resolve-Path ..\) 'Microservice_Monitoring\scripts\testdata.sql'
$containerSql = '/migrations/testdata.sql'

Write-Host "Running test scenarios SQL in container '$container'..."
$cleanup = @"
BEGIN;
DELETE FROM afwijking WHERE sensor_id IN (10001,10002,10003,10004);
DELETE FROM meting WHERE sensor_id IN (10001,10002,10003,10004);
DELETE FROM sensorconfiguratie WHERE sensor_id IN (10001,10002,10003,10004);
DELETE FROM sensor WHERE id IN (10001,10002,10003,10004);
DELETE FROM sensortype WHERE id IN (10001,10002,10003,10004);
DELETE FROM kunstwerk WHERE id IN (10001,10002,10003,10004);
COMMIT;
"@

Write-Host "Cleaning previous test rows..."
docker compose exec -T $container psql -U $dbUser -d $dbName -c "$cleanup" | Write-Host

Write-Host ('-' * 60) -ForegroundColor DarkGray
Write-Host " Running test scenarios SQL in container '$container' " -ForegroundColor Cyan
Write-Host ('-' * 60) -ForegroundColor DarkGray
if (Test-Path $hostSql) {
    Write-Host "Using local SQL: $hostSql" -ForegroundColor DarkCyan
    Get-Content $hostSql -Raw | docker compose exec -T $container psql -U $dbUser -d $dbName -f - | Write-Host
} else {
    Write-Host "Using container SQL: $containerSql" -ForegroundColor DarkCyan
    docker compose exec -T $container psql -U $dbUser -d $dbName -f $containerSql | Write-Host
}
Write-Host ('-' * 60) -ForegroundColor DarkGray

# Query results
$query = @"
SELECT sensor_id, COUNT(*) AS total, SUM(CASE WHEN is_warning THEN 1 ELSE 0 END) AS warnings
FROM afwijking
WHERE sensor_id IN (10001,10002,10003,10004)
GROUP BY sensor_id
ORDER BY sensor_id;
"@


Write-Host "Collecting results..." -ForegroundColor Cyan
$result = docker compose exec -T $container psql -U $dbUser -d $dbName -t -A -F"," -c $query
if ($LASTEXITCODE -ne 0) { Write-Error "Failed to query DB."; exit 2 }

Write-Host "Cleaning created test rows..." -ForegroundColor Magenta
docker compose exec -T $container psql -U $dbUser -d $dbName -c "$cleanup"

# Parse results into map
$map = @{}
$result -split "`n" | ForEach-Object {
    if ([string]::IsNullOrWhiteSpace($_)) { return }
    $parts = $_ -split ","
    $sid = [int]$parts[0]
    $total = [int]$parts[1]
    $warnings = [int]$parts[2]
    $map[$sid] = @{ total = $total; warnings = $warnings }
}

# Expected values
$expected = @{
    10001 = @{ total = 3; warnings = 1 }
    10002 = @{ total = 3; warnings = 2 }
    10003 = @{ total = 1; warnings = 0 }
    10004 = @{ total = 4; warnings = 3 }
}

$allPassed = $true
foreach ($k in $expected.Keys) {
    $exp = $expected[$k]
    if ($map.ContainsKey($k)) {
        $got = $map[$k]
        if ($got.total -eq $exp.total -and $got.warnings -eq $exp.warnings) {
            Write-Host "Sensor ${k}: PASS (total=$($got.total), warnings=$($got.warnings))" -ForegroundColor Green
        } else {
            Write-Host "Sensor ${k}: FAIL - expected total=$($exp.total),warnings=$($exp.warnings) but got total=$($got.total),warnings=$($got.warnings)" -ForegroundColor Red
            $allPassed = $false
        }
    } else {
        Write-Host "Sensor ${k}: FAIL - no anomalies found (expected total=$($exp.total))" -ForegroundColor Red
        $allPassed = $false
    }
}

Write-Host ('-' * 60) -ForegroundColor DarkGray
if ($allPassed) {
    Write-Host "All tests PASSED" -ForegroundColor Green
    Write-Host ('=' * 60) -ForegroundColor DarkGreen
    exit 0
} else {
    Write-Host "Some tests FAILED" -ForegroundColor Red
    Write-Host ('=' * 60) -ForegroundColor DarkRed
    exit 1
}

$ScriptLocation= Split-Path -Parent $PSCommandPath
Set-Location (Split-Path -Parent $ScriptLocation)

$standalone= Resolve-Path (Join-Path $ScriptLocation '../standalone')
$ffs_windows = $standalone -replace '\\','/'
$ffs_windows = $ffs_windows -replace 'C:','/mnt/c'

Set-Location $standalone
go build .\main.go
bash -c "rsync -avzh $ffs_windows dev@145.239.95.36:/home/dev/rental-saas"

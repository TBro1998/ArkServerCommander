# 开发用：在 Windows 上直接运行 Go 后端和 Next 开发服务器，不构建镜像（dev.sh 的 Windows 版本）
# 用法：在 PowerShell 里 .\dev.ps1，Ctrl+C 同时停止两者
#       执行策略不允许运行脚本时：powershell -ExecutionPolicy Bypass -File .\dev.ps1
# 本文件保存为带 BOM 的 UTF-8，否则 Windows PowerShell 5.1 会把中文读成乱码
# 需要：Go 1.27、Node >= 20.9、Docker Desktop（Linux 容器）
# 环境变量可在运行前设置覆盖，例如 PowerShell 里：$env:PUBLIC_HOST = '1.2.3.4'; .\dev.ps1
# ARK 服务器文件放在 data-dev 下（Windows 磁盘，经 Docker Desktop 文件共享挂进容器），读写比 Linux 上慢
$ErrorActionPreference = 'Stop'

$Root = $PSScriptRoot

function Set-Default([string]$Name, [string]$Value) {
    if (-not [Environment]::GetEnvironmentVariable($Name)) {
        Set-Item -Path "env:$Name" -Value $Value
    }
}

foreach ($cmd in 'go', 'node', 'docker') {
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
        Write-Host "[dev] 找不到 $cmd，请先安装并加入 PATH" -ForegroundColor Red
        exit 1
    }
}

# 数据放在 data-dev，不和正式部署的 ./data 混用
Set-Default 'DB_PATH' (Join-Path $Root 'data-dev\ark_server.db')
Set-Default 'SERVER_PORT' '8081'
Set-Default 'LISTEN_HOST' '127.0.0.1'
Set-Default 'GIN_MODE' 'debug'
Set-Default 'LOG_LEVEL' 'debug'
Set-Default 'LOG_FORMAT' 'console'
Set-Default 'API_BASE' "http://127.0.0.1:$env:SERVER_PORT/api"
$UiPort = if ($env:UI_PORT) { $env:UI_PORT } else { '3000' }

# 后端的 Docker SDK 只认 DOCKER_HOST，不读 docker context；让它和 docker 命令连同一个本机引擎
if (-not $env:DOCKER_HOST) {
    $dockerHost = (docker context inspect --format '{{.Endpoints.docker.Host}}' 2>$null | Out-String).Trim()
    if ($dockerHost -like 'npipe:*') {
        $env:DOCKER_HOST = $dockerHost
    }
}

New-Item -ItemType Directory -Force -Path (Split-Path -Parent $env:DB_PATH) | Out-Null

if (-not (Test-Path (Join-Path $Root 'ui\node_modules'))) {
    Write-Host '[dev] 安装前端依赖 (npm ci)'
    Push-Location (Join-Path $Root 'ui')
    try {
        & npm.cmd ci
        if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    } finally {
        Pop-Location
    }
}

# 两个进程共用当前控制台，Ctrl+C 会同时发给它们。
# 前端直接用 node 运行 next，不经过 npx.cmd，避免 Ctrl+C 后 cmd 询问“终止批处理操作吗”
$backend = $null
$frontend = $null
try {
    Write-Host "[dev] 后端: http://$($env:LISTEN_HOST):$($env:SERVER_PORT)  (Swagger: /swagger/index.html)"
    $backend = Start-Process -FilePath 'go' -ArgumentList 'run', '.' `
        -WorkingDirectory (Join-Path $Root 'server') -NoNewWindow -PassThru
    $null = $backend.Handle # 缓存句柄，进程退出后才能读到 ExitCode

    Write-Host "[dev] 前端: http://localhost:$UiPort"
    $frontend = Start-Process -FilePath 'node' `
        -ArgumentList 'node_modules/next/dist/bin/next', 'dev', '--turbopack', '--port', $UiPort `
        -WorkingDirectory (Join-Path $Root 'ui') -NoNewWindow -PassThru
    $null = $frontend.Handle

    # 任一进程退出就结束另一个
    while (-not $backend.HasExited -and -not $frontend.HasExited) {
        Start-Sleep -Milliseconds 500
    }
    $exited = if ($backend.HasExited) { $backend } else { $frontend }
    $code = $exited.ExitCode
    if ($code -eq 0) { $code = 1 }
    exit $code
} finally {
    # 按进程树结束：go run 编译出的子进程、Next 的工作进程都要一起停掉
    foreach ($p in $backend, $frontend) {
        if ($p -and -not $p.HasExited) {
            & taskkill.exe /T /F /PID $p.Id 2>&1 | Out-Null
        }
    }
}

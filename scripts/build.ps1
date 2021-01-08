cd ../

$root = "github.com/qianjunakasumi/project-shizuku/"
$commitid = git rev-parse --short master

"开始检查代码错误..."

go vet cmd/shizuku/main.go
scripts/golangci-lint.exe run

"请确认，按任意键开始编译"
[Console]::ReadKey() | Out-Null

"开始编译..."

$env:GOOS="linux"
go build -ldflags "-w -X ${root}configs.CommitId=${commitid}" -o build/shizuku cmd/shizuku/main.go

"编译完成，按任意键退出..."
[Console]::ReadKey() | Out-Null

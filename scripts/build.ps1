cd ../

$root = "github.com/qianjunakasumi/project-shizuku/"
$commitid = git rev-parse --short master

go build -ldflags "-w -X ${root}configs.CommitId=${commitid}" -o shizuku.exe cmd/shizuku/main.go

$env:GOOS="linux"
go build -ldflags "-w -X ${root}configs.CommitId=${commitid}" -o shizuku cmd/shizuku/main.go

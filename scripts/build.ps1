cd ../

$root = "github.com/qianjunakasumi/project-shizuku/"
$commitid = git rev-parse --short master
$time = Get-Date -Format "yyMMddHH"

go build -ldflags "-w -X ${root}configs.CommitId=${commitid} -X ${root}configs.BuildTime=${time}" -o shizuku.exe cmd/shizuku/main.go

./upx.exe shizuku.exe
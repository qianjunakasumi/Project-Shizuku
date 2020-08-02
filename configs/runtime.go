package configs

import "github.com/qianjunakasumi/project-shizuku/internal/uehara/messagechain"

var (
	QuoteJob = make(map[uint32]*func(call string, info *messagechain.MessageInfo) (*messagechain.MessageChain, error))
)

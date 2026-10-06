package restic

import (
	sshaccess "handfree-work/octo-backup/internal/plugins/access/ssh"
	"handfree-work/octo-backup/internal/svc"
)

// ClientConfig 保存 ResticClient 的执行环境和仓库配置。
type ClientConfig struct {
	Environment        ResticEnvironment
	ExecutionAccess    *sshaccess.SshAccess
	RepositoryAccess   *sshaccess.SshAccess
	RepositoryPath     string
	RepositoryPassword string
	Version            string
	RemoteDirectory    string
	ResticBinary       []byte
	ResticBinaryPath   string
}

// BackupRequest 保存一次备份所需的来源数据参数。
type BackupRequest struct {
	RepositoryTag string
	SourcePaths   []string
	ExcludePaths  []string
}

// RestoreRequest 保存还原一个快照所需的参数。
type RestoreRequest struct {
	SnapshotId string
	TargetPath string
}

// BackupResult 表示一次备份执行的最终状态。
type BackupResult struct {
	Status string `json:"status"`
	Output string `json:"output,omitempty"`
}

// ResticEnvironment 指定 Restic 命令的执行位置。
type ResticEnvironment string

const (
	// LocalResticEnvironment 表示在当前进程所在主机执行 Restic。
	LocalResticEnvironment ResticEnvironment = "local"
	// RemoteResticEnvironment 表示通过 SSH 在远程主机执行 Restic。
	RemoteResticEnvironment ResticEnvironment = "remote"
)

// NewResticClient 创建拥有完整运行环境的 Restic 客户端。
func NewResticClient(svcCtx *svc.ServiceContext, config ClientConfig) *ResticClient {
	if config.Version == "" && svcCtx != nil {
		config.Version = svcCtx.Restic.Version
	}
	client := &ResticClient{config: config}
	if config.Environment == LocalResticEnvironment {
		client.executor = &LocalResticExecutor{}
	} else {
		client.executor = &RemoteResticExecutor{Access: config.ExecutionAccess}
	}
	return client
}

package restic

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
	"handfree-work/octo-backup/internal/base/error_"
	"handfree-work/octo-backup/internal/base/log_"
	sshaccess "handfree-work/octo-backup/internal/plugins/access/ssh"
)

// RemoteResticExecutor 通过已连接的 SSH 主机执行 Restic 命令。
type RemoteResticExecutor struct {
	Access     *sshaccess.SshAccess
	BinaryPath string
}

func repositoryURL(config *ClientConfig) string {
	return "sftp:" + config.RepositoryAccess.Username + "@" + config.RepositoryAccess.Host + ":" + config.RepositoryPath
}

// Prepare 建立远程 SSH 执行连接，并确认远程主机可以执行基础命令。
func (e *RemoteResticExecutor) Prepare(ctx context.Context, config *ClientConfig) error {
	if e.Access == nil {
		return error_.NewTextError("远程 Restic 执行主机不能为空")
	}
	if err := e.Access.Connect(ctx); err != nil {
		return error_.NewWrapError("连接远程 Restic 执行主机失败", err)
	}
	if _, err := e.Access.Execute(ctx, "printf restic-remote-ready"); err != nil {
		_ = e.Access.Close()
		return error_.NewWrapError("远程 Restic 执行环境不可用", err)
	}
	remoteDirectory := strings.TrimSpace(config.RemoteDirectory)
	if remoteDirectory == "" {
		remoteDirectory = ".octo_backup"
	}
	binaryPath := path.Join(remoteDirectory, "restic")
	infoPath := path.Join(remoteDirectory, "restic.info")
	infoOutput, infoErr := e.Access.Execute(ctx, "cat "+shellQuote(infoPath))
	var info resticRemoteInfo
	infoValid := infoErr == nil && json.Unmarshal(infoOutput, &info) == nil && resticVersionMatches(info.Version, config.Version)
	binaryOutput, binaryErr := e.Access.Execute(ctx, "if [ -x "+shellQuote(binaryPath)+" ]; then printf yes; else printf no; fi")
	if infoValid && binaryErr == nil && strings.TrimSpace(string(binaryOutput)) == "yes" {
		e.BinaryPath = binaryPath
		config.ResticBinaryPath = binaryPath
		log_.Logger.Info("复用远程 restic 二进制", zap.String("path", binaryPath), zap.String("version", info.Version))
		return nil
	}
	osName, arch, err := e.Access.Platform(ctx)
	if err != nil {
		return error_.NewWrapError("探测远程 Restic 执行主机平台失败", err)
	}
	config.ResticBinary, err = downloadResticBinary(osName+"\n"+arch, config.Version)
	if err != nil {
		return error_.NewWrapError("准备远程 Restic 二进制失败", err)
	}
	if err := e.Access.Upload(ctx, binaryPath, config.ResticBinary, 0700); err != nil {
		return error_.NewWrapError("上传远程 Restic 二进制失败", err)
	}
	infoData, err := json.Marshal(resticRemoteInfo{Version: strings.TrimPrefix(config.Version, "v"), Size: len(config.ResticBinary), UploadedAt: time.Now().Format(time.RFC3339)})
	if err != nil {
		return error_.NewWrapError("生成远程 Restic 信息文件失败", err)
	}
	if err := e.Access.Upload(ctx, infoPath, infoData, 0600); err != nil {
		return error_.NewWrapError("写入远程 Restic 信息文件失败", err)
	}
	config.ResticBinaryPath = binaryPath
	e.BinaryPath = binaryPath
	return nil
}

// PrepareRepository 上传仓库环境配置，并在仓库不存在时初始化仓库。
func (e *RemoteResticExecutor) PrepareRepository(ctx context.Context, config *ClientConfig) error {
	if e.Access == nil {
		return error_.NewTextError("远程 Restic 执行主机不能为空")
	}
	configPath := path.Join(strings.TrimSpace(config.RemoteDirectory), "restic.env")
	if strings.TrimSpace(config.RemoteDirectory) == "" {
		configPath = path.Join(".octo_backup", "restic.env")
	}
	env := fmt.Sprintf("RESTIC_REPOSITORY=%s\nRESTIC_PASSWORD=%s\n", shellQuote(repositoryURL(config)), shellQuote(config.RepositoryPassword))
	if err := e.Access.Upload(ctx, configPath, []byte(env), 0600); err != nil {
		return error_.NewWrapError("上传远程 Restic 仓库配置失败", err)
	}
	binaryPath := config.ResticBinaryPath
	if binaryPath == "" {
		binaryPath = path.Join(path.Dir(configPath), "restic")
	}
	command := "set -a; . " + shellQuote(configPath) + "; set +a; " + shellQuote(binaryPath) + " -o sftp.command=" + shellQuote(repositorySftpCommand(config)) + " snapshots >/dev/null 2>&1 || " + shellQuote(binaryPath) + " -o sftp.command=" + shellQuote(repositorySftpCommand(config)) + " init"
	if _, err := e.Access.Execute(ctx, command); err != nil {
		return error_.NewWrapError("初始化远程 Restic 仓库失败", err)
	}
	return nil
}

func repositorySftpCommand(config *ClientConfig) string {
	return "ssh -o BatchMode=yes -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new -p " + fmt.Sprint(config.RepositoryAccess.Port) + " " + shellQuote(config.RepositoryAccess.Username+"@"+config.RepositoryAccess.Host) + " -s sftp"
}

// Execute 在远程主机执行 Restic 命令，并通过 shell 环境赋值传入变量。
func (e RemoteResticExecutor) Execute(ctx context.Context, args []string, environment []string) ([]byte, error) {
	command := make([]string, 0, len(environment)+len(args)+1)
	for _, value := range environment {
		key, variable, found := strings.Cut(value, "=")
		if found {
			command = append(command, key+"="+shellQuote(variable))
		}
	}
	command = append(command, shellQuote(e.BinaryPath))
	for _, argument := range args {
		command = append(command, shellQuote(argument))
	}
	return e.Access.Execute(ctx, strings.Join(command, " "))
}

// Close 关闭远程 SSH 执行连接。
func (e *RemoteResticExecutor) Close() error {
	if e.Access == nil {
		return nil
	}
	return e.Access.Close()
}

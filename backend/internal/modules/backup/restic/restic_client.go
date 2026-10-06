package restic

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"handfree-work/octo-backup/internal/base/error_"
	"handfree-work/octo-backup/internal/base/log_"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// ResticClient 通过 Restic CLI 实现备份和还原。
type ResticClient struct {
	config   ClientConfig
	executor ResticExecutor
}

// Snapshots 查询仓库快照，统一由 provider 准备本地 Restic 环境并执行命令。
func (p *ResticClient) Snapshots(ctx context.Context) (any, error) {
	defer p.executor.Close()
	if err := p.prepareEnvironment(ctx); err != nil {
		return nil, err
	}
	repository := "sftp:" + p.config.RepositoryAccess.Username + "@" + p.config.RepositoryAccess.Host + ":" + p.config.RepositoryPath
	args := []string{"-r", repository, "snapshots", "--json"}
	output, err := p.executor.Execute(ctx, args, []string{"RESTIC_PASSWORD=" + p.config.RepositoryPassword})
	if err != nil {
		return nil, error_.NewTextError("查询 Restic 仓库快照失败: output=%s error=%v", strings.TrimSpace(string(output)), err)
	}
	var snapshots any
	if err := yaml.Unmarshal(output, &snapshots); err != nil {
		return nil, error_.NewWrapError("解析 Restic 快照失败", err)
	}
	return snapshots, nil
}

// prepareEnvironment 让具体执行器准备本地或远程 Restic 环境。
func (p *ResticClient) prepareEnvironment(ctx context.Context) error {
	if p.config.RepositoryAccess == nil || strings.TrimSpace(p.config.RepositoryPath) == "" || p.config.RepositoryPassword == "" {
		return error_.NewTextError("Restic 仓库 SSH 配置、路径和密码不能为空")
	}
	return p.executor.Prepare(ctx, &p.config)
}

type resticRemoteInfo struct {
	Version    string `json:"version"`
	Platform   string `json:"platform,omitempty"`
	Arch       string `json:"arch,omitempty"`
	Size       int    `json:"size"`
	UploadedAt string `json:"uploadedAt"`
}

// StringSlice 将插件配置值转换为字符串列表。
func StringSlice(value any) []string {
	var out []string
	switch values := value.(type) {
	case []any:
		for _, item := range values {
			out = append(out, fmt.Sprint(item))
		}
	case []string:
		out = values
	case string:
		for _, item := range strings.Split(values, "\n") {
			if strings.TrimSpace(item) != "" {
				out = append(out, strings.TrimSpace(item))
			}
		}
	}
	return out
}

// Restore 将快照还原到指定目标路径。
func (p *ResticClient) Restore(ctx context.Context, request RestoreRequest) (string, error) {
	defer p.executor.Close()
	if strings.TrimSpace(p.config.RepositoryPath) == "" || p.config.RepositoryPassword == "" || strings.TrimSpace(request.SnapshotId) == "" || strings.TrimSpace(request.TargetPath) == "" {
		return "", error_.NewTextError("restic 还原参数不完整")
	}
	if err := p.prepareEnvironment(ctx); err != nil {
		return "", err
	}
	executor := p.executor
	args := []string{"restore", request.SnapshotId, "--target", request.TargetPath}
	log_.Logger.Info("执行 Restic 快照还原", zap.String("snapshotId", request.SnapshotId), zap.String("targetPath", request.TargetPath))
	output, err := executor.Execute(ctx, args, []string{"RESTIC_REPOSITORY=" + p.repositoryURL(), "RESTIC_PASSWORD=" + p.config.RepositoryPassword})
	if err != nil {
		return strings.TrimSpace(string(output)), err
	}
	return strings.TrimSpace(string(output)), nil
}

// Backup 准备远程 Restic 运行环境，按需初始化仓库并执行备份。
func (p *ResticClient) Backup(ctx context.Context, request BackupRequest) (*BackupResult, error) {
	defer p.executor.Close()
	started := time.Now()
	log_.Logger.Info("开始执行 restic 备份", zap.Int("sourcePathCount", len(request.SourcePaths)), zap.Int("excludePathCount", len(request.ExcludePaths)))
	if len(request.SourcePaths) == 0 || p.config.RepositoryPassword == "" {
		return nil, error_.NewTextError("restic 备份参数不完整")
	}
	if err := p.prepareEnvironment(ctx); err != nil {
		return nil, err
	}
	if err := p.executor.PrepareRepository(ctx, &p.config); err != nil {
		return nil, err
	}
	executor := p.executor
	args := []string{"backup", "--json"}
	if p.sftpCommand() != "" {
		args = append([]string{"-o", "sftp.command=" + p.sftpCommand()}, args...)
	}
	if strings.TrimSpace(request.RepositoryTag) != "" {
		args = append(args, "--tag", request.RepositoryTag)
	}
	for _, excludePath := range request.ExcludePaths {
		if strings.TrimSpace(excludePath) != "" {
			args = append(args, "--exclude", excludePath)
		}
	}
	args = append(args, request.SourcePaths...)
	log_.Logger.Info("执行 Restic 文件备份", zap.Int("sourcePathCount", len(request.SourcePaths)))
	output, err := executor.Execute(ctx, args, []string{"RESTIC_REPOSITORY=" + p.repositoryURL(), "RESTIC_PASSWORD=" + p.config.RepositoryPassword})
	if err != nil {
		return nil, err
	}
	log_.Logger.Info("restic 备份完成", zap.Duration("duration", time.Since(started)))
	return &BackupResult{Status: "success", Output: strings.TrimSpace(string(output))}, nil
}

func (p *ResticClient) repositoryURL() string {
	return "sftp:" + p.config.RepositoryAccess.Username + "@" + p.config.RepositoryAccess.Host + ":" + p.config.RepositoryPath
}
func (p *ResticClient) sftpCommand() string {
	return "ssh -o BatchMode=yes -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new -p " + fmt.Sprint(p.config.RepositoryAccess.Port) + " " + shellQuote(p.config.RepositoryAccess.Username+"@"+p.config.RepositoryAccess.Host) + " -s sftp"
}

func resticVersionMatches(output, expected string) bool {
	expected = strings.TrimPrefix(strings.TrimSpace(expected), "v")
	for _, field := range strings.Fields(output) {
		field = strings.TrimPrefix(field, "v")
		if field == expected {
			return true
		}
	}
	return false
}

func resticPlatform(output string) (string, string, error) {
	parts := strings.Fields(output)
	if len(parts) < 2 {
		return "", "", error_.NewTextError("无法识别远程主机系统类型: %s", output)
	}
	osName := strings.ToLower(parts[0])
	arch := strings.ToLower(parts[1])
	if osName != "linux" && osName != "windows" && osName != "darwin" {
		return "", "", error_.NewTextError("暂不支持远程系统: %s", parts[0])
	}
	switch arch {
	case "x86_64", "amd64":
		arch = "amd64"
	case "aarch64", "arm64":
		arch = "arm64"
	case "i386", "i686", "386":
		arch = "386"
	default:
		return "", "", error_.NewTextError("暂不支持远程架构: %s", parts[1])
	}
	return osName, arch, nil
}

func shellQuote(value string) string { return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'" }

func LoadResticBinary() ([]byte, error) {
	file := strings.TrimSpace(os.Getenv("OCTO_RESTIC_BINARY"))
	if file == "" {
		return nil, nil
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, error_.NewWrapError("读取 restic 二进制失败", err)
	}
	return data, nil
}

// DownloadLocalResticBinary 按当前系统平台下载匹配的 Restic 二进制。
func DownloadLocalResticBinary(version string) ([]byte, error) {
	platform := runtime.GOOS + "\n" + runtime.GOARCH
	cacheDir := filepath.Join("data", "restic")
	cacheFile := filepath.Join(cacheDir, "restic-"+strings.TrimPrefix(version, "v")+"-"+runtime.GOOS+"-"+runtime.GOARCH)
	infoFile := filepath.Join(cacheDir, "restic.info")
	var info resticRemoteInfo
	if infoData, readErr := os.ReadFile(infoFile); readErr == nil && json.Unmarshal(infoData, &info) == nil &&
		info.Version == strings.TrimPrefix(version, "v") && info.Platform == runtime.GOOS && info.Arch == runtime.GOARCH {
		if data, readErr := os.ReadFile(cacheFile); readErr == nil && len(data) > 0 {
			log_.Logger.Info("复用本地 restic 二进制", zap.String("path", cacheFile), zap.String("version", info.Version), zap.String("platform", runtime.GOOS+"/"+runtime.GOARCH))
			return data, nil
		}
	}
	log_.Logger.Info("本地 restic 缓存不存在或不匹配，开始下载", zap.String("path", cacheFile), zap.String("version", strings.TrimPrefix(version, "v")), zap.String("platform", runtime.GOOS+"/"+runtime.GOARCH))
	data, err := downloadResticBinary(platform, version)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(cacheDir, 0o755); err == nil {
		_ = os.WriteFile(cacheFile, data, 0o755)
		infoData, _ := json.Marshal(resticRemoteInfo{Version: strings.TrimPrefix(version, "v"), Platform: runtime.GOOS, Arch: runtime.GOARCH, Size: len(data), UploadedAt: time.Now().Format(time.RFC3339)})
		_ = os.WriteFile(infoFile, infoData, 0o644)
	}
	return data, nil
}

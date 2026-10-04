package ssh

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"handfree-work/octo-backup/internal/base/error_"
	"handfree-work/octo-backup/internal/modules/plugin"

	"golang.org/x/crypto/ssh"
)

//go:embed metadata.yaml
var metadata []byte

type Provider struct{}

func (Provider) Definition() (*plugin.Definition, error) {
	return plugin.NewDefinition(metadata, NewPluginInstance)
}

func NewPluginInstance(config map[string]any) (plugin.ActionExecutor, error) {
	return NewSshAccess(config), nil
}

type SshAccess struct {
	config  map[string]any
	client  *ssh.Client
}

func NewSshAccess(config map[string]any) *SshAccess { return &SshAccess{config: config} }

func (s *SshAccess) ExecuteAction(_ context.Context, action string, _ map[string]any) (any, error) {
	if action == "onTest" {
		return map[string]any{"ok": true}, nil
	}
	return nil, plugin.ErrActionNotFound
}

func (s *SshAccess) Connect(ctx context.Context) error {
	host, user := strings.TrimSpace(fmt.Sprint(s.config["host"])), strings.TrimSpace(fmt.Sprint(s.config["username"]))
	if host == "" || user == "" {
		return error_.NewTextError("SSH 主机和用户名不能为空")
	}
	port, timeout := intValue(s.config["port"], 22), time.Duration(intValue(s.config["timeout"], 30))*time.Second
	auth, err := authMethods(s.config)
	if err != nil {
		return err
	}
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := (&net.Dialer{Timeout: timeout}).DialContext(ctx, "tcp", address)
	if err != nil {
		return error_.NewTextError(fmt.Sprintf("连接 SSH 主机失败: %v", err))
	}
	clientConn, chans, reqs, err := ssh.NewClientConn(conn, address, &ssh.ClientConfig{User: user, Auth: auth, HostKeyCallback: ssh.InsecureIgnoreHostKey(), Timeout: timeout})
	if err != nil {
		_ = conn.Close()
		return error_.NewTextError(fmt.Sprintf("SSH 认证失败: %v", err))
	}
	s.client = ssh.NewClient(clientConn, chans, reqs)
	return nil
}

func (c *SshAccess) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

func (c *SshAccess) ListDirectories(path string) ([]map[string]any, error) {
	if c == nil || c.client == nil {
		return nil, error_.NewTextError("SSH 尚未连接")
	}
	if strings.TrimSpace(path) == "" {
		path = "/"
	}
	session, err := c.client.NewSession()
	if err != nil {
		return nil, error_.NewTextError(fmt.Sprintf("创建 SSH 会话失败: %v", err))
	}
	defer session.Close()
	command := "find " + shellQuote(path) + " -mindepth 1 -maxdepth 1 -type d -print"
	output, err := session.CombinedOutput(command)
	if err != nil {
		detail := strings.TrimSpace(string(output))
		if detail == "" {
			return nil, error_.NewTextError(fmt.Sprintf("读取远程目录失败: %v", err))
		}
		return nil, error_.NewTextError(fmt.Sprintf("读取远程目录失败: %v: %s", err, detail))
	}
	paths := make([]map[string]any, 0)
	for _, line := range strings.Split(string(output), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			paths = append(paths, map[string]any{"title": line, "key": line, "isLeaf": false})
		}
	}
	return paths, nil
}

func authMethods(config map[string]any) ([]ssh.AuthMethod, error) {
	if password := fmt.Sprint(config["password"]); password != "" {
		return []ssh.AuthMethod{ssh.Password(password)}, nil
	}
	key := strings.TrimSpace(fmt.Sprint(config["privateKey"]))
	if key == "" {
		return nil, error_.NewTextError("SSH 密码或私钥不能为空")
	}
	var signer ssh.Signer
	var err error
	if passphrase := fmt.Sprint(config["privateKeyPassword"]); passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(key), []byte(passphrase))
	} else {
		signer, err = ssh.ParsePrivateKey([]byte(key))
	}
	if err != nil {
		return nil, error_.NewTextError(fmt.Sprintf("解析 SSH 私钥失败: %v", err))
	}
	return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
}

func intValue(value any, fallback int) int {
	if parsed, err := strconv.Atoi(strings.TrimSpace(fmt.Sprint(value))); err == nil && parsed > 0 {
		return parsed
	}
	return fallback
}
func shellQuote(value string) string { return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'" }

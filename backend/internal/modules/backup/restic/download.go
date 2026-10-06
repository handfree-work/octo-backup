package restic

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"handfree-work/octo-backup/internal/base/error_"
	"handfree-work/octo-backup/internal/base/log_"

	"github.com/dsnet/compress/bzip2"
	"go.uber.org/zap"
)

type resticAsset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

type resticRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []resticAsset `json:"assets"`
}

func downloadResticBinary(platform, version string) ([]byte, error) {
	var release resticRelease
	version = strings.TrimSpace(version)
	response, err := (&http.Client{Timeout: 30 * time.Second}).Get(resticReleaseURL(version))
	if err != nil {
		return nil, error_.NewWrapError("获取 restic 版本失败", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, error_.NewTextError("获取 restic 版本失败: HTTP %d", response.StatusCode)
	}
	if err := json.NewDecoder(response.Body).Decode(&release); err != nil {
		return nil, error_.NewWrapError("解析 restic 版本失败", err)
	}
	osName, arch, err := resticPlatform(platform)
	if err != nil {
		return nil, err
	}
	githubAsset, ok := findResticAsset(release.Assets, osName, arch)
	if !ok {
		return nil, error_.NewTextError("找不到匹配的 restic 下载包: version=%s platform=%s/%s", release.TagName, osName, arch)
	}
	selected := githubAsset
	atomgitRelease, atomgitErr := fetchAtomGitRelease()
	atomgitAsset, atomgitOK := findResticAsset(atomgitRelease.Assets, osName, arch)
	if atomgitErr == nil && atomgitOK && atomgitRelease.TagName == release.TagName {
		selected = chooseFasterAsset(githubAsset, atomgitAsset)
	}
	log_.Logger.Info("下载 restic", zap.String("version", release.TagName), zap.String("url", selected.URL), zap.String("source", assetSource(selected, githubAsset, atomgitAsset)))
	data, err, retryable := downloadResticAsset(selected)
	if err != nil && retryable && selected.URL != githubAsset.URL {
		data, err, _ = downloadResticAsset(githubAsset)
	}
	return data, err
}

func downloadResticAsset(asset resticAsset) ([]byte, error, bool) {
	response, err := (&http.Client{Timeout: 30 * time.Second}).Get(asset.URL)
	if err != nil {
		return nil, error_.NewWrapError("下载 restic 失败", err), true
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, error_.NewTextError("下载 restic 失败: HTTP %d", response.StatusCode), true
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, error_.NewWrapError("读取 restic 下载内容失败", err), false
	}
	if strings.HasSuffix(strings.ToLower(asset.Name), ".zip") {
		archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return nil, error_.NewWrapError("解压 restic ZIP 安装包失败", err), false
		}
		for _, file := range archive.File {
			if file.FileInfo().IsDir() {
				continue
			}
			reader, err := file.Open()
			if err != nil {
				return nil, error_.NewWrapError("读取 restic ZIP 安装包失败", err), false
			}
			binary, readErr := io.ReadAll(reader)
			_ = reader.Close()
			if readErr != nil {
				return nil, error_.NewWrapError("读取 restic ZIP 压缩包失败", readErr), false
			}
			return binary, nil, false
		}
		return nil, error_.NewTextError("restic 压缩包为空"), false
	}
	bz, err := bzip2.NewReader(bytes.NewReader(data), nil)
	if err != nil {
		return nil, error_.NewWrapError("解压 restic BZ2 安装包失败", err), false
	}
	binary, err := io.ReadAll(bz)
	if err != nil {
		return nil, error_.NewWrapError("读取 restic BZ2 压缩包失败", err), false
	}
	if len(binary) == 0 {
		return nil, error_.NewTextError("restic 压缩包为空"), false
	}
	return binary, nil, false
}

func resticReleaseURL(version string) string {
	if version == "" {
		return "https://api.github.com/repos/restic/restic/releases/latest"
	}
	return "https://api.github.com/repos/restic/restic/releases/tags/" + url.PathEscape("v"+strings.TrimPrefix(version, "v"))
}

func findResticAsset(assets []resticAsset, osName, arch string) (resticAsset, bool) {
	for _, asset := range assets {
		if strings.HasPrefix(asset.Name, "restic_") && strings.Contains(asset.Name, "_"+osName+"_") && strings.Contains(asset.Name, "_"+arch+".") && (strings.HasSuffix(asset.Name, ".bz2") || (osName == "windows" && strings.HasSuffix(asset.Name, ".zip"))) {
			return asset, true
		}
	}
	return resticAsset{}, false
}

func fetchAtomGitRelease() (resticRelease, error) {
	response, err := (&http.Client{Timeout: 5 * time.Second}).Get("https://api.atomgit.com/api/v5/repos/handfree-work/restic/releases/latest")
	if err != nil {
		return resticRelease{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return resticRelease{}, error_.NewTextError("HTTP %d", response.StatusCode)
	}
	var release resticRelease
	err = json.NewDecoder(response.Body).Decode(&release)
	return release, err
}

func probeResticSource(sourceURL string) time.Duration {
	started := time.Now()
	response, err := (&http.Client{Timeout: 5 * time.Second}).Get(sourceURL)
	if err != nil {
		log_.Logger.Warn("测速 restic 发布信息接口失败", zap.String("url", sourceURL), zap.Error(err))
		return time.Hour
	}
	_, _ = io.Copy(io.Discard, response.Body)
	_ = response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 400 {
		return time.Hour
	}
	duration := time.Since(started)
	log_.Logger.Info("restic 发布信息接口测速完成", zap.String("url", sourceURL), zap.Duration("duration", duration))
	return duration
}

func chooseFasterAsset(github, atomgit resticAsset) resticAsset {
	var githubDuration, atomgitDuration time.Duration
	var wait sync.WaitGroup
	wait.Add(2)
	go func() { defer wait.Done(); githubDuration = probeResticSource("https://github.com/manifest.json") }()
	go func() {
		defer wait.Done()
		atomgitDuration = probeResticSource("https://api.atomgit.com/api/v5/repos/handfree-work/restic/releases/latest")
	}()
	wait.Wait()
	if atomgitDuration < githubDuration {
		return atomgit
	}
	return github
}

func assetSource(selected, github, atomgit resticAsset) string {
	if selected.URL == atomgit.URL {
		return "atomgit"
	}
	return "github"
}

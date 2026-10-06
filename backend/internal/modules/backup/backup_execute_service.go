package logic

import (
	"context"
	"go.uber.org/zap"
	"time"

	"handfree-work/octo-backup/internal/base/db_"
	"handfree-work/octo-backup/internal/base/error_"
	"handfree-work/octo-backup/internal/base/log_"
	"handfree-work/octo-backup/internal/models"
	restic "handfree-work/octo-backup/internal/modules/backup/restic"
	"handfree-work/octo-backup/internal/modules/plugin"
	pluginservice "handfree-work/octo-backup/internal/modules/plugin/service"
	sshaccess "handfree-work/octo-backup/internal/plugins/access/ssh"
	sftprepository "handfree-work/octo-backup/internal/plugins/repository/sftp"
	sshsource "handfree-work/octo-backup/internal/plugins/source/ssh"
	"handfree-work/octo-backup/internal/svc"
)

type BackupExecuteService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type BackupExecuteResult struct {
	Id     int64  `json:"id"`
	Status string `json:"status"`
}

func NewBackupExecuteService(ctx context.Context, s *svc.ServiceContext) *BackupExecuteService {
	return &BackupExecuteService{ctx: ctx, svcCtx: s}
}

func (s *BackupExecuteService) Execute(id int64, logId int64) (*BackupExecuteResult, error) {
	currentProgress := 0
	updateLog := func(status string, nextProgress int, stage, result, message string) {
		if status == "failed" && nextProgress == 100 {
			nextProgress = currentProgress
		}
		currentProgress = nextProgress
		now := time.Now().Unix()
		fields := map[string]any{"status": status, "progress": nextProgress, "stage": stage, "result": result, "error": message}
		if status == "success" || status == "failed" {
			fields["finished_at"] = now
		}
		_ = s.svcCtx.Db.Model(&models.BackupLog{}).Where("id = ?", logId).Updates(fields).Error
	}
	updateLog("running", 5, "读取备份计划", "", "")
	log_.Logger.Info("读取备份计划", zap.Int64("planId", id))
	dao := db_.New[models.BackupPlan](db_.NewCtx(s.ctx, s.svcCtx.Db))
	plan, err := dao.GetById(id)
	if err != nil {
		updateLog("failed", 100, "读取备份计划", "", err.Error())
		return nil, error_.NewApiError(err)
	}
	if plan == nil {
		updateLog("failed", 100, "读取备份计划", "", "备份计划不存在")
		return nil, error_.NewTextError("备份计划不存在")
	}
	if !plan.Enabled {
		updateLog("failed", 100, "校验备份计划", "", "已禁用的备份计划不能执行")
		return nil, error_.NewTextError("已禁用的备份计划不能执行")
	}
	setStatus := func(status, message string) {
		plan.LastStatus, plan.LastError = status, message
		now := time.Now().Unix()
		plan.LastRunAt = &now
		_, _ = dao.UpdateById(id, plan)
	}
	setStatus("running", "")
	log_.Logger.Info("备份计划状态更新为 running", zap.Int64("planId", id))
	instance := pluginservice.NewPluginInstanceService(s.ctx, s.svcCtx)
	source, err := instance.Instance(plan.SourceId)
	if err != nil {
		setStatus("failed", err.Error())
		return nil, err
	}
	result, err := s.doBackup(plan, source)
	if err != nil {
		updateLog("failed", 100, "执行备份", "", err.Error())
		setStatus("failed", err.Error())
		return nil, err
	}
	setStatus("success", "")
	updateLog("success", 100, "完成", result.Status, "")
	return &BackupExecuteResult{Id: id, Status: result.Status}, nil
}

func (s *BackupExecuteService) doBackup(plan *models.BackupPlan, source plugin.PluginInstance) (*restic.BackupResult, error) {
	sourceConfig, ok := source.(*sshsource.SshSource)
	if !ok {
		return nil, error_.NewTextError("备份来源插件不支持备份配置")
	}
	instance := pluginservice.NewPluginInstanceService(s.ctx, s.svcCtx)
	repositoryPlugin, err := instance.Instance(plan.RepositoryId)
	if err != nil {
		return nil, err
	}
	repositoryConfig, ok := repositoryPlugin.(*sftprepository.SftpRepository)
	if !ok {
		return nil, error_.NewTextError("存储仓库插件不支持备份配置")
	}
	accessIDValue := int64(sourceConfig.AccessId)
	if accessIDValue <= 0 {
		return nil, error_.NewTextError("SSH 授权配置不存在")
	}
	accessPlugin, err := instance.Instance(accessIDValue)
	if err != nil {
		return nil, err
	}
	sshConfig, ok := accessPlugin.(*sshaccess.SshAccess)
	if !ok {
		return nil, error_.NewTextError("SSH 授权插件类型无效")
	}
	paths := restic.StringSlice(sourceConfig.Paths)
	excludePaths := restic.StringSlice(sourceConfig.ExcludePaths)
	binary, err := restic.LoadResticBinary()
	if err != nil {
		return nil, err
	}
	repositoryAccessID := int64(repositoryConfig.AccessId)
	if repositoryAccessID <= 0 {
		return nil, error_.NewTextError("存储仓库 SSH 授权配置无效: %d", repositoryAccessID)
	}
	repositoryAccessPlugin, err := instance.Instance(repositoryAccessID)
	if err != nil {
		return nil, err
	}
	repositoryAccess, ok := repositoryAccessPlugin.(*sshaccess.SshAccess)
	if !ok {
		return nil, error_.NewTextError("存储仓库 SSH 授权插件类型无效")
	}
	clientConfig := restic.ClientConfig{
		Environment:        restic.RemoteResticEnvironment,
		ExecutionAccess:    sshConfig,
		RepositoryAccess:   repositoryAccess,
		RepositoryPath:     repositoryConfig.Path,
		RepositoryPassword: repositoryConfig.Password,
		Version:            s.svcCtx.Restic.Version,
		RemoteDirectory:    ".octo_backup",
		ResticBinary:       binary,
	}
	backupRequest := restic.BackupRequest{
		RepositoryTag: plan.RepoTag,
		SourcePaths:   paths,
		ExcludePaths:  excludePaths,
	}
	client := restic.NewResticClient(s.svcCtx, clientConfig)
	result, err := client.Backup(s.ctx, backupRequest)
	if err != nil {
		return nil, error_.NewWrapError("执行备份任务失败", err)
	}
	return result, nil
}

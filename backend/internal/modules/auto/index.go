package auto

import (
	"context"
	"go.uber.org/zap"
	"handfree-work/octo-backup/internal/base/log_"
	logic "handfree-work/octo-backup/internal/modules/backup"
	cronservice "handfree-work/octo-backup/internal/modules/cron"
	"handfree-work/octo-backup/internal/svc"
)

// Start 启动 cron 服务并注册各业务模块的自动任务。
func Start(serviceContext *svc.ServiceContext) (func(), error) {
	cronservice.Start()
	log_.Logger.Info("自动任务服务已启动")
	if err := logic.NewBackupPlanService(context.Background(), serviceContext).RegisterAll(); err != nil {
		log_.Logger.Error("自动任务批量注册失败", zap.Error(err))
		cronservice.Stop()
		return nil, err
	}
	return cronservice.Stop, nil
}

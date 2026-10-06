package cron

import (
	"context"
	"fmt"
	cronlib "github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"handfree-work/octo-backup/internal/base/log_"
	"sync"
)

// Service 提供按 type 和 id 管理的通用 cron 任务。
type Service struct {
	cron    *cronlib.Cron
	mu      sync.Mutex
	entries map[string]cronlib.EntryID
}

var defaultService *Service

// New 创建并启动 cron 服务。
func New() *Service {
	service := &Service{cron: cronlib.New(), entries: make(map[string]cronlib.EntryID)}
	service.cron.Start()
	log_.Logger.Info("cron 服务已启动")
	return service
}

// Start 启动全局 cron 服务。
func Start() { defaultService = New() }

// Register 注册全局 cron 任务。
func Register(taskType string, id int64, schedule string, handler func()) error {
	if defaultService == nil {
		return fmt.Errorf("cron 服务尚未启动")
	}
	return defaultService.Register(taskType, id, schedule, handler)
}

// Unregister 反注册全局 cron 任务。
func Unregister(taskType string, id int64) {
	if defaultService != nil {
		defaultService.Unregister(taskType, id)
	}
}

// Stop 停止全局 cron 服务。
func Stop() {
	if defaultService != nil {
		defaultService.Stop()
		defaultService = nil
	}
}

// ValidateSchedule 校验标准五字段 cron 表达式。
func ValidateSchedule(schedule string) error {
	_, err := cronlib.ParseStandard(schedule)
	if err != nil {
		return fmt.Errorf("cron 表达式无效: %w", err)
	}
	return nil
}

// Register 注册任务，同一 type 和 id 会先移除旧任务。
func (s *Service) Register(taskType string, id int64, schedule string, handler func()) error {
	key := fmt.Sprintf("%s:%d", taskType, id)
	s.Unregister(taskType, id)
	entryId, err := s.cron.AddFunc(schedule, handler)
	if err != nil {
		log_.Logger.Error("cron 任务注册失败", zap.String("taskType", taskType), zap.Int64("taskId", id), zap.String("schedule", schedule), zap.Error(err))
		return fmt.Errorf("注册 cron 任务 %s 失败: %w", key, err)
	}
	s.mu.Lock()
	s.entries[key] = entryId
	s.mu.Unlock()
	log_.Logger.Info("cron 任务已注册", zap.String("taskType", taskType), zap.Int64("taskId", id), zap.String("schedule", schedule), zap.Int("entryId", int(entryId)))
	return nil
}

// Unregister 移除指定任务。
func (s *Service) Unregister(taskType string, id int64) {
	key := fmt.Sprintf("%s:%d", taskType, id)
	s.mu.Lock()
	entryId, ok := s.entries[key]
	delete(s.entries, key)
	s.mu.Unlock()
	if ok {
		s.cron.Remove(entryId)
		log_.Logger.Info("cron 任务已反注册", zap.String("taskType", taskType), zap.Int64("taskId", id), zap.Int("entryId", int(entryId)))
	} else {
		log_.Logger.Info("cron 任务反注册跳过", zap.String("taskType", taskType), zap.Int64("taskId", id))
	}
}

// Stop 停止服务并等待正在执行的任务完成。
func (s *Service) Stop() context.Context {
	log_.Logger.Info("cron 服务正在停止")
	return s.cron.Stop()
}

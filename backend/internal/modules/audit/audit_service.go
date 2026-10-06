package audit

import (
	"context"
	"strings"

	"handfree-work/octo-backup/internal/base/error_"
	"handfree-work/octo-backup/internal/models"
	"handfree-work/octo-backup/internal/svc"
)

// AuditLogPageQuery 定义审计日志分页和筛选条件。
type AuditLogPageQuery struct {
	Offset, Limit  int64
	Username, Path string
}

// AuditLogPageResult 返回审计日志分页结果。
type AuditLogPageResult struct {
	Offset, Limit int64
	Records       []models.AuditLog
	Total         int64
}

// Service 提供审计日志查询能力。
type Service struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewService 创建审计日志服务。
func NewService(ctx context.Context, svcCtx *svc.ServiceContext) *Service {
	return &Service{ctx: ctx, svcCtx: svcCtx}
}

// Page 按用户和路径查询审计日志。
func (s *Service) Page(q *AuditLogPageQuery) (*AuditLogPageResult, error) {
	off, lim := int64(0), int64(20)
	if q != nil {
		off, lim = q.Offset, q.Limit
	}
	if off < 0 {
		off = 0
	}
	if lim <= 0 {
		lim = 20
	}
	if lim > 1000 {
		lim = 1000
	}
	db := s.svcCtx.Db.Model(&models.AuditLog{})
	if q != nil && strings.TrimSpace(q.Username) != "" {
		db = db.Where("username LIKE ?", "%"+strings.TrimSpace(q.Username)+"%")
	}
	if q != nil && strings.TrimSpace(q.Path) != "" {
		db = db.Where("path LIKE ?", "%"+strings.TrimSpace(q.Path)+"%")
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, error_.NewApiError(err)
	}
	var rows []models.AuditLog
	if err := db.Order("id desc").Offset(int(off)).Limit(int(lim)).Find(&rows).Error; err != nil {
		return nil, error_.NewApiError(err)
	}
	return &AuditLogPageResult{Offset: off, Limit: lim, Records: rows, Total: total}, nil
}

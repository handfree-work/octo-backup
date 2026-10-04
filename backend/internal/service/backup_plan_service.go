package logic

import (
	"context"
	"fmt"
	"path"
	"strings"

	"gorm.io/gorm"
	"handfree-work/octo-backup/internal/base/db_"
	"handfree-work/octo-backup/internal/base/error_"
	"handfree-work/octo-backup/internal/models"
	"handfree-work/octo-backup/internal/modules/plugin"
	"handfree-work/octo-backup/internal/svc"
)

type BackupPlanService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBackupPlanService(ctx context.Context, s *svc.ServiceContext) *BackupPlanService {
	return &BackupPlanService{ctx: ctx, svcCtx: s}
}

type BackupPlanInput struct {
	Name         string `json:"name"`
	SourceID     int64  `json:"sourceId"`
	RepositoryID int64  `json:"repositoryId"`
	RepoSubPath  string `json:"repoSubPath"`
	Schedule     string `json:"schedule"`
	Enabled      *bool  `json:"enabled"`
}
type BackupPlanPageQuery struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Name   string `json:"name"`
}
type BackupPlanPageResult struct {
	Offset  int64            `json:"offset"`
	Limit   int64            `json:"limit"`
	Records []map[string]any `json:"records"`
	Total   int64            `json:"total"`
}

func (s *BackupPlanService) validate(in *BackupPlanInput) error {
	if in == nil || strings.TrimSpace(in.Name) == "" || in.SourceID <= 0 || in.RepositoryID <= 0 || strings.TrimSpace(in.Schedule) == "" {
		return error_.NewTextError("备份计划名称、来源、仓库和调度不能为空")
	}
	normalized, err := normalizeRepoSubPath(in.RepoSubPath)
	if err != nil {
		return err
	}
	in.RepoSubPath = normalized
	dao := db_.New[models.PluginInstance](db_.NewCtx(s.ctx, s.svcCtx.Db))
	source, err := dao.GetById(in.SourceID)
	if err != nil || source == nil || source.PluginType != string(plugin.TypeSource) {
		return error_.NewTextError(fmt.Sprintf("备份来源不存在: %d", in.SourceID))
	}
	repo, err := dao.GetById(in.RepositoryID)
	if err != nil || repo == nil || repo.PluginType != string(plugin.TypeRepository) {
		return error_.NewTextError(fmt.Sprintf("存储仓库不存在: %d", in.RepositoryID))
	}
	return nil
}

func normalizeRepoSubPath(value string) (string, error) {
	value = strings.TrimSpace(strings.ReplaceAll(value, "\\", "/"))
	if value == "" {
		return "", nil
	}
	if strings.HasPrefix(value, "/") || strings.Contains(value, ":") {
		return "", error_.NewTextError("仓库目录必须是相对路径")
	}
	cleaned := path.Clean(value)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", error_.NewTextError("仓库目录不能越界")
	}
	return strings.Trim(cleaned, "/"), nil
}

func (s *BackupPlanService) checkRepoSubPath(id, repositoryID int64, subPath string) error {
	var count int64
	db := s.svcCtx.Db.Model(&models.BackupPlan{}).Where("repository_id = ? AND repo_sub_path = ?", repositoryID, subPath)
	if id > 0 {
		db = db.Where("id <> ?", id)
	}
	if err := db.Count(&count).Error; err != nil {
		return error_.NewTextError(fmt.Sprintf("检查仓库目录是否重复失败: %v", err))
	}
	if count > 0 {
		return error_.NewTextError("同一存储仓库下的目录已被其他备份计划使用")
	}
	return nil
}
func (s *BackupPlanService) Create(in *BackupPlanInput) (map[string]any, error) {
	if err := s.validate(in); err != nil {
		return nil, err
	}
	if err := s.checkRepoSubPath(0, in.RepositoryID, in.RepoSubPath); err != nil {
		return nil, err
	}
	enabled := true
	if in.Enabled != nil {
		enabled = *in.Enabled
	}
	row := &models.BackupPlan{Name: strings.TrimSpace(in.Name), SourceID: in.SourceID, RepositoryID: in.RepositoryID, RepoSubPath: in.RepoSubPath, Schedule: strings.TrimSpace(in.Schedule), Enabled: enabled, LastStatus: "pending"}
	if err := db_.New[models.BackupPlan](db_.NewCtx(s.ctx, s.svcCtx.Db)).Create(row); err != nil {
		return nil, err
	}
	return s.present(row), nil
}
func (s *BackupPlanService) Update(id int64, in *BackupPlanInput) (map[string]any, error) {
	dao := db_.New[models.BackupPlan](db_.NewCtx(s.ctx, s.svcCtx.Db))
	row, err := dao.GetById(id)
	if err != nil || row == nil {
		return nil, error_.NewTextError("备份计划不存在")
	}
	if err := s.validate(in); err != nil {
		return nil, err
	}
	if err := s.checkRepoSubPath(id, in.RepositoryID, in.RepoSubPath); err != nil {
		return nil, err
	}
	row.Name, row.SourceID, row.RepositoryID, row.RepoSubPath, row.Schedule = strings.TrimSpace(in.Name), in.SourceID, in.RepositoryID, in.RepoSubPath, strings.TrimSpace(in.Schedule)
	if in.Enabled != nil {
		row.Enabled = *in.Enabled
	}
	if _, err = dao.UpdateById(id, row); err != nil {
		return nil, err
	}
	return s.present(row), nil
}
func (s *BackupPlanService) Info(id int64) (map[string]any, error) {
	row, err := db_.New[models.BackupPlan](db_.NewCtx(s.ctx, s.svcCtx.Db)).GetById(id)
	if err != nil || row == nil {
		return nil, error_.NewTextError("备份计划不存在")
	}
	return s.present(row), nil
}
func (s *BackupPlanService) Delete(id int64) error {
	dao := db_.New[models.BackupPlan](db_.NewCtx(s.ctx, s.svcCtx.Db))
	n, err := dao.Delete(&id)
	if err != nil {
		return err
	}
	if n == 0 {
		return error_.NewTextError("备份计划不存在")
	}
	return nil
}
func (s *BackupPlanService) Page(q *BackupPlanPageQuery) (*BackupPlanPageResult, error) {
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
	page := &db_.Page{Start: off, Limit: lim}
	filters := []func(*gorm.DB){}
	if q != nil && strings.TrimSpace(q.Name) != "" {
		name := strings.TrimSpace(q.Name)
		filters = append(filters, func(db *gorm.DB) { db.Where("name LIKE ?", "%"+name+"%") })
	}
	rows, err := db_.New[models.BackupPlan](db_.NewCtx(s.ctx, s.svcCtx.Db)).FindPage(&db_.PageReq[models.BackupPlan]{Query: &models.BackupPlan{}, Page: page}, filters...)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(*rows))
	for i := range *rows {
		out = append(out, s.present(&(*rows)[i]))
	}
	return &BackupPlanPageResult{Offset: off, Limit: lim, Records: out, Total: page.Total}, nil
}
func (s *BackupPlanService) present(row *models.BackupPlan) map[string]any {
	return map[string]any{"id": row.Id, "name": row.Name, "sourceId": row.SourceID, "repositoryId": row.RepositoryID, "repoSubPath": row.RepoSubPath, "schedule": row.Schedule, "enabled": row.Enabled, "lastStatus": row.LastStatus, "lastError": row.LastError, "lastRunAt": row.LastRunAt, "createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt}
}

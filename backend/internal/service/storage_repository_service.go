package logic

import (
	"context"
	"fmt"
	"strings"

	"handfree-work/web-restic/internal/base/db_"
	"handfree-work/web-restic/internal/models"
	"handfree-work/web-restic/internal/svc"
)

var ErrRepositoryNotFound = fmt.Errorf("存储仓库不存在")

type StorageRepositoryInput struct {
	Name        string `json:"name"`
	Repository  string `json:"repository"`
	Password    string `json:"password"`
	Description string `json:"description"`
}
type UpdateStorageRepositoryInput struct {
	Name        *string `json:"name"`
	Repository  *string `json:"repository"`
	Password    *string `json:"password"`
	Description *string `json:"description"`
}
type StorageRepositoryPageQuery struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}
type StorageRepositoryPageResult struct {
	Offset  int64                       `json:"offset"`
	Limit   int64                       `json:"limit"`
	Records []*models.StorageRepository `json:"records"`
	Total   int64                       `json:"total"`
}

type StorageRepositoryService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStorageRepositoryService(ctx context.Context, svcCtx *svc.ServiceContext) *StorageRepositoryService {
	return &StorageRepositoryService{ctx: ctx, svcCtx: svcCtx}
}
func (s *StorageRepositoryService) Create(in *StorageRepositoryInput) (*models.StorageRepository, error) {
	if in == nil || strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Repository) == "" {
		return nil, fmt.Errorf("仓库名称和地址不能为空")
	}
	dao := db_.New[models.StorageRepository](db_.NewCtx(s.ctx, s.svcCtx.Db))
	row := &models.StorageRepository{Name: strings.TrimSpace(in.Name), Repository: strings.TrimSpace(in.Repository), Password: in.Password, Description: strings.TrimSpace(in.Description)}
	if err := dao.Create(row); err != nil {
		return nil, err
	}
	return row, nil
}
func (s *StorageRepositoryService) FindPage(q *StorageRepositoryPageQuery) (*StorageRepositoryPageResult, error) {
	offset, limit := int64(0), int64(20)
	if q != nil {
		offset, limit = q.Offset, q.Limit
	}
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 1000 {
		limit = 1000
	}
	dao := db_.New[models.StorageRepository](db_.NewCtx(s.ctx, s.svcCtx.Db))
	page := &db_.Page{Start: offset, Limit: limit}
	rows, err := dao.FindPage(&db_.PageReq[models.StorageRepository]{Query: &models.StorageRepository{}, Page: page})
	if err != nil {
		return nil, err
	}
	records := make([]*models.StorageRepository, 0, len(*rows))
	for i := range *rows {
		records = append(records, &(*rows)[i])
	}
	return &StorageRepositoryPageResult{Offset: offset, Limit: limit, Records: records, Total: page.Total}, nil
}
func (s *StorageRepositoryService) Update(id int64, in *UpdateStorageRepositoryInput) (*models.StorageRepository, error) {
	dao := db_.New[models.StorageRepository](db_.NewCtx(s.ctx, s.svcCtx.Db))
	row, err := dao.GetById(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrRepositoryNotFound
	}
	if in == nil {
		return row, nil
	}
	if in.Name != nil {
		row.Name = strings.TrimSpace(*in.Name)
	}
	if in.Repository != nil {
		row.Repository = strings.TrimSpace(*in.Repository)
	}
	if in.Password != nil {
		row.Password = *in.Password
	}
	if in.Description != nil {
		row.Description = strings.TrimSpace(*in.Description)
	}
	if _, err = dao.UpdateById(id, row); err != nil {
		return nil, err
	}
	return row, nil
}
func (s *StorageRepositoryService) Delete(id int64) error {
	dao := db_.New[models.StorageRepository](db_.NewCtx(s.ctx, s.svcCtx.Db))
	row, err := dao.GetById(id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrRepositoryNotFound
	}
	_, err = dao.Delete(&id)
	return err
}

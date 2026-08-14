package logic

import (
	"context"
	"handfree-work/web-restic/internal/base/conv_"
	"handfree-work/web-restic/internal/base/db_"
	"handfree-work/web-restic/internal/base/error_"
	"handfree-work/web-restic/internal/base/error_/code_"
	"handfree-work/web-restic/internal/base/log_"
	"handfree-work/web-restic/internal/models"
	"handfree-work/web-restic/internal/svc"
)

type UserService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserService(ctx context.Context, svcCtx *svc.ServiceContext) *UserService {
	return &UserService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserService) Add(in *models.User) (*models.User, error) {
	if in.Id == nil || *in.Id == 0 {
		return nil, error_.NewFormatError(code_.ParamIsBlank_, "user.id")
	}
	var entity models.User
	conv_.Convert(in, &entity)

	if in.Password != "" {
		if err := entity.EncryptPassword(in.Password); err != nil {
			log_.Error("EncryptPassword ", err)
			return nil, error_.NewTextError("EncryptPassword error")
		}
	}

	db_.NewCtx(l.ctx, l.svcCtx.Db)
	dao := db_.New[models.User](db_.NewCtx(l.ctx, l.svcCtx.Db))
	err := dao.Create(&entity)
	if err != nil {
		error_.LogError("add user ", err)
		return nil, error_.NewTextError("add user error")
	}
	return in, nil
}

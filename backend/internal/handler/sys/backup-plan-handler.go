package sys

import (
	"handfree-work/octo-backup/internal/base/error_"
	"handfree-work/octo-backup/internal/base/error_/code_"
	"handfree-work/octo-backup/internal/base/web_"
	logic "handfree-work/octo-backup/internal/service"
	"handfree-work/octo-backup/internal/svc"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func RegisterBackupPlan(app *fiber.App, s *svc.ServiceContext) {
	r := app.Group("/api/plan")
	web_.RegisterRoutes(r, s.Auth,
		web_.Route{Path: "/page", Permission: web_.Read, Handler: web_.HandleJSON(backupPlanPage(s))},
		web_.Route{Path: "/info", Permission: web_.Read, Handler: web_.Handle(backupPlanInfo(s))},
		web_.Route{Path: "/create", Permission: web_.Write, Handler: web_.HandleJSON(backupPlanCreate(s))},
		web_.Route{Path: "/update", Permission: web_.Write, Handler: web_.HandleJSON(backupPlanUpdate(s))},
		web_.Route{Path: "/delete", Permission: web_.Admin, Handler: web_.Handle(backupPlanDelete(s))},
	)
}
func backupPlanPage(s *svc.ServiceContext) func(fiber.Ctx, *logic.BackupPlanPageQuery) (any, error) {
	return func(c fiber.Ctx, q *logic.BackupPlanPageQuery) (any, error) {
		return logic.NewBackupPlanService(c.Context(), s).Page(q)
	}
}
func backupPlanInfo(s *svc.ServiceContext) func(fiber.Ctx) (any, error) {
	return func(c fiber.Ctx) (any, error) {
		id, err := strconv.ParseInt(c.Query("id"), 10, 64)
		if err != nil {
			return nil, error_.NewCodeTextError(code_.ParamError, "备份计划 ID 无效")
		}
		return logic.NewBackupPlanService(c.Context(), s).Info(id)
	}
}
func backupPlanCreate(s *svc.ServiceContext) func(fiber.Ctx, *logic.BackupPlanInput) (any, error) {
	return func(c fiber.Ctx, in *logic.BackupPlanInput) (any, error) {
		return logic.NewBackupPlanService(c.Context(), s).Create(in)
	}
}
func backupPlanUpdate(s *svc.ServiceContext) func(fiber.Ctx, *logic.BackupPlanInput) (any, error) {
	return func(c fiber.Ctx, in *logic.BackupPlanInput) (any, error) {
		id, err := strconv.ParseInt(c.Query("id"), 10, 64)
		if err != nil {
			return nil, error_.NewCodeTextError(code_.ParamError, "备份计划 ID 无效")
		}
		return logic.NewBackupPlanService(c.Context(), s).Update(id, in)
	}
}
func backupPlanDelete(s *svc.ServiceContext) func(fiber.Ctx) (any, error) {
	return func(c fiber.Ctx) (any, error) {
		id, err := strconv.ParseInt(c.Query("id"), 10, 64)
		if err != nil {
			return nil, error_.NewCodeTextError(code_.ParamError, "备份计划 ID 无效")
		}
		if err = logic.NewBackupPlanService(c.Context(), s).Delete(id); err != nil {
			return nil, err
		}
		return fiber.Map{}, nil
	}
}

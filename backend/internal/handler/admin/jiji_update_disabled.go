package admin

import (
	"context"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

// JijiUpdateDisabled 替换 SystemHandler 的更新服务：不访问 GitHub，
// 始终报告已是最新版本，并拒绝在线更新与回滚。
// 上报 build_type=source，前端据此不展示一键更新入口。
type JijiUpdateDisabled struct {
	version string
}

func NewJijiUpdateDisabled(updateSvc *service.UpdateService) *JijiUpdateDisabled {
	return &JijiUpdateDisabled{version: updateSvc.CurrentVersion()}
}

func (d *JijiUpdateDisabled) CheckUpdate(context.Context, bool) (*service.UpdateInfo, error) {
	return &service.UpdateInfo{
		CurrentVersion: d.version,
		LatestVersion:  d.version,
		HasUpdate:      false,
		BuildType:      "source",
	}, nil
}

func (d *JijiUpdateDisabled) PerformUpdate(context.Context) error {
	return service.ErrOnlineUpdateDisabled
}

func (d *JijiUpdateDisabled) Rollback() error {
	return service.ErrOnlineUpdateDisabled
}

func (d *JijiUpdateDisabled) ListRollbackVersions(context.Context) ([]service.RollbackVersion, error) {
	return nil, service.ErrOnlineUpdateDisabled
}

func (d *JijiUpdateDisabled) RollbackToVersion(context.Context, string) error {
	return service.ErrOnlineUpdateDisabled
}

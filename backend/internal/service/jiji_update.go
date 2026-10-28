package service

import infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"

// JIJI: 二开版本禁用在线更新/回滚——上游的更新源固定为 Wei-Shaw/sub2api，
// 在线更新会把二开版本替换为上游原版。升级一律通过发布新镜像完成。
var ErrOnlineUpdateDisabled = infraerrors.Forbidden(
	"ONLINE_UPDATE_DISABLED",
	"online update and rollback are disabled in this build; upgrade by deploying a new image",
)

// CurrentVersion returns the running build version.
func (s *UpdateService) CurrentVersion() string {
	return s.currentVersion
}

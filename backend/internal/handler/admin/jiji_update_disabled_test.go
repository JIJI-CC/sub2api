package admin

import (
	"context"
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestJijiUpdateDisabled(t *testing.T) {
	d := NewJijiUpdateDisabled(service.NewUpdateService(nil, nil, "0.2.18-jiji.1", "release"))
	ctx := context.Background()

	info, err := d.CheckUpdate(ctx, true)
	if err != nil {
		t.Fatalf("CheckUpdate: %v", err)
	}
	if info.HasUpdate || info.CurrentVersion != "0.2.18-jiji.1" || info.LatestVersion != "0.2.18-jiji.1" || info.BuildType != "source" {
		t.Fatalf("unexpected info: %+v", info)
	}

	_, listErr := d.ListRollbackVersions(ctx)
	for name, err := range map[string]error{
		"PerformUpdate":        d.PerformUpdate(ctx),
		"Rollback":             d.Rollback(),
		"ListRollbackVersions": listErr,
		"RollbackToVersion":    d.RollbackToVersion(ctx, "0.2.17"),
	} {
		if !errors.Is(err, service.ErrOnlineUpdateDisabled) {
			t.Errorf("%s: got %v, want ErrOnlineUpdateDisabled", name, err)
		}
	}
}

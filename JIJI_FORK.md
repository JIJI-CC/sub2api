# JIJI 二开维护规范

本仓库 fork 自 [Wei-Shaw/sub2api](https://github.com/Wei-Shaw/sub2api)，在上游**正式发布版本**之上做前后端逻辑二开，不新增数据库迁移。

## 分支与 tag

| 名称 | 含义 |
| --- | --- |
| `main` | 上游 `main` 的镜像，只做快进同步，**不提交任何二开代码** |
| `jiji/main` | 二开主线，基于上游正式版 tag，所有二开改动最终合到这里 |
| `jiji/feat-xxx` / `jiji/fix-xxx` | 功能/修复分支，从 `jiji/main` 拉出，完成后合回 |
| `upstream/vX.Y.Z` | 上游 tag（本地命名空间，只读，不推送到 origin） |
| `jiji-vX.Y.Z-N` | 我们的发布 tag，`X.Y.Z` 为所基于的上游版本，`N` 为该版本上的第 N 次发布 |

> ⚠️ 不要向 origin 推送 `v*` tag（包括 `git push --tags`）——会触发上游的 `release.yml`。

## 新克隆后的一次性配置

```bash
git remote add upstream https://github.com/Wei-Shaw/sub2api.git
git config remote.upstream.pushurl DISABLED                       # 禁止误推上游
git config remote.upstream.tagOpt --no-tags                       # 上游 tag 不直接落到 v* ...
git config --add remote.upstream.fetch '+refs/tags/*:refs/tags/upstream/*'  # ... 而是放到 upstream/ 下
git fetch upstream
```

## 日常开发

```bash
git switch jiji/main && git pull
git switch -c jiji/feat-xxx
# ... 开发、提交 ...
git switch jiji/main && git merge --no-ff jiji/feat-xxx && git push
```

提交信息沿用上游风格（`feat: ...` / `fix: ...`），二开提交可加 `jiji` scope，如 `feat(jiji): 新增 xxx`，方便同步上游时识别。

### 降低同步冲突的约定

- **新增优于修改**：新逻辑尽量放在新文件 / 新函数 / 新组件中，原文件里只留最小的接入点（一行调用、一个路由注册、一个组件引用）。
- 后端新文件可用 `jiji_` 前缀（如 `internal/service/jiji_xxx.go`），前端放在 `frontend/src/jiji/` 下，一眼可区分二开代码。
- 不改格式、不顺手重构上游代码；不升级依赖（跟随上游）。
- 不新增数据库迁移；如确需存储，优先复用现有 settings 等键值能力。
- 修改上游文件时，在改动处加 `// JIJI:` 注释说明原因。
- `wire_gen.go`、`ent/` 等生成文件冲突时，解决源文件后重新生成，不手工合并。

## 同步上游新版本

只跟进上游**正式 tag**，不跟上游 `main` 的未发布提交。

```bash
git fetch upstream
git tag -l 'upstream/v*' --sort=-v:refname | head -3      # 查看最新正式版

git switch main && git merge --ff-only upstream/main && git push   # 镜像分支同步（可选）

git switch jiji/main
git switch -c jiji/sync-v0.2.19
git merge upstream/v0.2.19          # 解决冲突
# 验证：后端 go test ./...，前端 pnpm build / pnpm test
git switch jiji/main && git merge --no-ff jiji/sync-v0.2.19 && git push
```

同步时重点检查：

- 上游是否新增了数据库迁移（`backend/migrations/` 等）——上游迁移会在启动时正常执行，升级前备份数据库。
- `git log --oneline upstream/v0.2.18..upstream/v0.2.19` 浏览上游变更，留意与二开逻辑重叠的模块。
- 用 `merge` 而非 `rebase`：`jiji/main` 已推送且可能多人协作，merge 不改写历史。

## 发布

```bash
git switch jiji/main
git tag -a jiji-v0.2.18-1 -m "jiji release 0.2.18-1"
git push origin jiji-v0.2.18-1      # 只推这一个 tag
```

推送后 `.github/workflows/jiji-release.yml` 会构建镜像并推送到：

- `ghcr.io/jiji-cc/sub2api:0.2.18-jiji.1`
- `ghcr.io/jiji-cc/sub2api:latest`

也可在 Actions 页面手动触发（可选 `linux/arm64`、或仅构建不推送）。

## 部署注意事项

- **后台「在线更新 / 回滚」已禁用**：上游更新源固定为 `Wei-Shaw/sub2api`，会把二开版本替换成上游原版。现由 `admin.JijiUpdateDisabled` 替换更新服务（接入点在 `handler/wire.go` 的 `ProvideSystemHandler`），检查更新恒为最新、不访问 GitHub，更新/回滚接口返回 403 `ONLINE_UPDATE_DISABLED`；前端隐藏回滚入口。升级一律通过拉取新镜像完成。同步上游时若 `ProvideSystemHandler` 或 `systemUpdateService` 接口有变，需同步调整。
- `deploy/install.sh`、`deploy/docker-deploy.sh` 默认指向上游仓库与镜像，部署时请改用 `ghcr.io/jiji-cc/sub2api` 镜像。
- GHCR 包默认私有，服务器拉取需 `docker login ghcr.io`，或在 GitHub Packages 设置中将其改为公开。

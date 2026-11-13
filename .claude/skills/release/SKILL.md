---
name: release
description: 发布新版本：根据上一个 tag 以来的提交生成中文更新说明，写入附注 tag 并推送，CI 构建镜像后用 tag 正文创建 GitHub Release。也用于给已发布的 Release 补写更新说明。用户说"发布""发版""打 tag""release""写更新说明"时使用。
---

# 发布新版本

更新说明由 Claude 生成，写在附注 tag 的正文里。推送 tag 后，`.github/workflows/docker-build.yml` 先构建并推送镜像，再取 tag 正文作为 Release 内容，后面接上镜像名、digest 和 Full Changelog 链接。仓库直接提交到 main，没有 PR，GitHub 自动生成的 "What's Changed" 是空的，所以说明必须写在 tag 里。

## 1. 检查状态

```bash
git status --short                 # 应为空
git fetch origin --tags
git rev-parse HEAD origin/main     # 两者应相同，提交要先推送，否则 CI 拿不到
git describe --tags --abbrev=0     # 上一个 tag
```

有未提交的改动或未推送的提交时先告诉用户，不要自行提交或推送。

## 2. 确定版本号

用户没有指定时，按上一个 tag 递增末位（v0.0.3 → v0.0.4），作为建议告诉用户。带 `-` 后缀的（如 `v0.0.4-b1`）会被 CI 标为预发布。

## 3. 收集变更

```bash
PREV=<上一个 tag>
git log "$PREV"..HEAD --format='%h %s%n%n%b%n----'
git diff --stat "$PREV" HEAD | tail -5
git diff "$PREV" HEAD -- ASE-Server-Docker        # 子模块指针变化
git -C ASE-Server-Docker log --oneline <旧指针>..<新指针>
```

提交说明不清楚某个改动对用户有什么影响时，看对应的代码差异，不要猜。

## 4. 写更新说明

读者是部署和使用 ARK Server Commander 的服主，不是开发者。用中文，按下面的结构，没有内容的小节省略：

```markdown
## 新功能

- **功能名**：用户能做什么、在界面哪里。

## 改进

- ...

## 修复

- ...

## 升级注意

- 需要用户手动处理的事项。
```

要求：

- 写用户看得到的行为，不写文件名、函数名、内部实现。API 路径只在确实影响外部调用时提。
- "升级注意"要特别检查：ARK 服务器镜像（子模块 `ASE-Server-Docker`）是否有改动、需要在首页"镜像管理"更新；新增或改名的环境变量；数据目录或配置文件的变化，旧服务器是否需要手动迁移；不兼容的改动。
- 只改文档、CI、测试、重构的提交通常不写，除非影响用户（例如 README 里的部署方式变了）。
- 只依据提交和代码写，不夸大，不写没有做的事。已知限制如实写。
- 不写镜像小节（CI 会加），不写 Co-Authored-By。

## 5. 给用户确认

把版本号和完整的更新说明发给用户，等用户确认或修改后再继续。推送 tag 会触发构建并公开发布，不能跳过这一步。

## 6. 打 tag 并推送

tag 说明的第一行是 `vX.Y.Z: 一句话概括`，空一行后是上面的 Markdown 正文。说明文件放在仓库外面：

```bash
NOTES=$(mktemp)
# 写入 NOTES：第一行概括，空行，正文
git tag -a vX.Y.Z -F "$NOTES" --cleanup=whitespace
git tag -l --format='%(contents:body)' vX.Y.Z    # 检查 ## 标题都还在
git push origin vX.Y.Z
rm "$NOTES"
```

必须用 `--cleanup=whitespace`：默认的 `strip` 会把 `#` 开头的行当注释删掉，Markdown 标题就没了。

## 7. 跟踪发布结果

```bash
gh run list --workflow=docker-build.yml -L 1
gh run watch <run-id> --exit-status
gh release view vX.Y.Z
```

失败时把出错的步骤和日志告诉用户。

## 给已发布的 Release 补写说明

tag 已经推送过、不想重打时，直接编辑 Release，把更新说明放在原有内容（镜像、Full Changelog）前面：

```bash
NOTES=$(mktemp)
# 写入更新说明（不带概括行）
{ cat "$NOTES"; echo; gh release view vX.Y.Z --json body -q .body; } > "$NOTES.full"
gh release edit vX.Y.Z --notes-file "$NOTES.full"
rm "$NOTES" "$NOTES.full"
```

同样先给用户确认内容再执行。

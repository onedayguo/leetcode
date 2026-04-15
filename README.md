# LeetCode Local Workflow (Top100 + Submit)

这个项目提供一个本地 CLI：

1. 拉取 LeetCode Top100 题目到 `problems/`
2. 在本地编辑 `solution.go`
3. 直接提交到 LeetCode 并轮询结果

## Commands

```bash
go run . pull top100 --limit 100
go run . submit two-sum
go run . status 1234567890
```

## Auth (for submit/status)

提交和查结果需要登录态（Cookie）：

- `LEETCODE_SESSION`
- `LEETCODE_CSRF_TOKEN`

可用环境变量：

```bash
export LEETCODE_SESSION="<your_session_cookie>"
export LEETCODE_CSRF_TOKEN="<your_csrf_token>"
export LC_SITE="https://leetcode.com"
```

也可以在项目根目录创建 `.leetcode.json`：

```json
{
  "site": "https://leetcode.com",
  "session": "<your_session_cookie>",
  "csrfToken": "<your_csrf_token>"
}
```

> 建议把 `.leetcode.json` 加入 `.gitignore`，避免泄露登录态。

## Generated Structure

执行 `pull top100` 后会生成：

- `problems/<frontendId>-<title-slug>/README.md`
- `problems/<frontendId>-<title-slug>/meta.json`
- `problems/<frontendId>-<title-slug>/solution.go`

`submit` 默认提交该目录下的 `solution.go`。

## Notes

- Top100 数据源使用 `GET /api/problems/top-100-liked/`（LeetCode 站点接口）
- 题目详情使用 GraphQL `question(titleSlug)`
- 提交使用 `POST /problems/{slug}/submit/`
- 判题结果使用 `GET /submissions/detail/{id}/check/`


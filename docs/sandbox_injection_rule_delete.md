# 简介
`sandbox injection-rule delete`（别名 `dl`）删除一个或多个注入规则。支持变参和交互式多选。

# 格式

```bash
qshell sandbox injection-rule delete [ruleIDs...] [-y] [-s]
qshell sbx ir dl [ruleIDs...] [-y] [-s]
```

# 帮助文档

```bash
$ qshell sandbox injection-rule delete -h
$ qshell sandbox injection-rule delete --doc
```

# 参数

- `ruleIDs...`：一个或多个注入规则 ID
- `-y, --yes`：跳过确认
- `-s, --select`：交互式选择规则进行删除

# 示例

删除单个规则：

```bash
$ qshell sandbox injection-rule delete rule-xxxxxxxxxxxx
```

不加 `-y` 时，命令会交互式询问确认。

跳过确认直接删除：

```bash
$ qshell sandbox injection-rule delete rule-xxxxxxxxxxxx -y
```

删除多个规则：

```bash
$ qshell sandbox injection-rule delete rule-aaa rule-bbb -y
```

交互式选择删除：

```bash
$ qshell sandbox injection-rule delete -s
```

# 非交互式调用（CI / AI Agent / 管道）

当 stdin 不是终端时（CI、`docker run` 不带 `-it`、`exec` 调用、管道等），交互式分支会立即报错并退出，不会卡住进程：

- 必须显式传入规则 ID（不能用 `--select`）
- 必须传 `-y` / `--yes` 跳过确认

```bash
$ qshell sandbox injection-rule delete rule-xxxxxxxxxxxx -y
```

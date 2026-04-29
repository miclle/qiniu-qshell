# 简介
`sandbox injection-rule delete`（别名 `dl`）删除一个或多个注入规则。

# 格式

```bash
qshell sandbox injection-rule delete <ruleIDs...> [-y]
qshell sbx ir dl <ruleIDs...> [-y]
```

# 帮助文档

```bash
$ qshell sandbox injection-rule delete -h
$ qshell sandbox injection-rule delete --doc
```

# 参数

- `ruleIDs...`：一个或多个注入规则 ID（必填）
- `-y, --yes`：跳过确认

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

# 非交互式调用（CI / AI Agent / 管道）

当 stdin 不是终端时，缺省的确认提示会立即报错并退出，不会卡住进程。
自动化场景必须传 `-y` / `--yes` 跳过确认：

```bash
$ qshell sandbox injection-rule delete rule-xxxxxxxxxxxx -y
```

如需批量"列出再删"，请用：

```bash
$ qshell sandbox injection-rule list --format json | jq -r '.[].rule_id' | xargs qshell sandbox injection-rule delete -y
```

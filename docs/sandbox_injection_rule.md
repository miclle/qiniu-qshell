# 简介
`sandbox injection-rule`（别名 `ir`）命令用于管理沙箱注入规则，支持列出、查看、创建、更新和删除注入规则。

# 格式

```bash
qshell sandbox injection-rule <子命令>
qshell sbx ir <子命令>
```

# 帮助文档

```bash
$ qshell sandbox injection-rule -h
$ qshell sandbox injection-rule --doc
```

# 鉴权
注入规则属于用户级资源，使用 AK/SK 鉴权（与七牛对象存储等命令一致）。优先级：

1. `qshell user` 配置的当前账号
2. 环境变量 `QINIU_ACCESS_KEY` / `QINIU_SECRET_KEY`

未配置 AK/SK 时执行子命令会得到明确报错，提示通过 `qshell user` 添加账号或设置环境变量。

# 子命令

`injection-rule` 的子命令有：

- `list`（`ls`）：列出所有注入规则
- `get`（`gt`）：查看指定注入规则详情
- `create`（`cr`）：创建新的注入规则
- `update`（`up`）：更新已有注入规则
- `delete`（`dl`）：删除一个或多个注入规则

# 示例

列出所有注入规则：

```bash
qshell sandbox injection-rule list
```

查看指定注入规则：

```bash
qshell sandbox injection-rule get rule-xxxxxxxxxxxx
```

创建 OpenAI 注入规则：

```bash
qshell sandbox injection-rule create --name openai-default --type openai --api-key sk-xxx
```

更新自定义 HTTP 注入规则：

```bash
qshell sandbox injection-rule update rule-xxxxxxxxxxxx --type http --base-url https://api.example.com --headers "Authorization=Bearer newtoken"
```

删除注入规则：

```bash
qshell sandbox injection-rule delete rule-xxxxxxxxxxxx -y
```

# railblock

区域调度 **闭塞分区占用报文** 的编解码与冲突检测服务：接入车站/轨道电路发来的二进制帧，维护占用位图，检测「双端同时占用」等冲突，输出许可通行/封锁建议。

## 功能

- **二进制帧编解码** — 17 字节定长帧，魔数 `RBLK`、IEEE CRC-32 校验
- **占用位图** — 按分区记录多来源占用状态
- **冲突检测** — R1 双端同时占用；R2 签名 `FORCE_CLEAR` 强制清占用
- **放行判定** — 进路覆盖分区全部空闲且无冲突 → Allow，否则 Reject
- **HTTP API** — 帧接入、分区查询、放行检查
- **嵌入式 Web UI** — 调度员操作台查看位图与冲突

## 快速开始

```bash
cd D:\lzg\railblock
set GOTOOLCHAIN=local
go build -o railblock.exe ./cmd/railblock
railblock.exe
```

默认监听 `:8080`。环境变量：

| 变量 | 说明 | 默认 |
|------|------|------|
| `RAILBLOCK_ADDR` | 监听地址 | `:8080` |
| `RAILBLOCK_FORCE_SECRET` | FORCE_CLEAR 签名密钥 | `railblock-dev-secret` |
| `RAILBLOCK_MAX_BODY` | 最大请求体字节 | `1048576` |
| `RAILBLOCK_CORS` | 启用 CORS | `false` |
| `RAILBLOCK_LOG` | 请求日志 | `false` |

## 帧格式

```
offset 0:  magic   uint32 BE  0x52424C4B ("RBLK")
offset 4:  version uint8      1
offset 5:  flags   uint8
offset 6:  blockID uint16 BE
offset 8:  occupied uint8     0/1
offset 9:  seq     uint32 BE
offset 13: crc32   uint32 BE  IEEE CRC-32 of bytes[0:13]
总长 17 字节
```

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/frames` | 上报二进制帧（`application/octet-stream` 或 base64） |
| POST | `/v1/frames/encode` | 编码演示帧（UI 辅助） |
| GET | `/v1/blocks` | 全部分区快照 |
| GET | `/v1/blocks/{id}` | 单分区查询 |
| POST | `/v1/clearance/check` | 进路放行检查 |
| GET | `/healthz` | 健康检查 |
| GET | `/` | 嵌入式调度台 |

上报帧需携带请求头 `X-Source-ID`（如 `station-a`）。`FORCE_CLEAR` 帧需额外携带 `X-Force-Sign`（HMAC-SHA256 签名）。

### 放行检查示例

```bash
curl -X POST http://localhost:8080/v1/clearance/check \
  -H "Content-Type: application/json" \
  -d '{"route_id":"train-1","blocks":[12,13,14]}'
```

## 项目结构

```
cmd/railblock/          主程序入口
internal/codec/       帧编解码与 CRC
internal/bitmap/        占用位图
internal/conflict/      冲突规则
internal/clearance/     放行判定
internal/ingest/        帧接入
internal/query/         查询 API
internal/app/           服务装配
internal/web/           嵌入式 Web UI
web/                    静态资源源码（嵌入到 internal/web/assets）
```

## 开发与测试

```bash
set GOTOOLCHAIN=local
go build ./...
go test ./... -count=1
```

## 业务场景

车站 A 上报分区 12「占用」，车站 B 随后也报分区 12「占用」→ 冲突模块判定双端占用，`clearance.Check` 对该进路返回 Reject。CRC 错误帧不会污染位图。调度员在操作台看到分区 12 标红。

## 许可

MIT

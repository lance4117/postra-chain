# Postra Chain（postra-chain）

Postra Chain 是 Postra 项目的区块链核心组件，基于 Cosmos SDK 与 Ignite CLI 构建。

Postra Chain 用于在链上记录“可验证的发布元数据”，而不是直接存储文章正文。
每一条发布记录都会锚定一个内容地址（content_uri）与内容哈希（content_hash），从而为链下内容提供公开、可审计、不可篡改的时间与作者证明。

- 项目名称：Postra
- 链仓库：postra-chain
- 区块链程序（daemon）：postrad
- 模块：x/posts
- 技术栈：Cosmos SDK + Ignite CLI

---

## Postra Chain 解决什么问题？

传统博客与内容平台是中心化的：
平台可以修改、删除内容，或者在不通知作者的情况下改变展示方式。

Postra Chain 提供了一种不同的模式：
它不负责“托管内容”，而是提供一个公开可信的记录层，用于证明：

- 某个账户
- 在某个时间点
- 发布了一个指向特定内容的声明

通过将内容哈希与时间戳写入区块链，Postra Chain 确保发布历史无法被事后篡改或抹除。

---

## 数据模型（v0.1）

### Post

```text
Post {
  id: int64             // 自增 ID，由链生成
  creator: string       // 发布者地址
  title: string         // 标题，长度 <= 140
  content_uri: string   // 正文内容地址（IPFS / HTTPS 等）
  content_hash: string  // 正文内容哈希（sha256:<hex> 或 <hex>）
  created_at: int64     // 发布时间（区块时间，Unix 秒，由链端写入）
}
````

---

## 内容哈希规范（建议）

为了保证不同客户端对同一内容计算出一致的哈希值，建议在客户端侧进行如下规范化处理：

* 使用 UTF-8 编码
* 将换行符统一为 `\n`（即将 `\r\n` 转换为 `\n`）
* 不包含 BOM

推荐的哈希格式为：

```
content_hash = "sha256:" + hex(sha256(canonical_bytes(content)))
```

---

## 架构设计说明

* 链上存储：

    * 发布元数据（作者、时间、标题）
    * 内容完整性锚点（content_hash）
    * 内容定位信息（content_uri）

* 链下存储：

    * 实际文章正文（Markdown / 文本 / JSON / 图片等）

* 客户端职责：

    * 根据 content_uri 拉取内容
    * 对内容计算哈希
    * 与链上的 content_hash 进行比对
    * 若不一致，应明确提示内容可能被篡改或来源不可用

这种设计可以在保证内容可验证性的同时，避免将大量正文数据写入区块链，从而保持链的轻量与可扩展性。

---

## 仓库结构（推荐）

```text
postra-chain/
├── app/                  # Cosmos SDK 应用组装（module manager、keepers 等）
├── proto/                # protobuf 定义（权威 API 描述）
├── x/
│   └── posts/            # 发布模块（x/posts）
├── cmd/
│   └── postrad/          # postrad 命令行程序
├── scripts/              # 本地网络、测试辅助脚本
├── docs/                 # 链相关文档（参数说明、升级说明等）
└── README.md
```

说明：

* `proto/` 目录应作为权威接口定义来源
* 其他仓库（前端、SDK 等）应从此处生成客户端代码

---

## 开发环境要求

* Ignite CLI
* Go（版本以所使用的 Cosmos SDK 要求为准）

---

## 快速开始（本地开发）

### 1）启动本地区块链（单节点）

```bash
ignite chain serve
```

该命令会编译链代码、初始化本地数据目录，并以开发模式启动节点。

---

### 2）查看 CLI 命令

```bash
postrad --help
postrad tx posts --help
postrad query posts --help
```

---

### 3）创建一条 Post

具体参数顺序以实际 scaffold 生成的命令为准，
请优先使用 `--help` 查看。

常见形式示例（位置参数）：

```bash
postrad tx posts create-post \
  "Hello Postra" \
  "https://example.com/content/hello.md" \
  "sha256:<hex>" \
  --from <key-name> \
  --chain-id postra-1 \
  --node http://localhost:26657 \
  --yes
```

说明：`created_at` 由链端写入，客户端不需要传入该字段。

---

### 4）更新 Post

仅允许更新 `title`、`content_uri`、`content_hash`。

```bash
postrad tx posts update-post \
  <id> \
  "Updated Title" \
  "https://example.com/content/updated.md" \
  "sha256:<hex>" \
  --from <key-name> \
  --chain-id postra-1 \
  --node http://localhost:26657 \
  --yes
```

---

### 5）查询 Post

Ignite 默认会生成如下查询命令：

```bash
postrad query posts list-post --limit 20 --node http://localhost:26657
postrad query posts show-post <id> --node http://localhost:26657
```

---

## x/posts 模块说明

`x/posts` 是 Postra Chain 的核心模块，定义最基础的发布原语。

当前 MVP 范围：

* Post 对象（content_uri + content_hash）
* 创建 Post
* 更新 Post（仅 title/content_uri/content_hash）
* 删除 Post
* 列表与分页查询
* 使用区块时间作为 created_at
* 基础字段校验（title 长度、URI scheme、content_hash 格式）

未来可扩展方向（非 v0.1 必需）：

* 评论（Comment）
* 反应（Reaction / Like）
* 关注关系（Follow）
* 参数模块（限制标题长度、URI 长度、哈希格式等）

---

## 版本与兼容性

postra-chain 使用语义化版本号（Semantic Versioning）。

前端与产品仓库（postra）应明确声明其所兼容的链版本。
具体对应关系请参见 postra 仓库中的“兼容矩阵”。

---

## 安全模型（简述）

* 区块链是发布元数据与时间的唯一可信来源
* 正文内容的可用性取决于链下存储方式（IPFS、对象存储等）
* 客户端必须验证 content_hash，而不能盲目信任内容来源

---

## 参与贡献

欢迎提交 Issue 与 Pull Request。

建议：

* 尽量保持模块接口的向后兼容
* 状态变更需提供清晰的迁移与升级说明
* 为消息校验与 keeper 逻辑补充测试

---

## 许可证

本项目采用 Apache-2.0 许可证，详见 [LICENSE](LICENSE) 文件。

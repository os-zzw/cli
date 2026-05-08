---
name: lark-slides
version: 1.0.0
description: "飞书幻灯片：创建和编辑幻灯片，接口通过 XML 协议通信。创建演示文稿、读取幻灯片内容、管理幻灯片页面（创建、删除、读取、局部替换）。当用户需要创建或编辑幻灯片、读取或修改单个页面时使用。"
metadata:
  requires:
    bins: ["lark-cli"]
  cliHelp: "lark-cli slides --help"
---

# slides (v1)

**CRITICAL — 开始前 MUST 先用 Read 工具读取 [`../lark-shared/SKILL.md`](../lark-shared/SKILL.md)，其中包含认证、权限处理**

**CRITICAL — 生成任何 XML 之前，MUST 先用 Read 工具读取 [xml-schema-quick-ref.md](references/xml-schema-quick-ref.md)，禁止凭记忆猜测 XML 结构。**

**CRITICAL — 新建演示文稿或大幅改写页面时，MUST 先生成 `.lark-slides/plan/<deck-or-task-id>/slide_plan.json`，再生成 XML。先创建对应目录，规划层规则见 [planning-layer.md](references/planning-layer.md)。仅替换一个标题、插入一个块等小型已有页编辑可豁免。**

**CRITICAL — 新建演示文稿或大幅改写页面时，生成 XML 前 MUST 读取 [visual-planning.md](references/visual-planning.md)，确保 `layout_type`、`visual_focus`、`text_density` 实际改变页面几何、主视觉和文本量。**

**CRITICAL — 新建演示文稿或大幅改写页面时，规划 `asset_need` MUST 遵循 [asset-planning.md](references/asset-planning.md)：只做元数据规划，必须有 `fallback_if_missing`，不得要求真实搜索、下载或上传素材。**

**CRITICAL — 创建或大幅改写后，MUST 按 [validation-checklist.md](references/validation-checklist.md) 做显式验证：回读全文 XML、核对页数和关键元素、检查空白/破损页、明显溢出、布局风险；本地 XML 可用时运行 `layout_lint.py`。**

**CRITICAL — 创建前自检或失败排障时，MUST 按 [troubleshooting.md](references/troubleshooting.md) 检查 XML 转义、结构、shell 截断、图片 token、3350001 和布局风险。**

**CRITICAL — 如果用户提到“模板”“套用模板”“参考某种主题/风格/版式”，或用户需求明显落在已有场景模板内（如工作汇报、产品介绍、商业计划书、培训、晋升汇报等），MUST 先用 [`scripts/template_tool.py`](scripts/template_tool.py) 的 `search` 做模板检索；默认给出 2-3 个最匹配模板候选供用户选择。锁定模板后用 `summarize` 获取主题和布局摘要；只有需要布局骨架时才用 `extract` 裁切目标页型 XML。不要直接读取完整模板 XML。**

> [!NOTE]
> `scripts/template_tool.py` 需要 Python 3。`references/template-index.json` 是脚本缓存/轻量路由索引，不是默认给 agent 阅读的文档；`assets/templates/*.xml` 是机器资源，只应通过脚本摘要或裁切，不要全文读取。

**CRITICAL — 使用模板生成或改写页面时，MUST 先 `summarize` 目标页型；只有需要具体布局骨架时才 `extract`。生成本地 XML 后，如可运行 Python，MUST 先用 [`scripts/layout_lint.py`](scripts/layout_lint.py) 检查 XML well-formed、重叠/越界/文本高度风险，再创建或追加页面。它不是完整 XSD schema 校验。**

**编辑已有幻灯片页面**：优先用 [`+replace-slide`](references/lark-slides-replace-slide.md)（块级替换/插入，不动页序）；选择 action 和完整读-改-写流程见 [`lark-slides-edit-workflows.md`](references/lark-slides-edit-workflows.md)。

## 身份选择

飞书幻灯片通常是用户自己的内容资源。**默认应优先显式使用 `--as user`（用户身份）执行 slides 相关操作**，始终显式指定身份。

- **`--as user`（推荐）**：以当前登录用户身份创建、读取、管理演示文稿。执行前先完成用户授权：

```bash
lark-cli auth login --domain slides
```

- **`--as bot`**：仅在用户明确要求以应用身份操作，或需要让 bot 持有/创建资源时使用。使用 bot 身份时，要额外确认 bot 是否真的有目标演示文稿的访问权限。

**执行规则**：

1. 创建、读取、增删 slide、按用户给出的链接继续编辑已有 PPT，默认都先用 `--as user`。
2. 如果出现权限不足，先检查当前是否误用了 bot 身份；不要默认回退到 bot。
3. 只有在用户明确要求"用应用身份 / bot 身份操作"，或当前工作流就是 bot 创建资源后再做协作授权时，才切换到 `--as bot`。

## 快速开始

一条命令创建包含页面内容的 PPT（推荐）：

```bash
lark-cli slides +create --title "演示文稿标题" --slides '[
  "<slide xmlns=\"http://www.larkoffice.com/sml/2.0\"><style><fill><fillColor color=\"rgb(245,245,245)\"/></fill></style><data><shape type=\"text\" topLeftX=\"80\" topLeftY=\"80\" width=\"800\" height=\"100\"><content textType=\"title\"><p>页面标题</p></content></shape><shape type=\"text\" topLeftX=\"80\" topLeftY=\"200\" width=\"800\" height=\"200\"><content textType=\"body\"><p>正文内容</p><ul><li><p>要点一</p></li><li><p>要点二</p></li></ul></content></shape></data></slide>"
]'
```

也可以分两步（先创建空白 PPT，再逐页添加），详见 [+create 参考文档](references/lark-slides-create.md)。

> [!WARNING]
> `--slides '[...]'` 适合简单页面批量创建，但并不等同于“10 页以内都安全”。如果 slide XML 含中文、大段文本、复杂布局、嵌套引号或较多特殊字符，shell 传参时可能出现转义或截断问题，导致内容丢失、页面空白或布局异常。遇到复杂页面时，优先改用“两步创建法”。

> [!IMPORTANT]
> `slides +create --slides` 底层是“先创建空白 PPT，再逐页调用 `xml_presentation.slide.create`”。这不是原子操作；中途某一页失败时，前面已创建成功的页面会保留。skill 必须把这种“部分成功”风险提前告诉用户，并在失败后先记录 `xml_presentation_id`，回读确认当前状态，再决定是否在现有 PPT 上继续修复或追加。

> 以上是最小可用示例。更丰富的页面效果（渐变背景、卡片、图表、表格等），参考下方 Workflow 和 XML 模板。

## 执行前必做

> **重要**：`references/slides_xml_schema_definition.xml` 是此 skill 唯一正确的 XML 协议来源；其他 md 仅是对它和 CLI schema 的摘要。

### 必读（每次创建前）

| 文档 | 说明 |
|------|------|
| [xml-schema-quick-ref.md](references/xml-schema-quick-ref.md) | **XML 元素和属性速查，必读** |
| [planning-layer.md](references/planning-layer.md) | **新建 PPT / 大幅改写前的持久化规划层，必读** |
| [visual-planning.md](references/visual-planning.md) | **新建 PPT / 大幅改写时的版式智能规则，必读** |
| [asset-planning.md](references/asset-planning.md) | **新建 PPT / 大幅改写时的轻量资产规划规则，必读** |
| [validation-checklist.md](references/validation-checklist.md) | **创建后 / 大幅改写后的显式验证清单，必读** |

### 选读（需要时查阅）

| 场景 | 文档 |
|------|------|
| 需要了解详细 XML 结构 | [xml-format-guide.md](references/xml-format-guide.md) |
| 需要新建 PPT 或大幅改写页面 | [planning-layer.md](references/planning-layer.md) |
| 需要把 plan 转成有差异的页面版式 | [visual-planning.md](references/visual-planning.md) |
| 需要规划图、图标、截图、图表或兜底视觉 | [asset-planning.md](references/asset-planning.md) |
| 需要创建后验收或交付验证记录 | [validation-checklist.md](references/validation-checklist.md) |
| 需要 XML 自检、失败排障、错误码处理 | [troubleshooting.md](references/troubleshooting.md) |
| 需要快速筛模板、做低成本路由 | [`scripts/template_tool.py search`](scripts/template_tool.py) |
| 需要匹配 PPT 模板/主题风格 | [template-catalog.md](references/template-catalog.md) |
| 需要按页型抽摘要或裁切 XML 片段 | [`scripts/template_tool.py`](scripts/template_tool.py) |
| 需要做本地布局风险检查 | [`scripts/layout_lint.py`](scripts/layout_lint.py) |
| 需要 CLI 调用示例 | [examples.md](references/examples.md) |
| 需要参考真实 PPT 的 XML | [slides_demo.xml](references/slides_demo.xml) |
| 需要用 table/chart 等复杂元素 | [slides_xml_schema_definition.xml](references/slides_xml_schema_definition.xml)（完整 Schema） |
| 需要编辑已有 PPT 的单个页面 | [lark-slides-edit-workflows.md](references/lark-slides-edit-workflows.md) |
| 需要了解某个命令的详细参数 | 对应命令的 reference 文档（见下方参考文档章节） |

## Workflow

> **这是演示文稿，不是文档。** 每页 slide 是独立的视觉画面，信息密度要低，排版要留白。

### 创建方式选择

| 场景 | 推荐方式 |
|------|----------|
| 简单 XML（1-3 页、结构简单、几乎无复杂中文和特殊字符） | `slides +create --slides '[...]'` 一步创建 |
| 复杂 XML（多页、含中文、大段文本、复杂布局、嵌套引号、特殊字符较多） | **两步创建**：先 `slides +create` 创建空白 PPT，再用 `xml_presentation.slide create` 逐页添加 |
| 已有 PPT 继续追加或插入页面 | 使用 `xml_presentation.slide create`，必要时配合 `before_slide_id` |

> [!WARNING]
> `--slides '[...]'` 的风险点主要在 shell 参数传递，而不是单纯页数。即使只有 1 页，只要 XML 足够复杂，也建议使用两步创建法。

### 模板与脚本优先流程

模板细则见 [template-catalog.md](references/template-catalog.md)。主流程只记住：先 `search`，锁定后 `summarize`，需要骨架时才 `extract`；不要直接读取完整模板 XML 或照搬占位文案。

```bash
python3 skills/lark-slides/scripts/template_tool.py search --query "<用户需求原文>" --limit 3
python3 skills/lark-slides/scripts/template_tool.py summarize --template <template-id> --label <封面|目录|分节|内容|结尾>
python3 skills/lark-slides/scripts/template_tool.py extract --template <template-id> --label <页型> --out /tmp/template-slice.xml
python3 skills/lark-slides/scripts/layout_lint.py --input /tmp/presentation.xml
```

```text
Step 1: 需求澄清 & 读取知识
  - 澄清主题、受众、页数、风格；模板需求按“模板与脚本优先流程”处理
  - 读取 xml-schema-quick-ref.md；新建 / 大幅改写时还要读取 planning-layer.md、visual-planning.md、asset-planning.md

Step 2: 生成大纲 → 用户确认 → 写入 slide_plan.json
  - 生成结构化大纲供用户确认；如使用模板，标明基于哪个模板改写
  - 新建 / 大幅改写必须先创建目录并写入 `.lark-slides/plan/<deck-or-task-id>/slide_plan.json`
  - plan 字段、路径命名、模板边界和 `asset_need` 结构按 planning-layer.md / asset-planning.md 执行

Step 3: 按 slide_plan.json 生成 XML → 创建
  - 逐页消费 plan：key_message 定主结论，layout_type 定几何，visual_focus 定主视觉，text_density 定文本量
  - 缺少真实素材时必须用 `fallback_if_missing` 生成 XML-native 兜底视觉；不要留空
  - 创建方式按“创建方式选择”判断；图片、复杂 XML、转义和 3350001 排查按 lark-slides-create.md、media-upload.md、troubleshooting.md 执行

Step 4: 审查 & 交付
  - 创建完成后，必须用 xml_presentations.get 读取全文 XML，并按 validation-checklist.md 做显式验证记录
  - 失败或部分成功按 troubleshooting.md 处理；局部问题优先用 `+replace-slide` 修正
  - 没问题 → 交付：告知用户演示文稿 ID 和访问方式
```

### jq 命令模板（编辑已有 PPT 时使用）

新建 PPT 推荐用 `+create --slides`。以下 jq 模板适用于向已有演示文稿追加页面的场景，可以避免手动转义双引号：

```bash
# 追加到末尾
lark-cli slides xml_presentation.slide create \
  --as user \
  --params '{"xml_presentation_id":"YOUR_ID"}' \
  --data "$(jq -n --arg content '<slide xmlns="http://www.larkoffice.com/sml/2.0">
  <style><fill><fillColor color="BACKGROUND_COLOR"/></fill></style>
  <data>
    在这里放置 shape、line、table、chart 等元素
  </data>
</slide>' '{slide:{content:$content}}')"

# 插到指定页之前：before_slide_id 必须在 --data body 里，与 slide 同级
# ⚠️ 不要把 before_slide_id 写进 --params —— CLI 会当未知 query 参数静默下发，服务端忽略，新页跑到末尾
lark-cli slides xml_presentation.slide create \
  --as user \
  --params '{"xml_presentation_id":"YOUR_ID"}' \
  --data "$(jq -n --arg content '<slide ...>...</slide>' --arg before 'TARGET_SLIDE_ID' \
    '{slide:{content:$content}, before_slide_id:$before}')"
```

### 风格快速判断表

> **注意**：渐变色必须使用 `rgba()` 格式并带百分比停靠点，如 `linear-gradient(135deg,rgba(15,23,42,1) 0%,rgba(56,97,140,1) 100%)`。使用 `rgb()` 或省略停靠点会导致服务端回退为白色。

| 场景/主题 | 推荐风格 | 背景 | 主色 | 文字色 |
|----------|---------|------|------|-------|
| 科技/AI/产品 | 深色科技风 | 深蓝渐变 `linear-gradient(135deg,rgba(15,23,42,1) 0%,rgba(56,97,140,1) 100%)` | 蓝色系 `rgb(59,130,246)` | 白色 |
| 商务汇报/季度总结 | 浅色商务风 | 浅灰 `rgb(248,250,252)` | 深蓝 `rgb(30,60,114)` | 深灰 `rgb(30,41,59)` |
| 教育/培训 | 清新明亮风 | 白色 `rgb(255,255,255)` | 绿色系 `rgb(34,197,94)` | 深灰 `rgb(51,65,85)` |
| 创意/设计 | 渐变活力风 | 紫粉渐变 `linear-gradient(135deg,rgba(88,28,135,1) 0%,rgba(190,24,93,1) 100%)` | 粉紫色系 | 白色 |
| 周报/日常汇报 | 简约专业风 | 浅灰 `rgb(248,250,252)` + 顶部彩色渐变条 | 蓝色 `rgb(59,130,246)` | 深色 `rgb(15,23,42)` |
| 用户未指定 | 默认简约专业风 | 同上 | 同上 | 同上 |

### 页面布局建议

| 页面类型 | 布局要点 |
|---------|---------|
| 封面页 | 居中大标题 + 副标题 + 底部信息，背景用渐变或深色 |
| 数据概览页 | 指标卡片横排（rect 背景 + 大号数字 + 小号说明），下方列表或图表 |
| 内容页 | 左侧竖线装饰 + 标题，下方分栏或列表 |
| 对比/表格页 | table 元素或并列卡片，表头深色背景白字 |
| 图表页 | chart 元素（column/line/pie），配合文字说明 |
| 结尾页 | 居中感谢语 + 装饰线，风格与封面呼应 |

### 大纲模板

生成大纲时使用以下格式，交给用户确认：

```text
[PPT 标题] — [定位描述]，面向 [目标受众]

模板：[未使用模板 / <category>/<template>.xml（推荐原因）]

页面结构（N 页）：
1. 封面页：[标题文案]
2. [页面主题]：[要点1]、[要点2]、[要点3]
3. [页面主题]：[要点描述]
...
N. 结尾页：[结尾文案]

风格：[配色方案]，[排版风格]
```

### 常用 Slide XML 模板

可直接复制使用的模板（封面页、内容页、数据卡片页、结尾页）：[slide-templates.md](references/slide-templates.md)

---

## 核心概念

### URL 格式与 Token

| URL 格式 | 示例 | Token 类型 | 处理方式 |
|----------|------|-----------|----------|
| `/slides/` | `https://example.larkoffice.com/slides/xxxxxxxxxxxxx` | `xml_presentation_id` | URL 路径中的 token 直接作为 `xml_presentation_id` 使用 |
| `/wiki/` | `https://example.larkoffice.com/wiki/wikcnxxxxxxxxx` | `wiki_token` | ⚠️ **不能直接使用**，需要先查询获取真实的 `obj_token` |

> `+replace-slide` 和 `+media-upload` shortcut 会自动解析以上两种 URL；直接调用原生 API 时仍需手动解析 wiki 链接。

### Wiki 链接特殊处理（关键！）

知识库链接（`/wiki/TOKEN`）不能直接当 `xml_presentation_id`。直接调用原生 API 前，先查询 wiki 节点，确认 `node.obj_type == "slides"`，再用 `node.obj_token` 作为真实 presentation ID。

```bash
lark-cli wiki spaces get_node --as user --params '{"token":"wiki_token"}'
```

Shortcut `+replace-slide` 和 `+media-upload` 会自动解析 `/wiki/` URL；手动调用 `xml_presentations.*` / `xml_presentation.slide.*` 时才需要自己做这一步。

### 资源关系

```text
Wiki Space (知识空间)
└── Wiki Node (知识库节点, obj_type: slides)
    └── obj_token → xml_presentation_id

Slides (演示文稿)
├── xml_presentation_id (演示文稿唯一标识)
├── revision_id (版本号)
└── Slide (幻灯片页面)
    └── slide_id (页面唯一标识)
```

## Shortcuts（推荐优先使用）

Shortcut 是对常用操作的高级封装（`lark-cli slides +<verb> [flags]`）。有 Shortcut 的操作优先使用。

| Shortcut | 说明 |
|----------|------|
| [`+create`](references/lark-slides-create.md) | 创建 PPT（可选 `--slides` 一步添加页面，支持 `<img src="@./local.png">` 占位符自动上传），bot 模式自动授权 |
| [`+media-upload`](references/lark-slides-media-upload.md) | 上传本地图片到指定演示文稿，返回 `file_token`（用作 `<img src="...">`），最大 20 MB |
| [`+replace-slide`](references/lark-slides-replace-slide.md) | 对已有幻灯片页面进行块级替换/插入（`block_replace` / `block_insert`），自动注入 id 和 `<content/>`，不改变页序 |

## API Resources

```bash
lark-cli schema slides.<resource>.<method>   # 调用 API 前必须先查看参数结构
lark-cli slides <resource> <method> [flags] # 调用 API
```

> **重要**：使用原生 API 时，必须先运行 `schema` 查看 `--data` / `--params` 参数结构，不要猜测字段格式。

### xml_presentations

  - `get` — 读取演示文稿全文信息，XML 格式返回

### xml_presentation.slide

  - `create` — 在指定 XML 演示文稿下创建页面
  - `delete` — 在指定 XML 演示文稿下删除页面
  - `get` — 获取指定 XML 演示文稿的单个页面 XML 内容
  - `replace` — 对指定 XML 演示文稿页面进行元素级别的局部替换

## 核心规则

1. **先规划再写 XML**：新建演示文稿或大幅改写页面时，必须先写入 `.lark-slides/plan/<deck-or-task-id>/slide_plan.json`；模板、风格和大纲只能作为规划输入，不能绕过规划层
2. **创建流程**：简单短 XML（1-3 页、结构简单、特殊字符少）可用 `slides +create --slides '[...]'` 一步创建；复杂内容、含图片/中文大段文本/嵌套引号/较多特殊字符，或超过 10 页时，默认先 `slides +create` 创建空白 PPT，再用 `xml_presentation.slide.create` 逐页添加
3. **`<slide>` 直接子元素只有 `<style>`、`<data>`、`<note>`**：文本和图形必须放在 `<data>` 内
4. **文本通过 `<content>` 表达**：必须用 `<content><p>...</p></content>`，不能把文字直接写在 shape 内
5. **保存关键 ID**：后续操作需要 `xml_presentation_id`、`slide_id`、`revision_id`
6. **删除谨慎**：删除操作不可逆，且至少保留一页幻灯片
7. **编辑已有页面优先块级替换**：修改单个 shape/img 用 `+replace-slide`（`block_replace` / `block_insert`），不要整页重建；只有需要替换整页结构时才用 `slide.delete` + `slide.create`
8. **`<img src>` 只能用上传到飞书 drive 的 `file_token`，禁止使用 http(s) 外链 URL**：飞书 slides 渲染端不会代理外链图片，外链 src 在 PPT 里通常不显示或显示破图。流程必须是「先把图存到本地 → 用 `slides +media-upload` 上传或 `+create --slides` 的 `@./path` 占位符自动上传 → 拿 `file_token` 写进 `<img src>`」。如果用户给了网图链接，先 `curl`/下载到 CWD 内再走上传流程，不要直接把外链 URL 塞进 `src`。**图片最大 20 MB**（slides upload API 不支持分片上传）。

## 权限表

| 方法 | 所需 scope |
|------|-----------|
| `slides +create` | `slides:presentation:create`, `slides:presentation:write_only`（含 `@` 占位符时还需 `docs:document.media:upload`） |
| `slides +media-upload` | `docs:document.media:upload`（wiki URL 解析还需 `wiki:node:read`） |
| `slides +replace-slide` | `slides:presentation:update`（wiki URL 解析还需 `wiki:node:read`） |
| `xml_presentations.get` | `slides:presentation:read` |
| `xml_presentation.slide.create` | `slides:presentation:update` 或 `slides:presentation:write_only` |
| `xml_presentation.slide.delete` | `slides:presentation:update` 或 `slides:presentation:write_only` |
| `xml_presentation.slide.get` | `slides:presentation:read` |
| `xml_presentation.slide.replace` | `slides:presentation:update` |

## 参考文档

| 文档 | 说明 |
|------|------|
| [lark-slides-create.md](references/lark-slides-create.md) | **+create Shortcut：创建 PPT（支持 `--slides` 一步添加页面，含 `@` 占位符自动上传图片）** |
| [lark-slides-media-upload.md](references/lark-slides-media-upload.md) | **+media-upload Shortcut：上传本地图片，返回 `file_token`** |
| [lark-slides-replace-slide.md](references/lark-slides-replace-slide.md) | **+replace-slide Shortcut：块级替换/插入，含合法根元素速查与 3350001 排错** |
| [lark-slides-edit-workflows.md](references/lark-slides-edit-workflows.md) | 编辑已有页面的读-改-写流程与 action 决策树 |
| [template-index.json](references/template-index.json) | **脚本缓存/轻量路由索引：由 `template_tool.py search` 使用，不是默认阅读入口** |
| [template-catalog.md](references/template-catalog.md) | **按场景/色调匹配现成 PPT 模板，并定位到页型范围** |
| [`scripts/template_tool.py`](scripts/template_tool.py) | **可选 Python 辅助脚本：`search` / `summarize` / `extract`，支持 `--layout-tag` 与 `extract --with-summary`** |
| [`scripts/layout_lint.py`](scripts/layout_lint.py) | **本地预检脚本：先检查 XML well-formed，再检测重叠、越界、页脚碰撞、文本高度风险；不是完整 XSD schema 校验** |
| [planning-layer.md](references/planning-layer.md) | 新建 PPT / 大幅改写前的持久化规划层 |
| [validation-checklist.md](references/validation-checklist.md) | 创建后 / 大幅改写后的显式验证清单 |
| [troubleshooting.md](references/troubleshooting.md) | XML 自检、失败排障、错误码和症状修复 |
| [xml-schema-quick-ref.md](references/xml-schema-quick-ref.md) | **XML Schema 精简速查（必读）** |
| [slide-templates.md](references/slide-templates.md) | 可复制的 Slide XML 模板 |
| [xml-format-guide.md](references/xml-format-guide.md) | XML 详细结构与示例 |
| [examples.md](references/examples.md) | CLI 调用示例 |
| [slides_demo.xml](references/slides_demo.xml) | 真实 PPT 的完整 XML |
| [slides_xml_schema_definition.xml](references/slides_xml_schema_definition.xml) | **完整 Schema 定义**（唯一协议依据） |
| [lark-slides-xml-presentations-get.md](references/lark-slides-xml-presentations-get.md) | 读取 PPT 命令详情 |
| [lark-slides-xml-presentation-slide-create.md](references/lark-slides-xml-presentation-slide-create.md) | 添加幻灯片命令详情 |
| [lark-slides-xml-presentation-slide-delete.md](references/lark-slides-xml-presentation-slide-delete.md) | 删除幻灯片命令详情 |
| [lark-slides-xml-presentation-slide-get.md](references/lark-slides-xml-presentation-slide-get.md) | 读取单个幻灯片命令详情 |
| [lark-slides-xml-presentation-slide-replace.md](references/lark-slides-xml-presentation-slide-replace.md) | 原生 slide.replace API 命令详情 |

> **注意**：如果 md 内容与 `slides_xml_schema_definition.xml` 或 `lark-cli schema slides.<resource>.<method>` 输出不一致，以后两者为准。

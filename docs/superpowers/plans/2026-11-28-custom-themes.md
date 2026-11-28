# 自定义主题实施计划

> **For agentic workers:** 使用 subagent-driven-development 执行独立模块并进行审查。

**Goal:** 提供可视化终端主题编辑、预览、持久化和可联动的界面配色。

**Architecture:** 保留 colorScheme 和 palette 渲染结构；用户主题独立存储并使用稳定 ID。编辑草稿与生效配置分离，界面主题解析和终端颜色解析避免循环。

**Tech Stack:** Go、Viva ui/components/gfx、现有配置存储。

## Global Constraints

- 在 feat/custom-themes 实施，不修改已有暂存 deb。
- 每次创建或编辑不超过 100 行，新文件 git add。
- 使用中文交流及注释，正常业务场景测试，保留现有配置兼容。
- 独立模块负责不同文件，协调共享 State、persist 和 scheme 接口。

### Task 1: 界面主题和联动

文件：新建 application_theme*.go；修改 theme.go、workbench.go 外观解析、persist.go 配置字段、settings.go 模式。
接口：State.applicationTheme 保存用户配色；refreshApplicationTheme 更新界面；resolveAppearance 保持单向来源。
- [x] 定义持久化界面配色、独立 followScheme 来源。
- [x] 添加自定义及跟随终端模式，派生完整 palette 和组件 Theme。
- [x] 暴露界面设置及颜色选择器、背景不透明度。
- [x] 正常保存恢复、模式切换、亮暗来源测试。
- [x] 审查配置迁移及已有深浅色行为。

### Task 2: 终端主题管理与编辑

文件：新建 custom_scheme*.go、theme_editor*.go；修改 scheme.go、settings_paint.go、State。
接口：saveCustomScheme(colorScheme)、deleteCustomScheme(id)、themeEditorState。
- [x] 稳定 ID 用户主题模型和原子保存，保留 xresources 加载。
- [x] 创建/复制/编辑/重命名/删除，内置只读。
- [x] 独立草稿及色盘、HEX/RGB、选区和 ANSI 色。
- [x] 编辑页面及保存、取消、重置操作。
- [x] 正常复制隔离、保存恢复、编辑应用测试。

### Task 3: 预览及集成验证

文件：新建 theme_preview.go、custom_theme_test.go；终端选区和背景渲染集成。
- [x] 固定示例终端渲染，实时预览，宽屏并排/窄屏页签。
- [x] 展开预览及界面配色预览。
- [x] 接入现有终端刷新、选区及背景混合。
- [x] 运行 go test ./internal/workbench 和相关包测试。
- [x] 运行原生 headless 控件交互测试、构建及整体代码审查。

验证命令：GOROOT=/soft/dev/go GOPATH=/home/tudou/go PATH=/soft/dev/go/bin:$PATH go test ./internal/workbench
预期：PASS；如仓库其他测试已有失败，独立复现并记录。

## 验证记录

- 主题/外观/背景/终端透明度/主题目录相关专项通过。
- 新增真实指针与键盘测试：主题名称、色块、颜色选择器应用、展开预览、窄屏预览、保存。
- 保存恢复、复制隔离、重命名稳定 ID、多窗口刷新与跟随模式切换通过。
- 最小窗口布局与配置字段归属检查通过。
- `go build -o /tmp/motty-custom-themes ./cmd/motty` 通过。
- workbench 测试按 A–H、I–P、Q–Z 分组执行；本次引入的四处失败均已修复并复测通过。
- 另外 45 个顶层测试失败均在原始 `677b8312` 独立检出中复现，涉及既有远程工具、设置搜索、文件管理等。本次未修改这些无关行为。
- 审查修复：保存冲突保留草稿并报错，避免通用冲突恢复误报成功；所有预览共用实时色盘草稿。

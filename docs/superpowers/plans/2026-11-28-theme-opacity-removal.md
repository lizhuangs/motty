# 统一终端背景可见度实施计划

目标：消除终端主题不透明度与图片可见度的叠加。
方案：删除主题 opacity 数据链，保留全局背景设置和独立界面设置。
技术：Go、Viva、现有配置存储。

- [x] 添加旧主题配置加载与编辑后背景可见度保持一致的成功场景测试。
- [x] 移除 scheme、custom_scheme、theme_editor 中的透明度字段和控件。
- [x] 更新 background、theme_preview_terminal、application_theme 的透明度来源。
- [x] 调整既有测试与使用文档，执行专项测试（含编译检查）。

验证：旧配置回归测试先复现 alpha=0.25，再修复为 0.5；主题、背景与透明度专项测试通过。
广泛匹配 Background 时包含已记录的 SFTP 菜单失败，未修改该模块。

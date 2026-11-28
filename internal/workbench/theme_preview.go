package workbench

import (
	"fmt"
	"github.com/vivaui/viva/gfx"
	"github.com/vivaui/viva/ui"
)

const themePreviewSample = "\x1b[32muser@motty\x1b[0m:~/project $ ls\r\n" +
	"\x1b[34msrc/\x1b[0m  \x1b[36mdocs/\x1b[0m  README.md\r\n\r\n" +
	"\x1b[32m+ added: custom themes\x1b[0m\r\n\x1b[31m- removed: old colors\x1b[0m\r\n" +
	"\x1b[33mWARN\x1b[0m  Configuration changed\r\n\x1b[36mINFO\x1b[0m  Theme saved successfully\r\n" +
	"Selected text / 选中文字\r\n\r\nuser@motty:~ $ "

func (state *State) declareThemePreview(c ui.C) {
	e := &state.themeEditor
	preview := e.previewScheme()
	ui.Column(c).FillWidth().Gap(12).Children(func() {
		ui.Row(c).Items(ui.Center).Children(func() {
			ui.Text(c, state.themeText("Live preview", "实时预览", "即時預覽")).Grow(1).Color(state.pal.text).FontSize(15)
			label := state.themeText("Expand", "展开", "展開")
			if e.expanded {
				label = state.themeText("Back to editor", "返回编辑", "返回編輯")
			}
			ui.Button(c, label).Outline().Hit(hitThemePreview, label)
		})
		ui.Slot(c).FillWidth().Height(280).Clip().Draw(func(d *ui.DrawContext) { state.paintThemePreviewTerminal(d) })
		ui.Row(c).Gap(3).Children(func() {
			for i := 0; i < 16; i++ {
				color := schemeRGB(preview.ansi[i])
				ui.Column(c).Grow(1).Basis(0).Height(28).Background(color).Center().Children(func() {
					ink := gfx.RGB(255, 255, 255)
					if luminance(rgb8(color)) > .5 {
						ink = gfx.RGB(0, 0, 0)
					}
					ui.Text(c, fmt.Sprint(i)).FontSize(10).Color(ink)
				})
			}
		})
		p := applicationPaletteForScheme(preview)
		ui.Text(c, state.themeText("Interface following this theme", "界面跟随此主题的效果", "介面跟隨此主題的效果")).Color(state.pal.muted).FontSize(12)
		ui.Column(c).FillWidth().Background(p.workspace).Radius(8).Border(1, p.line).Children(func() {
			ui.Row(c).FillWidth().Height(36).PadX(12).Items(ui.Center).Background(p.shell).Children(func() {
				ui.Text(c, "Motty · Terminal").FontSize(13).Color(p.text)
			})
			ui.Column(c).Pad(16).Gap(10).Children(func() {
				ui.Text(c, state.themeText("Connections & settings", "连接与设置", "連線與設定")).FontSize(14).Color(p.text)
				ui.Text(c, state.themeText("Secondary text / description", "次要文字 / 说明", "次要文字 / 說明")).FontSize(12).Color(p.muted)
				ui.Row(c).Height(34).PadX(12).Items(ui.Center).Background(p.surface).Radius(6).Children(func() {
					ui.Text(c, "example.server").FontSize(12).Color(p.text)
				})
				ui.Row(c).Height(32).PadX(12).Items(ui.Center).Background(p.accent).Radius(6).Children(func() {
					ui.Text(c, state.themeText("Connect", "连接", "連線")).FontSize(12).Color(p.onAccent)
				})
			})
		})
	})
}

package workbench

import (
	"maps"
	"strconv"

	"github.com/vivaui/viva/gfx"
	"github.com/vivaui/viva/ui"
)

func (state *State) paintApplicationThemePreview(page *settingsPage) {
	c, pal, label := page.c, state.applicationPreviewPalette(), state.themeText
	page.separator()
	ui.Column(c).MarginSides(12, 18, 12, 18).Radius(radius.control).Border(1, pal.line).Background(pal.workspace).Children(func() {
		ui.Row(c).Pad(8).Gap(6).Background(pal.shell).Children(func() {
			ui.Text(c, label("Preview", "预览", "預覽")).Color(pal.text).FontSize(metrics.FontSizeS).
				PadSides(3, 10, 3, 10).Radius(radius.item).Background(pal.tabActive)
			ui.Text(c, label("Inactive tab", "非活动标签页", "非使用中分頁")).Color(pal.muted).FontSize(metrics.FontSizeS).
				PadSides(3, 10, 3, 10).Radius(radius.item).Background(pal.tab)
			ui.Text(c, label("Hovered tab", "悬停标签页", "懸停分頁")).Color(pal.text).FontSize(metrics.FontSizeS).
				PadSides(3, 10, 3, 10).Radius(radius.item).Background(pal.tabHover)
		})
		ui.Row(c).Margin(10).Pad(10).Gap(8).Items(ui.Center).Radius(radius.control).Border(1, pal.line).Background(pal.surface).Children(func() {
			ui.Column(c).Grow(1).Shrink(1).Children(func() {
				ui.Text(c, label("Sample panel", "面板示例", "面板範例")).Color(pal.text).Bold().Ellipsis()
				ui.Text(c, label("Secondary text", "次要文字", "次要文字")).Color(pal.muted).FontSize(metrics.FontSizeS).Ellipsis()
			})
			ui.Text(c, label("Accent", "强调色", "強調色")).Color(pal.onAccent).FontSize(metrics.FontSizeS).
				PadSides(4, 10, 4, 10).Radius(radius.control).Background(pal.accent)
		})
		ui.Row(c).MarginSides(0, 10, 10, 10).Gap(6).Children(func() {
			ui.Text(c, label("Hover", "悬停", "懸停")).Color(pal.text).FontSize(metrics.FontSizeS).
				PadSides(3, 10, 3, 10).Radius(radius.control).Background(pal.surfaceSoft)
			ui.Text(c, label("Pressed", "按下", "按下")).Color(pal.text).FontSize(metrics.FontSizeS).
				PadSides(3, 10, 3, 10).Radius(radius.control).Background(pal.surfaceHot)
		})
	})
}

func (state *State) applicationPreviewPalette() *palette {
	if state.appearanceMode == appearanceCustom {
		for index, token := range applicationColors {
			if color, open := state.applicationThemeUI.pickers[index].Draft(); open {
				draft := state.applicationTheme
				draft.Colors = maps.Clone(draft.Colors)
				draft.setColor(token.key, color)
				return draft.palette()
			}
		}
	}
	return state.pal
}

func (state *State) paintApplicationOpacity(page *settingsPage) {
	c, slider := page.c, &state.applicationThemeUI.opacity
	label := state.themeText("Interface background opacity", "界面背景不透明度", "介面背景不透明度")
	description := state.themeText("Applies over a background image or on supported transparent windows.", "在背景图片上或支持透明的窗口中生效。", "在背景圖片上或支援透明的視窗中生效。")
	page.lineFor(hitApplicationOpacity, label, description, 182, func() {
		ui.Row(c).Grow(1).Gap(8).Items(ui.Center).Children(func() {
			wasDragging := slider.Dragging()
			control := ui.Slider(c, slider).Range(0, 100).Step(1).Label(label).Hit(hitApplicationOpacity, label).Grow(1).Shrink(1)
			slider.ID = hitApplicationOpacity
			changed := control.Changed()
			opacity := int(slider.Value)
			if opacity != state.applicationOpacity() {
				state.applicationTheme.Opacity = &opacity
				if opacity == 100 {
					state.applicationTheme.Opacity = nil
				}
			}
			if !slider.Dragging() && (wasDragging || changed) {
				state.save()
			}
			ui.Text(c, strconv.Itoa(opacity)+"%").FontSize(metrics.FontSizeS).Color(state.pal.muted).Width(36).Shrink(0).TextEnd()
		})
	})
}

func (state *State) applicationOpacity() int {
	if opacity := state.applicationTheme.Opacity; opacity != nil {
		return min(max(*opacity, 0), 100)
	}
	return 100
}

func (state *State) applicationGround(color gfx.Color) gfx.Color {
	if state.terminalTransparent() || state.showsBackground() {
		color.Alpha *= float32(state.applicationOpacity()) / 100
	}
	return color
}

package workbench

import (
	"strconv"
	"strings"

	"github.com/vivaui/viva/components"
	"github.com/vivaui/viva/gfx"
	"github.com/vivaui/viva/ui"
)

func (state *State) declareThemeColorFields(c ui.C) {
	e := &state.themeEditor
	ui.Row(c).Gap(6).Children(func() {
		hex := state.themeColorField(c, &e.hex, "HEX", "#RRGGBB", 2)
		if hex.Changed() {
			color, ok := components.ParseColorHex(strings.TrimSpace(e.hex.Text()), false)
			if ok && len(strings.TrimSpace(e.hex.Text())) == 7 {
				e.colors[e.colorIndex] = color
				channels := rgb8(color)
				e.red.SetText(strconv.Itoa(int(channels[0])))
				e.green.SetText(strconv.Itoa(int(channels[1])))
				e.blue.SetText(strconv.Itoa(int(channels[2])))
				e.error = ""
				e.colorInvalid = false
			} else {
				e.colorInvalid = true
				e.error = state.themeText("Enter #RRGGBB", "请输入 #RRGGBB", "請輸入 #RRGGBB")
			}
		}
		changed := false
		for i, field := range []*ui.Field{&e.red, &e.green, &e.blue} {
			box := state.themeColorField(c, field, []string{"R", "G", "B"}[i], "0–255", 1)
			changed = box.Changed() || changed
		}
		if changed {
			var channels [3]uint8
			valid := true
			for i, value := range []string{e.red.Text(), e.green.Text(), e.blue.Text()} {
				n, err := strconv.Atoi(value)
				if err != nil || n < 0 || n > 255 {
					valid = false
				} else {
					channels[i] = uint8(n)
				}
			}
			if valid {
				e.colors[e.colorIndex] = schemeRGB(channels)
				e.syncColorFields()
			} else {
				e.colorInvalid = true
				e.error = state.themeText("RGB values must be 0–255", "RGB 值须为 0–255", "RGB 值須為 0–255")
			}
		}
	})
}

func (state *State) paintThemeManagement(page *settingsPage) {
	page.group(state.themeText("Custom themes", "自定义主题", "自訂主題"), func() {
		ui.Column(page.c).Pad(16).Gap(10).Children(func() {
			ui.Grid(page.c, max((page.width-42)/2, 100)).FillWidth().Gap(8).Children(func() {
				for _, action := range []struct {
					id    gfx.ID
					label string
				}{
					{hitThemeNew, state.themeText("New", "新建", "新增")},
					{hitThemeCopy, state.themeText("Duplicate", "复制并编辑", "複製並編輯")},
					{hitThemeEdit, state.themeText("Edit / rename", "编辑 / 重命名", "編輯 / 重新命名")},
				} {
					page.button(action.id, action.label, aiButtonOutline, false)
				}
				page.button(hitThemeDelete, state.themeText("Delete", "删除", "刪除"), aiButtonOutline,
					!strings.HasPrefix(state.activeScheme(), customSchemePrefix))
			})
			ui.Paragraph(page.c, state.themeText("Built-in and imported themes are edited as a copy. Choose a theme above to manage it.",
				"内置和导入主题会复制后编辑。先在上方选择要管理的主题。", "內建與匯入主題會複製後編輯。先在上方選擇要管理的主題。"), 3).
				FontSize(12).Color(state.pal.muted).LineHeight(20)
		})
	})
}

func (state *State) themeColorField(c ui.C, field *ui.Field, label, placeholder string, grow float32) ui.Box {
	var box ui.Box
	ui.Column(c).Grow(grow).Basis(0).Gap(4).Children(func() {
		ui.Text(c, label).Color(state.pal.muted).FontSize(11)
		box = ui.TextField(c, field, label, placeholder).FillWidth()
	})
	return box
}

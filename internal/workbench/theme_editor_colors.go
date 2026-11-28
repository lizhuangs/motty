package workbench

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vivaui/viva/components"
	"github.com/vivaui/viva/gfx"
	"github.com/vivaui/viva/ui"
)

func (e *themeEditorState) loadDraft() {
	s := &e.draft
	e.name.SetText(s.name)
	e.colors[0], e.colors[1], e.colors[2] = schemeRGB(s.background), schemeRGB(s.foreground), schemeRGB(s.cursor)
	if !s.customSelection {
		s.selection = rgb8(e.colors[0].Mix(e.colors[1], .25))
		s.selectionText = s.foreground
	}
	e.colors[3], e.colors[4] = schemeRGB(s.selection), schemeRGB(s.selectionText)
	for i, color := range s.ansi {
		e.colors[i+5] = schemeRGB(color)
	}
	e.opacity.Value = 100
	if s.opacity != nil {
		e.opacity.Value = float32(*s.opacity)
	}
	e.syncColorFields()
}

func schemeRGB(c [3]uint8) gfx.Color { return gfx.RGB(c[0], c[1], c[2]) }

func (e *themeEditorState) syncColorFields() {
	c := rgb8(e.colors[e.colorIndex])
	e.hex.SetText(fmt.Sprintf("#%02x%02x%02x", c[0], c[1], c[2]))
	e.red.SetText(strconv.Itoa(int(c[0])))
	e.green.SetText(strconv.Itoa(int(c[1])))
	e.blue.SetText(strconv.Itoa(int(c[2])))
	e.error = ""
	e.colorInvalid = false
}

func (e *themeEditorState) collectDraft() colorScheme {
	s := e.draft
	s.name = strings.TrimSpace(e.name.Text())
	s.background, s.foreground, s.cursor = rgb8(e.colors[0]), rgb8(e.colors[1]), rgb8(e.colors[2])
	s.selection, s.selectionText, s.customSelection = rgb8(e.colors[3]), rgb8(e.colors[4]), true
	for i := range s.ansi {
		s.ansi[i] = rgb8(e.colors[i+5])
	}
	opacity := int(e.opacity.Value)
	s.opacity = &opacity
	return s
}

func (state *State) themeColorLabel(index int) string {
	labels := []string{
		state.themeText("Background", "背景", "背景"), state.themeText("Text", "文字", "文字"),
		state.themeText("Cursor", "光标", "游標"), state.themeText("Selection", "选区背景", "選取背景"),
		state.themeText("Selection text", "选区文字", "選取文字"),
	}
	if index < len(labels) {
		return labels[index]
	}
	return fmt.Sprintf("ANSI %d", index-5)
}

func (state *State) declareThemeColors(c ui.C) {
	e := &state.themeEditor
	ui.Grid(c, 64).FillWidth().Gap(8).Children(func() {
		for i := range e.colors {
			label := state.themeColorLabel(i)
			ui.Column(c).Gap(3).Items(ui.Center).Children(func() {
				swatch := ui.ColorSwatch(c, e.colors[i], e.colorIndex == i, label).
					Hit(hitThemeColors.At(i), label).Size(32, 32).DrawOver(func(d *ui.DrawContext) { e.colorAnchors[i] = d.Rect })
				ui.Text(c, label).FontSize(10).Color(state.pal.muted).LineHeight(16)
				if swatch.Clicked() {
					e.colorIndex = i
					e.syncColorFields()
					e.openPicker = true
				}
			})
		}
	})

	state.declareThemeColorFields(c)
}

// 弹层脱离滚动裁剪，但锚点始终来自被点击的色块。
func (state *State) declareThemeColorPicker(c ui.C) {
	e := &state.themeEditor
	picker := &e.pickers[e.colorIndex]
	if e.openPicker {
		picker.SetDraftValue(components.ColorPickerValue{}.WithColor(e.colors[e.colorIndex]))
		picker.SetOpenState(true)
		e.openPicker = false
	}

	if picker.Applied() {
		if e.colors[e.colorIndex].Alpha == 0 {
			original := themeEditorState{draft: e.original}
			original.loadDraft()
			e.colors[e.colorIndex] = original.colors[e.colorIndex]
		}
		e.syncColorFields()
	}
	if !picker.Open() {
		return
	}
	ui.Overlay(c, ui.OverlaySpec{Band: ui.Band(dialogLayer + 2), Bounds: e.colorAnchors[e.colorIndex]}).Children(func() {
		ui.ColorPicker(c, picker, &e.colors[e.colorIndex], state.themeColorLabel(e.colorIndex)).PanelLabels(components.ColorPickerLabels{
			Color: state.messages.ColorPickerChoose, Saturation: state.messages.ColorPickerSV,
			Hue: state.messages.ColorPickerHue, Hex: "HEX", Clear: state.messages.FindSettings.Reset,
			Cancel: state.messages.Cancel, Apply: state.messages.ColorPickerApply,
		})
	})
}

func (e *themeEditorState) useTemplate(look appearance) {
	name := e.name.Text()
	p := look.palette()
	e.draft = colorScheme{name: name, background: rgb8(p.terminal), foreground: rgb8(p.termFG), cursor: rgb8(p.cursor), ansi: p.ansi}
	e.loadDraft()
}

// 所有预览读取同一份临时色盘值；取消色盘后自然恢复已确认的草稿。
func (e *themeEditorState) previewScheme() colorScheme {
	s := e.collectDraft()
	if color, ok := e.pickers[e.colorIndex].Draft(); ok {
		rgb := rgb8(color)
		switch e.colorIndex {
		case 0:
			s.background = rgb
		case 1:
			s.foreground = rgb
		case 2:
			s.cursor = rgb
		case 3:
			s.selection = rgb
		case 4:
			s.selectionText = rgb
		default:
			s.ansi[e.colorIndex-5] = rgb
		}
	}
	return s
}

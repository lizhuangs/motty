package workbench

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vivaui/viva/app"
	"github.com/vivaui/viva/components"
	"github.com/vivaui/viva/gfx"
	"github.com/vivaui/viva/headless"
	"motty/internal/i18n"
)

func TestCustomThemeEditorPointerTypingAndPreview(t *testing.T) {
	state := newTestState(t, testOptions(t))
	state.setLanguage(i18n.Chinese)
	state.setTerminalScheme("dracula")
	state.openThemeEditor(true, false)
	view, err := headless.Open(state.WindowConfig(app.Config{Title: "Theme editor", Width: 1280, Height: 900, Tree: Tree(),
		Paint: state.Paint, OnPointer: state.Pointer, OnKey: state.Key, OnScroll: state.Scroll}), nil, headless.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer view.Close()
	if err := view.Wait(30 * time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if err := view.ClickLabel("最大化"); err != nil {
		t.Fatal(err)
	}
	if !state.themeEditor.maximized {
		t.Fatal("editor did not maximize")
	}
	if err := view.ClickLabel("还原"); err != nil {
		t.Fatal(err)
	}
	if state.themeEditor.maximized {
		t.Fatal("editor did not restore")
	}
	if err := view.ClickLabel("主题名称"); err != nil {
		t.Fatal(err)
	}
	if err := view.Type(" edited"); err != nil {
		t.Fatal(err)
	}
	if state.themeEditor.name.Text() != "Dracula 副本 edited" {
		t.Fatalf("name=%q", state.themeEditor.name.Text())
	}
	if err := view.ClickLabel("ANSI 4"); err != nil {
		t.Fatal(err)
	}
	if state.themeEditor.colorIndex != 9 {
		t.Fatal("swatch did not select ANSI 4")
	}
	if err := view.Rest(); err != nil {
		t.Fatal(err)
	}
	if !state.themeEditor.pickers[9].Open() {
		t.Fatal("clicking a color did not open the picker")
	}
	chosen := gfx.RGB(96, 160, 240)
	state.themeEditor.pickers[9].SetDraftValue(components.ColorPickerValue{}.WithColor(chosen))
	if err := view.Rest(); err != nil {
		t.Fatal(err)
	}
	if state.themeEditor.previewScheme().ansi[4] != rgb8(chosen) {
		t.Fatal("picker draft did not reach preview")
	}
	if err := view.ClickLabel(state.messages.ColorPickerApply); err != nil {
		t.Fatal(err)
	}
	if err := view.Rest(); err != nil {
		t.Fatal(err)
	}
	if rgb8(state.themeEditor.colors[9]) != rgb8(chosen) || state.themeEditor.hex.Text() != "#60a0f0" {
		t.Fatalf("color picker: color=%v HEX=%q", rgb8(state.themeEditor.colors[9]), state.themeEditor.hex.Text())
	}
	shots := os.Getenv("MOTTY_THEME_SCREENSHOTS")
	if shots != "" {
		if err := view.Wait(300 * time.Millisecond); err != nil {
			t.Fatal(err)
		}
		if err := view.Screenshot(filepath.Join(shots, "theme-editor.png")); err != nil {
			t.Fatal(err)
		}
	}
	if err := view.ClickLabel("展开"); err != nil {
		t.Fatal(err)
	}
	if !state.themeEditor.expanded {
		t.Fatal("preview did not expand")
	}
	if shots != "" {
		if err := view.Wait(300 * time.Millisecond); err != nil {
			t.Fatal(err)
		}
		if err := view.Screenshot(filepath.Join(shots, "theme-preview.png")); err != nil {
			t.Fatal(err)
		}
	}
	if err := view.ClickLabel("返回编辑"); err != nil {
		t.Fatal(err)
	}
	if err := view.Resize(720, 720); err != nil {
		t.Fatal(err)
	}
	if err := view.ClickLabel("预览"); err != nil {
		t.Fatal(err)
	}
	if !state.themeEditor.previewOnly {
		t.Fatal("narrow preview tab did not open")
	}
	if shots != "" {
		if err := view.Wait(300 * time.Millisecond); err != nil {
			t.Fatal(err)
		}
		if err := view.Screenshot(filepath.Join(shots, "theme-editor-narrow.png")); err != nil {
			t.Fatal(err)
		}
	}
	if err := view.ClickLabel("保存并应用"); err != nil {
		t.Fatal(err)
	}
	if state.modal != modalNone || len(state.managedSchemes) != 1 {
		t.Fatal("save did not persist editor")
	}
}

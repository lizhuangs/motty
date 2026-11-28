package workbench

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestCustomThemeUsesGlobalBackgroundVisibility(t *testing.T) {
	picked := filepath.Join(t.TempDir(), "wall.png")
	writePNG(t, picked, 64, 64)
	state, _ := backgroundTestState(t, picked)
	state.Click(hitBgChoose)
	waitFor(t, state, state.backgroundOn)
	inTerminal(t, state)
	var saved []savedCustomScheme
	if err := json.Unmarshal([]byte(`[{"id":"custom:legacy","name":"Legacy","background":[24,12,40],"opacity":50}]`), &saved); err != nil {
		t.Fatal(err)
	}
	state.restoreCustomSchemes(saved)
	state.setTerminalScheme("custom:legacy")
	state.bg.visibility.Value = 50
	if got := state.terminalGround().Alpha; got != .5 {
		t.Fatalf("terminal background alpha = %v, want 0.5", got)
	}
	state.chooseAppearanceMode(appearanceFollowTerminal)
	if got := state.applicationOpacity(); got != 100 {
		t.Fatalf("follow opacity = %v, want 100", got)
	}
	state.openThemeEditor(false, false)
	state.themeEditor.name.SetText("Edited legacy")
	state.clickThemeEditor(hitThemeApply)
	if got := state.terminalGround().Alpha; got != .5 {
		t.Fatalf("edited terminal background alpha = %v, want 0.5", got)
	}
	data, err := json.Marshal(state.savedCustomSchemes())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `"opacity"`) {
		t.Fatalf("theme still saves opacity: %s", data)
	}
}

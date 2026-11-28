package workbench

import (
	"github.com/vivaui/viva/components"
	"github.com/vivaui/viva/gfx"
)

type applicationColorToken struct {
	key, en, zh, tw string
	color           func(*palette) *gfx.Color
}

var applicationColors = [...]applicationColorToken{
	{"shell", "Title bar and sidebar", "标题栏和侧栏", "標題列和側欄", func(p *palette) *gfx.Color { return &p.shell }},
	{"workspace", "Workspace", "工作区", "工作區", func(p *palette) *gfx.Color { return &p.workspace }},
	{"surface", "Panels and controls", "面板和控件", "面板和控制項", func(p *palette) *gfx.Color { return &p.surface }},
	{"text", "Text", "正文", "正文", func(p *palette) *gfx.Color { return &p.text }},
	{"muted", "Secondary text", "次要文字", "次要文字", func(p *palette) *gfx.Color { return &p.muted }},
	{"accent", "Accent", "强调色", "強調色", func(p *palette) *gfx.Color { return &p.accent }},
	{"border", "Border", "边框", "邊框", func(p *palette) *gfx.Color { return &p.line }},
	{"hover", "Hover", "悬停", "懸停", func(p *palette) *gfx.Color { return &p.surfaceSoft }},
	{"pressed", "Pressed", "按下", "按下", func(p *palette) *gfx.Color { return &p.surfaceHot }},
	{"tab", "Tab", "标签页", "分頁", func(p *palette) *gfx.Color { return &p.tab }},
	{"tabHover", "Hovered tab", "悬停标签页", "懸停分頁", func(p *palette) *gfx.Color { return &p.tabHover }},
	{"tabActive", "Active tab", "活动标签页", "使用中分頁", func(p *palette) *gfx.Color { return &p.tabActive }},
}

const applicationPrimaryColors = 7

func (state *State) resolveApplicationPalette(look appearance) *palette {
	switch state.appearanceMode {
	case appearanceFollowTerminal:
		if scheme, ok := state.schemeByID(state.followScheme()); ok {
			pal := applicationPaletteForScheme(*scheme)
			return &pal
		}
	case appearanceCustom:
		return state.applicationTheme.palette()
	}
	return look.palette()
}

func applicationPaletteForScheme(scheme colorScheme) palette {
	pal := *applicationColorAppearance(scheme.background).palette()
	pal.workspace, pal.text = applicationRGB(scheme.background), applicationRGB(scheme.foreground)
	if contrastRatio(scheme.foreground, scheme.background) < 4.5 {
		pal.text = applicationTextAgainst(scheme.background)
	}
	pal.shell = readableApplicationColor(pal.text, blendApplicationColor(pal.workspace, pal.text, 0.045), pal.workspace, 4.5)
	pal.surface = readableApplicationColor(pal.text, blendApplicationColor(pal.workspace, pal.text, 0.085), pal.workspace, 4.5)
	pal.muted = readableApplicationColor(pal.surface, blendApplicationColor(pal.workspace, pal.text, 0.72), pal.text, 4.5)
	pal.line = blendApplicationColor(pal.workspace, pal.text, 0.22)
	pal.accent = readableApplicationColor(pal.surface, applicationRGB(scheme.cursor), pal.text, 4.5)
	pal.success = readableApplicationColor(pal.surface, applicationRGB(scheme.ansi[2]), pal.text, 3)
	pal.warning = readableApplicationColor(pal.surface, applicationRGB(scheme.ansi[3]), pal.text, 3)
	pal.danger = readableApplicationColor(pal.surface, applicationRGB(scheme.ansi[1]), pal.text, 3)
	pal.deriveApplicationColors()
	return pal
}

func applicationTextAgainst(background [3]uint8) gfx.Color {
	black, white := [3]uint8{}, [3]uint8{255, 255, 255}
	if contrastRatio(background, white) > contrastRatio(background, black) {
		return applicationRGB(white)
	}
	return applicationRGB(black)
}

func readableApplicationColor(background, color, fallback gfx.Color, ratio float32) gfx.Color {
	for _, fraction := range [...]float32{0, 0.25, 0.5, 0.75} {
		candidate := blendApplicationColor(color, fallback, fraction)
		if contrastRatio(rgb8(background), rgb8(candidate)) >= ratio {
			return candidate
		}
	}
	return fallback
}

func (config applicationTheme) palette() *palette {
	pal := *config.baseAppearance().palette()
	pal.dark = config.colorAppearance() == appearanceDark
	apply := func(tokens []applicationColorToken) bool {
		overridden := false
		for _, token := range tokens {
			if color, ok := components.ParseColorHex(config.Colors[token.key], false); ok {
				*token.color(&pal) = color
				overridden = true
			}
		}
		return overridden
	}
	if apply(applicationColors[:applicationPrimaryColors]) {
		pal.deriveApplicationColors()
	}
	apply(applicationColors[applicationPrimaryColors:])
	if _, ok := components.ParseColorHex(config.Colors["tabActive"], false); ok {
		pal.tabCloseHot = blendApplicationColor(pal.tabActive, pal.text, 0.10)
	}
	return &pal
}

func (pal *palette) deriveApplicationColors() {
	pal.customApplicationColors = true
	pal.surfaceSoft = blendApplicationColor(pal.surface, pal.text, 0.08)
	pal.surfaceHot = blendApplicationColor(pal.surface, pal.text, 0.18)
	pal.tab = blendApplicationColor(pal.shell, pal.text, 0.05)
	pal.tabHover = blendApplicationColor(pal.shell, pal.text, 0.13)
	pal.tabActive = blendApplicationColor(pal.shell, pal.text, 0.21)
	pal.tabCloseHot = blendApplicationColor(pal.tabActive, pal.text, 0.10)
	pal.faint = blendApplicationColor(pal.workspace, pal.muted, 0.76)
	pal.accentHot = blendApplicationColor(pal.accent, pal.text, 0.22)
	pal.accentSoft = blendApplicationColor(pal.surface, pal.accent, 0.13)
	pal.windowStroke = blendApplicationColor(pal.line, pal.text, 0.10)
	pal.onAccent = applicationTextAgainst(rgb8(pal.accent))
}

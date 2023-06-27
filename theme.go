package main

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var darkScheme = map[fyne.ThemeColorName]color.Color{
	theme.ColorBlue:                  color.RGBA{0x35, 0x84, 0xe4, 0xff}, // Adwaita color name @blue_3
	theme.ColorBrown:                 color.RGBA{0x98, 0x6a, 0x44, 0xff}, // Adwaita color name @brown_3
	theme.ColorGray:                  color.RGBA{0x5e, 0x5c, 0x64, 0xff}, // Adwaita color name @dark_2
	theme.ColorGreen:                 color.RGBA{0x26, 0xa2, 0x69, 0xff}, // Adwaita color name @green_5
	theme.ColorNameBackground:        color.RGBA{0x24, 0x24, 0x24, 0xff}, // Adwaita color name @window_bg_color
	theme.ColorNameButton:            color.RGBA{0x30, 0x30, 0x30, 0xff}, // Adwaita color name @headerbar_bg_color
	theme.ColorNameError:             color.RGBA{0xc0, 0x1c, 0x28, 0xff}, // Adwaita color name @error_bg_color
	theme.ColorNameForeground:        color.RGBA{0xef, 0xef, 0xef, 0xff}, // Adwaita color name @window_fg_color
	theme.ColorNameInputBackground:   color.RGBA{0x1e, 0x1e, 0x1e, 0xff}, // Adwaita color name @view_bg_color
	theme.ColorNameMenuBackground:    color.RGBA{0x1e, 0x1e, 0x1e, 0xff}, // Adwaita color name @view_bg_color
	theme.ColorNameOverlayBackground: color.RGBA{0x1e, 0x1e, 0x1e, 0xff}, // Adwaita color name @view_bg_color
	theme.ColorNamePrimary:           color.RGBA{0x35, 0x84, 0xe4, 0xff}, // Adwaita color name @accent_bg_color
	theme.ColorNameSelection:         color.RGBA{0x35, 0x84, 0xe4, 0xff}, // Adwaita color name @accent_bg_color
	theme.ColorNameShadow:            color.RGBA{0x00, 0x00, 0x00, 0x5b}, // Adwaita color name @shade_color
	theme.ColorNameSuccess:           color.RGBA{0x26, 0xa2, 0x69, 0xff}, // Adwaita color name @success_bg_color
	theme.ColorNameWarning:           color.RGBA{0xcd, 0x93, 0x09, 0xff}, // Adwaita color name @warning_bg_color
	theme.ColorOrange:                color.RGBA{0xff, 0x78, 0x00, 0xff}, // Adwaita color name @orange_3
	theme.ColorPurple:                color.RGBA{0x91, 0x41, 0xac, 0xff}, // Adwaita color name @purple_3
	theme.ColorRed:                   color.RGBA{0xc0, 0x1c, 0x28, 0xff}, // Adwaita color name @red_4
	theme.ColorYellow:                color.RGBA{0xf6, 0xd3, 0x2d, 0xff}, // Adwaita color name @yellow_3
	theme.ColorNameSeparator:         color.RGBA{0x00, 0x00, 0x00, 0x00},
}

var lightScheme = map[fyne.ThemeColorName]color.Color{
	theme.ColorBlue:                  color.RGBA{0x35, 0x84, 0xe4, 0xff}, // Adwaita color name @blue_3
	theme.ColorBrown:                 color.RGBA{0x98, 0x6a, 0x44, 0xff}, // Adwaita color name @brown_3
	theme.ColorGray:                  color.RGBA{0x5e, 0x5c, 0x64, 0xff}, // Adwaita color name @dark_2
	theme.ColorGreen:                 color.RGBA{0x2e, 0xc2, 0x7e, 0xff}, // Adwaita color name @green_4
	theme.ColorNameBackground:        color.RGBA{0xfa, 0xfa, 0xfa, 0xff}, // Adwaita color name @window_bg_color
	theme.ColorNameButton:            color.RGBA{0xeb, 0xeb, 0xeb, 0xff}, // Adwaita color name @headerbar_bg_color
	theme.ColorNameError:             color.RGBA{0xe0, 0x1b, 0x24, 0xff}, // Adwaita color name @error_bg_color
	theme.ColorNameForeground:        color.RGBA{0x3d, 0x3d, 0x3d, 0xff}, // Adwaita color name @window_fg_color
	theme.ColorNameInputBackground:   color.RGBA{0xff, 0xff, 0xff, 0xff}, // Adwaita color name @view_bg_color
	theme.ColorNameMenuBackground:    color.RGBA{0xff, 0xff, 0xff, 0xff}, // Adwaita color name @view_bg_color
	theme.ColorNameOverlayBackground: color.RGBA{0xff, 0xff, 0xff, 0xff}, // Adwaita color name @view_bg_color
	theme.ColorNamePrimary:           color.RGBA{0x35, 0x84, 0xe4, 0xff}, // Adwaita color name @accent_bg_color
	theme.ColorNameShadow:            color.RGBA{0x00, 0x00, 0x00, 0x11}, // Adwaita color name @shade_color
	theme.ColorNameSuccess:           color.RGBA{0x2e, 0xc2, 0x7e, 0xff}, // Adwaita color name @success_bg_color
	theme.ColorNameWarning:           color.RGBA{0xe5, 0xa5, 0x0a, 0xff}, // Adwaita color name @warning_bg_color
	theme.ColorOrange:                color.RGBA{0xff, 0x78, 0x00, 0xff}, // Adwaita color name @orange_3
	theme.ColorPurple:                color.RGBA{0x91, 0x41, 0xac, 0xff}, // Adwaita color name @purple_3
	theme.ColorRed:                   color.RGBA{0xe0, 0x1b, 0x24, 0xff}, // Adwaita color name @red_3
	theme.ColorYellow:                color.RGBA{0xf6, 0xd3, 0x2d, 0xff}, // Adwaita color name @yellow_3
	theme.ColorNameSeparator:         color.RGBA{0x00, 0x00, 0x00, 0x00},
	// theme.ColorNameSelection: color.RGBA{0xff, 0x00, 0x00, 0xff},
}

type CustomTheme struct{}

func NewCustomTheme() *CustomTheme {
	t := CustomTheme{}
	return &t
}

func (t CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch variant {
	case theme.VariantLight:
		if c, ok := lightScheme[name]; ok {
			return c
		}
	case theme.VariantDark:
		if c, ok := darkScheme[name]; ok {
			return c
		}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (t CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	switch name {
	// custom icon names here...
	default:
		return theme.DefaultTheme().Icon(name)
	}
}

func (t CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

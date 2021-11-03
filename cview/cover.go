package cview

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const logo = `
_________                            .__         .___   _____   
\_   ___ \  ____   ____   __________ |  |   ____ |   | /     \  
/    \  \/ /  _ \ /    \ /  ___/  _ \|  | _/ __ \|   |/  \ /  \ 
\     \___(  <_> )   |  \\___ (  <_> )  |_\  ___/|   /    Y    \
 \______  /\____/|___|  /____  >____/|____/\___  >___\____|__  /
        \/            \/     \/                \/            \/ 
`

const (
	subtitle   = `ConsoleIM - 一款基于 tview 构建的简易聊天程序`
	navigation = `Ctrl-N: 下一项    Ctrl-P: 上一项    Ctrl-C: 退出`
	mouse      = `(也可以用鼠标点击)`
)

// Cover returns the cover page.
func Cover(nextSlide func()) (title string, content tview.Primitive) {
	// What's the size of the logo?
	lines := strings.Split(logo, "\n")
	logoWidth := 0
	logoHeight := len(lines)
	for _, line := range lines {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}
	logoBox := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetDoneFunc(func(key tcell.Key) {
			nextSlide()
		})
	fmt.Fprint(logoBox, logo)

	// Create a frame for the subtitle and navigation infos.
	frame := tview.NewFrame(tview.NewBox()).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(subtitle, true, tview.AlignCenter, tcell.ColorWhite).
		AddText("", true, tview.AlignCenter, tcell.ColorWhite).
		AddText(navigation, true, tview.AlignCenter, tcell.ColorDarkMagenta).
		AddText(mouse, true, tview.AlignCenter, tcell.ColorDarkMagenta)

	// Create a Flex layout that centers the logo and subtitle.
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 7, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(logoBox, logoWidth, 1, true).
			AddItem(tview.NewBox(), 0, 1, false), logoHeight, 1, true).
		AddItem(frame, 0, 10, false)

	return "Start", flex
}

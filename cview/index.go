/*
A presentation of the tview package, implemented with tview.

Navigation

The presentation will advance to the next slide when the primitive demonstrated
in the current slide is left (usually by hitting Enter or Escape). Additionally,
the following shortcuts can be used:

  - Ctrl-N: Jump to next slide
  - Ctrl-P: Jump to previous slide
*/
package cview

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CView struct {
	IsLogin       bool
	App           *tview.Application
	Pages         *tview.Pages
	CurrentSlide  Slide
	PreviousSlide Slide
	NextSlide     Slide
	Slides        []Slide
}

// Slide is a function which returns the slide's main primitive and its title.
// It receives a "nextSlide" function which can be called to advance the
// presentation to the next slide.
type Slide func(nextSlide func()) (title string, content tview.Primitive)

func NewCView() *CView {
	f := &CView{
		IsLogin:       false,
		App:           tview.NewApplication(),
		Pages:         tview.NewPages(),
		CurrentSlide:  nil,
		PreviousSlide: nil,
		NextSlide:     nil,
		Slides:        make([]Slide, 10),
	}

	return f
}

// Starting point for the presentation.
func (c *CView) Index() {
	// The presentation slides.
	c.Slides[0] = Cover
	//如果登录成功则不再显示登录页面
	if !c.IsLogin {
		c.Slides[1] = c.Login
	}
	c.Slides[2] = c.Chat
	//初始化app
	c.App = tview.NewApplication()
	//初始化页面
	c.Pages = tview.NewPages()

	// The bottom row has some info on where we are.
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			c.Pages.SwitchToPage(added[0])
		})

	// Create the pages for all slides.
	previousSlide := func() {
		slide, _ := strconv.Atoi(info.GetHighlights()[0])
		slide = (slide - 1 + len(c.Slides)) % len(c.Slides)
		info.Highlight(strconv.Itoa(slide)).
			ScrollToHighlight()
	}
	nextSlide := func() {
		slide, _ := strconv.Atoi(info.GetHighlights()[0])
		slide = (slide + 1) % len(c.Slides)
		info.Highlight(strconv.Itoa(slide)).
			ScrollToHighlight()
	}
	for index, slide := range c.Slides {
		if slide == nil {
			//跳过无效的
			continue
		}
		title, primitive := slide(nextSlide)
		c.Pages.AddPage(strconv.Itoa(index), primitive, true, index == 0)
		fmt.Fprintf(info, `%d ["%d"][darkcyan]%s[white][""]  `, index+1, index, title)
	}
	info.Highlight("0")

	// Create the main layout.
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(c.Pages, 0, 1, true).
		AddItem(info, 1, 1, false)

	// Shortcuts to navigate the slides.
	c.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			nextSlide()
			return nil
		} else if event.Key() == tcell.KeyCtrlP {
			previousSlide()
			return nil
		}
		return event
	})

	// Start the application.
	if err := c.App.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
func (c *CView) CheckLogin() {
	if !c.IsLogin {
		c.Pages.SwitchToPage("1")
	}
}
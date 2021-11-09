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
	"container/list"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/gtck520/ConsoleIM/common/http"
	"github.com/gtck520/ConsoleIM/common/util/ext"
	"github.com/gtck520/ConsoleIM/websocket"
	"github.com/rivo/tview"
)

type CView struct {
	IsLogin       bool
	App           *tview.Application
	Pages         *tview.Pages
	CurrentSlide  string
	PreviousSlide string
	Slides        []Slide
	TabInfo       *tview.TextView
	Api           *http.Api
	FriendsList   []interface{}
	UserInfo      map[string]interface{}
	ChatToUserId  int                //当前聊天用户
	MessageList   map[int]*list.List //用户消息列表
	Ws            websocket.Ws
	TextView      *tview.TextView
}

//不需要检查登录的页面
var Exclude_Check []string

// Slide is a function which returns the slide's main primitive and its title.
// It receives a "nextSlide" function which can be called to advance the
// presentation to the next slide.
type Slide func(nextSlide func()) (title string, content tview.Primitive)

//发送内容content 载体
type TypeMessage struct {
	Type interface{} `json:"type"` //内容分发类型：ping、init、message 等等
	Data interface{} `json:"data"`
}

//TypeMessage Data
type ClientMessage struct {
	Name     string `json:"name"`
	Avator   string `json:"avator"`
	Id       string `json:"id"`
	Group    string `json:"group"`
	Time     string `json:"time"`
	ToId     string `json:"to_id"`
	Content  string `json:"content"`
	City     string `json:"city"`
	ClientIp string `json:"client_ip"`
	Refer    string `json:"refer"`
}

//初始化一个显示客户端
func NewCView() *CView {
	//不需要检查登录的页面
	Exclude_Check = []string{"0", "1"}
	f := &CView{
		IsLogin:       false,
		App:           tview.NewApplication(),
		Pages:         tview.NewPages(),
		CurrentSlide:  "0",
		PreviousSlide: "0",
		Slides:        make([]Slide, 10),
		TabInfo:       tview.NewTextView(), //底部切换栏
		Api:           http.NewApi(),
		FriendsList:   make([]interface{}, 0),
		UserInfo:      make(map[string]interface{}),
		MessageList:   make(map[int]*list.List),
		Ws:            websocket.Ws{},
		TextView:      tview.NewTextView(),
	}

	return f
}

// Starting point for the presentation.
func (c *CView) Index() {
	// The presentation slides.
	c.Slides[0] = c.Cover //此处key值 跳转页面时需要用到 需要对应
	//如果登录成功则不再显示登录页面
	if !c.IsLogin {
		c.Slides[1] = c.Login
	}
	c.Slides[2] = c.Chat
	//初始化app
	c.App = tview.NewApplication()
	//初始化页面
	c.Pages = tview.NewPages()

	// The bottom row has some c.TabInfo on where we are.
	c.TabInfo.SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			c.PreviousSlide = c.CurrentSlide
			c.CurrentSlide = added[0]
			c.CheckJump(func() {
				c.Pages.SwitchToPage(added[0])
			}, added[0])

		})

	// Create the pages for all slides.
	previousSlide := func() {
		slide, _ := strconv.Atoi(c.TabInfo.GetHighlights()[0])
		slide = (slide - 1 + len(c.Slides)) % len(c.Slides)
		c.PreviousSlide = c.CurrentSlide
		c.CurrentSlide = strconv.Itoa(slide)
		c.CheckJump(func() {
			c.TabInfo.Highlight(strconv.Itoa(slide)).
				ScrollToHighlight()
		}, strconv.Itoa(slide))

	}
	nextSlide := func() {

		slide, _ := strconv.Atoi(c.TabInfo.GetHighlights()[0])
		slide = (slide + 1) % len(c.Slides)
		c.PreviousSlide = c.CurrentSlide
		c.CurrentSlide = strconv.Itoa(slide)
		c.CheckJump(func() {
			c.TabInfo.Highlight(strconv.Itoa(slide)).
				ScrollToHighlight()
		}, strconv.Itoa(slide))

	}
	for index, slide := range c.Slides {
		if slide == nil {
			//跳过无效的
			continue
		}
		title, primitive := slide(nextSlide)
		c.Pages.AddPage(strconv.Itoa(index), primitive, true, index == 0)
		fmt.Fprintf(c.TabInfo, `%d ["%d"][darkcyan]%s[white][""]  `, index+1, index, title)
	}
	c.TabInfo.Highlight("0")

	// Create the main layout.
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(c.Pages, 0, 1, true).
		AddItem(c.TabInfo, 1, 1, false)

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
func (c *CView) CheckLogin() bool {
	if c.IsLogin {
		result := c.Api.Info()
		if result.Code == 200 {
			Datas := result.Data.(map[string]interface{})
			c.UserInfo = Datas["userinfo"].(map[string]interface{})
			c.FriendsList = Datas["friends"].([]interface{})
			c.Reload(2, c.Chat)
		} else {
			c.Api.Header["X-Token"] = ""
			c.IsLogin = false
		}
	}
	return c.IsLogin
}

//跳转页面
func (c *CView) JumpTo(pagename string) {
	c.Pages.SwitchToPage(pagename)
	c.TabInfo.Highlight(pagename).ScrollToHighlight()
}

//检查登录与跳转页面 passDone 检查通过后执行
func (c *CView) CheckJump(passDone func(), slide string) {
	if c.CheckLogin() {
		passDone()
	} else {
		if !ext.In(c.CurrentSlide, Exclude_Check) {
			c.alert(c.Pages, "alert-dialog", "error", "请先登录", "1")
		} else {
			passDone()
		}
	}
}

//重载页面 index 页面编号， slide 定义的视图方法
func (c *CView) Reload(index int, slide Slide) {
	_, primitive := slide(func() {
		c.JumpTo(strconv.Itoa(index + 1))
	})
	c.Pages.RemovePage(strconv.Itoa(index))
	c.Pages.AddPage(strconv.Itoa(index), primitive, true, index == 0)
}

//从socket读取消息到视图
func (c *CView) ReadMessage() {

	go func() {
		for {
			_, message, err := c.Ws.Conn.ReadMessage()
			if err != nil {
				c.Ws.Log.Error("read:", err)
				return
			}
			typemessage := TypeMessage{}
			json.Unmarshal(message, &typemessage)
			c.Ws.Log.Infof("typemessage: %s", typemessage)
			reflectType := reflect.TypeOf(typemessage.Data)
			c.Ws.Log.Infof("reflecttype: %s", reflectType.Name())
			if reflectType.Name() == "" {
				clientmessage := typemessage.Data.(map[string]interface{})
				c.Ws.Log.Infof("clientmessage: %s", clientmessage)
				if clientmessage["to_id"] != nil {
					fromuserid, _ := strconv.Atoi(clientmessage["id"].(string))
					c.ScreenAndSave(fromuserid, clientmessage["name"].(string), clientmessage["time"].(string), clientmessage["content"].(string))
				}
			}

		}
	}()
}

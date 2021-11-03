// Demo code for the Flex primitive.
package cview

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (c *CView) Chat(nextSlide func()) (title string, content tview.Primitive) {
	//消息窗
	textView := tview.NewTextView().
		SetTextColor(tcell.ColorYellow).
		SetScrollable(false).
		SetDoneFunc(func(key tcell.Key) {
			nextSlide()
		})
	textView.SetChangedFunc(func() {
		if textView.HasFocus() {
			c.App.Draw()
		}
	})
	go func() {
		var n int
		for {
			if textView.HasFocus() {
				n++
				if n > 512 {
					n = 1
					textView.SetText("")
				}

				fmt.Fprintf(textView, "%d ", n)
				time.Sleep(200 * time.Millisecond)
			} else {
				time.Sleep(time.Second)
			}
		}
	}()
	textView.SetBorder(true).SetTitle("  消息")
	//输入窗口
	inputform := tview.NewForm()
	inputform.AddInputField("请输入消息:", "", 50, nil, nil).
		AddButton("发送", func() {
			message := inputform.GetFormItem(0).(*tview.InputField).GetText()
			fmt.Fprintln(textView, message)
			inputform.GetFormItem(0).(*tview.InputField).SetText("")
		}).
		SetHorizontal(true)
	inputform.SetBorder(true).SetTitle("")

	//整体框架
	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("  好友列表 "), 0, 1, false). //Left (1/2 x width of Top)
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("  信息栏 "), 0, 1, false). //Top
			AddItem(textView, 0, 3, false).                                          //Middle (3 x height of Top)
			AddItem(inputform, 5, 1, false), 0, 2, false)                            //Bottom (5 rows)
		//AddItem(tview.NewBox().SetBorder(true).SetTitle("待定"), 20, 1, false)  //Right (20 cols)
	return "Chat", flex
}

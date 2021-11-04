package cview

import (
	"fmt"

	"github.com/rivo/tview"
)

// Form demonstrates forms.
func (c *CView) Login(nextSlide func()) (title string, content tview.Primitive) {
	c.PreviousSlide = c.CurrentSlide
	c.CurrentSlide = "1"
	f := tview.NewForm()
	f.AddInputField("用户名:", "", 20, nil, nil).
		//AddInputField("Last name:", "", 20, nil, nil).
		//AddDropDown("Role:", []string{"Engineer", "Manager", "Administration"}, 0, nil).
		//AddCheckbox("On vacation:", false, nil).
		AddPasswordField("密码:", "", 10, '*', nil).
		AddButton("登录", func() {
			c.ApiLogin(f)
		}).
		AddButton("退出", func() {
			c.JumpTo("0")

		})
	f.SetBorder(true).SetTitle(" 账户登录")
	return "Login", f
}
func (c *CView) ApiLogin(form *tview.Form) {
	userName := form.GetFormItem(0).(*tview.InputField).GetText()
	userPwd := form.GetFormItem(1).(*tview.InputField).GetText()

	c.alert(c.Pages, "alert-dialog", fmt.Sprintf("保存成功，%s %s！", userName, userPwd))
}

// alert shows a confirmation dialog.
func (c *CView) alert(pages *tview.Pages, id string, message string) *tview.Pages {
	return pages.AddPage(
		id,
		tview.NewModal().
			SetText(message).
			AddButtons([]string{"确定"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.HidePage(id).RemovePage(id)
				c.JumpTo("2") //跳转聊天页面
			}),
		false,
		true,
	)
}

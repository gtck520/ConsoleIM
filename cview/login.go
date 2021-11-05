package cview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Form demonstrates forms.
func (c *CView) Login(nextSlide func()) (title string, content tview.Primitive) {
	f := tview.NewForm()
	f.AddInputField("用户名:", "18695732895", 20, nil, nil).
		//AddInputField("Last name:", "", 20, nil, nil).
		//AddDropDown("Role:", []string{"Engineer", "Manager", "Administration"}, 0, nil).
		//AddCheckbox("On vacation:", false, nil).
		AddPasswordField("密码:", "18695732895", 10, '*', nil).
		AddButton("登录", func() {
			c.ApiLogin(f)
		}).
		AddButton("退出", func() {
			c.Api.Header["X-Token"] = ""
			c.IsLogin = false
			c.JumpTo("0")

		})
	f.SetBorder(true).SetTitle(" 账户登录")
	return "Login", f
}
func (c *CView) ApiLogin(form *tview.Form) {
	userName := form.GetFormItem(0).(*tview.InputField).GetText()
	userPwd := form.GetFormItem(1).(*tview.InputField).GetText()

	result := c.Api.Login(userName, userPwd)

	if result.Code == 200 {
		data := result.Data.(map[string]interface{})
		c.Api.Header["X-Token"] = data["token"].(string)
		c.IsLogin = true
		c.alert(c.Pages, "alert-dialog", "success", "登录成功！", "2")
	} else {
		c.alert(c.Pages, "alert-dialog", "error", result.Message, "")
	}
}

// alert shows a confirmation dialog.
func (c *CView) alert(pages *tview.Pages, id string, alerttype string, message string, jumpkey string) *tview.Pages {
	modal := tview.NewModal()
	modal.SetText(message).
		AddButtons([]string{"确定"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.HidePage(id).RemovePage(id)
			if jumpkey != "" {
				c.JumpTo(jumpkey) //跳转聊天页面
			}
		})
	if alerttype == "error" {
		modal.SetBackgroundColor(tcell.ColorRed)
	}

	return pages.AddPage(
		id,
		modal,
		false,
		true,
	)
}

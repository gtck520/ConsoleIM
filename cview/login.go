package cview

import (
	"github.com/rivo/tview"
)

// Form demonstrates forms.
func (c *CView) Login(nextSlide func()) (title string, content tview.Primitive) {
	f := tview.NewForm().
		AddInputField("用户名:", "", 20, nil, nil).
		//AddInputField("Last name:", "", 20, nil, nil).
		//AddDropDown("Role:", []string{"Engineer", "Manager", "Administration"}, 0, nil).
		//AddCheckbox("On vacation:", false, nil).
		AddPasswordField("密码:", "", 10, '*', nil).
		AddButton("保存", nextSlide).
		AddButton("退出", nextSlide)
	f.SetBorder(true).SetTitle(" 账户登录")
	return "Login", f
}

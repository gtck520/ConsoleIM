// Demo code for the Flex primitive.
package cview

import (
	"container/list"
	"fmt"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//打开好友聊天
func (c *CView) ChatToUser(userid int) {
	c.TextView.SetTitle(" " + strconv.Itoa(userid) + "消息")
	c.ChatToUserId = userid //重置当前聊天id
	//将消息列表打印到屏幕
	c.TextView.SetText("") //先清空
	for e := c.MessageList[userid].Front(); e != nil; e = e.Next() {
		fmt.Fprintln(c.TextView, e.Value)
	}

}

func (c *CView) Chat(nextSlide func()) (title string, content tview.Primitive) {

	//消息窗
	c.TextView.
		SetTextColor(tcell.ColorYellow).
		SetScrollable(false).
		SetDoneFunc(func(key tcell.Key) {
			nextSlide()
		})
	c.TextView.SetChangedFunc(func() {
		// if c.TextView.HasFocus() {
		c.App.Draw()
		// }
	})
	c.TextView.SetBorder(true).SetTitle("  消息")
	//好友列表
	type node struct {
		text     string
		expand   bool
		selected func()
		children []*node
	}

	tree := tview.NewTreeView()

	var rootNode = &node{
		text: "全部分组",
		children: []*node{
			{text: "展开全部", selected: func() { tree.GetRoot().ExpandAll() }},
			{text: "收起全部", selected: func() {
				for _, child := range tree.GetRoot().GetChildren() {
					child.CollapseAll()
				}
			}},
		}}
	for _, groupI := range c.FriendsList {
		if groupI == nil {
			continue
		}
		group := groupI.(map[string]interface{})
		var newcode node
		newcode.text = group["group_name"].(string)
		newcode.expand = true
		for _, friendI := range group["group_members"].([]interface{}) {
			friend := friendI.(map[string]interface{})
			//初始化聊天记录
			list := list.New()
			c.MessageList[int(friend["friend_id"].(float64))] = list

			subcode := &node{
				text: strconv.Itoa(int(friend["friend_id"].(float64))),
				//expand: true,
				selected: func() {
					c.ChatToUser(int(friend["friend_id"].(float64)))
				},
			}
			newcode.children = append(newcode.children, subcode)
		}
		rootNode.children = append(rootNode.children, &newcode)

	}
	tree.SetBorder(true).
		SetTitle("好友列表")

	// Add nodes.
	var add func(target *node) *tview.TreeNode
	add = func(target *node) *tview.TreeNode {
		node := tview.NewTreeNode(target.text).
			SetSelectable(target.expand || target.selected != nil).
			SetExpanded(target == rootNode).
			SetReference(target)
		if target.expand {
			node.SetColor(tcell.ColorGreen)
		} else if target.selected != nil {
			node.SetColor(tcell.ColorRed)
		}
		for _, child := range target.children {
			node.AddChild(add(child))
		}
		return node
	}
	root := add(rootNode)
	tree.SetRoot(root).
		SetCurrentNode(root).
		SetSelectedFunc(func(n *tview.TreeNode) {
			original := n.GetReference().(*node)
			if original.expand {
				n.SetExpanded(!n.IsExpanded())
			} else if original.selected != nil {
				original.selected()
			}
		})

	// list := tview.NewList()
	// list.ShowSecondaryText(false).
	// 	AddItem("Basic table", "", 'b', func() {
	// 		textView.SetText("Basic table")
	// 	}).
	// 	AddItem("Table with separator", "", 's', func() {
	// 		textView.SetText("Table with separator")
	// 	}).
	// 	AddItem("Table with borders", "", 'o', func() {
	// 		textView.SetText("Table with borders")
	// 	}).
	// 	AddItem("Selectable rows", "", 'r', func() {
	// 		textView.SetText("Selectable rows")
	// 	})
	// list.SetBorderPadding(1, 1, 2, 2)
	// list.SetBorder(true).SetTitle("  好友列表 ")
	//输入窗口
	inputform := tview.NewForm()
	inputform.AddInputField("请输入消息:", "", 50, nil, nil).
		AddButton("发送", func() {
			message := inputform.GetFormItem(0).(*tview.InputField).GetText()
			result := c.Api.SendMessage(int(c.UserInfo["id"].(float64)), c.ChatToUserId, "user", message)
			timestr := time.Now().Format("2006-01-02 15:04:05")
			if result.Code == 200 {
				c.ScreenAndSave(c.ChatToUserId, "我", timestr, message)
				inputform.GetFormItem(0).(*tview.InputField).SetText("")
			}

		}).
		SetHorizontal(true)
	inputform.SetBorder(true).SetTitle("")

	//整体框架
	flex := tview.NewFlex().
		AddItem(tree, 0, 1, false). //Left (1/2 x width of Top)
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("  信息栏 "), 0, 1, false). //Top
			AddItem(c.TextView, 0, 3, false).                                        //Middle (3 x height of Top)
			AddItem(inputform, 5, 1, false), 0, 2, false)                            //Bottom (5 rows)
		//AddItem(tview.NewBox().SetBorder(true).SetTitle("待定"), 20, 1, false)  //Right (20 cols)
	return "Chat", flex
}

//发送到屏幕 并保存到本地
func (c *CView) ScreenAndSave(userid int, name string, time string, message string) {
	//消息接收时，接收对象为当前聊天对象才打印到屏幕
	if c.ChatToUserId == userid {
		sendmessage := name + ":" + message + "  at " + time
		fmt.Fprintln(c.TextView, sendmessage)
	}
	//往消息列表后插入一条
	c.MessageList[userid].PushBack(message)
}

// Demo code for the Flex primitive.
package cview

import (
	"fmt"
	"strconv"
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
			var subcode node
			subcode.text = strconv.Itoa(int(friend["friend_id"].(float64)))
			subcode.expand = true
			subcode.selected = func() {
			}
			newcode.children = append(newcode.children, &subcode)
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
			fmt.Fprintln(textView, message)
			inputform.GetFormItem(0).(*tview.InputField).SetText("")
		}).
		SetHorizontal(true)
	inputform.SetBorder(true).SetTitle("")

	//整体框架
	flex := tview.NewFlex().
		AddItem(tree, 0, 1, false). //Left (1/2 x width of Top)
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("  信息栏 "), 0, 1, false). //Top
			AddItem(textView, 0, 3, false).                                          //Middle (3 x height of Top)
			AddItem(inputform, 5, 1, false), 0, 2, false)                            //Bottom (5 rows)
		//AddItem(tview.NewBox().SetBorder(true).SetTitle("待定"), 20, 1, false)  //Right (20 cols)
	return "Chat", flex
}

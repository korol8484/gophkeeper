package bubble

//import (
//	"github.com/charmbracelet/bubbles/list"
//	tea "github.com/charmbracelet/bubbletea"
//	"github.com/charmbracelet/lipgloss"
//)
//
//type item string
//
//func (i item) Title() string       { return string(i) }
//func (i item) Description() string { return "" }
//func (i item) FilterValue() string { return string(i) }
//
//type NavigateModel struct {
//	list   list.Model
//	choice string
//	auth   *authModel
//
//	showAuth bool
//}
//
//func NewNavigateModel() *NavigateModel {
//	items := []list.Item{
//		item("Auth"),
//		item("Register"),
//	}
//
//	delegate := list.NewDefaultDelegate()
//	delegate.ShowDescription = false
//	delegate.SetSpacing(0)
//
//	l := list.New(items, delegate, 0, 0)
//	l.SetShowStatusBar(false)
//	l.SetShowFilter(false)
//	l.SetFilteringEnabled(false)
//	l.Title = "GophKeeper"
//
//	//m := &NavigateModel{
//	//	list: l,
//	//}
//
//	return &NavigateModel{
//		list: l,
//		auth: initAuthModel(),
//	}
//}
//
//func (n *NavigateModel) Init() tea.Cmd {
//	return nil
//}
//
////type authMsg struct{}
////
////func auth() tea.Msg {
////	return authMsg{}
////}
//
//func (n *NavigateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.WindowSizeMsg:
//		h, v := lipgloss.NewStyle().GetFrameSize()
//		n.list.SetSize(msg.Width-h, msg.Height-v)
//
//		return n, nil
//	case tea.KeyMsg:
//		if msg.String() == containsKey(msg.String(), []string{"q", "ctrl+c"}) {
//			return n, tea.Quit
//		}
//
//		if n.showAuth {
//			var cmd tea.Cmd
//
//			n.auth, cmd = n.auth.Update(msg)
//			return n, cmd
//		}
//
//		switch keypress := msg.String(); keypress {
//		case "enter":
//			i, ok := n.list.SelectedItem().(item)
//			if ok {
//				n.choice = string(i)
//			}
//
//			//var cmd tea.Cmd
//
//			//n.auth, cmd = n.auth.Update(msg)
//			//n.showAuth = true
//
//			return n, tea.Batch(auth)
//		}
//	case authMsg:
//		var cmd tea.Cmd
//
//		n.auth, cmd = n.auth.Update(msg)
//		n.showAuth = true
//		return n, cmd
//	}
//
//	var cmd tea.Cmd
//	n.list, cmd = n.list.Update(msg)
//
//	return n, cmd
//}
//
//func (n *NavigateModel) View() string {
//	var views []string
//	views = append(views, n.list.View())
//
//	if n.showAuth {
//		views = append(views, n.auth.View())
//	}
//
//	//views = append(views, "121`2`12`12`12")
//
//	return lipgloss.JoinHorizontal(lipgloss.Top, views...)
//
//	//return docStyle.Render(n.list.View())
//}
//
//func containsKey(v string, a []string) string {
//	for _, i := range a {
//		if i == v {
//			return v
//		}
//	}
//	return ""
//}

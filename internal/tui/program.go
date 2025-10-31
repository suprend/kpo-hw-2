package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"kpo-hw-2/internal/application/facade"
	"kpo-hw-2/internal/tui/styles"
)

// Model coordinates screen stack and routes Bubble Tea messages.
type Model struct {
	ctx   *programContext
	stack []Screen
}

// NewProgram constructs Bubble Tea model with provided root screen.
func NewProgram(
	account facade.AccountFacade,
	category facade.CategoryFacade,
	operation facade.OperationFacade,
	root Screen,
) *Model {
	m := &Model{
		ctx: &programContext{
			account:   account,
			category:  category,
			operation: operation,
		},
	}

	if root != nil {
		m.stack = append(m.stack, root)
	}

	return m
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	if current := m.current(); current != nil {
		return current.Init(m.ctx)
	}
	return tea.Quit
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	current := m.current()
	if current == nil {
		return m, tea.Quit
	}

	res := current.Update(msg, m.ctx)

	var cmds []tea.Cmd
	if res.Cmd != nil {
		cmds = append(cmds, res.Cmd)
	}

	if res.Pop {
		m.pop()
	}

	if res.Replace != nil {
		m.replace(res.Replace, &cmds)
	}

	if res.Push != nil {
		m.push(res.Push, &cmds)
	}

	if len(m.stack) == 0 {
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model.
func (m *Model) View() string {
	current := m.current()
	if current == nil {
		return ""
	}

	body := current.View()
	title := current.Name()
	if title == "" {
		return body
	}

	titleStyled := styles.TitleStyle.Render(title)
	return fmt.Sprintf("%s\n\n%s", titleStyled, body)
}

func (m *Model) current() Screen {
	if len(m.stack) == 0 {
		return nil
	}
	return m.stack[len(m.stack)-1]
}

func (m *Model) pop() {
	if len(m.stack) == 0 {
		return
	}
	m.stack = m.stack[:len(m.stack)-1]
}

func (m *Model) replace(screen Screen, cmds *[]tea.Cmd) {
	if len(m.stack) == 0 {
		m.stack = append(m.stack, screen)
	} else {
		m.stack[len(m.stack)-1] = screen
	}

	if initCmd := screen.Init(m.ctx); initCmd != nil {
		*cmds = append(*cmds, initCmd)
	}
}

func (m *Model) push(screen Screen, cmds *[]tea.Cmd) {
	m.stack = append(m.stack, screen)
	if initCmd := screen.Init(m.ctx); initCmd != nil {
		*cmds = append(*cmds, initCmd)
	}
}

type programContext struct {
	account   facade.AccountFacade
	category  facade.CategoryFacade
	operation facade.OperationFacade
}

func (c *programContext) Accounts() facade.AccountFacade     { return c.account }
func (c *programContext) Categories() facade.CategoryFacade  { return c.category }
func (c *programContext) Operations() facade.OperationFacade { return c.operation }

var _ ScreenContext = (*programContext)(nil)

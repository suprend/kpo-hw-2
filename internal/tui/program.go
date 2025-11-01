package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	accountcmd "kpo-hw-2/internal/application/command/account"
	categorycmd "kpo-hw-2/internal/application/command/category"
	exportcmd "kpo-hw-2/internal/application/command/export"
	fileimportcmd "kpo-hw-2/internal/application/command/import"
	operationcmd "kpo-hw-2/internal/application/command/operation"
	"kpo-hw-2/internal/tui/styles"
)

// Model coordinates screen stack and routes Bubble Tea messages.
type Model struct {
	ctx   *programContext
	stack []Screen
}

// NewProgram constructs Bubble Tea model with provided root screen.
func NewProgram(
	baseCtx context.Context,
	accountCommands *accountcmd.Service,
	categoryCommands *categorycmd.Service,
	operationCommands *operationcmd.Service,
	exportCommands *exportcmd.Service,
	importCommands *fileimportcmd.Service,
	root Screen,
) *Model {
	if baseCtx == nil {
		baseCtx = context.Background()
	}
	m := &Model{
		ctx: &programContext{
			ctx:               baseCtx,
			accountCommands:   accountCommands,
			categoryCommands:  categoryCommands,
			operationCommands: operationCommands,
			exportCommands:    exportCommands,
			importCommands:    importCommands,
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
	ctx               context.Context
	accountCommands   *accountcmd.Service
	categoryCommands  *categorycmd.Service
	operationCommands *operationcmd.Service
	exportCommands    *exportcmd.Service
	importCommands    *fileimportcmd.Service
}

func (c *programContext) Context() context.Context {
	return c.ctx
}

func (c *programContext) AccountCommands() *accountcmd.Service {
	return c.accountCommands
}
func (c *programContext) CategoryCommands() *categorycmd.Service {
	return c.categoryCommands
}
func (c *programContext) OperationCommands() *operationcmd.Service {
	return c.operationCommands
}
func (c *programContext) ExportCommands() *exportcmd.Service {
	return c.exportCommands
}
func (c *programContext) ImportCommands() *fileimportcmd.Service {
	return c.importCommands
}

var _ ScreenContext = (*programContext)(nil)

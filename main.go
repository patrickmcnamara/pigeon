package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/patrickmcnamara/bird"
	"github.com/patrickmcnamara/bird/seed"
)

type model struct {
	vp    viewport.Model
	ti    textinput.Model
	ready bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			if m.ti.Focused() {
				m.ti.Blur()
			} else {
				m.ti.Focus()
			}
			return m, nil
		case tea.KeyEnter:
			m.vp.GotoTop()
			m.ti.Blur()
			if m.ti.Value() == "" {
				return m, nil
			}
			return m, func() tea.Msg {
				sr, err := bird.Fetch(m.ti.Value())
				if err != nil {
					buf := new(bytes.Buffer)
					sw := seed.NewWriter(buf)
					sw.Header(1, "ERROR from pigeon")
					sw.Break()
					sw.Code()
					sw.Text(err.Error())
					return seed.NewReader(buf)
				}
				return sr
			}
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.vp = viewport.Model{Width: msg.Width, Height: msg.Height - 5, YPosition: 2}
			m.ready = true
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = msg.Height - 5
		}
	case *seed.Reader:
		m.vp.SetContent(seedToText(msg))
		return m, nil
	}
	var cmd tea.Cmd
	if m.ti.Focused() {
		m.ti, cmd = m.ti.Update(msg)
	} else {
		m.vp, cmd = m.vp.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	// no rendering until ready
	if !m.ready {
		return ""
	}

	// header
	header := m.ti.View() + " " + strings.Repeat("─", m.vp.Width-m.ti.Width-1)

	// body
	body := m.vp.View()
	if strings.Trim(body, "\n") != "" {
		body += "\n"
	}

	// footer
	footer := strings.Repeat("─", m.vp.Width-5) + " " + fmt.Sprintf("%3.f%%", m.vp.ScrollPercent()*100)

	// all together now
	return strings.Join([]string{header, body, footer}, "\n\n")
}

func main() {
	// setup textbox
	ti := textinput.NewModel()
	ti.Placeholder = "bird://"

	prog := tea.NewProgram(model{ti: ti})
	prog.Start()
}

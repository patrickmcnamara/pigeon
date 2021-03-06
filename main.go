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
				rawurl := m.ti.Value()
				sr, c, err := bird.Fetch(rawurl)
				if err != nil {
					return err
				}
				return birdResponse{sr: sr, close: c}
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
	case birdResponse:
		m.vp.SetContent(seedToText(msg.sr))
		return m, func() tea.Msg {
			return msg.close()
		}
	case error:
		buf := new(bytes.Buffer)
		sw := seed.NewWriter(buf)
		sw.Header(1, "ERROR from pigeon")
		sw.Text("")
		sw.Text(msg.Error())
		sr := seed.NewReader(buf)
		m.vp.SetContent(seedToText(sr))
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
	header := "> bird:" + m.ti.View() + " " + strings.Repeat("─", m.vp.Width-m.ti.Width-8)

	// body
	body := m.vp.View()
	body += strings.Repeat("\n", m.vp.Height-strings.Count(body, "\n"))

	// footer
	footer := strings.Repeat("─", m.vp.Width-5) + " " + fmt.Sprintf("%3.f%%", m.vp.ScrollPercent()*100)

	// all together now
	return strings.Join([]string{header, body, footer}, "\n\n")
}

func main() {
	// setup textbox
	ti := textinput.NewModel()
	ti.Prompt = ""
	ti.Placeholder = "//"

	prog := tea.NewProgram(model{ti: ti})
	prog.Start()
}

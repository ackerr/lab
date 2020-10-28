package internal

import (
	"fmt"

	"github.com/ackerr/lab/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xanzy/go-gitlab"
)

func NewJobUI(client *gitlab.Client, pid interface{}, jobs []*gitlab.Job) JobModel {
	return JobModel{
		client:   client,
		pid:      pid,
		choices:  jobs,
		selected: make(map[int]*gitlab.Job),
	}
}

type JobModel struct {
	client   *gitlab.Client
	pid      interface{}
	choices  []*gitlab.Job
	cursor   int
	selected map[int]*gitlab.Job
}

func (m JobModel) Init() tea.Cmd {
	return nil
}

func (m JobModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "tab", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = m.choices[m.cursor]
			}
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "o", "O":
			job := m.choices[m.cursor]
			err := utils.OpenBrowser(job.WebURL)
			utils.Check(err)
			return m, tea.Quit
		case "enter":
			println("running")
		}
	}
	return m, nil
}

func (m JobModel) View() string {
	s := "All job done!\n\n"

	for i, job := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %-20s\t%-10s\t%-10s\n", cursor, checked, job.Name, job.Status, job.Stage)
	}

	s += "\nPress o to open job page in browser."
	s += "\nPress q to quit.\n"
	return s
}

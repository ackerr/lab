package ui

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xanzy/go-gitlab"
)

type event int

const (
	maxLength       = 20
	intervalRefresh = 2 * time.Second
)

const (
	choice event = iota
	help
	trace
)

func NewJobUI(client *gitlab.Client, pid interface{}, jobs []*gitlab.Job, tailLine int64) JobModel {
	return JobModel{
		wg:       &sync.WaitGroup{},
		client:   client,
		pid:      pid,
		choices:  jobs,
		selected: make(map[int]*gitlab.Job),
		tailLine: tailLine,
		event:    0,
	}
}

type JobModel struct {
	wg       *sync.WaitGroup
	client   *gitlab.Client
	pid      interface{}
	choices  []*gitlab.Job
	cursor   int
	selected map[int]*gitlab.Job
	tailLine int64
	event    event
}

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(intervalRefresh)
	return tickMsg{}
}

func (m JobModel) Init() tea.Cmd {
	return tick
}

func (m JobModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.wg.Wait()
	if m.event > 0 {
		return eventCallback(msg, m)
	}
	switch msg := msg.(type) {
	case tickMsg:
		m.refreshJobStatus()
		return m, tick
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			m.event = 1
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
		case "tab", "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = m.choices[m.cursor]
			}
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "r":
			m.refreshJobStatus()
			return m, nil
		case "R":
			var err error
			job := m.choices[m.cursor]
			if job.Status == "manual" {
				job, _, err = m.client.Jobs.PlayJob(m.pid, job.ID)
			} else {
				job, _, err = m.client.Jobs.RetryJob(m.pid, job.ID)
			}
			if err != nil {
				return m, tea.Quit
			}
			m.choices[m.cursor] = job
		case "o":
			job := m.choices[m.cursor]
			err := utils.OpenBrowser(job.WebURL)
			if err != nil {
				return m, tea.Quit
			}
		case "V":
			m.event = 2
			return m, nil
		}
	}
	return m, nil
}

func (m JobModel) View() (s string) {
	switch m.event {
	case trace:
		s = viewTrace(m)
	case help:
		s = viewHelp()
	case choice:
		s = viewDefault(m)
	}
	return
}

func (m *JobModel) refreshJobStatus() {
	for i, job := range m.choices {
		if internal.IsRunning(job.Status) {
			go func(index, jobID int) {
				job, _, err := m.client.Jobs.GetJob(m.pid, jobID)
				if err != nil {
					log.Println(err)
				}
				m.choices[index] = job
			}(i, job.ID)
		}
	}
}

func (m *JobModel) resetEvent() {
	m.event = 0
}

func eventCallback(msg tea.Msg, m JobModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.resetEvent()
	return m, nil
}

func viewDefault(m JobModel) (s string) {
	s = utils.ColorFg("Pipeline Jobs! (type ? for help)\n\n", internal.MainConfig.ThemeColor)
	for i, job := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = ">>"
		}

		checked := false
		if _, ok := m.selected[i]; ok {
			checked = true
		}

		length := len(job.Name)
		if len(job.Name) > maxLength {
			length = maxLength
		}
		o := fmt.Sprintf("%s %-25s\t%-10s\t%-10s\n", cursor, job.Name[:length], job.Status, job.Stage)
		if checked {
			o = utils.ColorFg(o, internal.MainConfig.ThemeColor)
		}
		s += o
	}
	return
}

func viewTrace(m JobModel) string {
	m.wg.Add(len(m.selected))
	for _, job := range m.selected {
		go func(job *gitlab.Job) {
			_ = internal.DoTrace(m.client, m.pid, job, m.tailLine)
			m.wg.Done()
		}(job)
	}
	return utils.ColorFg("Trace job ...", internal.MainConfig.ThemeColor)
}

func viewHelp() string {
	var help = `
   ? : toggle help
   q : quit
   j : move up
   k : move down
   o : open job page in browser
   r : refresh job status
   R : retry current job or run the manual job
   V : view the chosen job trace
   <tab> or <enter> : select current job
	`

	s := utils.ColorFg("Lab ci quickhelp ~\n", internal.MainConfig.ThemeColor)
	s += help
	return s
}

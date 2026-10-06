package admin

import (
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

type tuiModel struct {
	db                 *gorm.DB
	selected           int
	username, password textinput.Model
	form               bool
	quitting           bool
	status             string
	labels             uiLabels
}

type uiLabels struct {
	menuTitle, reset, quit string
	username, password     string
	failed, succeeded      string
	controls               string
}

func labelsForLocale(locale string) uiLabels {
	if strings.HasPrefix(strings.ToLower(locale), "en") {
		return uiLabels{
			menuTitle: "OctoBackup Admin", reset: "Reset user password", quit: "Quit",
			username: "Username: ", password: "New password: ",
			failed: "Operation failed: ", succeeded: "Password reset successfully. Press Esc to return.",
			controls: "↑/k up  ↓/j down  Enter select  Ctrl+C quit",
		}
	}
	return uiLabels{
		menuTitle: "OctoBackup 管理菜单", reset: "重置用户密码", quit: "退出",
		username: "用户名: ", password: "新密码: ",
		failed: "操作失败: ", succeeded: "密码重置成功，按 Esc 返回菜单",
		controls: "↑/k 上移  ↓/j 下移  Enter 选择  Ctrl+C 退出",
	}
}

func newTUIModel(db *gorm.DB) tuiModel {
	u, p := textinput.New(), textinput.New()
	labels := labelsForLocale(systemLocale())
	u.Prompt, p.Prompt = labels.username, labels.password
	p.EchoMode, p.EchoCharacter = textinput.EchoPassword, '*'
	return tuiModel{db: db, username: u, password: p, labels: labels}
}

func (m tuiModel) Init() tea.Cmd { return nil }
func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	if key.String() == "ctrl+c" {
		m.quitting = true
		return m, tea.Quit
	}
	if !m.form {
		switch key.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < 1 {
				m.selected++
			}
		case "enter":
			if m.selected == 1 {
				m.quitting = true
				return m, tea.Quit
			}
			m.form = true
			m.username.Focus()
		}
		return m, nil
	}
	if key.String() == "esc" {
		m.form = false
		m.username.Blur()
		m.password.Blur()
		return m, nil
	}
	if key.String() == "enter" && m.password.Focused() {
		if err := ResetPassword(m.db, m.username.Value(), m.password.Value()); err != nil {
			m.status = m.labels.failed + err.Error()
		} else {
			m.status = m.labels.succeeded
			m.username.Reset()
			m.password.Reset()
			m.username.Blur()
			m.password.Blur()
		}
		return m, nil
	}
	var cmd tea.Cmd
	if m.username.Focused() {
		m.username, cmd = m.username.Update(msg)
		if key.String() == "enter" {
			m.username.Blur()
			m.password.Focus()
		}
	} else {
		m.password, cmd = m.password.Update(msg)
	}
	return m, cmd
}
func (m tuiModel) View() string {
	if m.quitting {
		return ""
	}
	if !m.form {
		reset, quit := "  "+m.labels.reset, "  "+m.labels.quit
		if m.selected == 0 {
			reset = "> " + m.labels.reset
		} else {
			quit = "> " + m.labels.quit
		}
		return lipgloss.NewStyle().Margin(1, 2).Render(m.labels.menuTitle + "\n\n" + reset + "\n" + quit + "\n\n" + m.labels.controls)
	}
	formTitle, controls := "重置用户密码", "Enter 提交    Esc 返回"
	if m.labels.reset == "Reset user password" {
		formTitle, controls = m.labels.reset, "Enter submit    Esc back"
	}
	return lipgloss.NewStyle().Margin(1, 2).Render(formTitle + "\n\n" + m.username.View() + "\n" + m.password.View() + "\n\n" + m.status + "\n\n" + controls)
}
func runTUI(db *gorm.DB, input io.Reader, output io.Writer) error {
	_, err := tea.NewProgram(newTUIModel(db), tea.WithInput(input), tea.WithOutput(output)).Run()
	return err
}

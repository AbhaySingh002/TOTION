package styles

import "github.com/charmbracelet/lipgloss"

var (
	WelcomeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("219")).
			Margin(2, 0, 0, 2).
			Padding(0, 1)

	TotionLogostyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffd505ff")).
			Bold(true).
			Align(lipgloss.Center).
			Margin(2,0).
			Width(80)

	CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd505ff"))

	DocStyle = lipgloss.NewStyle().Margin(1, 2)

	ListTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("16")).Background(lipgloss.Color("#ffd505ff")).Padding(0, 1).Margin(1, 1)

	DescriptionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#888")).
				Italic(true).
				Align(lipgloss.Center).
				Width(80).
				Margin(0, 0, 1, 0)
)

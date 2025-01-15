package secrets

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/korol8484/gophkeeper/internal/client/service"
)

type Model struct {
	service *service.Client
}

func NewSecretsScreen(service *service.Client) *Model {
	return &Model{service: service}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {

	return nil
}

func (m *Model) View() string {
	return "table"
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model for our TUI: Holds the anatomy of *every* view we want to see
type model struct {
	username     string
	connection   net.Conn
	serverReader io.Reader
	viewport     viewport.Model
	messages     []string
	textarea     textarea.Model
	senderStyle  lipgloss.Style
	err          error
}

// Custom types for our Model to use
type errMsg error
type serverMsg string

// Initalize and return a base model
func CreateModel(c net.Conn, name string) model {

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		username:     name,
		connection:   c,
		serverReader: bufio.NewReader(c),
		textarea:     ta,
		messages:     []string{},
		viewport:     vp,
		senderStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:          nil,
	}
}

// Init: Start commands for model
func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.RecieveFromServer)
}

// Update: Main update function.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return UpdateChannel(msg, m)
}

// UpdateChannel: Updates the ViewChannel
func UpdateChannel(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	// Check incomming message
	switch msg := msg.(type) {
	case serverMsg:
		m.messages = append(m.messages, m.senderStyle.Render(string(msg)))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()

	// if its a keyboard input
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.SendToServer(m.username + ": " + m.textarea.Value())
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd, m.RecieveFromServer)
}

// View: Views selected view
func (m model) View() string {
	return ViewChannel(m)
}

// ViewChannel: Displays Channel TUI
func ViewChannel(m model) string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

// ReceiveFromServer: tea.Cmd to get messages continuesly from the selected server each cycle through the UpdateChannel()
func (m model) RecieveFromServer() tea.Msg {
	inbuf := make([]byte, 1024)
	n, err := m.serverReader.Read(inbuf[:])
	if err != nil {
		log.Println(err)
	}
	return serverMsg(string(inbuf[:n]))
}

// SentToServer: Send a given message to the selected server to be displayed in the current ViewChannel()
func (m model) SendToServer(msg string) {
	_, err := m.connection.Write([]byte(msg + "\n"))
	if err != nil {
		log.Println(err)
	}
}

// Runner for client
// Connects to specified server, Creates a Scanner and Reader, and continuesly scans and reads
func main() {
	// Connection
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	var selectedName string
	fmt.Printf("Please enter your username: ")
	fmt.Scanln(&selectedName)

	clientRunner := tea.NewProgram(CreateModel(conn, selectedName))
	if _, err := clientRunner.Run(); err != nil {
		log.Fatal(err)
	}
}

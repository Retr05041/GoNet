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

// model: Holds the anatomy of every view we want to see in our TUI
type model struct {
	username     string         // Given username
	connection   net.Conn       // Connection interface
	serverReader io.Reader      // Connction used as Reader from server
	viewport     viewport.Model // Viewport model
	messages     []string       // Holds the messages for the viewport
	textarea     textarea.Model // Text area for user input
	senderStyle  lipgloss.Style // Styles
	err          error          // If there is an error
}

type errMsg error     // errorMSg: to wrap error messages for the model to display
type serverMsg string // serverMsg: to wrap server messages for the model to display

// CreateModel: Creates the base model for our TUI
func CreateModel(c net.Conn, name string) model {

	// model.textarea setup
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

	// model.viewport setup
	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	// Return base model
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

// View: Views selected view
func (m model) View() string {
	return ViewChannel(m)
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

// main: Runner for Client
func main() {
	// Connection
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// Username
	var selectedName string
	fmt.Printf("Please enter your username: ")
	fmt.Scanln(&selectedName)

	// Run model
	clientRunner := tea.NewProgram(CreateModel(conn, selectedName))
	if _, err := clientRunner.Run(); err != nil {
		log.Fatal(err)
	}
}

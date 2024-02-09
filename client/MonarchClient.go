package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChannelModel struct {
	// username string
	connection net.Conn
	// serverReader io.Reader
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}
type errMsg error
type serverMsg string

func (m ChannelModel) sendToServer(msg string) {
	// fmt.Println("Sending: " + m)
	_, err := m.connection.Write([]byte(msg + "\n"))
	if err != nil {
		log.Println(err)
	}
	// fmt.Println("Message sent!")
}

// Initalize and return a base model
func initialChannelModel(c net.Conn) ChannelModel {

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

	return ChannelModel{
		connection:  c,
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

// Returns a text area tp type in
func (m ChannelModel) Init() tea.Cmd {
	return textarea.Blink
}

// Updates TUI Model every run
func (m ChannelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	// Check incomming message
	switch msg := msg.(type) {
	// if its a keyboard input
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.sendToServer(m.textarea.Value())
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

// Displays TUI model
func (m ChannelModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
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

	currentModel := initialChannelModel(conn)

	clientRunner := tea.NewProgram(currentModel)
	if _, err := clientRunner.Run(); err != nil {
		log.Fatal(err)
	}
}

// inbuf := make([]byte, 1024)
// for {
// 	n, err := serverReader.Read(inbuf[:])
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
// 	return string(inbuf[:n])

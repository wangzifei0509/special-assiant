package tui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/sashabaranov/go-openai"

	"github.com/chunfengshili/sa/pkg/llm"
)

// Model TUI 模型
type Model struct {
	client   *llm.Client
	messages []openai.ChatCompletionMessage
	input    string
	output   string
	err      string
	loading  bool
	renderer *glamour.TermRenderer
	width    int
	height   int
}

// NewModel 创建新的 TUI 模型
func NewModel(client *llm.Client) *Model {
	r, _ := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(80))
	return &Model{
		client:   client,
		renderer: r,
		messages: make([]openai.ChatCompletionMessage, 0),
	}
}

// Init 初始化
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update 更新
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case aiResponseMsg:
		m.loading = false
		if msg.err != nil {
			m.err = fmt.Sprintf("❌ %v", msg.err)
		} else {
			m.messages = append(m.messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: msg.text,
			})
			rendered, _ := m.renderer.Render(msg.text)
			m.output = rendered
			m.err = ""
		}
	}
	return m, nil
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	case "enter":
		if !m.loading && m.input != "" {
			m.messages = append(m.messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: m.input,
			})
			m.loading = true
			m.output = ""
			m.err = ""
			input := m.input
			m.input = ""
			return m, m.fetchAI(input)
		}
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.input += msg.String()
		}
	}
	return m, nil
}

// View 渲染
func (m *Model) View() string {
	var s strings.Builder

	s.WriteString("🤖 Agent CLI - 输入问题后回车发送 (Ctrl+C 退出)\n")
	s.WriteString("════════════════════════════════════════════════\n\n")

	// 显示最近几条对话
	if len(m.messages) > 0 {
		maxShow := min(4, len(m.messages))
		start := len(m.messages) - maxShow
		s.WriteString("📜 最近对话:\n")
		for i := start; i < len(m.messages); i++ {
			msg := m.messages[i]
			role := "👤 你"
			if msg.Role == openai.ChatMessageRoleAssistant {
				role = "🤖 Agent"
			}
			content := msg.Content
			if len(content) > 100 {
				content = content[:100] + "..."
			}
			s.WriteString(fmt.Sprintf("%s: %s\n", role, content))
		}
		s.WriteString("\n")
	}

	if m.loading {
		s.WriteString("💭 Agent 思考中...\n\n")
	}

	if m.output != "" {
		s.WriteString(m.output)
		s.WriteString("\n")
	}

	if m.err != "" {
		s.WriteString(m.err + "\n\n")
	}

	s.WriteString("────────────────────────────────────────────────\n")
	s.WriteString("👉 ")
	s.WriteString(m.input)
	s.WriteString("▌")

	return s.String()
}

// aiResponseMsg AI 响应消息
type aiResponseMsg struct {
	text string
	err  error
}

func (m *Model) fetchAI(input string) tea.Cmd {
	return func() tea.Msg {
		// 添加用户消息
		messages := append(m.messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})

		stream, err := m.client.CreateChatCompletionStream(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    m.client.Model(),
				Messages: messages,
			},
		)
		if err != nil {
			return aiResponseMsg{err: err}
		}
		defer stream.Close()

		var fullText strings.Builder
		for {
			response, err := stream.Recv()
			if err != nil {
				break
			}
			if len(response.Choices) > 0 {
				fullText.WriteString(response.Choices[0].Delta.Content)
			}
		}
		return aiResponseMsg{text: fullText.String()}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
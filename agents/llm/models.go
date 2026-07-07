package llm

import "encoding/json"

type Message struct {
	Role    string  `json:"role"`
	Content string  `json:"content"`
	Id      *string `json:"tool_call_id"`
	Tools   *Tools
}

type ToolReply struct {
	Content string `json:"content" binding:"required"`
	Id      string `json:"tool_call_id" binding:"required"`
}
type Arg struct {
	Name        string
	Description string
	Type        string
	Required    bool
}

type Messages []Message
type Tools []Tool
type Tool struct {
	Name        string
	Description string
	Args        []Arg
	fn          func(id string, args ToolInput) (any, error)
}
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	Id       string           `json:"id"`
	Type     string           `json:"type"`
	Function ToolCallFunction `json:"function"`
}

type ChoiceMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls"`
}

type Choice struct {
	Index   int           `json:"index"`
	Message ChoiceMessage `json:"message"`
}

type RawLLMResponse struct {
	Id      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Response struct {
	Content    string
	ToolName   string
	ToolCallId string
	Arguments  map[string]any
}

type ToolInput map[string]any

func (r RawLLMResponse) getResponse() (Response, error) {
	var response Response
	if len(r.Choices) == 0 {
		return response, nil
	}
	message := r.Choices[0].Message
	response.Content = message.Content
	if len(message.ToolCalls) == 0 {
		return response, nil
	}
	call := message.ToolCalls[0]
	response.ToolName = call.Function.Name
	response.ToolCallId = call.Id
	if err := json.Unmarshal([]byte(call.Function.Arguments), &response.Arguments); err != nil {
		return response, err
	}
	return response, nil
}

func (t ToolInput) get(name string) (string, bool) {
	rawValue, ok := t[name]
	if !ok {
		return "", false
	}
	value, ok := rawValue.(string)
	if !ok {
		return "", false
	}
	return value, true
}

func (a Arg) property() map[string]any {
	return map[string]any{
		"type":        a.Type,
		"description": a.Description,
	}
}

func (t Tool) Exec(id string, args map[string]any) (any, error) {
	return t.fn(id, args)
}

func (t Tool) format() map[string]any {
	properties := make(map[string]any, len(t.Args))
	required := make([]string, 0, len(t.Args))
	for _, arg := range t.Args {
		properties[arg.Name] = arg.property()
		if arg.Required {
			required = append(required, arg.Name)
		}
	}
	return map[string]any{
		"type": "function",
		"function": map[string]any{
			"name":        t.Name,
			"description": t.Description,
			"parameters": map[string]any{
				"type":       "object",
				"properties": properties,
				"required":   required,
			},
		},
	}
}

func (t Tools) List() []map[string]any {
	list := make([]map[string]any, 0, len(t))
	for _, tool := range t {
		list = append(list, tool.format())
	}
	return list
}

func (t Tools) Exec(id string, name string, args map[string]any) (any, error) {
	for _, existingTool := range t {
		if existingTool.Name != name {
			continue
		}
		return existingTool.Exec(id, args)
	}
	return nil, nil
}

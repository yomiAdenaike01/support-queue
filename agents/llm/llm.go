package llm

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	CUSTOMER_AGENT_PERSONALITY = `You are a customer support agent responding to client messages. Each incoming message includes a sentiment score for that message and an average sentiment score for the conversation so far.
Rules:
- Never make promises to the customer that you cannot guarantee (refunds, timelines, exceptions to policy, etc.).
- Always be professional and helpful, even if the customer is frustrated or rude.
- Never state facts about an order, refund, or policy from memory — use a tool to check first if one exists for that purpose.
- Use at most one tool per response.
- You can hand a conversation off to a specific human by their personId when it's outside your scope — see the escalate_to_human tool for how to pick the right person.
- Never write your reply to the customer as plain response content. Once you've decided what to say, call send_response_to_customer with that message — that is the only way the customer receives it.`
)

func Do(messages Messages, tools Tools) (Response, error) {
	var r Response
	chatMessages := []map[string]any{
		{
			"role":    "system",
			"content": CUSTOMER_AGENT_PERSONALITY,
		},
	}

	for _, message := range messages {
		chatMessages = append(chatMessages, map[string]any{
			"role":    message.Role,
			"content": message.Content,
		})
	}
	body := map[string]any{
		"model":    "llama3.2",
		"tools":    tools.List(),
		"messages": chatMessages,
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		return r, err
	}

	req, err := http.NewRequest("POST", "http://localhost:11434/v1/chat/completions", bytes.NewReader(rawBody))
	if err != nil {
		return r, err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return r, err
	}
	var llmResponse RawLLMResponse
	if err := json.NewDecoder(response.Body).Decode(&llmResponse); err != nil {
		return r, err
	}
	return llmResponse.getResponse()
}

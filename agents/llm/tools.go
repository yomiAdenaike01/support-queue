package llm

import (
	"fmt"
	"log"
)

const TEAM_ROSTER = `team_billing_01 | Billing | online
- team_legal_04   | Legal/Compliance | online
- team_eng_02     | Engineering/Support | offline
- team_general_01 | General | online`

var pendingStore map[string]any

func escalate(personId string, pendingId string, reason string) any {
	message := fmt.Sprintf("personId=%s pendingId=%s reason=%s", personId, pendingId, reason)
	log.Println(message)
	return nil
}

func escalateToHuman(ticketId string) func(id string, args ToolInput) (any, error) {
	return func(id string, args ToolInput) (any, error) {
		log.Printf("ticketId=%s", ticketId)
		reason, ok := args.get("reason")
		if !ok {
			return nil, nil
		}
		personId, ok := args.get("personId")
		if !ok {
			return nil, nil
		}
		return escalate(personId, id, reason), nil
	}

}

func sendResponseToCustomer(ticketId string) func(id string, args ToolInput) (any, error) {
	return func(id string, args ToolInput) (any, error) {
		log.Printf("ticketid=%s", ticketId)
		log.Println(id, args)
		return nil, nil
	}

}

func NewTools(ticketId string) Tools {
	return Tools{
		Tool{
			fn:          sendResponseToCustomer(ticketId),
			Name:        "send_response_to_customer",
			Description: "Send your reply directly to the customer. Use this to deliver the message you've composed once you've decided what to say — do not put your reply in plain text without calling this tool.",
			Args: []Arg{
				{
					Name:        "message",
					Type:        "string",
					Description: "The message you want to send to the customer",
					Required:    true,
				},
			},
		},
		Tool{
			fn:          escalateToHuman(ticketId),
			Name:        "escalate_to_human",
			Description: "Hand the conversation off to a human agent immediately. Use this for angry/abusive customers, legal or safety-related complaints, or anything explicitly outside your scope.",
			Args: []Arg{
				{
					Name:        "reason",
					Type:        "string",
					Description: "Why this conversation needs a human, in one short sentence",
					Required:    true,
				},
				{
					Name:        "personId",
					Type:        "string",
					Description: fmt.Sprintf("The ID of the team best suited to handle this, chosen from the current roster based on department/role match to the issue:\n%s\nMatch billing/refund disputes to a Billing role, legal/safety/abuse complaints to Legal/Compliance, product or technical issues to Engineering/Support. If nothing matches clearly, pick the most senior general role. Only use IDs from this list.", TEAM_ROSTER),
					Required:    true,
				},
			},
		},
	}
}

package integrationsdomain

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	inputsourcesdomain "github.com/yomiAdenaike01/support-queue/internals/domains/inputsources"
	"github.com/yomiAdenaike01/support-queue/internals/domains/integrations/oauth"
)

const gmailAPIBase = "https://gmail.googleapis.com/gmail/v1/users/me"

// Excludes Gmail's own promotions/social/forums categories and spam so we never
// even fetch the full content of mail that isn't a real customer request.
const ticketWorthyQuery = "-category:promotions -category:social -category:forums -in:spam"

var errNotTicketWorthy = errors.New("message is not ticket worthy")

type GmailClient struct {
	inputSourceRepository inputsourcesdomain.Repository
	out                   chan<- TicketPayload
	providers             *oauth.Providers
}

type GmailMessage struct {
	Id      string
	From    string
	Subject string
	Body    string
}

func (g GmailMessage) GetSubject() string {
	return g.Subject
}

func (g GmailMessage) Message() string {
	return g.Body
}

func (g GmailMessage) Email() string {
	return g.From
}

type gmailHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type gmailMessagePart struct {
	MimeType string        `json:"mimeType"`
	Headers  []gmailHeader `json:"headers"`
	Body     struct {
		Data string `json:"data"`
	} `json:"body"`
	Parts []gmailMessagePart `json:"parts"`
}

type gmailListResponse struct {
	Messages []struct {
		Id string `json:"id"`
	} `json:"messages"`
}

type gmailMessageResponse struct {
	Id      string           `json:"id"`
	Snippet string           `json:"snippet"`
	Payload gmailMessagePart `json:"payload"`
}

type TicketPayload interface {
	Email() string
	GetSubject() string
	Message() string
}

func gmailGet(ctx context.Context, accessToken string, path string, query url.Values) ([]byte, error) {
	reqURL := gmailAPIBase + path
	if len(query) > 0 {
		reqURL = fmt.Sprintf("%s?%s", reqURL, query.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gmail api request to %s failed with status %d: %s", path, response.StatusCode, string(body))
	}
	return body, nil
}

func listLatestMessageIds(ctx context.Context, accessToken string, maxResults int) ([]string, error) {
	body, err := gmailGet(ctx, accessToken, "/messages", url.Values{
		"maxResults": []string{fmt.Sprintf("%d", maxResults)},
		"labelIds":   []string{"INBOX"},
		"q":          []string{ticketWorthyQuery},
	})
	if err != nil {
		return nil, err
	}
	var listResponse gmailListResponse
	if err := json.Unmarshal(body, &listResponse); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(listResponse.Messages))
	for _, message := range listResponse.Messages {
		ids = append(ids, message.Id)
	}
	return ids, nil
}

func decodeBase64URL(data string) string {
	if data == "" {
		return ""
	}
	decoded, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		log.Printf("Failed to decode gmail message body error=%s", err.Error())
		return ""
	}
	return string(decoded)
}

func extractBody(part gmailMessagePart) string {
	if part.MimeType == "text/plain" && part.Body.Data != "" {
		return decodeBase64URL(part.Body.Data)
	}
	for _, child := range part.Parts {
		if body := extractBody(child); body != "" {
			return body
		}
	}
	return ""
}

func getHeader(headers []gmailHeader, name string) string {
	for _, header := range headers {
		if strings.EqualFold(header.Name, name) {
			return header.Value
		}
	}
	return ""
}

func getFrom(headers []gmailHeader) string {
	from := getHeader(headers, "From")
	for _, key := range []string{"noreply", "donotreply", "no-reply"} {
		if strings.Contains(from, key) {
			return ""
		}
	}
	return from
}

func getMessage(ctx context.Context, accessToken string, id string) (GmailMessage, error) {
	body, err := gmailGet(ctx, accessToken, "/messages/"+id, url.Values{"format": []string{"full"}})
	if err != nil {
		return GmailMessage{}, err
	}
	var message gmailMessageResponse
	if err := json.Unmarshal(body, &message); err != nil {
		return GmailMessage{}, err
	}
	// Bulk/marketing senders that dodge Gmail's own category tagging almost always
	if getHeader(message.Payload.Headers, "List-Unsubscribe") != "" {
		return GmailMessage{}, errNotTicketWorthy
	}
	messageBody := extractBody(message.Payload)
	from := getFrom(message.Payload.Headers)
	if len(from) == 0 {
		return GmailMessage{}, errNotTicketWorthy
	}
	return GmailMessage{
		Id:      message.Id,
		From:    from,
		Subject: getHeader(message.Payload.Headers, "Subject"),
		Body:    messageBody,
	}, nil
}

func fetchLatestMessages(ctx context.Context, accessToken string, maxResults int) []GmailMessage {
	ids, err := listLatestMessageIds(ctx, accessToken, maxResults)
	if err != nil {
		log.Printf("Failed to list gmail messages error=%s", err.Error())
		return nil
	}
	messages := make([]GmailMessage, 0, len(ids))
	for _, id := range ids {
		message, err := getMessage(ctx, accessToken, id)
		if errors.Is(err, errNotTicketWorthy) {
			continue
		}
		if err != nil {
			log.Printf("Failed to fetch gmail message id=%s error=%s", id, err.Error())
			continue
		}
		messages = append(messages, message)
	}
	return messages
}

func (g *GmailClient) hasGmailSource(ctx context.Context) bool {
	srcType := inputsourcesdomain.InputSourceTypeEmail
	inputSources, err := g.inputSourceRepository.FindMany(ctx, inputsourcesdomain.FindInput{
		SourceType: &srcType,
	})
	if err != nil {
		log.Printf("Failed to find email input sources error=%s", err.Error())
		return false
	}
	for _, inputSource := range inputSources {
		if strings.Contains(inputSource.ConnectionValue, "gmail") {
			return true
		}
	}
	return false
}

func (g *GmailClient) fetchAndForward(ctx context.Context) {
	if !g.hasGmailSource(ctx) {
		return
	}

	accessToken, err := g.providers.GetAccessToken(oauth.GOOGLE)
	if err != nil {
		log.Printf("Failed to fetch Google OAuth tokens error=%s", err.Error())
		return
	}
	if accessToken == "" {
		log.Println("No valid Google access token available, skipping gmail poll")
		return
	}

	for _, message := range fetchLatestMessages(ctx, accessToken, 10) {
		g.out <- message
	}
}

func (g *GmailClient) poll(ctx context.Context) {
	g.fetchAndForward(ctx)
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			g.fetchAndForward(ctx)
		}
	}
}

type emailIntegrationDeps struct {
	repository     inputsourcesdomain.Repository
	onNewEmail     chan<- TicketPayload
	oauthProviders *oauth.Providers
}

func newGmailClient(ctx context.Context, deps emailIntegrationDeps) *GmailClient {
	client := &GmailClient{
		inputSourceRepository: deps.repository,
		out:                   deps.onNewEmail,
		providers:             deps.oauthProviders,
	}
	go client.poll(ctx)
	return client
}

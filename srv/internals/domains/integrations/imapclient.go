package integrationsdomain

import (
	"bytes"
	"context"
	"io"
	"log"
	"mime"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/charset"
	inputsourcesdomain "github.com/yomiAdenaike01/support-queue/internals/domains/inputsources"
)

type ImapClient struct {
	in                    chan *imapclient.Client
	inputSourceRepository inputsourcesdomain.Repository
	out                   chan any
}

func authenticateWithIMAP(inputSource inputsourcesdomain.InputSource) (*imapclient.Client, error) {
	credentials, err := inputSource.GetCredentials()
	if err != nil {
		log.Printf("Failed to get credentials error=%s", err.Error())
		return nil, err
	}

	options := &imapclient.Options{
		WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader},
	}
	imapAddress := ""
	if strings.Contains(inputSource.ConnectionValue, "gmail") {
		imapAddress = "imap.gmail.com:993"
	}
	client, err := imapclient.DialTLS(imapAddress, options)
	cmd := client.Login(credentials.Principal, credentials.Password)

	if err = cmd.Wait(); err != nil {
		log.Printf("Failed to authenticate with creds username=%s password=%s err=%s", credentials.Principal, credentials.Password, err.Error())
		return nil, err
	}
	return client, nil

}

func fetchNewMessages(client *imapclient.Client) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		data, err := client.Select("INBOX", nil).Wait()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Printf("NUM_MESSAGE= %d", data.NumMessages)
		bodySection := &imap.FetchItemBodySection{}
		messages, err := client.Fetch(imap.SeqSetNum(data.NumMessages), &imap.FetchOptions{
			UID:         true,
			Envelope:    true,
			Flags:       true,
			BodySection: []*imap.FetchItemBodySection{bodySection},
		}).Collect()

		if err != nil {
			log.Println("Failed to fetch mailbox messages")
			continue
		}
		for _, message := range messages {
			sub := message.Envelope.Subject
			body := message.FindBodySection(&imap.FetchItemBodySection{})
			reader := bytes.NewReader(body)
			text, err := io.ReadAll(reader)
			if err != nil {
				panic(err)
			}
			log.Printf("subject=%s", text, sub)

		}
	}
}

func (i *ImapClient) poll(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			break
		case client, ok := <-i.in:
			if ok {
				go fetchNewMessages(client)
			}
		}
	}
}

func (i *ImapClient) authenticate(ctx context.Context) {
	log.Println("Starting polling....")
	srcType := inputsourcesdomain.InputSourceTypeEmail

	inputSources, err := i.inputSourceRepository.FindMany(ctx, inputsourcesdomain.FindInput{
		SourceType: &srcType,
	})

	if err != nil {
		panic(err)
	}

	for _, inputSource := range inputSources {
		_, err := inputSource.GetCredentials()

		if err != nil {
			continue
		}

		imapClient, err := authenticateWithIMAP(inputSource)

		if err != nil {
			continue
		}
		i.in <- imapClient

	}

}

func newImapClient(ctx context.Context, repository inputsourcesdomain.Repository, onNewEmail chan any) *ImapClient {
	client := &ImapClient{
		inputSourceRepository: repository,
		in:                    make(chan *imapclient.Client),
		out:                   onNewEmail,
	}
	go client.poll(ctx)
	go client.authenticate(ctx)
	return client
}

package gmail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kravitz/contra_mailer/mailer"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Mailer implemets Mailer interface for GMAIL api
type Mailer struct {
	service *gmail.Service
}

func getToken(credentialsFilename string) (*oauth2.Token, error) {
	f, err := os.Open(credentialsFilename)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func encodeWeb64String(b []byte) string {

	s := base64.URLEncoding.EncodeToString(b)

	var i = len(s) - 1
	for s[i] == '=' {
		i--
	}

	return s[0 : i+1]
}

// CreateMailer is a constructor for gmail.Mailer
func CreateMailer(clientSecretFilename string, credentialsFilename string) (m *Mailer, err error) {
	// ctx := context.Background()
	//
	// b, err := ioutil.ReadFile("client_secret.json")
	// if err != nil {
	//   log.Fatalf("Unable to read client secret file: %v", err)
	// }
	//
	// // If modifying these scopes, delete your previously saved credentials
	// // at ~/.credentials/gmail-go-quickstart.json
	// config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	// if err != nil {
	//   log.Fatalf("Unable to parse client secret file to config: %v", err)
	// }
	// client := getClient(ctx, config)
	//
	// srv, err := gmail.New(client)
	// if err != nil {
	//   log.Fatalf("Unable to retrieve gmail Client %v", err)
	// }
	//
	// // user := "me"
	// from := mail.Address{Name: "CONTRA service", Address: "cubadvfu@gmail.com"}
	// to := mail.Address{Name: "idkravitz", Address: "idkravitz@gmail.com"}
	b, err := ioutil.ReadFile(clientSecretFilename)
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		return nil, err
	}
	token, err := getToken(credentialsFilename)
	client := config.Client(context.Background(), token)
	srv, err := gmail.New(client)
	if err != nil {
		return nil, err
	}
	m = &Mailer{
		service: srv,
	}

	return m, err
}

// Send message via Gmail API
func (m *Mailer) Send(msg *mailer.Message) (err error) {
	header := map[string]string{}
	header["From"] = msg.From.String()
	header["To"] = msg.To.String()
	header["Subject"] = msg.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	var rawMsg string
	for k, v := range header {
		rawMsg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	rawMsg += "\r\n" + msg.Body

	gmsg := &gmail.Message{Raw: encodeWeb64String([]byte(rawMsg))}
	_, err = m.service.Users.Messages.Send("me", gmsg).Do()
	if err != nil {
		return err
	}

	return nil
}

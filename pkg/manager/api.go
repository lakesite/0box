// api.go
// This file defines all the web service handlers, and the web service API
// endpoints.

package manager

import (
	"bufio"
	// "crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/emersion/go-mbox"
	"github.com/gorilla/mux"

	"github.com/lakesite/ls-config/pkg/config"
	"github.com/lakesite/ls-fibre/pkg/service"
)

var api_key string

// struct for holding message header and body
type Message struct {
	Header mail.Header
	Body   string
}

// Handle requests to post mail
func (ms *ManagerService) PostMailHandler(w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from") // &mail.Address{r.FormValue("from_name"), r.FormValue("from")}
	to := r.FormValue("to")     // &mail.Address{r.FormValue("to_name"), r.FormValue("to")}
	subject := r.FormValue("subject")
	body := r.FormValue("body")
	// template (last)

	// handle cleaning mail addresses up;
	replacer := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	if from /*.String()*/ == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Missing form data: from")
		return
	}

	if to /*.String()*/ == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Missing form data: to")
		return
	}

	if subject == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Missing form data: subject")
		return
	}

	if body == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Missing form data: body")
		return
	}

	fmt.Printf("from [%s], to [%s], subject [%s], body: %s\n", from, to, subject, body)

	headers := make(map[string]string)
	headers["From"] = from // .String()
	headers["To"] = to     // .String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Handle plain jane susie mary localhost authenticary
	c, err := smtp.Dial("127.0.0.1:25")
	if err != nil {
		// 500
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer c.Close()

	// from
	if err = c.Mail(replacer.Replace(from /*.String() */)); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	// handle range of recipients
	// for i := range to {
	if err = c.Rcpt(replacer.Replace(to /*.String() */)); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	// }

	// Data
	cw, err := c.Data()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	_, err = cw.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	err = cw.Close()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	c.Quit()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Message sent.")
}

// Handle requests to return all a user's mail
func (ms *ManagerService) GetMailUserHandler(w http.ResponseWriter, r *http.Request) {
	var messages []*Message

	vars := mux.Vars(r)
	root := ms.GetSectionPropertyOrDefault("0box", "mboxroot", "/var/mail")

	// get the username
	mailbox := root + "/" + vars["user"]
	// if root + /username does not exist return json no user
	if _, err := os.Stat(mailbox); os.IsNotExist(err) {
		// no user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("No such mailbox user.")
		return
	} else {
		// read;
		box, _ := os.Open(mailbox)
		mr := mbox.NewReader(box)
		for {
			nm, err := mr.NextMessage()

			if err == io.EOF {
				break
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)
				return
			}

			msg, err := mail.ReadMessage(nm)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)
				return
			}
			body, _ := ioutil.ReadAll(msg.Body)
			messages = append(messages, &Message{Header: msg.Header, Body: string(body)})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

// Handle requests to return a user's message # from mail.
func (ms *ManagerService) GetMailUserNumberHandler(w http.ResponseWriter, r *http.Request) {
	var message *Message

	vars := mux.Vars(r)
	root := ms.GetSectionPropertyOrDefault("0box", "mboxroot", "/var/mail")

	mailbox := root + "/" + vars["user"]
	if _, err := os.Stat(mailbox); os.IsNotExist(err) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("No such mailbox user.")
	} else {
		box, _ := os.Open(mailbox)
		mr := mbox.NewReader(box)
		n, _ := vars["number"]
		number, _ := strconv.Atoi(n)
		counter := 1
		for {
			nm, err := mr.NextMessage()
			if err == io.EOF {
				break
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)
				return
			}
			if counter == number {
				msg, err := mail.ReadMessage(nm)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(err)
					return
				}
				body, _ := ioutil.ReadAll(msg.Body)
				message = &Message{Header: msg.Header, Body: string(body)}
			}
			if counter > number {
				break
			}
			counter = counter + 1
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if message == nil {
		json.NewEncoder(w).Encode("No such message number.")
	} else {
		json.NewEncoder(w).Encode(message)
	}
}

// Handle requests to delete a user's message # from mail.
// Todo: This method truncates and re-writes from structures in memory,
// this is not efficient and could cause problems.
func (ms *ManagerService) DeleteMailUserNumberHandler(w http.ResponseWriter, r *http.Request) {
	var messages []*Message
	status := "No such message to delete."

	vars := mux.Vars(r)
	root := ms.GetSectionPropertyOrDefault("0box", "mboxroot", "/var/mail")

	mailbox := root + "/" + vars["user"]
	if _, err := os.Stat(mailbox); os.IsNotExist(err) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("No such mailbox user.")
	} else {
		box, _ := os.OpenFile(mailbox, os.O_RDWR, 0600)
		defer box.Close()
		mr := mbox.NewReader(box)
		n, _ := vars["number"]
		number, _ := strconv.Atoi(n)
		counter := 1
		for {
			nm, err := mr.NextMessage()
			if err == io.EOF {
				break
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)
				return
			}
			if counter != number {
				msg, err := mail.ReadMessage(nm)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(err)
					return
				}
				body, _ := ioutil.ReadAll(msg.Body)
				messages = append(messages, &Message{Header: msg.Header, Body: string(body)})
			} else {
				status = "Message deleted."
			}

			counter = counter + 1
		}

		// truncate box, panic on errors.
		box.Truncate(0)
		box.Seek(0, 0)
		bw := bufio.NewWriter(box)

		// write it back out, header first (From: )
		for _, message := range messages {
			// header first
			_, err = fmt.Fprintf(bw, "From %s  %s\n", message.Header.Get("From"), message.Header.Get("Date"))
			ms.PanicCheck(err)
			for k, v := range message.Header {
				_, err = fmt.Fprintf(bw, "%s: %s\n", k, v[0])
				ms.PanicCheck(err)
			}
			_, err = fmt.Fprintf(bw, "\n%s\n", message.Body)
			ms.PanicCheck(err)
			bw.Flush()
			// fmt.Printf("Got message[%v]: %s\n%s\n", i, message.Header, message.Body)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// Handle requests to delete a user mailbox
func (ms *ManagerService) DeleteMailUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	root := ms.GetSectionPropertyOrDefault("0box", "mboxroot", "/var/mail")

	mailbox := root + "/" + vars["user"]
	if _, err := os.Stat(mailbox); os.IsNotExist(err) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("No such mailbox.")
	} else {
		box, _ := os.OpenFile(mailbox, os.O_RDWR, 0600)
		defer box.Close()

		// truncate box, panic on errors.
		box.Truncate(0)
		box.Seek(0, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted.")
}

// Handle requests to list available mailbox users:
func (ms *ManagerService) GetMailboxesHandler(w http.ResponseWriter, r *http.Request) {
	var files []string
	var boxes []string

	root := ms.GetSectionPropertyOrDefault("0box", "mboxroot", "/var/mail")

	files, err := ms.GetFilesInPath(root)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		boxes = append(boxes, filepath.Base(file))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(boxes); err != nil {
		panic(err)
	}
}

// setupRoutes defines and associates routes to handlers.
func (ms *ManagerService) setupRoutes(ws *service.WebService) {
	// [1, 2]: POST /api/0box/v1/mail/
	ws.Router.HandleFunc("/api/0box/v1/mail/", ms.PostMailHandler).Methods("POST")

	// [3]: GET /api/0box/v1/mail/ - return list of boxes
	ws.Router.HandleFunc("/api/0box/v1/mail/", ms.GetMailboxesHandler).Methods("GET")

	// [3]: GET /api/0box/v1/mail/<username from /var/spool/mail>[/mail #]
	ws.Router.HandleFunc("/api/0box/v1/mail/{user}", ms.GetMailUserHandler).Methods("GET")
	ws.Router.HandleFunc("/api/0box/v1/mail/{user}/{number}", ms.GetMailUserNumberHandler).Methods("GET")

	// [4]: DELETE /api/0box/v1/mail/<username>/<mail #>
	ws.Router.HandleFunc("/api/0box/v1/mail/{user}", ms.DeleteMailUserHandler).Methods("DELETE")
	ws.Router.HandleFunc("/api/0box/v1/mail/{user}/{number}", ms.DeleteMailUserNumberHandler).Methods("DELETE")

	// [6] Secure the API with an API key for all operations.
	ws.Apikey = config.Getenv("0BOX_API_KEY", ms.GetSectionPropertyOrDefault("0box", "apikey", ""))
	ws.Router.Use(ws.LogMiddleware)
	ws.Router.Use(ws.APIKeyMiddleware)
}

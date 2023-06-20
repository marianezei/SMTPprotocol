package main

import (
	"fmt"
	"net"
	"net/textproto"
	"time"
)

func main() {
	// Configurações do servidor SMTP
	addr := fmt.Sprintf(":%d", 8080)

	// Define o endereço de email do remetente
	from := "seu_email@example.com"

	// Define o(s) endereço(s) de email do(s) destinatário(s)
	to := []string{"destinatario1@example.com"}

	to2 := []string{"destinatario2@example.com"}

	to3 := []string{"destinatario3@example.com"}

	// Define o conteúdo do email
	subject := "Assunto do email"
	body := "Corpo do email"

	// Monta o cabeçalho do email
	msg := "From: " + from + "\r\n" +
		"To: " + to[0]
	for i := 1; i < len(to); i++ {
		msg += "," + to[i]
	}
	msg += "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n"

	
	for i := 0; i < 5; i++ {
		
		go sendMail(addr, from, to, msg)
		go sendMail(addr, from, to2, msg)
		go sendMail(addr, from, to3, msg)
		
		time.Sleep(2 * time.Second)
	}

	time.Sleep(2 * time.Second)
	fmt.Println("Todos os emails foram enviados!")
	
}

func sendMail(addr string, from string, to []string, msg string) error {
	// validação de emails se estão no DB

	client, err := dial(addr)
	if err != nil {
		return err
	}
	
	defer client.conn.Close()

	// Envia o email
	err = client.mail(from)
	
	if err != nil {
		fmt.Println("Erro ao definir o remetente:", err)
		return err
	}

	for _, addr := range to {
		err = client.rcpt(addr)
		if err != nil {
			fmt.Println("Erro ao definir o destinatário:", err)
			return err
		} 
	}

	err1 := client.data()
	if err1 != nil {
		fmt.Println("Erro ao iniciar a transferência de dados:", err1)
		return err1
	}


	fmt.Println("Email enviado com sucesso!")

	err = client.quit()
	if err != nil {
		fmt.Println("Erro ao sair do servidor SMTP:", err)
		return err
	}

	return err
}

// A Client represents a client connection to an SMTP server.
type Client struct {
	// Text is the textproto.Conn used by the Client. It is exported to allow for
	// clients to add extensions.
	Text *textproto.Conn
	// keep a reference to the connection so it can be used to create a TLS
	// connection later
	conn net.Conn

	serverName string
	// map of supported extensions
	ext map[string]string
	
	localName  string // the name to use in HELO/EHLO
	didHello   bool   // whether we've said HELO/EHLO
	helloError error  // the error from the hello
}

// NewClient returns a new Client using an existing connection and host as a
// server name to be used when authenticating.
func newClient(conn net.Conn, host string) (*Client, error) {
	text := textproto.NewConn(conn)
	_, _, err := text.ReadResponse(220)

	if err != nil {
		text.Close()
		return nil, err
	}
	c := &Client{Text: text, conn: conn, serverName: host, localName: "localhost"}
	return c, nil
}

// cmd is a convenience function that sends a command and returns the response
func (c *Client) cmd(expectCode int, format string) (int, string, error) {
	id, err := c.Text.Cmd(format)
	if err != nil {
		return 0, "", err
	}
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
	code, msg, err := c.Text.ReadResponse(expectCode)
	fmt.Println("Servidor: ", msg)
	return code, msg, err
}

func dial(addr string) (*Client, error) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	
	return newClient(conn, host)
}

func (c *Client) helo() error {
	_, _, err := c.cmd(250, "EHLO")

	return err
}

func (c *Client) mail(from string) error {

	if err := c.helo(); err != nil {
		return err
	}
	
	_, _, err := c.cmd(250, "MAIL")

	return err
}

func (c *Client) rcpt(to string) error {
	_, _, err := c.cmd(250, "RCPT " + to)

	return err
}

func (c *Client) data() error {
	_, _, err := c.cmd(354, "DATA")

	if err != nil {
		return err
	}
	return err
}

func (c *Client) quit() error {

	_, _, err := c.cmd(221, "QUIT")
	if err != nil {
		return err
	}
	return c.Text.Close()
}
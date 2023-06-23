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
	to1 := []string{"destinatario1@example.com"}

	to2 := []string{"destinatario2@example.com"}

	to3 := []string{"destinatario3@example.com"}

	// Define o conteúdo do email
	subject := "Assunto do email"
	body1 := "Corpo do email 1"
	body2 := "Corpo do email 2"
	body3 := "Corpo do email 3"

	// Monta o cabeçalho do email
	msg1 := "From: " + from + "\r\n" +
		"To: " + to1[0]
	for i := 1; i < len(to1); i++ {
		msg1 += "," + to1[i]
	}
	msg1+= "subject: " + subject +"\r\n" + "body: " + body1
	msg2 := "From: " + from + "\r\n" +
		"To: " + to2[0]
	for i := 1; i < len(to2); i++ {
		msg2 += "," + to2[i]
	}
	msg2+= "subject: " + subject +"\r\n" + "body: " + body2
	msg3 := "From: " + from + "\r\n" +
		"To: " + to2[0]
	for i := 1; i < len(to3); i++ {
		msg3 += "," + to3[i]
	}
	msg3+= "subject: " + subject +"\r\n" + "body: " + body3

	
	for i := 0; i < 5; i++ {
		
		go sendMail(addr, from, to1, msg1)
		go sendMail(addr, from, to2, msg2)
		go sendMail(addr, from, to3, msg3)
		
		time.Sleep(2 * time.Second)
	}

	time.Sleep(2 * time.Second)
	fmt.Println("Todos os emails foram enviados!")
	
}

func sendMail(addr string, from string, to []string, msg string) error {

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

	err1 := client.data(msg)
	if err1 != nil {
		fmt.Println("Erro ao iniciar a transferência de dados:", err1)
		return err1
	}


	fmt.Println("Email enviado com sucesso!")

	err = client.quit(from)
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

func (c *Client) helo(from string) error {
	_, _, err := c.cmd(250, "EHLO")
	fmt.Println("EHLO -> from: " + from)

	return err
}

func (c *Client) mail(from string) error {

	if err := c.helo(from); err != nil {
		return err
	}
	
	_, _, err := c.cmd(250, "MAIL")
	fmt.Println("MAIL -> from: " + from)
	
	return err
}

func (c *Client) rcpt(to string) error {
	_, _, err := c.cmd(250, "RCPT")

	fmt.Println("RCPT -> to " + to)

	return err
}

func (c *Client) data(msg string) error {
	_, _, err := c.cmd(354, "DATA")
	fmt.Println("DATA -> to " + msg)

	if err != nil {
		return err
	}
	return err
}

func (c *Client) quit(from string) error {

	_, _, err := c.cmd(221, "QUIT")
	fmt.Println("QUIT -> from " + from)

	if err != nil {
		return err
	}
	return c.Text.Close()
}
package main

import (
		"fmt"
		"time"
    "log"
    "net/smtp"
)

func sendEmail(emails []string, msg []byte) {
	smtpHost := "smtp.gmail.com"
	username := "projetoconcorrente@gmail.com"
	password := "kesbrjebqslmkxos"

  // Set up authentication information.
  auth := smtp.PlainAuth("", username, password, smtpHost)


	err := smtp.SendMail(smtpHost+":587", auth, username, emails, msg)
    if err != nil {
        log.Fatal(err)
    }
}

func calma() {
	for {
			time.Sleep(time.Second * 1)
			fmt.Printf("calma ezequias\n")
	}
}


func main() {

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	mails := []string{"h.parcelly@gmail.com"}
	msg := []byte("Este é o PRIMEIRO teste de e-mail enviado usando Go e o protocolo SMTP.")

	mails1 := []string{"ezequias.bernardo@ccc.ufcg.edu.br", "mariane.zeitouni@ccc.ufcg.edu.br", "huggo.silva@ccc.ufcg.edu.br", "h.parcelly@gmail.com"}
	msg1 := []byte("Este é o SEGUNDO email de teste usando Go e o protocolo SMTP.")

	go sendEmail(mails, msg)
	go calma()
	sendEmail(mails1, msg1)
	

	fmt.Println("E-mails enviados com sucesso!")
}

package main

import (
		"fmt"
    "log"
    "net/smtp"
)

func main() {

	smtpHost := "smtp.gmail.com"
	username := "projetoconcorrente@gmail.com"
	password := "kesbrjebqslmkxos"

    // Set up authentication information.
    auth := smtp.PlainAuth("", username, password, smtpHost)

    // Connect to the server, authenticate, set the sender and recipient,
    // and send the email all in one step.
    err := smtp.SendMail(smtpHost+":587", auth, username, []string{"h.parcelly@gmail.com"}, []byte("Este é um teste de e-mail enviado usando Go e o protocolo SMTP."))
    if err != nil {
        log.Fatal(err)
    }

		fmt.Println("E-mail enviado com sucesso!")
}
package main

import (
    "log"
    "net"
    "fmt"
    "bufio"
    "strings"
)

func server() {
    // Configuração do servidor SMTP
	smtpPort := 8080

	// Inicia o servidor na porta especificada
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", smtpPort))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Servidor SMTP iniciado na porta %d", smtpPort)

	// Aceita conexões de clientes SMTP
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Erro ao aceitar conexão:", err)
			continue
		}

		// Lida com a conexão em uma goroutine separada
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Cria um scanner para ler as linhas enviadas pelo cliente
	scanner := bufio.NewScanner(conn)

	// Envia a saudação inicial
	conn.Write([]byte("220 Servidor SMTP de exemplo\r\n"))

	// Processa as linhas enviadas pelo cliente
	for scanner.Scan() {
		line := scanner.Text()
		log.Println("Cliente:", line)

		// Verifica o comando enviado pelo cliente
		cmdParts := strings.SplitN(line, " ", 2)
		cmd := strings.ToUpper(cmdParts[0])

		switch cmd {
		case "EHLO":
			conn.Write([]byte("250-Hello\r\n"))
			conn.Write([]byte("250 Servidor SMTP OK\r\n"))
		case "MAIL":
			conn.Write([]byte("250 MAIL OK\r\n"))
		case "RCPT":
			conn.Write([]byte("250 RCPT OK\r\n"))
		case "DATA":
			conn.Write([]byte("354 DATA OK\r\n"))
		case "QUIT":
			conn.Write([]byte("221 Bye\r\n"))
			return
		default:
			conn.Write([]byte("500 Comando não suportado!\r\n"))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Erro ao ler do cliente:", err)
	}
}

func main() {
    server()
   
}

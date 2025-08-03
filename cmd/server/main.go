// WazMeow API
//
//	@title			WazMeow API
//	@version		1.0.0
//	@description	API para gerenciamento de sessões WhatsApp usando whatsmeow library. Permite criar, conectar e gerenciar múltiplas sessões WhatsApp, enviar mensagens e realizar operações de automação.
//	@termsOfService	https://github.com/wazmeow/wazmeow/blob/main/LICENSE
//
//	@contact.name	WazMeow API Support
//	@contact.url	https://github.com/wazmeow/wazmeow
//	@contact.email	support@wazmeow.com
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:8080
//	@BasePath	/
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-API-Key
//	@description				API Key para autenticação. Configure AUTH_ENABLED=true no .env para habilitar.
//
//	@securityDefinitions.basic	BasicAuth
//	@description				Autenticação básica HTTP. Configure AUTH_TYPE=basic no .env para habilitar.
//
//	@schemes	http https
//	@produce	json
//	@accept		json
//
//	@tag.name			Sessions
//	@tag.description	Operações de gerenciamento de sessões WhatsApp
//
//	@tag.name			Health
//	@tag.description	Endpoints de monitoramento e saúde da aplicação
package main

import (
	"log"

	"wazmeow/internal/app"
)

func main() {

	// Initialize and start the application
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start the application (this handles graceful shutdown internally)
	if err := application.Start(); err != nil {
		log.Printf("Application stopped: %v", err)
	}

	// Stop the application (cleanup)
	if err := application.Stop(); err != nil {
		log.Printf("Error stopping application: %v", err)
	}

	log.Println("Application stopped gracefully")
}

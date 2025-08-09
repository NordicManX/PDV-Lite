// pdv-lite/handlers/auth_handler.go
package handlers

import (
	"net/http"

	"github.com/nordicmanx/pdv-lite/models"

	"github.com/gin-gonic/gin"
)

// LoginPayload define a estrutura que esperamos receber no corpo do login.
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler processa a tentativa de login.
func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload LoginPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}

		// --- LÓGICA DE LOGIN SIMPLIFICADA ---
		// No futuro, aqui você buscaria o usuário no banco de dados
		// e compararia a senha com um hash seguro (bcrypt).
		// Por enquanto, vamos usar um usuário "fixo" para teste.

		var user models.Usuario // Reutilizando a struct que já deve existir em models
		// Supondo que você tenha um usuário de teste
		// db.QueryRow("SELECT ... FROM usuarios WHERE email = ?", payload.Email).Scan(...)

		if payload.Email == "carlos@cafe.com" && payload.Password == "123456" {
			// Login bem-sucedido
			user.ID = 1
			user.NomeUsuario = "Carlos do Café"
			user.Email = payload.Email

			// Em um app real, você geraria um JWT (JSON Web Token) aqui.
			// Por enquanto, vamos retornar uma mensagem de sucesso e os dados do usuário.
			c.JSON(http.StatusOK, gin.H{
				"message": "Login bem-sucedido!",
				"user":    user,
				"token":   "fake-jwt-token-for-testing", // Token falso
			})
		} else {
			// Login falhou
			c.JSON(http.StatusUnauthorized, gin.H{"error": "E-mail ou senha inválidos."})
		}
	}
}

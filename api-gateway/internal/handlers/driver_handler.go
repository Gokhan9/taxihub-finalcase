package handlers

import (
	"bitaksi-finalcase/api-gateway/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// ProxyHandler tek bir driverServiceURL ile tüm gelen driver isteklerini ilgili servise yönlendirir.
// Logic, utils paketindeki genel ForwardRawRequest fonksiyonuna devredilmiştir.
type ProxyHandler struct {
	handler fiber.Handler
}

// Hedef URL ile birlikte utils'den genel bir proxy handler oluşturur.
func NewProxyHandler(targetRaw string) *ProxyHandler {
	return &ProxyHandler{
		handler: utils.ForwardRawRequest(targetRaw),
	}
}

/*
Handle: gelen isteği oluşturulan proxy handler'a iletir.
*/
func (h *ProxyHandler) Handle(c *fiber.Ctx) error {
	return h.handler(c)
}

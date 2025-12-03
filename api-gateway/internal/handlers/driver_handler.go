package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ProxyHandler tek bir driverServiceURL ile tüm gelen driver isteklerini ilgili servise yönlendirir.
// Bu yaklaşım, api-gateway'in minimal proxy görevi görmesi için uygundur.
type ProxyHandler struct {
	target *url.URL
	client *http.Client
}

// Hedef URL ile birlikte http.Client oluşturur.
func NewProxyHandler(targetRaw string) *ProxyHandler {
	u, _ := url.Parse(targetRaw)
	return &ProxyHandler{
		target: u,
		client: &http.Client{
			Timeout: 15 * time.Second, // isteğin zaman aşımına uğrayacağı süre
		},
	}
}

// Handle: fiber'den gelen isteği alır, driver-service'e aynı method & path ile gönderir,
// response'u okuyup client'a iletir.
func (h *ProxyHandler) Handle(c *fiber.Ctx) error {

	originalPath := c.OriginalURL()

	targetURL := h.target.Scheme + "://" + h.target.Host + originalPath //* https "://" + "example.com:8080" + "path"

	req, err := http.NewRequestWithContext(
		context.Background(),
		c.Method(),
		targetURL,
		bytes.NewReader(c.Body()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "proxy request oluşturulamadı", "details": err.Error()})
	}
	defer req.Body.Close()

	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})

	resp, err := h.client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "driver-service erişilemedi", "details": err.Error()})
	}
	defer resp.Body.Close()

	//TODO: Servisin döndüğü headlerları response'a kopyala
	for k, v := range resp.Header {
		for _, v := range v {
			c.Set(k, v)
		}
	}

	//TODO: DURUM KODU AYARLAR VE BODY'Yİ STREAM OLARAK İLET
	c.Status(resp.StatusCode)
	_, err = io.Copy(c, resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "response stream hatası", "details": err.Error()})
	}
	return nil
}

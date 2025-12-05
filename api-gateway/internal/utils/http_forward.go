package utils

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
)

/*
Yönlendirme amaçlı kullanılan bir katman.
Router içerisinde yer alan /_forward/* metodu gibi yollar için kullanılır.
Daha çok "raw" yönlendiricidir. Yani http frameworklerinin (fiber,gin,echo) requestine dokunulmamış bir haliyle başka bir servise iletilmesi.
*/
func ForwardRawRequest(targetBase string) fiber.Handler {

	u, _ := url.Parse(targetBase)
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	return func(c *fiber.Ctx) error {
		//*Path yolu /_forward/{rest}
		targetURL := u.Scheme + "://" + u.Host + c.OriginalURL()

		req, err := http.NewRequestWithContext(
			context.Background(),
			c.Method(),
			targetURL,
			bytes.NewReader(c.Body()))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "forward request oluşturma işlevi başarılamadı",
				"details": err.Error(),
			})
		}

		//todo: header'ı kopyalıyoruz
		c.Request().Header.VisitAll(func(key, val []byte) {
			req.Header.Set(string(key), string(val))
		})

		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error":   "istenen servise erişim sağlanamadı",
				"details": err.Error(),
			})
		}

		defer resp.Body.Close()

		for key, v := range resp.Header {
			for _, val := range v {
				c.Set(key, val)
			}
		}

		c.Status(resp.StatusCode)
		_, err = io.Copy(c, resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "forward resp. stream hatası",
				"details": err.Error(),
			})
		}
		return nil

	}
}

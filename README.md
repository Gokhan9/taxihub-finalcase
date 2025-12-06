# BiTaksi : TaxiHub Final Case - Driver System

Proje mikroservis mimarisi ile geliştirilmiş, ölçeklenebilir ve yüksek performanslı bir backend uygulamasıdır. Sistemin temel amacı, sürücülerin konum bilgilerini (lat-lon) cinsinden alması ve haritalama sistemi üzerinde coğrafi konuma dayalı sorgulama yapmasıdır.

## Mimari Yapı 

Proje, modern bir mikroservis mimarisi ile temellendirilmeye çalışılmıştır.

```marmaid
  Client - HTTP Request - API Gateway[8080]
  API Gateway - JWT Auth ve Rate Limit - GatewayLogic
  GatewayLogic - Proxy - Driver Service[8081]
  Driver Service - Read&Write - Database[Mongo]
```

* API GATEWAY, proje özelinde tek giriş noktasıdır. Kimlik doğrulama(validation), yük dengeleme ve istek yönlendirme yapıyor.
* Driver Service katmanı ise sürücü ile ilgili verileri yönetmeye yarıyor.

### Kullandığım Teknolojiler
* Dil : Go
* Web Framework : Fiber (Diğer frameworklere göre hızlı ve yönetilmesi daha kolay olduğu için seçildi.)
* Veritabanı : MongoDB
* Authentication: JWT (JSON Web Tokens)
* Conteinerization: Docker - Docker Compose


## Proje'nin Mimari Yapısı
* bitaksi.finalcase/api-gateway/internal/client-config-handlers-middleware-router-utils
* bitaksi.finalcase/driver-service/config-docs-postman/internal/config-db-dto-handlers-models-repository-router-services-utils
* bitaksi.finalcase/docker-compose.yaml

## Kurulum ve Çalıştırma
* Projeyi çalıştırmak için 2 farklı yöntem kullanılıyor.

### 1 - Docker ile 
Mevcut bilgisayarınızda Docker ve Docker Compose kuruluysa, tek bir komut kullanarak projeyi ayağa kaldırabilir veya sonlandırabilirsiniz.
```bash
# Servis & Sevisleri Başlatmak için:
docker-compose up --build
# Servis & Servisleri Kapatmak için:
docker-compose down
```

### 2 - Manuel
Servisleri ilgili klasör içerisine giderek tek tek çalıştırabilirsiniz. Bunu yapmak isterseniz:

```
1. Gereksinimler:
MongoDB'nin localinizde çalışıyor olması gerekir.

2. API Gateway:

cd api-gateway
go run cmd/api-gateway/main.go

3. Driver Service:

cd driver-service
go run cmd/driver-service/main.go
```

## Konfigürasyon (.env)

Her servis kendine ait olan (.env) dosyasını kullanıyor.

*API Gateway (.env)
`PORT` = `8080`
`DRIVER_SERVICE_URL` = `http://driver-service:8081`
`JWT_SECRET` = `gizli-anahtar-123`


*Driver Service (.env)
`PORT` = `8081`
`MONGO_URI` = `mongodb://localhost:27017`
`MONGO_DB_NAME` = `bitaksi_db`


# API Dökümantasyonu

Projeyi çalıştırdıktan sonra dökümantasyonlara aşağıdaki adreslerden erişim sağlayabilirsiniz.

* **Swagger UI:** `http://localhost:8081/swagger/index.html` (Driver Service endpointleri için)
* **Postman Collection:** Proje içerisinde `docs/postman/` klasöründe ki JSON dosyalarını Postman'a import edebilirsiniz.

---

![Swagger](https://github.com/user-attachments/assets/f09dd11a-2a12-49a4-aa5f-ffcd53a71b22)
![Post](https://github.com/user-attachments/assets/0bfdcd3f-b9ce-4144-bada-63f9a43dfe4c)
![GET-1 - updatedAt ve Status](https://github.com/user-attachments/assets/08bf4f91-99d5-4163-80ec-2e80b06c923b)
![driverspage=1 pageSize=20](https://github.com/user-attachments/assets/58104f90-c5c0-40e0-9e75-f99d04dd3dc1)
![driversnearbylat=41 00 lon=28 99 taxiType=sari radius=6](https://github.com/user-attachments/assets/1c126ddb-cbfc-457d-88ac-c26ba3389421)
![status(available,busy,offline)](https://github.com/user-attachments/assets/bd9b1a2e-f9f3-4cbf-95f1-720dfe5399a1)
![score](https://github.com/user-attachments/assets/61e526ba-64c5-4265-bc62-34ca72cd2adf)
![Location](https://github.com/user-attachments/assets/15ace614-9fcc-4408-abcf-683448e4af73)
![Delete](https://github.com/user-attachments/assets/3dedd8d8-223f-440d-a2a3-ca148e06eca5)
![Benzer Plaka sonrası err ](https://github.com/user-attachments/assets/1479b231-8ba3-4f48-b51c-862190f28c00)
![Geçersiz Plaka sonrası err](https://github.com/user-attachments/assets/213e99f2-666d-4719-ad89-40080bc19435)
![GetDriverByID](https://github.com/user-attachments/assets/b0fdaf74-0daf-4ec9-a2aa-f870989ae327)

"updatedAt" - Before ![updatedAt-before](https://github.com/user-attachments/assets/06e7d400-e78b-4406-ae94-c88bb750c98f)

"updatedAt" - after ![updatedAt-after](https://github.com/user-attachments/assets/74d17b6b-9f67-4e0e-85d7-4dc62d414d9d)



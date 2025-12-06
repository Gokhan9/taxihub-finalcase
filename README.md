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

### 1: Docker ile 
Mevcut bilgisayarınızda Docker ve Docker Compose kuruluysa, tek bir komut kullanarak projeyi ayağa kaldırabilir veya sonlandırabilirsiniz.
```bash
# Servis & Sevisleri Başlatmak için:
docker-compose up --build
# Servis & Servisleri Kapatmak için:
docker-compose down
```

### 2: Manuel
Servisleri ilgili klasör içerisine giderek tek tek çalıştırabilirsiniz. Bunu yapmak isterseniz:

```

1. Gereksinimler:**
MongoDB'nin localinizde çalışıyor olması gerekir.

2. API Gateway:**
```bash
cd api-gateway
go run cmd/api-gateway/main.go

3. Driver Service:**
cd driver-service
go run cmd/driver-service/main.go
```

## Konfigürasyon (.env)

Her servis kendine ait olan (.env) dosyasını kullanıyor.

**API Gateway (.env)



**Driver Service (.env)

`PORT` = `8081`
`MONGO_URI` = `mongodb://localhost:27017`
`MONGO_DB_NAME` = `bitaksi`

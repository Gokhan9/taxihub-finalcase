# BiTaksi : TaxiHub Final Case - Driver System

Proje mikroservis mimarisi ile geliştirilmiş, ölçeklenebilir ve yüksek performanslı bir backend uygulamasıdır. Sistemin temel amacı, sürücülerin konum bilgilerini (lat-lon) cinsinden alması ve haritalama sistemi üzerinde coğrafi konuma dayalı sorgulama yapmasıdır.

## Mimari Yapı 

Proje, modern bir mikroservis mimarisi ile temellendirilmeye çalışılmıştır.

```marmaid
graph LR
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
bitaksi.finalcase/api-gateway/internal/client-config-handlers-middleware-router-utils
bitaksi.finalcase/driver-service/config-docs-postman/internal/config-db-dto-handlers-models-repository-router-services-utils
bitaksi.finalcase/docker-compose.yaml

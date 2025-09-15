E-Commerce API

Bu proje, MongoDB üzerinde çalışan, Go dilinde geliştirilmiş bir E-ticaret REST API’sidir. Servisler, Echo framework ile oluşturulmuş ve mikroservisler Fast HTTP ile birbirleriyle haberleşmektedir. Proje, güvenli bir kullanıcı yönetimi ve sipariş yönetimi altyapısı sunar.


İçerik

Mikroservisler: Customer ve Order
Veritabanı: MongoDB
Framework: Echo (Go)
Kimlik Doğrulama: JWT
Middleware: Authentication, Authorization, Correlation ID, Error Handling, Logging, Recovery, Role-Based Routing
Özel Hata Yapısı: Kendi Error paketi
Validation: Request doğrulama

Servisler

Customer Servisi

Customer servisinin endpoint’leri:

HTTP	Endpoint	Açıklama
POST	/customer/create	Yeni kullanıcı oluştur
POST	/customer/login	Kullanıcı girişi
GET	/customer/:id	ID’ye göre kullanıcı bilgisi al
GET	/customer/email/:email	Email’e göre kullanıcı bilgisi al (yetkili)
PUT	/customer/:id	Kullanıcıyı güncelle
DELETE	/customer/:id	Kullanıcıyı sil
GET	/customer/list	Tüm kullanıcıları listele
GET	/customer/verify	JWT doğrulama



Order Servisi

Order servisinin endpoint’leri:

HTTP	Endpoint	Açıklama
POST	/order	Yeni sipariş oluştur
GET	/order/:id	ID’ye göre sipariş getir
PATCH	/order/:id/ship	Siparişi kargoya ver
PATCH	/order/:id/deliver	Siparişi teslim et
DELETE	/order/cancel/:id	Siparişi iptal et
GET	/order/list	Tüm siparişleri listele

Auth & Authorization

JWT Tabanlı Authentication
Role-Based Authorization
Middleware ile korunmuş endpoint’ler
JWT token doğrulama ve erişim kontrolü

Middleware

Projede kullanılan middleware’ler:
Authentication – Kullanıcı doğrulaması
Authorization – Role bazlı yetki kontrolü
Correlation ID – Request takibi
Error Handling – Özel hata yönetimi
Logging – Tüm request ve error log’ları
Recovery – Panic durumlarını yakalama
Role Routing – Roller üzerinden endpoint erişimi

Validation

Request veri doğrulama uygulanmıştır
Kullanıcı giriş, kayıt ve sipariş işlemlerinde validation kontrolleri yapılır

Özel Hata Yapısı

Proje, kendi error paketi ile yapılandırılmıştır
Hatalar, detaylı mesaj ve HTTP durum kodları ile döner
Merkezi error handling middleware ile yönetilir

Mikroservis Haberleşmesi

Servisler Fast HTTP kullanarak birbirleriyle haberleşir
Performans odaklı ve düşük gecikmeli istek yönetimi sağlar

Kurulum

Repo’yu klonlayın:

git clone <repo-url>
cd <repo-folder>
Gerekli Go modüllerini yükleyin:
go mod tidy
.env veya config dosyasını yapılandırın (MongoDB bağlantısı, JWT secret vb.)
Sunucuyu başlatın:
go run main.go

Kullanım

Postman veya herhangi bir API client ile endpoint’leri test edebilirsiniz.
Öncelikle /customer/create ile kullanıcı oluşturun ve /customer/login ile JWT token alın.
Token’ı kullanarak yetkilendirilmiş endpoint’lere erişin.

Teknolojiler

Go (Echo framework)
MongoDB
JWT
Fast HTTP
Middleware & Validation

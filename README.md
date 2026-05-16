# KARE-REHBER

Koçluk takip ve değerlendirme platformu. Onaylı veri akışı ile öğrencilerin haftalık gelişim sürecini takip eder. Admin / Koç / Öğrenci / Veli / Koordinatör (Vakıf) rolleri için ayrı paneller.

## İçindekiler

- [Stack](#stack)
- [Dizin Yapısı](#dizin-yapısı)
- [Hızlı Başlangıç](#hızlı-başlangıç)
- [Ortam Değişkenleri](#ortam-değişkenleri)
- [Railway Deploy](#railway-deploy)
- [Roller ve Yetkiler](#roller-ve-yetkiler)
- [API Özeti](#api-özeti)
- [İş Mantığı Notları](#iş-mantığı-notları)

## Stack

- **Backend:** Go 1.22+, Fiber v2, GORM, PostgreSQL 16, golang-migrate, JWT (golang-jwt/jwt), bcrypt, validator/v10
- **Frontend:** Vue 3 + Vite + TypeScript, Element Plus, Pinia, Vue Router, vue-i18n, ECharts
- **DB:** PostgreSQL 16

## Dizin Yapısı

```
.
├── backend/
│   ├── cmd/
│   │   ├── api/           # HTTP server entry
│   │   ├── migrate/       # golang-migrate CLI wrapper
│   │   └── seed/          # admin + 36 hafta seed
│   ├── internal/
│   │   ├── auth/          # JWT + bcrypt + password gen
│   │   ├── config/        # env, DB connect
│   │   ├── domain/        # GORM modelleri (User, Evaluation, ...)
│   │   ├── dto/           # request/response yapıları
│   │   ├── handler/       # Fiber HTTP handler'lar
│   │   ├── middleware/    # JWTAuth, RequireRole
│   │   ├── repository/    # DB erişim katmanı
│   │   ├── router/        # route kayıtları
│   │   ├── service/       # iş mantığı
│   │   └── sms/           # Provider interface + MockProvider
│   ├── migrations/        # SQL migration dosyaları
│   ├── pkg/               # logger, errors, pagination
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── api/           # axios client'ları
│   │   ├── components/    # paylaşılan UI
│   │   ├── layouts/       # AdminLayout, RoleLayout
│   │   ├── locales/       # vue-i18n (TR)
│   │   ├── router/        # rol bazlı guard'lar
│   │   ├── stores/        # Pinia (auth)
│   │   ├── types/
│   │   └── views/
│   │       ├── auth/      # Login, kayıt formları
│   │       ├── admin/     # Yönetim ekranları
│   │       ├── coach/
│   │       ├── coordinator/
│   │       ├── parent/
│   │       └── student/
│   └── package.json
├── docker-compose.yml     # Postgres (yalnızca dev)
├── Makefile
└── README.md
```

## Hızlı Başlangıç

### 1. Veritabanı

**Seçenek A — Docker (önerilen)**
```bash
docker compose up -d postgres
# veya
make db-up
```

Docker konteyneri **5433** portunda yayınlar (yerelde 5432'de çalışan başka bir Postgres'le çakışmasın diye). `backend/.env` zaten 5433'e ayarlı.

**Seçenek B — Yerel Postgres** (Homebrew ile kurulu ise)
```bash
psql -d postgres -c "CREATE USER kareuser WITH PASSWORD 'karepass';"
psql -d postgres -c "CREATE DATABASE karerehber OWNER kareuser;"
# .env içindeki DATABASE_URL'in port'unu 5432 yap
```

### 2. Backend

```bash
cd backend
cp .env.example .env
go mod tidy
go run ./cmd/migrate up   # tabloları oluştur
go run ./cmd/seed         # admin + 36 hafta seed et
go run ./cmd/api          # API'yi başlat
```

API: http://localhost:8080/api/v1/health

Hot-reload için [`air`](https://github.com/air-verse/air):
```bash
go install github.com/air-verse/air@latest
air
```

### 3. Frontend

```bash
cd frontend
npm install
npm run dev
```

UI: http://localhost:5173

Vite proxy ile `/api` istekleri backend'e (`localhost:8080`) yönlendirilir.

### Test Kullanıcıları

`go run ./cmd/seed` aşağıdaki test verisini oluşturur. Login ekranında **tek tıkla giriş butonları** (Admin / Koç / Öğrenci / Veli / Koordinatör) bu hesapları kullanır.

| Rol | Telefon | Şifre |
|---|---|---|
| Admin | `05000000000` | `Admin123!` |
| Koç (Ahmet) | `05311111111` | `Test123!` |
| Koç (Zeynep) | `05311111112` | `Test123!` |
| Koç (Mehmet) | `05311111113` | `Test123!` |
| Öğrenci (Ali) | `05322222221` | `Test123!` |
| Öğrenci (Ayşe…Emir) | `05322222222`–`05322222226` | `Test123!` |
| Veli (Hasan) | `05333333331` | `Test123!` |
| Veli (Fatma, Selma) | `05333333332`–`05333333333` | `Test123!` |
| Koordinatör (Aslı) | `05344444441` | `Test123!` |
| Koordinatör (Murat) | `05344444442` | `Test123!` |

Seed ayrıca **3 koç ↔ 6 öğrenci ataması**, **veli–çocuk ilişkileri**, **koordinatör atamaları** ve **3 örnek değerlendirme** (1 onaylı, 1 pending, 1 onaylı) ekler. Bilerek bazı öğrencilere değerlendirme girilmez ki "eksik girişler raporu" test edilebilsin.

Tüm verileri sıfırlayıp baştan seed'lemek için:
```bash
cd backend && go run ./cmd/seed --wipe
```

⚠ Production'da bu hesapları silin veya şifrelerini değiştirin.

## Ortam Değişkenleri

Backend için (`backend/.env`):

| Değişken | Açıklama | Varsayılan |
|---|---|---|
| `APP_ENV` | `development` / `production` | `development` |
| `PORT` veya `HTTP_PORT` | API port (Railway `PORT`'u otomatik atar) | `8080` |
| `DATABASE_URL` | Postgres bağlantısı (sslmode yoksa otomatik `disable` eklenir) | `postgres://kareuser:karepass@localhost:5433/karerehber?sslmode=disable` |
| `JWT_SECRET` | Token imzalama | (production'da zorunlu) |
| `JWT_EXPIRES_HOURS` | Token süresi | `24` |
| `BCRYPT_COST` | Şifre hash cost | `12` |
| `CORS_ALLOWED_ORIGINS` | Virgülle ayrılmış izinli origin'ler veya `*` | `http://localhost:5173` |
| `AUTO_MIGRATE` | API başlarken migration'ları çalıştırsın mı? | `true` |
| `SEED_ADMIN` | Admin + 36 hafta yoksa eklensin mi? | `true` |
| `SEED_TEST_USERS` | Test koç/öğrenci/veli/koordinatör seed edilsin mi? | `false` |

Frontend için (build sırasında okunur):

| Değişken | Açıklama | Varsayılan |
|---|---|---|
| `VITE_API_URL` | Backend API base URL (örn. `https://api.kare.app`). Boşsa relative `/api/...` kullanır (dev'de Vite proxy yönlendirir). | `""` |

## Railway Deploy

Tek repo'dan iki servis çıkar: **backend** (Go API) ve **frontend** (Vue SPA). Railway otomatik olarak `railway.json` dosyasını okur.

### Adım 1 — PostgreSQL

1. Railway projesinde **New → Provision PostgreSQL**.
2. Postgres servisinin `Variables` sekmesinden `DATABASE_URL` (internal) kopyalayın. Format:
   `postgresql://postgres:****@postgres.railway.internal:5432/railway`

### Adım 2 — Backend Servisi

1. **New → GitHub Repo** ile `kare-rehber` repo'sunu seç.
2. Servisin **Settings → Source → Root Directory** ayarını `backend` olarak değiştir.
3. **Variables** sekmesinde:
   - `DATABASE_URL` → Postgres servisinden **Reference Variable** olarak ekle (şifre değişse bile bağlı kalır).
   - `JWT_SECRET` → uzun rastgele bir string (örn. `openssl rand -hex 32`).
   - `APP_ENV` → `production`
   - `CORS_ALLOWED_ORIGINS` → frontend'in Railway URL'i (Adım 3'ten sonra eklenir; başlangıçta `*` bırakılabilir).
   - (Opsiyonel) `SEED_TEST_USERS=true` ile demo kullanıcıları otomatik oluştur. Canlıya çıktıktan sonra `false` yap.
4. **Settings → Networking → Generate Domain** ile public URL al (örn. `kare-rehber-api.up.railway.app`).

`AUTO_MIGRATE=true` (default) olduğu için API ilk başlangıçta migration'ları çalıştırıp admin'i ve haftaları seed eder. Manuel adım gerekmez.

### Adım 3 — Frontend Servisi

1. **New → GitHub Repo** ile aynı repo'yu seç (Railway aynı repo'dan ikinci servisi sorunsuz açar).
2. **Settings → Source → Root Directory** = `frontend`.
3. **Variables**:
   - `VITE_API_URL` → backend servisinin URL'i (Adım 2'den, sondaki `/` olmadan, örn. `https://kare-rehber-api.up.railway.app`).
4. **Settings → Networking → Generate Domain** ile frontend için public URL al.
5. Backend servisine geri dön ve `CORS_ALLOWED_ORIGINS`'i frontend URL'i yap (örn. `https://kare-rehber.up.railway.app`).

### Önemli notlar

- **PORT**: Hem backend hem frontend Railway'in inject ettiği `PORT` ortam değişkenini kullanır; manuel port atama yapma.
- **Build**:
  - Backend: `go build -o api ./cmd/api`
  - Frontend: `npm install --no-audit --no-fund && npm run build`
- **Healthcheck**: Backend `/api/v1/health` endpoint'i Railway tarafından kontrol edilir.
- **Migration'lar binary'ye embed edilir** (`backend/internal/dbmigrate/sql/`). Yeni migration eklerken oraya `.sql` koy ve commit et — Railway redeploy'da otomatik uygulanır.
- **Frontend env'leri build-time**: `VITE_API_URL`'i değiştirirsen Railway redeploy gerekir.

## Roller ve Yetkiler

| Rol | Görebildiği |
|---|---|
| **Admin** | Her şey: kullanıcı yönetimi, başvuru onayları, eşleştirme, hafta yönetimi, tüm değerlendirmeler (ham + düzenli + tüm versiyonlar + admin_only_note), mesajlar, SMS, raporlar |
| **Koç** | Atandığı öğrenciler, haftalık değerlendirme yazma/düzenleme (yalnızca pending iken) |
| **Öğrenci** | Mesajlaşma (admin/koordinatör), geri bildirim ve şikayet açma |
| **Veli** | Çocuğunun yalnızca **onaylanmış** değerlendirmeleri (admin_only_note hariç), mesajlaşma |
| **Koordinatör (Vakıf)** | Sorumlu olduğu öğrencilerin **aktif** değerlendirmeleri (admin düzenlemişse düzenlenmiş hâli; admin_only_note görmez), mesajlaşma |

## API Özeti

Tüm endpoint'ler `/api/v1/...` altındadır.

```
GET    /health
POST   /auth/login
GET    /auth/me                                  (auth)
POST   /registrations/student                    (public)
POST   /registrations/coach                      (public)

GET    /weeks/open                               (auth)
GET    /evaluations                              (auth, role-filtered)
GET    /evaluations/:id                          (auth)

POST   /coach/evaluations                        (coach)
PATCH  /coach/evaluations/:id                    (coach, only pending)
GET    /coach/students                           (coach)

GET    /threads                                  (auth)
POST   /threads                                  (auth, opener: student/parent/admin/coordinator)
GET    /threads/:id/messages                     (auth)
POST   /threads/:id/messages                     (auth)
POST   /threads/:id/close                        (auth)

# Admin only
GET    /admin/registrations
POST   /admin/registrations/:id/decision

GET    /admin/users
POST   /admin/users
GET    /admin/users/:id
PATCH  /admin/users/:id
POST   /admin/users/:id/reset-password

GET    /admin/assignments
POST   /admin/assignments
POST   /admin/assignments/unassign
GET    /admin/assignments/students-without-coach
POST   /admin/assignments/set-coordinator
POST   /admin/assignments/set-parent

GET    /admin/weeks
POST   /admin/weeks
PATCH  /admin/weeks/:id
POST   /admin/weeks/:id/open
POST   /admin/weeks/:id/close
POST   /admin/weeks/generate

PATCH  /admin/evaluations/:id        (eski sürüm evaluation_versions'a yazılır)
POST   /admin/evaluations/:id/approve
POST   /admin/evaluations/:id/revoke
GET    /admin/evaluations/:id/versions

GET    /admin/sms/templates
POST   /admin/sms/send                            (mock provider — DB'ye log)
GET    /admin/sms/logs
GET    /admin/sms/missing-coaches?weekId=...

GET    /admin/reports/summary
GET    /admin/reports/coach-performance
GET    /admin/reports/city-distribution?role=...
GET    /admin/reports/week-stats
GET    /reports/student-trend?studentId=...
```

## İş Mantığı Notları

### Onay & Versiyon Akışı

1. **Koç** haftalık değerlendirme yazar → `status = pending`
2. **Admin** doğrudan düzenleyebilir → eski hâl `evaluation_versions`'a yazılır, aktif kayıt güncellenir, `status = edited_by_admin`
3. **Admin** onaylar → `status = approved`, `approved_by/approved_at` doldurulur
4. **Admin** onayı geri çekebilir → `status = pending`

### Hafta Yönetimi

- Akademik yıl için **36 hafta** seed edilir (Eylül ilk Pazartesi'den itibaren)
- Her istek başlangıcında / hafta endpoint çağrısında `is_open` durumu otomatik senkronlanır:
  - Bugünün dahil olduğu hafta + bir önceki hafta → açık
  - Daha eski haftalar → otomatik kapalı
  - Manuel açılmış (admin tarafından) haftalar elle kapatılana dek açık kalır
- Koç UI'sinde sadece `is_open=true` olan haftalar görünür

### Mock SMS

Faz 7'de mock provider kullanılır — gerçekten SMS göndermez, `sms_logs` tablosuna `status = mock_sent` ile yazar. Gerçek sağlayıcı (NetGSM / İletiMerkezi) eklemek için `internal/sms/Provider` interface'ini implemente edip `MockProvider` yerine inject edin.

### Görünürlük Kuralları

- **Veli:** Sadece `status = approved` değerlendirmeler. `admin_only_note` her zaman gizli.
- **Koordinatör:** Aktif kayıt (admin düzenlemişse düzenli hâl). `admin_only_note` gizli.
- **Admin:** Her şey + tüm versiyon geçmişi.

## Geliştirme Notları

- Soft-delete YOK; admin değişiklikleri `evaluation_versions` ve genel `audit_logs` ile takip edilir.
- API versiyonlama: `/api/v1/...` baştan.
- Token süresi 24h; refresh token MVP'de yok.
- Şifre hash: bcrypt cost 10 (dev) / 12 (prod).

## Lisans

İç kullanım.

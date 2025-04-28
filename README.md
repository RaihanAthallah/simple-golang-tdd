# Simple Golang TDD Project

## 🔄 Clone Project

```bash
git clone https://github.com/RaihanAthallah/simple-golang-tdd.git
cd simple-golang-tdd
```

## 🔧 Installation

```bash
go mod tidy
```

## 🌍 Running Locally

```bash
go run main.go
```

Aplikasi akan berjalan di `http://localhost:8080`.

Swagger dapat diakses di:

```text
http://localhost:8080/swagger/index.html#/
```

---

# 📦 Deploy menggunakan Docker

### 1. Build Docker Image

```bash
docker build -t simple-golang-tdd .
```

### 2. Run Docker Container

```bash
docker-compose up --build
```

Aplikasi akan otomatis berjalan di `http://localhost:8080` dalam container.

---

# 📊 Struktur Folder

```plaintext
├── config/
├── controller/
│   ├── auth/
│   └── customer/
├── data/
├── docs/
├── dto/
├── middleware/
├── model/
├── repository/
│   ├── customer/
│   ├── history/
│   └── merchant/
├── routes/
├── service/
│   ├── auth/
│   └── customer/
├── utils/
├── main.go
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── .env
```

## 📂 Penjelasan Folder

| Folder / File          | Penjelasan                                                           |
| :--------------------- | :------------------------------------------------------------------- |
| **config/**            | Konfigurasi aplikasi (database, environment).                        |
| **controller/**        | Menangani request & response API. Dibagi ke `auth/` dan `customer/`. |
| **data/**              | Menyimpan file JSON sebagai database sederhana.                      |
| **docs/**              | Dokumentasi project, termasuk file swagger.                          |
| **dto/**               | Data Transfer Object: format data request & response.                |
| **middleware/**        | Middleware untuk autentikasi dan logging.                            |
| **model/**             | Definisi struktur data utama (struct).                               |
| **repository/**        | Interaksi data: membaca/menulis file JSON atau database.             |
| **routes/**            | Mapping endpoint URL ke controller.                                  |
| **service/**           | Business logic aplikasi, dibagi untuk `auth/` dan `customer/`.       |
| **utils/**             | Helper function seperti token generator, hashing, validator.         |
| **main.go**            | Entry point aplikasi, menginisialisasi semua komponen.               |
| **Dockerfile**         | Instruksi untuk membuat Docker image.                                |
| **docker-compose.yml** | Mengatur service container.                                          |
| **.env**               | Menyimpan environment variables seperti PORT dan SECRET KEY.         |

---

# 📂 Contoh Environment (.env)

```env
ACCESS_SECRET=youraccesstokensecret
REFRESH_SECRET=yourrefreshtokensecret
```

---

# 📂 Membaca Data dari Folder Data

Semua file JSON yang ada di `./data/*.json` akan dibaca dan dipakai sebagai database sederhana untuk aplikasi ini.

Pastikan file JSON valid untuk menghindari error saat aplikasi membaca file.

Contoh file di `./data/history.json`:

```json
[
  {
    "id": "history-001",
    "customer_id": "cust-001",
    "action": "payment",
    "details": {
      "merchant_id": "merchant-001",
      "amount": 100.5
    },
    "timestamp": "2024-04-27T12:00:00Z"
  }
]
```

---

# 🚀 Tentang Pengembangan

Aplikasi ini dibangun menggunakan **Golang 1.23.0** sebagai bahasa pemrograman utama, dengan menggunakan **Gin** sebagai framework web. Gin adalah framework yang ringan dan cepat, cocok untuk membuat API dengan performa tinggi.

## 🧪 Pendekatan Test-Driven Development (TDD)

Pendekatan yang digunakan dalam pengembangan aplikasi ini adalah **Test-Driven Development (TDD)**. Dalam pendekatan ini, pengujian dilakukan terlebih dahulu (unit tests) sebelum kode aplikasi yang sesungguhnya ditulis. Berikut adalah alur umum dalam pengembangan aplikasi dengan TDD:

1. **Menulis Unit Test**: Pertama, kita menulis unit test yang mendefinisikan perilaku yang diinginkan dari fungsi atau fitur yang akan dibangun.
2. **Menjalankan Unit Test**: Setelah unit test ditulis, jalankan untuk memastikan bahwa tes gagal (karena implementasinya belum ada).
3. **Menulis Kode**: Kemudian, implementasikan kode untuk membuat unit test berhasil.
4. **Refactor**: Setelah kode implementasi berhasil, refactor kode untuk meningkatkan kualitasnya tanpa mengubah fungsionalitasnya.

---

## 🔒 Endpoint yang Membutuhkan Token

### Proteksi Token

- **POST** `/api/v1/customer/payment`

Untuk mendapatkan token:

- **POST** `/user/v1/auth/login`

Setelah login berhasil, gunakan token pada header Authorization:

```bash
Authorization: Bearer <your_token>
```

---

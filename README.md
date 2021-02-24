# RISA_RESTFULL

Restfull Backend for Risa Aplication Pelindo III using Golang (Fiber) and MongoDB

## Dependency lokal

- [ErruUtils](https://github.com/muchlist/erru_utils_go/)  
  Library ini digunakan untuk memformat response error dan logger sehingga response error memiliki format yang standart
  di setiap service (berguna jika akan mengimplementasikan microservice).

## Dependency pihak ketiga

- [Go Fiber Framework](https://github.com/gofiber/fiber/) : Web framework golang yang memiliki kemiripan dengan express
  js dan menggunakan fast-http (tidak berbeda jauh dengan gin dan echo).
- [Mongo go driver](https://go.mongodb.org/mongo-driver/) : Saat ini service ini full menggunakan MongoDB.
- [JWT go](https://github.com/dgrijalva/jwt-go/)
- [Ozzo validation](https://github.com/go-ozzo/ozzo-validation/) : Library yang digunakan untuk validasi request body
  dari user. (Karena Go Fiber tidak memiliki input validasi seperti Binding di Gin)

## LOG

- `gen_unit` domain. `gen_unit` digunakan untuk meng-collect semua perangkat dengan hanya menyimpan data umumnya saja
  dan meninggalkan data detil.
  `gen_unit` dibuat karena ada permintaan dari client agar semua perangkat dapat dicari menggunakan satu buah kolom
  pencarian tanpa harus memilih kategori. Semakin banyak data akan semakin lambat sehingga kedepan akan diganti
  menggunakan database elasticsearch. domain ini tidak bersentuhan secara langsung dengan user dari segi inputan.
  updatenya akan dilakukan dibelakang layar berasarkan : pembuatan perangkat pada kategori apapun, pengeditan jika nama,
  ip , category, cabang berubah. dan penghapusan. serta ada update pada history/incident.
- `history` digunakan untuk mencatat semua riwayat perangkat, riwayat ini memiliki status progress (1), persetujuan
  pending (2), pending (3), complete (4). Setiap penambahan `history` yang belum komplit akan mengupdate field `cases`
  pada domain `gen_unit` dan jika `history` diubah statusnya menjadi complete maka case di `gen_unit` akan dikurangi

## Kontrak Struktur
### Handler > Middleware > Service > Dao

- Handler digunakan untuk mengekstrak inputan dari user. params, query, json body, claims dari jwt serta validasi input
  ,termasuk memastikan dan menimpa huruf besar atau kecil.
- Service digunakan untuk bisnis logic, menggabungkan dua atau lebih dao atau utilitas pembantu lainnya, mengisi data
  yang dibutuhkan dao misalnya saat perpindahan dari requestData (data sedikit) ke Data (data banyak).
- Dao berkomunikasi langsung ke database. Sedikit juga memastikan inputan huruf besar dan kecil pada inputan database
  yang caseSensitif untuk memaksimalkan indexing.
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
- `history` digunakan untuk mencatat semua riwayat perangkat, riwayat ini memiliki status info (0), progress (1),
  persetujuan pending (2), pending (3), complete (4). Setiap penambahan `history` yang belum komplit akan mengupdate
  field `cases`
  pada domain `gen_unit` dan jika `history` diubah statusnya menjadi complete maka case di `gen_unit` akan dikurangi
- `cctv`
- `check` menggenerate daftar tempat atau perangkat yang harus di cek dengan menyesuaikan shifts inputan. Didalam check
  akan ada `check item` yang cocok dan juga `general unit cctv` dengan kriteria tertentu (dalam hal ini perangkat yang
  down dan field `cases` kosong maka harus dimunculkan di checklist). ketika mengupdate check item didalamnyam akan
  mengupdate juga entity `check item` kecuali yang berdasarkan general unit (misalnya cctv) akan diupdate ketika `check`
  ditandai sudah selesai (isFinish = true).
  `check item` yang ditandai mempunyai problem juga akan di munculkan pada saat pembuatan check berikutnya.
- `check item` sebagai template yang mana saja yang mau di cek. didalamnya ada slice shift untuk dimunculkan
  saat `check` dibuat.
- `stock` menyimpan stock sebagai satu buah dokumen saja , pemakaian dan penambahan stok dijadikan sebagai child didalam
  dokumen dan setiap perubahannya akan mempengaruhi field QTY pada stock.

## Kontrak Struktur

### Handler > Middleware > Service > Dao

- Handler digunakan untuk mengekstrak inputan dari user. params, query, json body, claims dari jwt serta validasi input
  ,termasuk memastikan dan menimpa huruf besar atau kecil.
- Service digunakan untuk bisnis logic, menggabungkan dua atau lebih dao atau utilitas pembantu lainnya, mengisi data
  yang dibutuhkan dao misalnya saat perpindahan dari requestData (data sedikit) ke Data (data banyak). termasuk merubah
  string menjadi ObjectID dan Pengecekan IP address.
- Dao berkomunikasi langsung ke database. Beberapa kasus juga memastikan inputan huruf besar dan kecil pada inputan
  database yang caseSensitif untuk memaksimalkan indexing, memastikan nilai yang di input array<T> apabila array nil.
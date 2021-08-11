# RISA_RESTFULL

RestfulApi Backend for Risa Aplication Pelindo III using Golang (Fiber) and MongoDB

## Fitur
- riwayat pemeliharaan
- penarikan laporan
- data inventaris alat
- stok manajemen
- monitoring speed test, cctv
- checklist pengecekan harian
- checklist maintenance cctv oleh vendor
- tugas perbaikan / kemajuan
- notifikasi perangkat bermasalah


## Playstore
- [Google Playstore](https://play.google.com/store/apps/details?id=dev.muchlis.risa2)

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
  dari user.
- [Go Cron](https://github.com/go-co-op/gocron/) : Scheduller
- [Maroto](https://github.com/johnfercher/maroto/) : Framework pembuatan PDF
- [Firebase](https://firebase.google.com/go/v4) : Notifikasi realtime ke Android

## LOG

- `gen_unit` domain. `gen_unit` digunakan untuk meng-collect semua perangkat dengan hanya menyimpan data umumnya saja
  dan meninggalkan data detil.
  `gen_unit` dibuat karena ada permintaan dari client agar semua perangkat dapat dicari menggunakan satu buah kolom
  pencarian tanpa harus memilih kategori. Semakin banyak data akan semakin lambat sehingga kedepan akan diganti
  menggunakan database elasticsearch. domain ini tidak bersentuhan secara langsung dengan user dari segi inputan.
  updatenya akan dilakukan dibelakang layar berdasarkan : pembuatan perangkat pada kategori apapun, pengeditan jika nama,
  ip , category, cabang berubah. dan penghapusan. serta ada update pada history/incident.  `gen_unit` juga memuat data ping alamat ip kghusus perangkat
yang memiliki ip address.
- `history` digunakan untuk mencatat semua riwayat perangkat, riwayat ini memiliki status info (0), progress (1),
  persetujuan pending (2), pending (3), complete (4). Setiap penambahan `history` yang belum komplit akan mengupdate
  field `cases` pada domain `gen_unit` dan jika `history` diubah statusnya menjadi complete maka case di `gen_unit` akan dikurangi.
  `history` memiliki `history` lagi didalamnya untuk keperluan tracking perubahan dan pembuatan laporan
  berdasarkan range waktu tertentu.
- `cctv`, `computer`, `application` dll yang serupa memuat data inventaris.
- `check` menggenerate daftar tempat atau perangkat yang harus di cek dengan menyesuaikan waktu shifts realtime.
  `check item` yang ditandai have problem juga akan di munculkan pada saat pembuatan check berikutnya.
- `check item` sebagai template item yang mana saja yang mau di cek. didalamnya ada slice shift untuk dimunculkan
  saat `check` dibuat.
- `checklist_cctv` membuat ceklist maintenance harian atau bulanan cctv oleh vendor cctv
- `stock` menyimpan stock sebagai satu buah dokumen saja , pemakaian dan penambahan stok dijadikan sebagai child didalam
  dokumen dan setiap perubahannya akan mempengaruhi field QTY pada stock.
- `scheduller` setiap satu jam sistem akan memeriksa cctv yang status pingnya down. status ping didapatkan dari inputan aplikasi lain bernama pingers.
hasil pemeriksaan dikirimkan ke user menggunakan `firebase`


## Kontrak Struktur

### Middleware > Handler > Service > Dao || Api

- Handler digunakan untuk mengekstrak inputan dari user. params, query, json body, claims dari jwt serta validasi input
  ,termasuk memastikan dan menimpa huruf besar atau kecil.
- Service digunakan untuk bisnis logic, menggabungkan dua atau lebih dao atau utilitas pembantu lainnya, mengisi data
  yang dibutuhkan dao misalnya saat perpindahan dari requestData (data sedikit) ke Data (data banyak). termasuk merubah
  string menjadi ObjectID dan Pengecekan IP address.
- Dao berkomunikasi langsung ke database. Beberapa kasus juga memastikan inputan huruf besar dan kecil pada inputan
  database yang caseSensitif untuk memaksimalkan indexing, memastikan nilai yang di input array<T> apabila array nil.
- Api (folder client) merupakan aplikasi pihak luar. aplikasi bisa berkomunikasi dengan api pihak luar menggunakan rest api.

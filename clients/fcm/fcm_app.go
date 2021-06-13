package fcm

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/muchlist/erru_utils_go/logger"
	"google.golang.org/api/option"
	"log"
	"os"
)

const (
	firebaseCred = "GOOGLE_APPLICATION_CREDENTIALS"
)

var (
	FCM     *firebase.App
	fcmCred string
)

// Init menginisiasi firebase app
// responsenya digunakan untuk memutus koneksi apabila main program dihentikan
func Init() error {
	if os.Getenv(firebaseCred) == "" {
		log.Fatal("firebase credensial tidak boleh kosong ENV: GOOGLE_APPLICATION_CREDENTIALS")
	}
	fcmCred = os.Getenv(firebaseCred)
	opt := option.WithCredentialsFile(fcmCred)

	fcm, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Error("gagal membuat app firebase", err)
		return err
	}
	FCM = fcm

	return nil
}

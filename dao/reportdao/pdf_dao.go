package reportdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

const (
	connectTimeout   = 3
	keyPdfCollection = "pdf"

	keyPdfCreatedAt = "created_at"
	keyPdfBranch    = "branch"
)

func NewPdfDao() PdfDaoAssumer {
	return &pdfDao{}
}

type pdfDao struct {
}

type PdfDaoAssumer interface {
	InsertPdf(input dto.PdfFile) (*string, rest_err.APIError)
	FindPdf(branch string) ([]dto.PdfFile, rest_err.APIError)
}

func (c *pdfDao) InsertPdf(input dto.PdfFile) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyPdfCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	input.Branch = strings.ToUpper(input.Branch)
	input.CreatedBy = strings.ToUpper(input.CreatedBy)

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan pdf ke database", err)
		logger.Error("Gagal menyimpan %s ke database (InsertPdf)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *pdfDao) FindPdf(branch string) ([]dto.PdfFile, rest_err.APIError) {
	coll := db.DB.Collection(keyPdfCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	branch = strings.ToUpper(branch)

	// filter
	filter := bson.M{
		keyPdfBranch: branch,
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyPdfCreatedAt, -1}}) //nolint:govet
	opts.SetLimit(10)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar pdf dari database (FindPDF)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.PdfFile{}, apiErr
	}

	var pdfList []dto.PdfFile
	if err = cursor.All(ctx, &pdfList); err != nil {
		logger.Error("Gagal decode pdfList cursor ke objek slice (FindPdf)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.PdfFile{}, apiErr
	}

	return pdfList, nil
}

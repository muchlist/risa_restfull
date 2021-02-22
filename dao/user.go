package dao

import (
	"errors"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"context"

	"github.com/muchlist/risa_restfull/db"
)

const (
	connectTimeout = 2

	keyUserColl = "user"

	keyUserID        = "_id"
	keyUserEmail     = "email"
	keyUserHashPw    = "hash_pw"
	keyUserName      = "name"
	keyUserRoles     = "roles"
	keyUserAvatar    = "avatar"
	keyUserTimeStamp = "timestamp"
)

func NewUserDao() UserDaoAssumer {
	return &userDao{}
}

type userDao struct {
}

type UserDaoAssumer interface {
	InsertUser(user dto.UserRequest) (*string, rest_err.APIError)
	GetUserByID(userID string) (*dto.UserResponse, rest_err.APIError)
	GetUserByIDWithPassword(userID string) (*dto.User, rest_err.APIError)
	FindUser() (dto.UserResponseList, rest_err.APIError)
	CheckIDAvailable(email string) (bool, rest_err.APIError)
	EditUser(userID string, userRequest dto.UserEditRequest) (*dto.UserResponse, rest_err.APIError)
	DeleteUser(userID string) rest_err.APIError
	PutAvatar(userID string, avatar string) (*dto.UserResponse, rest_err.APIError)
	ChangePassword(data dto.UserChangePasswordRequest) rest_err.APIError
}

//InsertUser menambahkan user dan mengembalikan insertedID, err
func (u *userDao) InsertUser(user dto.UserRequest) (*string, rest_err.APIError) {

	coll := db.Db.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	insertDoc := bson.D{
		{keyUserID, user.ID},
		{keyUserName, user.Name},
		{keyUserEmail, strings.ToLower(user.Email)},
		{keyUserRoles, user.Roles},
		{keyUserAvatar, user.Avatar},
		{keyUserHashPw, user.Password},
		{keyUserTimeStamp, user.Timestamp},
	}

	result, err := coll.InsertOne(ctx, insertDoc)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan user ke database", err)
		logger.Error("Gagal menyimpan user ke database", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

//GetUser mendapatkan user dari database berdasarkan userID, jarang digunakan
//pada case ini biasanya menggunakan email karena user yang digunakan adalah email
func (u *userDao) GetUserByID(userID string) (*dto.UserResponse, rest_err.APIError) {

	coll := db.Db.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	var user dto.UserResponse
	opts := options.FindOne()
	opts.SetProjection(bson.M{keyUserHashPw: 0})

	if err := coll.FindOne(ctx, bson.M{keyUserID: strings.ToUpper(userID)}, opts).Decode(&user); err != nil {

		if err == mongo.ErrNoDocuments {
			//apiErr := rest_err.NewNotFoundError(fmt.Sprintf("User dengan ID %v tidak ditemukan", userID.Hex()))
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("User dengan ID %s tidak ditemukan", userID))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan user (by ID) dari database", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan user dari database", err)
		return nil, apiErr
	}

	return &user, nil
}

//GetUserByIDWithPassword mendapatkan user dari database berdasarkan id dengan memunculkan passwordhash
//password hash digunakan pada endpoint login dan change password
func (u *userDao) GetUserByIDWithPassword(userID string) (*dto.User, rest_err.APIError) {

	coll := db.Db.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	var user dto.User

	if err := coll.FindOne(ctx, bson.M{keyUserID: strings.ToUpper(userID)}).Decode(&user); err != nil {

		if err == mongo.ErrNoDocuments {
			// karena sudah pasti untuk keperluan login maka error yang dikembalikan unauthorized
			apiErr := rest_err.NewUnauthorizedError("Username atau password tidak valid")
			return nil, apiErr
		}

		logger.Error("Gagal mendapatkan user dari database (GetUserByIDWithPassword)", err)
		apiErr := rest_err.NewInternalServerError("Error pada database", errors.New("database error"))
		return nil, apiErr
	}

	return &user, nil
}

//FindUser mendapatkan daftar semua user dari database
func (u *userDao) FindUser() (dto.UserResponseList, rest_err.APIError) {

	coll := db.Db.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	users := dto.UserResponseList{}
	opts := options.Find()
	opts.SetSort(bson.D{{keyUserID, -1}})
	sortCursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan user dari database", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.UserResponseList{}, apiErr
	}

	if err = sortCursor.All(ctx, &users); err != nil {
		logger.Error("Gagal decode usersCursor ke objek slice", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.UserResponseList{}, apiErr
	}

	return users, nil
}

//CheckEmailAvailable melakukan pengecekan apakah alamat email sdh terdaftar di database
//jika ada akan return false ,yang artinya email tidak available
func (u *userDao) CheckIDAvailable(userID string) (bool, rest_err.APIError) {

	coll := db.Db.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOne()
	opts.SetProjection(bson.M{keyUserID: 1})

	var user dto.UserResponse

	if err := coll.FindOne(ctx, bson.M{keyUserID: strings.ToUpper(userID)}, opts).Decode(&user); err != nil {

		if err == mongo.ErrNoDocuments {
			return true, nil
		}

		logger.Error("Gagal mendapatkan user dari database", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan user dari database", err)
		return false, apiErr
	}

	apiErr := rest_err.NewBadRequestError("ID tidak tersedia")
	return false, apiErr
}

//EditUser mengubah user, memerlukan timestamp int64 agar lebih safety pada saat pengeditan oleh dua user
func (u *userDao) EditUser(userID string, userRequest dto.UserEditRequest) (*dto.UserResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyUserID:        strings.ToUpper(userID),
		keyUserTimeStamp: userRequest.TimestampFilter,
	}
	update := bson.M{
		"$set": bson.M{
			keyUserName:      userRequest.Name,
			keyUserRoles:     userRequest.Roles,
			keyUserTimeStamp: time.Now().Unix(),
		},
	}

	var user dto.UserResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("User tidak diupdate karena ID atau timestamp tidak valid")
		}

		logger.Error("Gagal mendapatkan user dari database", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan user dari database", err)
		return nil, apiErr
	}

	return &user, nil
}

//DeleteUser menghapus user
func (u *userDao) DeleteUser(userID string) rest_err.APIError {
	coll := db.Db.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyUserID: strings.ToUpper(userID),
	}

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return rest_err.NewBadRequestError("User gagal dihapus, dokumen tidak ditemukan")
		}

		logger.Error("Gagal menghapus user dari database", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan user dari database", err)
		return apiErr
	}

	if result.DeletedCount == 0 {
		return rest_err.NewBadRequestError("User gagal dihapus, dokumen tidak ditemukan")
	}

	return nil
}

//PutAvatar hanya mengubah avatar berdasarkan filter email
func (u *userDao) PutAvatar(userID string, avatar string) (*dto.UserResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyUserID: strings.ToUpper(userID),
	}
	update := bson.M{
		"$set": bson.M{
			keyUserAvatar:    avatar,
			keyUserTimeStamp: time.Now().Unix(),
		},
	}

	var user dto.UserResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("User avatar gagal diupload, user dengan id %s tidak ditemukan", userID))
		}

		logger.Error("Gagal mendapatkan user dari database", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan user dari database", err)
		return nil, apiErr
	}

	return &user, nil
}

//ChangePassword merubah hash_pw dengan password baru sesuai masukan
func (u *userDao) ChangePassword(data dto.UserChangePasswordRequest) rest_err.APIError {
	coll := db.Db.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyUserID: strings.ToUpper(data.ID),
	}

	update := bson.M{
		"$set": bson.M{
			keyUserHashPw:    data.NewPassword,
			keyUserTimeStamp: time.Now().Unix(),
		},
	}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return rest_err.NewBadRequestError("Penggantian password gagal, ID salah")
		}

		logger.Error("Gagal mendapatkan user dari database (ChangePassword)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mengganti password user", err)
		return apiErr
	}

	if result.ModifiedCount == 0 {
		return rest_err.NewBadRequestError("Penggantian password gagal, kemungkinan ID salah")
	}

	return nil
}

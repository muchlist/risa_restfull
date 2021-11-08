package dto

// User struct lengkap dari document user di Mongodb
type User struct {
	ID        string   `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string   `json:"email" bson:"email"`
	Name      string   `json:"name" bson:"name"`
	Branch    string   `json:"branch" bson:"branch"`
	Roles     []string `json:"roles" bson:"roles"`
	Avatar    string   `json:"avatar" bson:"avatar"`
	HashPw    string   `json:"hash_pw,omitempty" bson:"hash_pw,omitempty"`
	Position  string   `json:"position" bson:"position"`
	Division  string   `json:"division" bson:"division"`
	FcmToken  string   `json:"fcm_token" bson:"fcm_token"`
	Timestamp int64    `json:"timestamp" bson:"timestamp"`
}

// UserResponseList tipe slice dari UserResponse
type UserResponseList []UserResponse

// UserResponse struct kembalian dari MongoDB dengan menghilangkan hashPassword
type UserResponse struct {
	ID        string   `json:"id" bson:"_id"`
	Email     string   `json:"email" bson:"email"`
	Name      string   `json:"name" bson:"name"`
	Branch    string   `json:"branch" bson:"branch"`
	Roles     []string `json:"roles" bson:"roles"`
	Avatar    string   `json:"avatar" bson:"avatar"`
	Position  string   `json:"position" bson:"position"`
	Division  string   `json:"division" bson:"division"`
	FcmToken  string   `json:"-" bson:"fcm_token"`
	Timestamp int64    `json:"timestamp" bson:"timestamp"`
}

// UserRequest input JSON untuk keperluan register, timestamp dapat diabaikan
type UserRequest struct {
	ID        string   `json:"id" bson:"_id"`
	Name      string   `json:"name" bson:"name"`
	Email     string   `json:"email" bson:"email"`
	Branch    string   `json:"branch" bson:"branch"`
	Roles     []string `json:"roles" bson:"roles"`
	Avatar    string   `json:"avatar" bson:"avatar"`
	Position  string   `json:"position" bson:"position"`
	Division  string   `json:"division" bson:"division"`
	Password  string   `json:"password" bson:"password"`
	Timestamp int64    `json:"timestamp" bson:"timestamp"`
}

// UserEditRequest input JSON oleh admin untuk mengedit user
type UserEditRequest struct {
	Name            string   `json:"name" bson:"name"`
	Branch          string   `json:"branch" bson:"branch"`
	Roles           []string `json:"roles" bson:"roles"`
	Position        string   `json:"position" bson:"position"`
	Division        string   `json:"division" bson:"division"`
	TimestampFilter int64    `json:"timestamp_filter" bson:"timestamp"`
}

// UserUpdateFcmRequest input fcm dari firebase client
type UserUpdateFcmRequest struct {
	FcmToken string `json:"fcm_token" bson:"fcm_token"`
}

// UserLoginRequest input JSON oleh client untuk keperluan login
type UserLoginRequest struct {
	ID       string `json:"id" bson:"_id"`
	Password string `json:"password" bson:"password"`
	Limit    int    `json:"limit"`
}

// UserChangePasswordRequest struck untuk keperluan change password dan reset password
// pada reset password hanya menggunakan NewPassword dan mengabaikan Password
type UserChangePasswordRequest struct {
	ID          string `json:"id" bson:"_id"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

// UserLoginResponse balikan user ketika sukses login dengan tambahan AccessToken
type UserLoginResponse struct {
	ID           string   `json:"id" bson:"_id"`
	Email        string   `json:"email" bson:"email"`
	Name         string   `json:"name" bson:"name"`
	Branch       string   `json:"branch" bson:"branch"`
	Roles        []string `json:"roles" bson:"roles"`
	Avatar       string   `json:"avatar" bson:"avatar"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Expired      int64    `json:"expired"`
}

type UserRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	Limit        int    `json:"limit"`
}

// UserRefreshTokenResponse mengembalikan token dengan claims yang
// sama dengan token sebelumnya dengan expired yang baru
type UserRefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	Expired     int64  `json:"expired"`
}

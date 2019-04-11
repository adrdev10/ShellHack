// Here, we register the Pusher client
var client = pusher.Client{
	AppId:   "600556",
	Key:     "1791bdb1a90b21e2e8ac",
	Secret:  "9c55dc8ac436be5be3db",
	Cluster: "us2",
	Secure:  true,
}


// Here, we define a user struct
type user struct {
	Username    string `json:"username" xml:"username" form:"username" query:"username"`
	Email       string `json:"email" xml:"email" form:"email" query:"email"`
	PhoneNumber string `json:"phonenumber"`
}
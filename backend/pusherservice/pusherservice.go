package pusherservice


// Here, we check if a user is logged in
func IsUserLoggedIn(rw http.ResponseWriter, req *http.Request) {
	if loggedInUser.Username != "" && loggedInUser.Email != "" {
		json.NewEncoder(rw).Encode(loggedInUser)
	} else {
		json.NewEncoder(rw).Encode("false")
	}
}



// -------------------------------------------------------
// Here, we receive a new user's details in a POST request and
// bind it to an instance of the user struct, we further use this
// user instance to check if a user is logged in or not
// -------------------------------------------------------
func NewUser(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &loggedInUser)
	if err != nil {
		panic(err)
	}
	json.NewEncoder(rw).Encode(loggedInUser)
}

// -------------------------------------------------------
// Here, we authorize users so that they can subscribe to the presence channel
// -------------------------------------------------------
func PusherAuth(res http.ResponseWriter, req *http.Request) {
	params, _ := ioutil.ReadAll(req.Body)
	presenceData := pusher.MemberData{
		UserId: loggedInUser.Username,
		UserInfo: map[string]string{
			"email": loggedInUser.Email,
		},
	}
	response, err := client.AuthenticatePresenceChannel(params, presenceData)

	if err != nil {
		panic(err)
	}
	fmt.Fprintf(res, string(response))
}

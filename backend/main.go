package main

// Here, we import the required packages (including Pusher)
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	pusher "github.com/pusher/pusher-http-go"
	"github.com/urfave/negroni"
)



var value int

// Here, we create a global user variable to hold user details for a session
var loggedInUser user

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
func pusherAuth(res http.ResponseWriter, req *http.Request) {

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

func createEmail(r http.ResponseWriter, w *http.Request) {
	var accountSid = "AC47589820a42731a6c3c4b00e4dc942f9"
	var AuthToken = "89806f248c7bfbd65a655a1b25fb6c31"
	var urlStr = "https://api.twilio.com/2010-04-01/Accounts/"

	strConv := fmt.Sprintf("Welcome %s, Your have chosen the perfect restaurant with the perfect person/s. Soon you will get to set the address and time. Thank You", loggedInUser.Username)

	msgData := url.Values{}
	msgData.Set("To", "+18135855228")
	msgData.Set("From", "+18505338424")
	msgData.Set("Body", strConv)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}
}

func main() {
	port := os.Getenv("PORT")
	newMux := http.NewServeMux()
	// Serve the static files and templates from the static directory
	newMux.Handle("/", http.FileServer(http.Dir("./static")))

	// Here, we determine if a user is logged in or not.
	newMux.HandleFunc("/isLoggedIn", isUserLoggedIn)

	// -------------------------------------------------------
	// Listen on these routes for new user registration and user authorization,
	// thereafter, handle each request using the matching handler function.
	// -------------------------------------------------------
	newMux.HandleFunc("/new/user", NewUser)
	newMux.HandleFunc("/success/send", createEmail)
	newMux.HandleFunc("/pusher/auth", pusherAuth)
	newMux.HandleFunc("/search/food", func(w http.ResponseWriter, r *http.Request) {
		var res *http.Response
		if r.Method == "GET" {
			city := r.FormValue("city")
			fmt.Println(city)
			strReq := fmt.Sprintf("https://developers.zomato.com/api/v2.1/locations?query=%s", url.QueryEscape(city))
			req, err := http.NewRequest(http.MethodGet, strReq, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			req.Header.Set("Accept", " application/json")
			req.Header.Set("user-key", "966b80a56bf2c010b2aef1a1f9c1b772")

			clients := http.Client{}
			if res, err = clients.Do(req); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			defer res.Body.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			var location Location

			data, _ := ioutil.ReadAll(res.Body)

			if err := json.Unmarshal(data, &location); err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			fmt.Println(location.LocationData)

			var locData = location.LocationData

			value = location.LocationData[0].EntityId

			encoder := json.NewEncoder(w)
			if err := encoder.Encode(locData); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		}
	})

	newMux.HandleFunc("/search/food/rest", func(w http.ResponseWriter, r *http.Request) {
		var res *http.Response
		var err error
		fmt.Println(value)
		stringValue := fmt.Sprintf("%v", value)
		strReq := fmt.Sprintf("https://developers.zomato.com/api/v2.1/search?entity_id=%s&entity_type=city&count=200", url.QueryEscape(stringValue))
		fmt.Println(strReq)
		req, err := http.NewRequest(http.MethodGet, strReq, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Accept", " application/json")
		req.Header.Set("user-key", "966b80a56bf2c010b2aef1a1f9c1b772")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		clients := http.Client{}
		if res, err = clients.Do(req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		defer res.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var rest Restaurants

		data, _ := ioutil.ReadAll(res.Body)

		if err := json.Unmarshal(data, &rest); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		fmt.Println(rest.Restaurant)

		var restData = rest.Restaurant

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(restData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	newMux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		var res *http.Response
		var err error
		req, err := http.NewRequest(http.MethodGet, "https://developers.zomato.com/api/v2.1/cuisines?city_id=280", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Accept", " application/json")
		req.Header.Set("user-key", "966b80a56bf2c010b2aef1a1f9c1b772")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		clients := http.Client{}
		if res, err = clients.Do(req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		defer res.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var cuis Cuisines

		data, _ := ioutil.ReadAll(res.Body)

		if err := json.Unmarshal(data, &cuis); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		fmt.Println(cuis.Cuisines)

		var cui = cuis.Cuisines

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(cui); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(newMux)

	// Start executing the application on port 8090
	http.ListenAndServe(":"+port, n)
}

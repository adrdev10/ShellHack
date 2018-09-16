package main

// Here, we import the required packages (including Pusher)
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	pusher "github.com/pusher/pusher-http-go"
	"github.com/urfave/negroni"
)

// Here, we register the Pusher client
var client = pusher.Client{
	AppId:   "600556",
	Key:     "1791bdb1a90b21e2e8ac",
	Secret:  "9c55dc8ac436be5be3db",
	Cluster: "us2",
	Secure:  true,
}

type Cuisines struct {
	Cuisines []Cui `json:"cuisines"`
}
type Cui struct {
	Cuisine Cuisine `json:"cuisine"`
}
type Cuisine struct {
	Cuisine_id   int    `json:"cuisine_id"`
	Cuisine_name string `json:"cuisine_name"`
}

type Location struct {
	LocationData []Data `json:"location_suggestions"`
}

type Data struct {
	EntityName string `json:"entity_type"`
	EntityId   int    `json:"entity_id"`
}

type Restaurants struct {
	Restaurant []Restaurant `json:"restaurants"`
}

type Restaurant struct {
	Rest Rest `json:"restaurant"`
}

type Rest struct {
	Name          string             `json:"name"`
	LocationRest  LocationRestaurant `json:"location"`
	Cuisine       string             `json:"cuisines"`
	AvgCostforTwo int                `json:"average_cost_for_two"`
	UserRating    UserRating         `json:"user_rating"`
}

type LocationRestaurant struct {
	Lat          string `json:"latitude"`
	Long         string `json:"longitude"`
	CuiseneLocal string `json:"locality_verbose"`
}

type UserRating struct {
	Rating     string `json:"aggregate_rating"`
	RatingText string `json:"rating_text"`
	Votes      string `json:"votes"`
}

var value int

// Here, we define a user struct
type user struct {
	Username string `json:"username" xml:"username" form:"username" query:"username"`
	Email    string `json:"email" xml:"email" form:"email" query:"email"`
}

// Here, we create a global user variable to hold user details for a session
var loggedInUser user

// Here, we check if a user is logged in
func isUserLoggedIn(rw http.ResponseWriter, req *http.Request) {
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

func main() {
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
	n.Run(":8080")
}

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
_
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
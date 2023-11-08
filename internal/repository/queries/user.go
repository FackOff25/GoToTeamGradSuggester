package queries

import (
	"fmt"
	"strconv"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
)

const (
	addUserQuery     = "INSERT INTO users (id) VALUES ($1);"
	getUserByIdQuery = `SELECT id, 
							username, 
							amusement_park_rating,
							aquarium_rating,
							art_gallery_rating,
							bar_rating,
							mosque_rating,
							movie_theater_rating,
							bowling_alley_rating,
							museum_rating,
							cafe_rating,
							night_club_rating,
							park_rating,
							casino_rating,
							cemetery_rating,
							church_rating,
							city_hall_rating,
							restaurant_rating,
							shopping_mall_rating,
							stadium_rating,
							synagogue_rating,
							tourist_attraction_rating,
							hindu_temple_rating,
							zoo_rating
						FROM users WHERE id = $1;`
)

func (q *Queries) GetUser(id string) (*domain.User, error) {
	row := q.Pool.QueryRow(q.Ctx, getUserByIdQuery, id)

	user := domain.User{PlaceTypePreferences: make(map[string]float64)}
	var username *string

	var amusement_park_rating, aquarium, art_gallery, bar, mosque, movie_theater, bowling_alley, 
	museum,	cafe, night_club, park, casino,	cemetery, church, city_hall, restaurant, shopping_mall, 
	stadium, synagogue,	tourist_attraction,	hindu_temple, zoo string

	err := row.Scan(
		&user.Id, 
		&username,
		&amusement_park_rating,
		&aquarium,
		&art_gallery,
		&bar,
		&mosque,
		&movie_theater,
		&bowling_alley,
		&museum,
		&cafe,
		&night_club,
		&park,
		&casino,
		&cemetery,
		&church,
		&city_hall,
		&restaurant,
		&shopping_mall,
		&stadium,
		&synagogue,
		&tourist_attraction,
		&hindu_temple,
		&zoo,
	)

	if username != nil {
		user.Username = *username
	}
	
	n, err := strconv.ParseFloat(amusement_park_rating, 64)
	user.PlaceTypePreferences["amusement_park_rating"] = n

	n, err = strconv.ParseFloat(aquarium, 64)
	user.PlaceTypePreferences["aquarium"] = n

	n, err = strconv.ParseFloat(art_gallery, 64)
	user.PlaceTypePreferences["art_gallery"] = n

	n, err = strconv.ParseFloat(bar, 64)
	user.PlaceTypePreferences["bar"] = n

	n, err = strconv.ParseFloat(mosque, 64)
	user.PlaceTypePreferences["mosque"] = n

	n, err = strconv.ParseFloat(movie_theater, 64)
	user.PlaceTypePreferences["movie_theater"] = n

	n, err = strconv.ParseFloat(bowling_alley, 64)
	user.PlaceTypePreferences["bowling_alley"] = n

	n, err = strconv.ParseFloat(cafe, 64)
	user.PlaceTypePreferences["cafe"] = n

	n, err = strconv.ParseFloat(night_club, 64)
	user.PlaceTypePreferences["night_club"] = n

	n, err = strconv.ParseFloat(park, 64)
	user.PlaceTypePreferences["park"] = n

	n, err = strconv.ParseFloat(casino, 64)
	user.PlaceTypePreferences["casino"] = n
	
	n, err = strconv.ParseFloat(cemetery, 64)
	user.PlaceTypePreferences["cemetery"] = n

	n, err = strconv.ParseFloat(church, 64)
	user.PlaceTypePreferences["church"] = n

	n, err = strconv.ParseFloat(city_hall, 64)
	user.PlaceTypePreferences["city_hall"] = n

	n, err = strconv.ParseFloat(restaurant, 64)
	user.PlaceTypePreferences["restaurant"] = n

	n, err = strconv.ParseFloat(shopping_mall, 64)
	user.PlaceTypePreferences["shopping_mall"] = n

	n, err = strconv.ParseFloat(stadium, 64)
	user.PlaceTypePreferences["stadium"] = n

	n, err = strconv.ParseFloat(synagogue, 64)
	user.PlaceTypePreferences["synagogue"] = n

	n, err = strconv.ParseFloat(tourist_attraction, 64)
	user.PlaceTypePreferences["tourist_attraction"] = n

	n, err = strconv.ParseFloat(hindu_temple, 64)
	user.PlaceTypePreferences["hindu_temple"] = n

	n, err = strconv.ParseFloat(zoo, 64)
	user.PlaceTypePreferences["zoo"] = n

	if err != nil {
		return nil, err
	}

	fmt.Println(user.PlaceTypePreferences)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (q *Queries) AddUser(id string) error {
	_, err := q.Pool.Query(q.Ctx, addUserQuery, id)
	return err
}

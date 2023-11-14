package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/FackOff25/GoToTeamGradSuggester/internal/domain"
	"github.com/FackOff25/GoToTeamGradSuggester/internal/usecase"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (pc *Controller) GetRoute(c echo.Context) error {
	defer c.Request().Body.Close()

	uuid, ok := c.Request().Header[uuidHeader]
	if !ok {
		return echo.ErrUnauthorized
	}

	u, _ := pc.Usecase.GetUser(uuid[0])
	if u == nil {
		return echo.ErrUnauthorized
	}

	var requestBody domain.RouteReq
	err := json.NewDecoder(c.Request().Body).Decode(&requestBody)
	if err != nil {
		return echo.ErrBadRequest
	}

	Greq, err := pc.Usecase.PrepareGreq(&requestBody)
	if err != nil {
		return echo.ErrBadRequest
	}

	route, err := pc.GetRouteFromG(Greq)
	if err != nil {
		return echo.ErrInternalServerError
	}

	resBody, err := pc.Usecase.GetRoute(route, Greq.TravelMode)
	if err != nil {
		log.Errorf("getting route error: %s", err.Error())
		return echo.ErrInternalServerError
	}

	resBodyBytes, err := json.Marshal(resBody)
	if err != nil {
		log.Errorf("marshaling error: %s", err.Error())
		return echo.ErrInternalServerError
	}

	return c.JSONBlob(http.StatusOK, resBodyBytes)
}

func (pc *Controller) GetRouteFromG(GreqBody *domain.GrouteRequest) (*domain.GrouteResp, error) {
	BytesGreqBody, err := json.Marshal(GreqBody)
	if err != nil {
		return nil, err
	}

	request := pc.Cfg.RoutesApiHost

	client := &http.Client{}
	Grequest, err := http.NewRequest(http.MethodPost, request, bytes.NewReader(BytesGreqBody))
	if err != nil {
		return nil, err
	}

	Grequest.Header.Set("Proxy-Header", "go-explore")
	Grequest.Header.Set("X-Goog-FieldMask", "routes.duration,routes.distanceMeters,routes.legs")

	gHttpResp, err := client.Do(Grequest)
	if err != nil {
		return nil, err
	}

	data, _ := io.ReadAll(gHttpResp.Body)
	var result domain.GrouteResp
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pc *Controller) SortPlaces(c echo.Context) error {
	defer c.Request().Body.Close()

	uuid, ok := c.Request().Header[uuidHeader]
	if !ok {
		return echo.ErrUnauthorized
	}

	u, _ := pc.Usecase.GetUser(uuid[0])
	if u == nil {
		return echo.ErrUnauthorized
	}

	var requestBody domain.SortPlacesReq
	err := json.NewDecoder(c.Request().Body).Decode(&requestBody)
	if err != nil {
		return echo.ErrBadRequest
	}

	path := make([]domain.ApiLocation, 0)

	path = append(path, requestBody.Start)
	path = append(path, requestBody.End)
	for _, v := range requestBody.Waypoints {
		path = append(path, v.Location)
	}

	path = usecase.SortPlacesForRoute(path)

	resp := domain.SortPlacesResp{
		Start:     path[0],
		End:       path[len(path)-1],
		Waypoints: path[1 : len(path)-1],
	}

	resBodyBytes := new(bytes.Buffer)
	encoder := json.NewEncoder(resBodyBytes)
	encoder.SetEscapeHTML(false)
	encoder.Encode(resp)

	return c.JSONBlob(http.StatusOK, resBodyBytes.Bytes())
}

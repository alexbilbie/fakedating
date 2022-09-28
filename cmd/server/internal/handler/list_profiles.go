package handler

import (
	"errors"
	"net/http"
	"strconv"

	"fakedating/pkg/middleware"
	"fakedating/pkg/model"
	"fakedating/pkg/payload"
	"fakedating/pkg/util"
)

func (h Handler) ListProfiles(w http.ResponseWriter, r *http.Request) {
	searchParams, validateErr := getSearchParametersFromRequest(r)
	if validateErr != nil {
		util.WriteErrorResponse("Invalid search parameters", validateErr, http.StatusBadRequest, w)
		return
	}

	// Fetch profiles
	matches, listErr := h.userRepository.ListMatches(
		middleware.GetUserIDFromContext(r.Context()),
		searchParams...,
	)
	if listErr != nil {
		util.WriteErrorResponse("Failed to list matches", listErr, http.StatusInternalServerError, w)
		return
	}

	util.WriteJSONResponse(
		payload.ListProfilesResponse{
			Matches: matches,
		},
		http.StatusOK,
		w,
	)
}

func getSearchParametersFromRequest(r *http.Request) ([]model.SearchParameterOpt, error) {
	var searchParams []model.SearchParameterOpt

	queryStringParams := r.URL.Query()

	// Look for filter parameters related to location
	latitude, latOk := queryStringParams["latitude"]
	longitude, longOk := queryStringParams["longitude"]
	radius, radiusOk := queryStringParams["radius"]

	if latOk && longOk && radiusOk {
		latVal, _ := strconv.ParseFloat(latitude[0], 64)
		longVal, _ := strconv.ParseFloat(longitude[0], 64)
		radiusVal, _ := strconv.ParseInt(radius[0], 10, 0)

		if latVal > 90 || latVal < -90 || longVal > 90 || longVal < -90 {
			return nil, errors.New("latitude/longitude out of bounds")
		}

		if radiusVal < 25 {
			radiusVal = 25
		}
		if radiusVal > 500 {
			radiusVal = 500
		}

		searchParams = append(
			searchParams, model.SearchWithLocationConstraint(latVal, longVal, uint(radiusVal)),
		)
	}

	// Look for filter parameters related to age
	ageLower, ageLowerOk := queryStringParams["age_lower"]
	ageUpper, ageUpperOk := queryStringParams["age_upper"]

	if ageLowerOk && ageUpperOk {
		ageLowerVal, _ := strconv.ParseInt(ageLower[0], 10, 0)
		ageUpperVal, _ := strconv.ParseInt(ageUpper[0], 10, 0)

		if ageLowerVal < 18 || ageLowerVal == 0 {
			ageLowerVal = 18
		}
		if ageLowerVal > ageUpperVal {
			return nil, errors.New("lower age constraint is greater than upper")
		}
		if ageUpperVal > 99 || ageUpperVal == 0 {
			ageUpperVal = 99
		}
		if ageUpperVal < ageLowerVal {
			return nil, errors.New("upper age constraint is less than lower")
		}

		searchParams = append(
			searchParams, model.SearchWithAgeConstraint(uint(ageLowerVal), uint(ageUpperVal)),
		)
	}

	// Offset
	offset, offsetOk := queryStringParams["offset"]
	if offsetOk {
		offsetVal, _ := strconv.ParseInt(offset[0], 10, 0)
		searchParams = append(searchParams, model.SearchWithOffset(uint(offsetVal)))
	}

	return searchParams, nil
}

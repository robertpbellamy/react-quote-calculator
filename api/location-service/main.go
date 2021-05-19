package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/go-playground/validator"
)

// Struct definitions
type Location struct {
    Postcode string `json:"postcode" validate:"required"`
}

type PostcodesIOAutoCompleteResult struct {
    Status int `json:"status" validate:"required"`
    Result []string `json:"result"`
}

type PostcodesIOLookupResult struct {
    Status int `json:"status" validate:"required"`
    Error string `json:"error"`
    Result PostcodesIOLookupResultInner `json:"result"`
}

type PostcodesIOLookupResultInner struct {
    Postcode string `json:"postcode"`
    Longitude float64 `json:"longitude"`
    Latitude float64 `json:"latitude"`
    AdminWard string `json:"admin_ward"`
    AdminDistrict string `json:"admin_district"`
}

// Global validator instance
var validate = validator.New()

// Setup router for our Location REST API
func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/api/location/autocomplete", locationAutocomplete).Methods("POST")
    myRouter.HandleFunc("/api/location/lookup", locationLookup).Methods("POST")
    log.Fatal(http.ListenAndServe(":80", myRouter))
}

// POST /api/location/autocomplete
func locationAutocomplete(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var location Location
    // Bad request if malformed JSON or does not meet structure validation rules.
    if json.Unmarshal(reqBody, &location) != nil ||
       validate.Struct(location) != nil {
         w.WriteHeader(http.StatusBadRequest)
    } else {
        // Success, query postcodes.io
        url := fmt.Sprintf("http://api.postcodes.io/postcodes/%s/autocomplete",
                            url.QueryEscape(location.Postcode))
        resp, err := http.Get(url)
        var respBody []byte
        if (err == nil) { respBody, err = ioutil.ReadAll(resp.Body) }
        if (err != nil) {
            w.WriteHeader(http.StatusServiceUnavailable)
        } else {
            // Received valid response, try to load into result object.
            var result PostcodesIOAutoCompleteResult
            // Fail if cannot load into object, received data is invalid,
            // or the remote API reports a status code which can't be handled.
            if json.Unmarshal(respBody, &result) != nil ||
               validate.Struct(result) != nil ||
               result.Status != 200 && result.Status != 404 {
                    w.WriteHeader(http.StatusInternalServerError)
            }  else if result.Status == 200 {
                // 200 success, process result object.
                if (len(result.Result) == 0) { result.Result = []string{} }
                output := struct { Postcodes []string `json:"postcodes"` }{ Postcodes: result.Result }
                w.WriteHeader(http.StatusOK)
                json.NewEncoder(w).Encode(output)
            } else {
                // 404: postcode not found, client made a bad request.
                w.WriteHeader(http.StatusBadRequest)
            }
        }
    }
}

// POST /api/location/lookup
func locationLookup(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var location Location
    validate := validator.New()
    // Bad request if malformed JSON or does not meet structure validation rules.
    if json.Unmarshal(reqBody, &location) != nil ||
       validate.Struct(location) != nil {
         w.WriteHeader(http.StatusBadRequest)
    } else {
        // Success, query postcodes.io
        url := fmt.Sprintf("http://api.postcodes.io/postcodes/%s",
                            url.QueryEscape(location.Postcode))
        resp, err := http.Get(url)
        var respBody []byte
        if (err == nil) { respBody, err = ioutil.ReadAll(resp.Body) }
        if (err != nil) {
            w.WriteHeader(http.StatusServiceUnavailable)
        } else {
            // Received valid response, try to load into result object.
            var result PostcodesIOLookupResult

            // Fail if cannot load into object, received data is invalid,
            // or the remote API reports a status code which can't be handled.
            if json.Unmarshal(respBody, &result) != nil ||
               validate.Struct(result) != nil ||
               result.Status != 200 && result.Status != 404 {
                    w.WriteHeader(http.StatusInternalServerError)
            } else if result.Status == 200 {
                // Success, process result object.
                townCity := "";
                if result.Result.AdminWard != "" {
                    townCity = result.Result.AdminWard
                }
                if result.Result.AdminDistrict != "" {
                    if townCity != "" { townCity += ", "}
                    townCity += result.Result.AdminDistrict
                }

                if townCity != "" { townCity += " "}
                townCity += result.Result.Postcode

                output := struct {
                    Latitude float64 `json:"latitude"`
                    Longitude float64 `json:"longitude"`
                    TownCity string `json:"town_city"`
                    }{ Latitude: result.Result.Latitude,
                        Longitude: result.Result.Longitude,
                        TownCity: townCity }
                w.WriteHeader(http.StatusOK)
                json.NewEncoder(w).Encode(output)
            } else {
                // 404: postcode not found, client made a bad request.
                w.WriteHeader(http.StatusBadRequest)
            }
        }
    }
}

func main() {
    handleRequests()
}

package main

import (
    "context"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/robertpbellamy/react-quote-calculator/api/db"
)

// Struct definitions
type PropertyType struct {
    Name string `json:"property_type"`
    Subtypes []PropertyTypeSubtype `json:"subtypes"`
    Floors []PropertyTypeFloor `json:"floors"`
}

type PropertyTypeSubtype struct {
    Name string `json:"name"`
    Multifloor bool `json:"multifloor"`
}

type PropertyTypeFloor struct {
    Name string `json:"name"`
    AboveGround bool `json:"above_ground"`
}

type RemovalService struct {
    Name string `json:"name"`
    ItemGroups []RemovalServiceItemGroup `json:"item_groups"`
}

type RemovalServiceItemGroup struct {
    Name string `json:"name"`
    Items []RemovalServiceItemGroupItem `json:"items"`
}

type RemovalServiceItemGroupItem struct {
    Name string `json:"name"`
    Minutes int `json:"minutes"`
}

// Setup router for our Location REST API
func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/api/quote/property-types", quotePropertyTypes).Methods("GET")
    myRouter.HandleFunc("/api/quote/removal-services", quoteRemovalServices).Methods("GET")
    log.Fatal(http.ListenAndServe(":80", myRouter))
}

// POST /api/quote/property-types
func quotePropertyTypes(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    client, err := db.GetClient(ctx)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    } else {
        database := client.Database("react_quote_calculator")
        collection := database.Collection("property_types")
        propertyTypes, err := collection.Find(ctx, bson.D{})
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
        } else {
            results := []PropertyType{}
            if err = propertyTypes.All(ctx, &results); err != nil {
                w.WriteHeader(http.StatusInternalServerError)
            } else {
                output := struct { PropertyTypes []PropertyType `json:"property_types"` } {
                                   PropertyTypes: results }
                w.WriteHeader(http.StatusOK)
                json.NewEncoder(w).Encode(output)
            }
        }
    }
}

// POST /api/quote/removal-services
func quoteRemovalServices(w http.ResponseWriter, r *http.Request) {
  ctx := context.Background()
  client, err := db.GetClient(ctx)
  if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
  } else {
      database := client.Database("react_quote_calculator")
      collection := database.Collection("removal_services")
      removalServices, err := collection.Find(ctx, bson.D{})
      if err != nil {
          w.WriteHeader(http.StatusInternalServerError)
      } else {
          results := []RemovalService{}
          if err = removalServices.All(ctx, &results); err != nil {
              log.Println(err)
              w.WriteHeader(http.StatusInternalServerError)
          } else {
              output := struct { RemovalServices []RemovalService `json:"removal_services"` } {
                                 RemovalServices: results }
              w.WriteHeader(http.StatusOK)
              json.NewEncoder(w).Encode(output)
          }
      }
  }
}

func main() {
    handleRequests()
}

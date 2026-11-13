package main

import (
	"log"

	"mapup-geofencing/internal/repository"
	"mapup-geofencing/internal/service"
	"mapup-geofencing/internal/handler"
	helper "mapup-geofencing/internal/helper"
	wsHandler "mapup-geofencing/internal/websocket"
	"net/http"
)

func main() {
	// ctx := context.Background()

	db, err := repository.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	alertHub := wsHandler.NewHub()

	repo := repository.NewGeofenceRepository(db)
	vehicleRepo := repository.NewVehicleRepository(db)
	vehicleLocationRepo := repository.NewVehicleLocationRepository(db)
	vehicleGeofenceStateRepo := repository.NewVehicleGeofenceStateRepository(db)
	violationRepo := repository.NewViolationRepository(db)
	alertConfigRepo := repository.NewAlertConfigRepository(db)
	geofenceRepo := repository.NewGeofenceRepository(db)


	// Service
	geofenceService := service.NewGeofenceService(repo)
	vehicleService := service.NewVehicleService(vehicleRepo)
	// vehicleLocationService := service.NewVehicleLocationService(vehicleLocationRepo)
	geofenceTransitionService := service.NewGeofenceTransitionService(vehicleGeofenceStateRepo)
	violationService := service.NewViolationService(violationRepo, alertConfigRepo, vehicleRepo, geofenceRepo, alertHub)
	alertConfigService := service.NewAlertConfigService(alertConfigRepo)
	vehicleLocationService := service.NewVehicleLocationService(vehicleLocationRepo, geofenceService, geofenceTransitionService, violationService)

	// Handler
	geofenceHandler := handler.NewGeofenceHandler(geofenceService)
	vehicleHandler := handler.NewVehicleHandler(vehicleService)
	vehicleLocationHandler := handler.NewVehicleLocationHandler(vehicleLocationService)
	violationHandler := handler.NewViolationHandler(violationService)
	webSocketHandler := handler.NewWebSocketHandler(alertHub)
	alertConfigHandler := handler.NewAlertConfigHandler(alertConfigService)

	// Routes
	http.HandleFunc("/geofences", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case http.MethodPost:
				geofenceHandler.CreateGeofence(w, r)

			case http.MethodGet:
				geofenceHandler.GetGeofences(w, r)

			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/vehicles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case http.MethodPost:
				vehicleHandler.CreateVehicle(w, r)

			case http.MethodGet:
				vehicleHandler.GetVehicles(w, r)

			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/vehicles/location", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case http.MethodPost:
				vehicleLocationHandler.CreateLocation(w, r)

			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/vehicles/location/{vehicle_id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case http.MethodGet:
				vehicleLocationHandler.GetLatestLocation(w, r)

			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})


	http.HandleFunc("/geofences/test-point", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			geofenceHandler.FindContainingGeofences(w, r)
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/violations/history", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case http.MethodGet:
				violationHandler.GetViolations(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ws/alerts", webSocketHandler.HandleAlerts)

	http.HandleFunc("/alerts/configure", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case http.MethodPost:
				alertConfigHandler.CreateConfig(w, r)

			default:
				http.Error(
					w,
					"method not allowed",
					http.StatusMethodNotAllowed,
				)
		}
	})

	http.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case http.MethodGet:
				alertConfigHandler.GetAlerts(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server running on :8080")

	if err := http.ListenAndServe(":8080", helper.EnableCORS(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}
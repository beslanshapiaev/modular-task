package main

import (
	"log"
	"modular-task/internal/api"
	"modular-task/internal/eventbus"
	"modular-task/internal/notifications"
	"modular-task/internal/products"
	"modular-task/internal/users"
	"net/http"
)

func main() {
	userRepo := users.NewUserRepository()
	productRepo := products.NewProductRepository()

	bus := eventbus.NewEventBus()

	userService := users.NewUserService(userRepo, bus)
	productService := products.NewProductService(productRepo, userService, bus)

	notifService := notifications.NewNotificationsService("notif", bus)
	notifService.Start()
	defer notifService.Stop()

	server := api.NewServer(userService, productService)
	server.Routes()

	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

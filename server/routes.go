package main

//import (
//	"Oauth/handler"
//	"Oauth/middleware"
//	"fmt"
//	"github.com/go-chi/chi"
//	"net/http"
//)
//
//type Server struct {
//	chi.Router
//}
//
//// SetupRoutes initializes the routes
//func SetupRoutes() *Server {
//	router := chi.NewRouter()
//	router.Use(middleware.CommonMiddlewares()...)
//
//	router.Route("/v1", func(audiophile chi.Router) {
//		audiophile.Route("/health", func(r chi.Router) {
//			r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
//				_, err := fmt.Fprintf(w, "health")
//				if err != nil {
//					return
//				}
//			})
//		})
//		audiophile.Post("/log-in", handler.Login)
//		audiophile.Post("/register", handler.Register)
//		audiophile.Route("/auth", func(auth chi.Router) {
//			auth.Use(middleware.AuthMiddleware)
//			auth.Put("/update-password", handler.UpdatePassword)
//			auth.Get("/products", handler.GetProducts)
//			auth.Post("/products", handler.BuyProduct)
//			auth.Post("/feedback", handler.CreateFeedback)
//			auth.Post("/log-out", handler.Logout)
//			auth.Route("/admin", func(auth chi.Router) {
//				auth.Use(middleware.AdminMiddleware)
//				auth.Post("/sell-product", handler.CreateProduct)
//			})
//		})
//	})
//	return &Server{router}
//}
//
//func (svc *Server) Run(port string) error {
//	return http.ListenAndServe(port, svc)
//}

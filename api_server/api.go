package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/omer1998/booking_api/controllers/doctor"
	"github.com/omer1998/booking_api/database"
	"github.com/omer1998/booking_api/services"
	"github.com/omer1998/booking_api/utils"
)

type ApiServer struct {
	listenAddr string
	db         database.Database
	cxt        context.Context
}

func NewApiServer(listenAddr string, db database.Database, cxt context.Context) *ApiServer {
	return &ApiServer{listenAddr: listenAddr, db: db, cxt: cxt}
}

func (s *ApiServer) Run() error {
	mux := mux.NewRouter()
	// mux.Use(LoggerMiddleWare)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	doctorRoutes := doctor.DoctorApi{Db: s.db, Auth: services.NewAuthenticationService(s.db)}
	// mux.HandleFunc("/doctor", utils.MakeHandleFunc(doctorRoutes.HandleAddDoctor)).Methods("POST")
	mux.HandleFunc("/doctor/all", utils.MakeHandleFunc(doctorRoutes.HandleGetAllDoctors)).Methods("GET")
	// mux.HandleFunc("/doctor/{email}", utils.CheckAuthState(utils.MakeHandleFunc(doctorRoutes.HandleGetDoctorByEmail))).Methods("GET")
	mux.HandleFunc("/doctor/auth/login", utils.MakeHandleFunc(doctorRoutes.HandleDoctorLogin)).Methods("POST")
	mux.HandleFunc("/doctor/auth/register", doctorRoutes.HandleDoctorRegister()).Methods("POST")
	mux.HandleFunc("/doctor/info", doctorRoutes.HandleDoctorInfo()).Methods("GET")
	// mux.HandleFunc("/doctor/{email}", utils.MakeHandleFunc(doctorRoutes.HandleDeleteDoctorByEmail)).Methods("DELETE")
	server := http.Server{
		Addr:    fmt.Sprintf("localhost%s", s.listenAddr),
		Handler: mux,
	}
	go func() {
		fmt.Printf("Server is running on port %s\n", s.listenAddr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Error starting server: %v", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		defer wg.Done()
		<-s.cxt.Done()
		shutdownCxt := context.Background()
		shutdownCxt, cancel := context.WithTimeout(shutdownCxt, 10*time.Second)
		fmt.Println("shutting down server ....")
		if err := server.Shutdown(shutdownCxt); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
		defer cancel()

	}()
	wg.Wait()
	return nil
}

package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/go-chi/chi"

	profile "github.com/DeepayanMallick/go-grpc/gunk/v1/profile"

	"google.golang.org/grpc"
)

const (
	ADDRESS = "localhost:50051"
)

type Server struct {
	Templates *template.Template
	Profile   profile.ProfileServiceClient
}

type ProfileTempData struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	Mobile      string
	Username    string
	Image       string
	FacebookURL string
	GithubURL   string
	TwitterURL  string
	Description string
}

func main() {
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}

	defer conn.Close()

	s := &Server{
		Templates: &template.Template{},
		Profile:   profile.NewProfileServiceClient(conn),
	}
	s.ParseTemplates()

	r := chi.NewRouter()
	r.Get("/", s.TestHandler)
	http.ListenAndServe(":5000", r)
}

func (s *Server) TestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := s.Profile.Profile(ctx, &profile.ProfileRequest{
		UserID: "1",
	})
	if err != nil {
		log.Println("Unable to get profile data")
		return
	}

	data := ProfileTempData{
		ID:          res.ID,
		FirstName:   res.FirstName,
		LastName:    res.LastName,
		Email:       res.Email,
		Mobile:      res.Mobile,
		Username:    res.Username,
		Image:       res.Image,
		FacebookURL: res.FacebookURL,
		GithubURL:   res.GithubURL,
		TwitterURL:  res.TwitterURL,
		Description: res.Description,
	}

	s.parseHomeTemplate(w, data)
}

func (s Server) parseHomeTemplate(w http.ResponseWriter, data any) {
	t := s.Templates.Lookup("home.html")
	if t == nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) ParseTemplates() error {
	templates := template.New("graphlogic-template").Funcs(template.FuncMap{
		"globalfunc": func(n string) string {
			return ""
		},
	}).Funcs(sprig.FuncMap())
	newFS := os.DirFS("assets/template")
	tmpl := template.Must(templates.ParseFS(newFS, "*.html"))
	if tmpl == nil {
		log.Fatalln("unable to parse templates")
	}

	s.Templates = tmpl
	return nil
}

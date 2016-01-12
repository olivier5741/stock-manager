package main

import (
	//"encoding/base64"
	"encoding/json"
	"github.com/auth0/go-jwt-middleware"
	"github.com/boltdb/bolt"
	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gocarina/gocsv"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"os"
)

type Prod struct {
	Id   string `csv:"id"`
	Name string `csv:"name"`
}

func SecuredPingHandler(w http.ResponseWriter, r *http.Request) {

	file, err := os.Open("stock.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	prods := []Prod{}

	if err := gocsv.UnmarshalFile(file, &prods); err != nil {
		panic(err)
	}

	db, err := bolt.Open("stock.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("products"))
		if err != nil {
			log.Print(err)
			return err
		}
		encoded, err := json.Marshal(prods)
		if err != nil {
			log.Print(err)
			return err
		}
		return b.Put([]byte("2015-01-01"), encoded)
	})

	decoded := []Prod{}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("products"))
		v := b.Get([]byte("2015-01-01"))
		err := json.Unmarshal(v, &decoded)

		if err != nil {
			log.Print(err)
			return err
		}
		return nil
	})

	re := render.New()
	re.JSON(w, http.StatusOK, decoded)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// log.Print(base64.StdEncoding.EncodeToString([]byte(os.Getenv("AUTH0_CLIENT_SECRET"))))

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("AUTH0_CLIENT_SECRET")), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	r := mux.NewRouter()

	r.Handle("/api/products", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(SecuredPingHandler)),
	))

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":3001")

}

package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(nil, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Serve the image on a HTTP endpoint
	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		// Find the document in the collection
		coll := client.Database("test").Collection("image")
		var doc bson.M
		err := coll.FindOne(nil, bson.M{}).Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}

		// Decode the base64 string and serve it as a JPEG image
		imgBase64 := doc["image"].(string)
		img, err := base64.StdEncoding.DecodeString(imgBase64)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(img)
	})

	// Serve the HTML page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
			<html>
			  <head>
				<title>Image from MongoDB</title>
			  </head>
			  <body>
				<img src="/image" alt="Image from MongoDB">
			  </body>
			</html>
		`
		fmt.Fprintln(w, html)
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

/*to save image to mongodb

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(nil, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Read image file and encode it as base64 string
	imgFile, err := ioutil.ReadFile("image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	imgBase64 := base64.StdEncoding.EncodeToString(imgFile)

	// Create a MongoDB document with the image as a string
	doc := bson.M{
		"image": imgBase64,
	}

	// Insert the document into a collection
	coll := client.Database("test").Collection("images")
	_, err = coll.InsertOne(nil, doc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Image saved in MongoDB")
}

*/

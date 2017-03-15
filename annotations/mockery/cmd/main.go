package main

import (
	"log"
	"os"

	"net/http"

	ann "github.com/Sotera/go_watchman/annotations"
	mock "github.com/Sotera/go_watchman/annotations/mockery"
	"github.com/Sotera/go_watchman/loogo"
	"gopkg.in/gin-gonic/gin.v1"
)

var eventsAPIRoot = os.Getenv("EVENTS_API_ROOT")

func main() {

	if eventsAPIRoot == "" {
		eventsAPIRoot = "http://localhost:3000/api/events"
	}

	r := gin.Default()

	parser := &loogo.HTTPRequestParser{
		Client: &loogo.HTTPClient{},
	}
	mocker := mock.Mockery{
		Parser: parser,
	}

	annotationsGroup := r.Group("/annotations")
	{
		annotationsGroup.GET("/refid/:refid", func(c *gin.Context) {
			// refid := c.Param("refid")
			// make time inputs look like standardized time
			options := ann.AnnotationOptions{
				StartTime:      c.Query("from_date") + "T00:00:00Z",
				EndTime:        c.Query("to_date") + "T00:00:00Z",
				AnnotationType: c.Query("type"),
				APIRoot:        eventsAPIRoot, // let's (mis)use this option.
			}
			annotations, err := mocker.GetAnnotations(options)
			if err != nil {
				log.Println("Error:", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			} else {
				c.JSON(http.StatusOK, &annotations)
			}
		})
	}

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}

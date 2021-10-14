package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	uuid "github.com/nu7hatch/gouuid"
)

type linkResponse struct {
    ID string `form:"ID" json:"ID" binding:"required"`
    OriginalLink string `form:"OriginalLink" json:"OriginalLink" binding:"required"`
    ShortLinkHash string `form:"ShortLinkHash" json:"ShortLinkHash" binding:"required"`
    NumberOfVisits int `form:"NumberOfVisits" json:"NumberOfVisits" binding:"required"`
}

type linkRequest struct {
    OriginalLink string `form:"OriginalLink" json:"OriginalLink" binding:"required"`
}

var client *redis.Client


func main() {

    client = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: "",
		DB: 0,
	})

    
    router := gin.Default()
    
    router.GET("/:hash", redirect)
    router.POST("/link", postLink)
    router.GET("/link/:hash", getLink)

    router.Run(":8080")
}


func redirect(c *gin.Context) {
    hash := c.Param("hash")
    val, err := client.Get(hash).Result()
    
    if err != nil {
        fmt.Println(err)
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "link not found"})
 
    }
    var linkData linkResponse
    
    json.Unmarshal([]byte(val), &linkData)

    //c.IndentedJSON(http.StatusOK, linkData)
    fmt.Println(linkData)
    c.Redirect(http.StatusMovedPermanently, linkData.OriginalLink )
    
    //Get link from DB
    //Increment number of visits
    //Fire event - visited
    //Redirect
    
    //c.IndentedJSON(http.StatusOK, albums)
}



func postLink(c *gin.Context) {
    
	
    //pong, err := client.Ping().Result()
    
    var link linkRequest

    // Call BindJSON to bind the received JSON to link

    if err := c.BindJSON(&link); err != nil {
        return
    }

    var linkDB linkResponse

    linkDB.NumberOfVisits = 0
    linkDB.OriginalLink = link.OriginalLink

    fmt.Println(link.OriginalLink)
    uuid, _ := uuid.NewV4()
    linkDB.ID = uuid.String()
    
    hasher := sha1.New()
    hasher.Write([]byte(link.OriginalLink))
    shortLinkHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

    linkDB.ShortLinkHash = shortLinkHash[0:7]
    
    b, _ := json.Marshal(linkDB)

    fmt.Println(linkDB)

    err := client.Set(linkDB.ShortLinkHash, b, 0).Err()
    if err != nil {
        panic(err)
    }


   c.IndentedJSON(http.StatusOK, linkDB)


}

func getLink(c *gin.Context) {
    hash := c.Param("hash")
    val, err := client.Get(hash).Result()

    if err == nil {
        

        var linkData linkResponse
        json.Unmarshal([]byte(val), &linkData)

        c.IndentedJSON(http.StatusOK, linkData)
        return
    }

    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "link not found"})

    
}



// getAlbums responds with the list of all albums as JSON.

/*
// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
    var newAlbum album

    // Call BindJSON to bind the received JSON to
    // newAlbum.
    if err := c.BindJSON(&newAlbum); err != nil {
        return
    }

    // Add the new album to the slice.
    albums = append(albums, newAlbum)
    c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
    id := c.Param("id")

    // Loop through the list of albums, looking for
    // an album whose ID value matches the parameter.
    for _, a := range albums {
        if a.ID == id {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
*/
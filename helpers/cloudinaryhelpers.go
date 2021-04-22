package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/admin/search"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

var CLOUDINARYNAME string
var CLOUDINARYKEY string
var CLOUDINARYSECRET string

const MAX_GAGS int = 20

func InitiCloudinary() {
	CLOUDINARYNAME = os.Getenv("CLOUDINARYNAME")
	CLOUDINARYKEY = os.Getenv("CLOUDINARYKEY")
	CLOUDINARYSECRET = os.Getenv("CLOUDINARYSECRET")

	log.Printf("%s:%s:%s", CLOUDINARYNAME, CLOUDINARYKEY, CLOUDINARYSECRET)
}

// Upload Image to Cloudinary
func UploadCloudinary(file, user_id string, tags []string) (string, error) {
	fmt.Println("Started Uploading file")
	cld, err := cloudinary.NewFromParams(CLOUDINARYNAME, CLOUDINARYKEY, CLOUDINARYSECRET)
	// cld, err := cloudinary.NewFromParams("doxkhafkv", "324734846843959", "P1QZWqGVSFaucLqTlWP-gLirdEY")
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("failed to upload gag")
	}
	var ctx = context.Background()
	tgs := api.CldAPIArray{}
	for _, t := range tags {
		tgs = append(tgs, t)
	}
	fmt.Printf("tags: %+v\n", tags)
	fmt.Printf("tags: %+v\n", tgs)
	res, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: user_id, Tags: tgs})
	if err != nil || res.PublicID == "" {
		fmt.Printf("Err: %+v\n", err.Error())
		fmt.Printf("Failed file finished: %+v\n", res)
		return "", errors.New("failed to upload gag")
	}
	fmt.Printf("Uploading file finished: %+v\n", res)
	return res.PublicID, nil
}

// Upload Image to Cloudinary
func FeedCloudinary(next_cursor string) (*admin.SearchResult, error) {
	fmt.Println("Started Searching for feed")
	cld, err := cloudinary.NewFromParams(CLOUDINARYNAME, CLOUDINARYKEY, CLOUDINARYSECRET)
	// cld, err := cloudinary.NewFromParams("doxkhafkv", "324734846843959", "P1QZWqGVSFaucLqTlWP-gLirdEY")
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("failed to create cloudinary connection")
	}
	var ctx = context.Background()

	searchQuery := search.Query{
		Expression: "resource_type:image AND type:upload",
		SortBy:     []search.SortByField{{"created_at": search.Descending}},
		MaxResults: MAX_GAGS,
		NextCursor: next_cursor,
		WithField:  []search.WithField{search.TagsField},
	}
	searchResult, err := cld.Admin.Search(ctx, searchQuery)

	if err != nil {
		log.Printf("Failed to search for gags, %v\n", err)
		return nil, err
	}

	log.Printf("Gags found: %v\n", searchResult.TotalCount)

	return searchResult, nil
}

// Get author gags from cloudinary
func GagsOfAuthor(user_id string, next_cursor string) (*admin.AssetsResult, error) {
	fmt.Println("Started Searching for author gags")
	cld, err := cloudinary.NewFromParams(CLOUDINARYNAME, CLOUDINARYKEY, CLOUDINARYSECRET)
	// cld, err := cloudinary.NewFromParams("doxkhafkv", "324734846843959", "P1QZWqGVSFaucLqTlWP-gLirdEY")
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("failed to create cloudinary connection")
	}
	var ctx = context.Background()
	prefix := user_id + "/"
	searchQuery := admin.AssetsParams{
		AssetType:  "image/upload",
		Prefix:     prefix,
		MaxResults: MAX_GAGS,
		NextCursor: next_cursor,
		Tags:       true,
	}

	searchResult, err := cld.Admin.Assets(ctx, searchQuery)

	if err != nil {
		log.Printf("Failed to get author for gags, %v\n", err.Error())
		return nil, err
	}

	return searchResult, nil
}

// Search for all images with tags
func GagsWithTags(tags string, next_cursor string) (*admin.SearchResult, error) {
	fmt.Println("Started Searching for tags")
	cld, err := cloudinary.NewFromParams(CLOUDINARYNAME, CLOUDINARYKEY, CLOUDINARYSECRET)
	// cld, err := cloudinary.NewFromParams("doxkhafkv", "324734846843959", "P1QZWqGVSFaucLqTlWP-gLirdEY")
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("failed to create cloudinary connection")
	}
	var ctx = context.Background()

	tagsQuery := ""
	tagSplits := strings.Split(tags, ",")
	for i, t := range tagSplits {
		tagsQuery += t
		if i != (len(tagSplits) - 1) {
			tagsQuery += " OR "
		}
	}

	searchQuery := search.Query{
		Expression: "resource_type:image AND type:upload AND ( " + tagsQuery + " )",
		SortBy:     []search.SortByField{{"created_at": search.Descending}},
		MaxResults: MAX_GAGS,
		NextCursor: next_cursor,
		WithField:  []search.WithField{search.TagsField},
	}
	log.Printf("%+v\n", searchQuery)
	searchResult, err := cld.Admin.Search(ctx, searchQuery)

	if err != nil {
		log.Printf("Failed to search for tags, %v\n", err)
		return nil, err
	}

	log.Printf("Gags found: %v\n", searchResult.TotalCount)

	return searchResult, nil
}

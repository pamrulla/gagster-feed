package helpers

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

// Upload Image to Cloudinary
func UploadCloudinary(file, user_id string, tags []string) (string, int, int, error) {
	fmt.Println("Started Uploading file")
	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUDINARYNAME"), os.Getenv("CLOUDINARYKEY"), os.Getenv("CLOUDINARYSECRET"))
	if err != nil {
		fmt.Println(err.Error())
		return "", 0, 0, errors.New("failed to upload gag")
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
		fmt.Printf("Failed file finished: %+v\n", res)
		return "", 0, 0, errors.New("failed to upload gag")
	}
	fmt.Printf("Uploading file finished: %+v\n", res)
	return res.PublicID, res.Width, res.Height, nil
}

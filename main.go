// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

const (
	availablePhotosPath = "/home/stefan/lib/daily-photo/photos/available"
	postedPhotosPath    = "/home/stefan/lib/daily-photo/photos/posted"
	failedPhotosPath    = "/home/stefan/lib/daily-photo/photos/failed"
)

type twitterCredentials struct {
	accessToken       string
	accessTokenSecret string
	consumerKey       string
	consumerSecret    string
}

func randomAvailablePhoto() (string, string, error) {
	paths, err := filepath.Glob(availablePhotosPath + "/*.txt")
	if err != nil {
		return "", "", err
	}

	if len(paths) == 0 {
		return "", "", nil
	}

	path := paths[rand.Intn(len(paths))]

	_, file := filepath.Split(path)
	name := file[0 : len(file)-4]

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", err
	}

	return name, string(b), nil
}

func availablePath(name string, ext string) string {
	return availablePhotosPath + "/" + name + "." + ext
}

func postedPath(name string, ext string) string {
	return postedPhotosPath + "/" + name + "." + ext
}

func failedPath(name string, ext string) string {
	return failedPhotosPath + "/" + name + "." + ext
}

//

func init() {
	rand.Seed(time.Now().UnixNano())
}

func postPhoto(path string, comment string, credentials twitterCredentials) (string, error) {
	twitter := anaconda.NewTwitterApiWithCredentials(credentials.accessToken, credentials.accessTokenSecret,
		credentials.consumerKey, credentials.consumerSecret)

	image, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	media, err := twitter.UploadMedia(base64.StdEncoding.EncodeToString(image))
	if err != nil {
		return "", err
	}

	tweet, err := postTweetWithMedia(twitter, comment, media)
	if err != nil {
		return "", err
	}

	// TODO Is there an API to get the URL?
	url := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IdStr)
	return url, nil
}

func moveToFailed(photo string) error {
	if err := os.Rename(availablePath(photo, "jpg"), postedPath(photo, "jpg")); err != nil {
		return err
	}
	if err := os.Rename(availablePath(photo, "txt"), postedPath(photo, "txt")); err != nil {
		return err
	}
	return nil
}

func moveToPosted(photo string) error {
	if err := os.Rename(availablePath(photo, "jpg"), postedPath(photo, "jpg")); err != nil {
		return err
	}
	if err := os.Rename(availablePath(photo, "txt"), postedPath(photo, "txt")); err != nil {
		return err
	}
	return nil
}

func postTweetWithMedia(twitter *anaconda.TwitterApi, status string, media anaconda.Media) (anaconda.Tweet, error) {
	values := url.Values{}
	values.Set("media_ids", strconv.FormatInt(media.MediaID, 10))
	return twitter.PostTweet(status, values)
}

func main() {
	photo, comment, err := randomAvailablePhoto()
	if err != nil {
		log.Fatal("Could not get a random photo: ", err)
	}

	if photo == "" {
		log.Println("No photos available")
		return
	}

	log.Printf("Found photo <%s> with comment <%s>\n", photo, comment)

	credentials := twitterCredentials{
		accessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		accessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		consumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		consumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
	}

	url, err := postPhoto(availablePath(photo, "jpg"), comment, credentials)
	if err != nil {
		log.Println("Could not post photo: ", err)
		if err := moveToFailed(photo); err != nil {
			log.Fatal("Could not move photo to failed/: ", err)
		}
	}

	log.Println("Uploaded photo to ", url)

	if err := moveToPosted(photo); err != nil {
		log.Fatal("Could not move photo to posted/: ", err)
	}
}

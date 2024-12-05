package google_test

import (
	"context"
	"fmt"
	"github.com/guneyin/disgo/internal/google"
	"github.com/guneyin/disgo/internal/utils"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/oauth2"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestOAuth2(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	var tokenCh = make(chan *oauth2.Token)

	Convey("Test Auth", t, func() {
		gp, err := google.NewTestApi(ctx, t)
		So(err, ShouldBeNil)
		url := gp.InitAuth(false)
		So(url, ShouldNotBeEmpty)
		t.Logf("URL: %s", url)

		go google.ServeHTTP(ctx, gp, tokenCh)
		go utils.OpenURL(url)

		token := &oauth2.Token{}
		select {
		case token = <-tokenCh:
		case <-ctx.Done():
		}

		So(token, ShouldNotBeEmpty)
		t.Logf("Token: %-v", token)

		err = google.WriteToken(token)
		So(err, ShouldBeNil)
	})
}

func TestAbout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	Convey("Test About", t, func() {
		gp, err := google.NewTestApi(ctx, t)
		So(err, ShouldBeNil)
		user, err := gp.About()
		So(err, ShouldBeNil)
		So(user, ShouldNotBeEmpty)
		t.Logf("User: %-v", user)
	})
}

func TestFileList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	parentId := "1Q6EXwfWKxNJuNLHefiwBViVg4WrcvgmJ"

	Convey("Test File List", t, func() {
		gp, err := google.NewTestApi(ctx, t)
		So(err, ShouldBeNil)
		fileList, err := gp.FileList(google.MimeTypeNone, parentId)
		So(err, ShouldBeNil)
		So(fileList, ShouldNotBeEmpty)

		for i, file := range fileList.Files {
			t.Logf("%-5d %-15s %-40s %-100s", i+1, file.MimeType, file.Id, file.Name)
		}
	})
}

func TestCreateDirectory(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	Convey("Test Create FileList", t, func() {
		gp, err := google.NewTestApi(ctx, t)
		So(err, ShouldBeNil)
		parentId := "1dFNmPyQTuszo4dcYyRr6cDkebbWyazer"
		dirName := fmt.Sprintf("my-test-dir-%d", rand.IntN(100))

		dir, err := gp.CreateDirectory(dirName, parentId)
		So(err, ShouldBeNil)
		So(dir, ShouldNotBeEmpty)
		t.Logf("FileList ID: %-v", dir)
	})
}

func TestDeleteDirectory(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	dirId := "1cWoA6o1wBWtfHciMhNFl1c3rc5DI8i3P"

	Convey("Test Delete FileList", t, func() {
		gp, err := google.NewTestApi(ctx, t)
		So(err, ShouldBeNil)
		err = gp.DeleteDirectory(dirId)
		So(err, ShouldBeNil)
	})
}

func TestDownloadFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	fileId := "1ytNbgLmYvHcTZr6u8Fc_zLPx21OjYuwv"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	Convey("Test Download File", t, func() {
		gp, err := google.NewTestApi(ctx, t)

		out, err := os.Create("downloaded.zip")
		So(err, ShouldBeNil)
		defer out.Close()

		err = gp.DownloadFile(fileId, out)
		So(err, ShouldBeNil)
	})
}

func TestApi_UploadFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	const testFile = "testdata/icons8-golang-480.png"
	Convey("Test Upload File", t, func() {
		gp, err := google.NewTestApi(ctx, t)
		So(err, ShouldBeNil)

		fileName := filepath.Base(testFile)
		file, err := os.Open(testFile)
		So(err, ShouldBeNil)

		uploaded, err := gp.UploadFile(fileName, "1dFNmPyQTuszo4dcYyRr6cDkebbWyazer", file)
		So(err, ShouldBeNil)
		So(uploaded, ShouldNotBeEmpty)
	})
}

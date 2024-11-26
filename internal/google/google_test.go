package google_test

import (
	"context"
	"fmt"
	"github.com/guneyin/disgo/internal/google"
	"github.com/guneyin/disgo/internal/utils"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand/v2"
	"testing"
	"time"
)

func TestOAuth2(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	var tokenCh = make(chan string)

	Convey("Test Auth", t, func() {
		gp := google.NewTestApi(t)
		url := gp.InitAuth()
		So(url, ShouldNotBeEmpty)
		t.Logf("URL: %s", url)

		go google.ServeHTTP(ctx, gp, tokenCh)
		go utils.OpenURL(url)

		token := ""
		select {
		case token = <-tokenCh:
		case <-ctx.Done():
		}

		So(token, ShouldNotBeEmpty)
		t.Logf("Token: %s", token)

		err := google.WriteToken(token)
		So(err, ShouldBeNil)
	})
}

func TestAbout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	Convey("Test About", t, func() {
		gp := google.NewTestApi(t)

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

	parentId := "1Swd9MRtY0ON-xI1jwR-NQXNQhMe0_4eY"

	Convey("Test File List", t, func() {
		gp := google.NewTestApi(t)

		fileList, err := gp.FileList(google.MimeTypeFolder, parentId)
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

	Convey("Test Create FileList", t, func() {
		gp := google.NewTestApi(t)

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

	dirId := "1tFY3gHz6Yc66oP7MWuztKYqS8LerUg6-"

	Convey("Test Delete FileList", t, func() {
		gp := google.NewTestApi(t)

		err := gp.DeleteDirectory(dirId)
		So(err, ShouldBeNil)
	})
}

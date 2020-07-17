package httpx

import (
	"github.com/motclub/common/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGET(t *testing.T) {
	var reply struct {
		Identity struct {
			ID    string `json:"id"`
			Login string `json:"login"`
		} `json:"identity"`
		Permissions struct {
			Roles []string `json:"roles"`
		} `json:"permissions"`
	}
	if err := GET("https://run.mocky.io/v3/c6652490-278f-458a-8499-82694398e3f0", &reply); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "b06cd03f-75d0-413a-b94b-35e155444d70", reply.Identity.ID)
	assert.Equal(t, "John Doe", reply.Identity.Login)
	assert.Equal(t, "moderator", reply.Permissions.Roles[0])
}

func TestPOST(t *testing.T) {
	var args = struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}{
		Code: "he",
		Name: "llo",
	}
	var reply struct {
		Name     string `json:"name"`
		Status   string `json:"status"`
		Url      string `json:"url"`
		ThumbUrl string `json:"thumbUrl"`
	}
	if err := POST("https://run.mocky.io/v3/fd257431-155a-4a8e-9d49-4ca233067206", &args, &reply); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "xxx.png", reply.Name)
	assert.Equal(t, "done", reply.Status)
	assert.Equal(t, "https://example.com/a/jkjgkEfvpUPVyRjUImniVslZfWPnJuuZ.png", reply.Url)
	assert.Equal(t, "https://example.com/a/jkjgkEfvpUPVyRjUImniVslZfWPnJuuZ.png", reply.ThumbUrl)

	var reply2 interface{}
	if err := POST("https://run.mocky.io/v3/52d0c8c5-7121-4f42-afb6-2244402229fb", nil, &reply2); err != nil {
		t.Fatal(err)
	}
	var results2 []map[string]interface{}
	if data, err := json.STD().Marshal(reply2); err != nil {
		t.Fatal(err)
	} else {
		if err := json.STD().Unmarshal(data, &results2); err != nil {
			t.Fatal(err)
		}
	}
	assert.Equal(t, 2, len(results2))
	assert.Equal(t, "aaa.png", results2[0]["name"].(string))
	assert.Equal(t, "https://example.com/a/bbb.png", results2[1]["thumbUrl"].(string))
}

func TestPUT(t *testing.T) {
	var reply struct {
		Name     string `json:"name"`
		Status   string `json:"status"`
		Url      string `json:"url"`
		ThumbUrl string `json:"thumbUrl"`
	}
	if err := PUT("https://run.mocky.io/v3/fd257431-155a-4a8e-9d49-4ca233067206", nil, &reply); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "xxx.png", reply.Name)
	assert.Equal(t, "done", reply.Status)
	assert.Equal(t, "https://example.com/a/jkjgkEfvpUPVyRjUImniVslZfWPnJuuZ.png", reply.Url)
	assert.Equal(t, "https://example.com/a/jkjgkEfvpUPVyRjUImniVslZfWPnJuuZ.png", reply.ThumbUrl)
}

func TestPATCH(t *testing.T) {
	var args = struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}{
		Code: "he",
		Name: "llo",
	}
	var reply struct {
		Name     string `json:"name"`
		Status   string `json:"status"`
		Url      string `json:"url"`
		ThumbUrl string `json:"thumbUrl"`
	}
	if err := PATCH("https://run.mocky.io/v3/fd257431-155a-4a8e-9d49-4ca233067206", &args, &reply); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "xxx.png", reply.Name)
	assert.Equal(t, "done", reply.Status)
	assert.Equal(t, "https://example.com/a/jkjgkEfvpUPVyRjUImniVslZfWPnJuuZ.png", reply.Url)
	assert.Equal(t, "https://example.com/a/jkjgkEfvpUPVyRjUImniVslZfWPnJuuZ.png", reply.ThumbUrl)
}

func TestDELETE(t *testing.T) {
	var reply struct {
		Name     string `json:"name"`
		Status   string `json:"status"`
		Url      string `json:"url"`
		ThumbUrl string `json:"thumbUrl"`
	}
	if err := DELETE("https://run.mocky.io/v3/fd257431-155a-4a8e-9d49-4ca233067206", &reply); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "xxx.png", reply.Name)
	assert.Equal(t, "done", reply.Status)
	assert.Equal(t, "https://example.com/a/jkjgkEfvpUPVyRjUImniVslZfWPnJuuZ.png", reply.Url)
	assert.Equal(t, "https://example.com/a/jkjgkEfvpUPVyRjUImniVslZfWPnJuuZ.png", reply.ThumbUrl)
}

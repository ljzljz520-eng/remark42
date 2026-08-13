package api

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrivateCreateCommentContentType(t *testing.T) {
	ts, _, teardown := startupT(t)
	defer teardown()

	resp, err := post(t, ts.URL+"/api/v1/comment", `{"text":"test","locator":{"url":"https://radio-t.com/blah1","site":"remark42"}}`)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	// The frontend's fetcher returns the raw text instead of a parsed object when the
	// content type is not application/json, which leaves comment.user undefined and breaks
	// rendering with "can't access property 'block', user is undefined". Enforce the type
	// so the body is parsed as JSON by the client.
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	// The created comment must be a JSON object carrying a user sub-object; without it the
	// client throws on accessing comment.user.block and the form stays greyed out.
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	created := map[string]interface{}{}
	require.NoError(t, json.Unmarshal(body, &created), "response body should be valid JSON: %s", string(body))
	user, ok := created["user"].(map[string]interface{})
	require.True(t, ok, "response should contain a user object: %s", string(body))
	assert.NotEmpty(t, user["id"], "user should have an id")
}

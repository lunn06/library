//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	genrepb "github.com/lunn06/library/book/internal/api/proto/genre"
)

func TestGenreCreateWithNoBooks(t *testing.T) {
	const (
		testTitle = "TestGenreCreateWithNoBooks"
		testDescription
	)
	// Put author
	req := genrepb.CreateRequest{
		Title:       testTitle,
		Description: testDescription,
	}
	data, err := proto.Marshal(&req)
	require.NoError(t, err)

	resMsg, err := nc.Request(genrePutSubj, data, reqTimeout)
	require.NoError(t, err)

	var resp genrepb.CreateResponse
	err = proto.Unmarshal(resMsg.Data, &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode))
	//////////////

	// Check putted author
	getReq := genrepb.GetRequest{
		GenreId: resp.GenreId,
	}
	data, _ = proto.Marshal(&getReq)

	validateMsg, err := nc.Request(genreGetSubj, data, reqTimeout)
	require.NoError(t, err)

	var getResp genrepb.GetResponse
	err = proto.Unmarshal(validateMsg.Data, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(getResp.StatusCode))
	assert.Equal(t, testTitle, getResp.Title)
	assert.Equal(t, testDescription, getResp.Description)
}

func TestGenreUpdateWithNoBooks(t *testing.T) {
	const (
		testTitle = "TestGenreUpdateWithNoBooks"
		testDescription

		testUpdatedTitle = "UpdatedTestGenreUpdateWithNoBooks"
		testUpdatedDescription
	)
	// Put author
	putReq := genrepb.CreateRequest{
		Title:       testTitle,
		Description: testDescription,
	}
	var putResp genrepb.CreateResponse
	err := request(genrePutSubj, &putReq, &putResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(putResp.StatusCode))
	//////////////

	// Update author
	updateReq := genrepb.UpdateRequest{
		Id:          putResp.GenreId,
		Title:       testUpdatedTitle,
		Description: testUpdatedDescription,
	}
	var updateResp genrepb.EmptyResponse
	err = request(genreUpdateSubj, &updateReq, &updateResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(updateResp.StatusCode))
	//////////////

	// Check updated author
	getReq := genrepb.GetRequest{
		GenreId: putResp.GenreId,
	}
	var getResp genrepb.GetResponse
	err = request(genreGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(getResp.StatusCode))
	assert.Equal(t, testUpdatedTitle, getResp.Title)
	assert.Equal(t, testUpdatedDescription, getResp.Description)
}

func TestGenreDelete(t *testing.T) {
	const (
		testTitle = "TestGenreDelete"
		testDescription
	)
	// Put author
	putReq := genrepb.CreateRequest{
		Title:       testTitle,
		Description: testDescription,
	}
	var putResp genrepb.CreateResponse
	err := request(genrePutSubj, &putReq, &putResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(putResp.StatusCode))
	//////////////

	// Check putted author
	getReq := genrepb.GetRequest{
		GenreId: putResp.GenreId,
	}
	var getResp genrepb.GetResponse
	err = request(genreGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(getResp.StatusCode))
	assert.Equal(t, testTitle, getResp.Title)
	assert.Equal(t, testDescription, getResp.Description)
	//////////////

	// Delete author
	deleteReq := genrepb.DeleteRequest{
		GenreId: putResp.GenreId,
	}
	var deleteResp genrepb.EmptyResponse
	err = request(genreDeleteSubj, &deleteReq, &deleteResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(deleteResp.StatusCode))
	//////////////

	// Check deleted author not found
	getReq = genrepb.GetRequest{
		GenreId: putResp.GenreId,
	}
	getResp = genrepb.GetResponse{}
	err = request(genreGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, int(getResp.StatusCode))
}

//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authorpb "github.com/lunn06/library/book/internal/api/proto/author"
)

func TestAuthorCreateWithNoBooks(t *testing.T) {
	const (
		testName = "TestAuthorCreateWithNoBooks"
		testDescription
	)
	// Put author
	req := authorpb.CreateRequest{
		Name:        testName,
		Description: testDescription,
	}
	var resp authorpb.CreateResponse
	err := request(authorPutSubj, &req, &resp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(resp.StatusCode))
	//////////////

	// Check putted author
	getReq := authorpb.GetRequest{
		AuthorId: resp.AuthorId,
	}
	var getResp authorpb.GetResponse
	err = request(authorGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(getResp.StatusCode))
	assert.Equal(t, testName, getResp.Name)
	assert.Equal(t, testDescription, getResp.Description)
}

func TestAuthorUpdateWithNoBooks(t *testing.T) {
	const (
		testName = "TestAuthorUpdateWithNoBooks"
		testDescription

		testUpdatedName = "UpdatedTestAuthorUpdateWithNoBooks"
		testUpdatedDescription
	)
	// Put author
	putReq := authorpb.CreateRequest{
		Name:        testName,
		Description: testDescription,
	}
	var putResp authorpb.CreateResponse
	err := request(authorPutSubj, &putReq, &putResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(putResp.StatusCode))
	//////////////

	// Update author
	updateReq := authorpb.UpdateRequest{
		Id:          putResp.AuthorId,
		Name:        testUpdatedName,
		Description: testUpdatedDescription,
	}
	var updateResp authorpb.EmptyResponse
	err = request(authorUpdateSubj, &updateReq, &updateResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(updateResp.StatusCode))
	//////////////

	// Check updated author
	getReq := authorpb.GetRequest{
		AuthorId: putResp.AuthorId,
	}
	var getResp authorpb.GetResponse
	err = request(authorGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(getResp.StatusCode))
	assert.Equal(t, testUpdatedName, getResp.Name)
	assert.Equal(t, testUpdatedDescription, getResp.Description)
}

func TestAuthorDelete(t *testing.T) {
	const (
		testName = "TestAuthorDelete"
		testDescription
	)
	// Put author
	putReq := authorpb.CreateRequest{
		Name:        testName,
		Description: testDescription,
	}
	var putResp authorpb.CreateResponse
	err := request(authorPutSubj, &putReq, &putResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(putResp.StatusCode))
	//////////////

	// Check putted author
	getReq := authorpb.GetRequest{
		AuthorId: putResp.AuthorId,
	}
	var getResp authorpb.GetResponse
	err = request(authorGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(getResp.StatusCode))
	assert.Equal(t, testName, getResp.Name)
	assert.Equal(t, testDescription, getResp.Description)
	//////////////

	// Delete author
	deleteReq := authorpb.DeleteRequest{
		AuthorId: putResp.AuthorId,
	}
	var deleteResp authorpb.EmptyResponse
	err = request(authorDeleteSubj, &deleteReq, &deleteResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(deleteResp.StatusCode))
	//////////////

	// Check deleted author not found
	getReq = authorpb.GetRequest{
		AuthorId: putResp.AuthorId,
	}
	getResp = authorpb.GetResponse{}
	err = request(authorGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, int(getResp.StatusCode))
}

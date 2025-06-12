//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	bookpb "github.com/lunn06/library/bookinfo/internal/api/proto/book"
)

func TestBookCreateWithNoAdditional(t *testing.T) {
	const (
		testUserID = 999
		testTitle  = "TestBookCreateWithNoAdditional"
		testDescription
		testBookUrl
	)
	// Put gateway
	req := bookpb.CreateRequest{
		UserId:      testUserID,
		Title:       testTitle,
		Description: testDescription,
		BookUrl:     testBookUrl,
	}
	data, err := proto.Marshal(&req)
	require.NoError(t, err)

	resMsg, err := nc.Request(bookPutSubj, data, reqTimeout)
	require.NoError(t, err)

	var resp bookpb.CreateResponse
	err = proto.Unmarshal(resMsg.Data, &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode))
	//////////////

	// Check putted gateway
	getReq := bookpb.GetRequest{
		BookId: resp.BookId,
	}
	data, _ = proto.Marshal(&getReq)

	validateMsg, err := nc.Request(bookGetSubj, data, reqTimeout)
	require.NoError(t, err)

	var getResp bookpb.GetResponse
	err = proto.Unmarshal(validateMsg.Data, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(getResp.StatusCode))
	assert.Equal(t, testUserID, int(getResp.UserId))
	assert.Equal(t, testTitle, getResp.Title)
	assert.Equal(t, testDescription, getResp.Description)
	assert.Equal(t, testBookUrl, getResp.BookUrl)
}

func TestBookUpdateWithNoAdditional(t *testing.T) {
	const (
		testUserID = 999
		testTitle  = "TestBookUpdateWithNoAdditional"
		testDescription
		testBookUrl

		testUpdatedUserID = 1000
		testUpdatedTitle  = "UpdatedTestBookUpdateWithNoAdditional"
		testUpdatedDescription
		testUpdatedBookUrl
	)
	// Put gateway
	putReq := bookpb.CreateRequest{
		UserId:      testUserID,
		Title:       testTitle,
		Description: testDescription,
		BookUrl:     testBookUrl,
	}
	var putResp bookpb.CreateResponse
	err := request(bookPutSubj, &putReq, &putResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(putResp.StatusCode))
	//////////////

	// Update gateway
	updateReq := bookpb.UpdateRequest{
		Id:          putResp.BookId,
		UserId:      testUpdatedUserID,
		Title:       testUpdatedTitle,
		Description: testUpdatedDescription,
		BookUrl:     testUpdatedBookUrl,
	}
	var updateResp bookpb.EmptyResponse
	err = request(bookUpdateSubj, &updateReq, &updateResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(updateResp.StatusCode))
	//////////////

	// Check updated gateway
	getReq := bookpb.GetRequest{
		BookId: putResp.BookId,
	}
	var getResp bookpb.GetResponse
	err = request(bookGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(getResp.StatusCode))
	assert.Equal(t, testUpdatedUserID, int(getResp.UserId))
	assert.Equal(t, testUpdatedTitle, getResp.Title)
	assert.Equal(t, testUpdatedDescription, getResp.Description)
	assert.Equal(t, testUpdatedBookUrl, getResp.BookUrl)
}

func TestBookDelete(t *testing.T) {
	const (
		testUserID = 999
		testTitle  = "TestBookDelete"
		testDescription
		testBookUrl
	)
	// Put gateway
	putReq := bookpb.CreateRequest{
		UserId:      testUserID,
		Title:       testTitle,
		Description: testDescription,
		BookUrl:     testBookUrl,
	}
	var putResp bookpb.CreateResponse
	err := request(bookPutSubj, &putReq, &putResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(putResp.StatusCode))
	//////////////

	// Check putted gateway
	getReq := bookpb.GetRequest{
		BookId: putResp.BookId,
	}
	var getResp bookpb.GetResponse
	err = request(bookGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(getResp.StatusCode))
	assert.Equal(t, testUserID, int(getResp.UserId))
	assert.Equal(t, testTitle, getResp.Title)
	assert.Equal(t, testDescription, getResp.Description)
	assert.Equal(t, testBookUrl, getResp.BookUrl)
	//////////////

	// Delete gateway
	deleteReq := bookpb.DeleteRequest{
		BookId: putResp.BookId,
	}
	var deleteResp bookpb.EmptyResponse
	err = request(bookDeleteSubj, &deleteReq, &deleteResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(deleteResp.StatusCode))
	//////////////

	// Check deleted gateway not found
	getReq = bookpb.GetRequest{
		BookId: putResp.BookId,
	}
	getResp = bookpb.GetResponse{}
	err = request(bookGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, int(getResp.StatusCode))
}

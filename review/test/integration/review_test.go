//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	reviewpb "github.com/lunn06/library/review/internal/api/proto/review"
)

func TestReviewCreate(t *testing.T) {
	const (
		testUserID = 999
		testBookID
		testTitle = "TestReviewCreate"
		testText
		testScore = 1
	)
	// Put review
	req := reviewpb.CreateRequest{
		UserId: testUserID,
		BookId: testBookID,
		Title:  testTitle,
		Text:   testText,
		Score:  testScore,
	}
	var resp reviewpb.CreateResponse
	err := request(reviewPutSubj, &req, &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode))
	//////////////

	// Check putted review
	getReq := reviewpb.GetRequest{
		ReviewId: resp.ReviewId,
	}
	var getResp reviewpb.GetResponse
	err = request(reviewGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(getResp.StatusCode))
	assert.Equal(t, testUserID, int(getResp.UserId))
	assert.Equal(t, testBookID, int(getResp.BookId))
	assert.Equal(t, testTitle, getResp.Title)
	assert.Equal(t, testText, getResp.Text)
	assert.Equal(t, testScore, int(getResp.Score))
}

func TestReviewUpdate(t *testing.T) {
	const (
		testUserID = 999
		testBookID
		testTitle = "TestReviewUpdate"
		testText
		testScore = 1

		testUpdatedTitle = "TestUpdatedReviewUpdate"
		testUpdatedText
		testUpdatedScore = 2
	)
	// Put review
	putReq := reviewpb.CreateRequest{
		UserId: testUserID,
		BookId: testBookID,
		Title:  testTitle,
		Text:   testText,
		Score:  testScore,
	}
	var putResp reviewpb.CreateResponse
	err := request(reviewPutSubj, &putReq, &putResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(putResp.StatusCode))
	//////////////

	// Update review
	updateReq := reviewpb.UpdateRequest{
		Id:    putResp.ReviewId,
		Title: testUpdatedTitle,
		Text:  testUpdatedText,
		Score: testUpdatedScore,
	}
	var updateResp reviewpb.EmptyResponse
	err = request(reviewUpdateSubj, &updateReq, &updateResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(updateResp.StatusCode))
	//////////////

	// Check updated review
	getReq := reviewpb.GetRequest{
		ReviewId: putResp.ReviewId,
	}
	var getResp reviewpb.GetResponse
	err = request(reviewGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(getResp.StatusCode))
	assert.Equal(t, testUserID, int(getResp.UserId))
	assert.Equal(t, testBookID, int(getResp.BookId))
	assert.Equal(t, testUpdatedTitle, getResp.Title)
	assert.Equal(t, testUpdatedText, getResp.Text)
	assert.Equal(t, testUpdatedScore, int(getResp.Score))
}

func TestReviewDelete(t *testing.T) {
	const (
		testUserID = 999
		testBookID
		testTitle = "TestReviewReview"
		testText
		testScore = 1
	)
	// Put review
	putReq := reviewpb.CreateRequest{
		UserId: testUserID,
		BookId: testBookID,
		Title:  testTitle,
		Text:   testText,
		Score:  testScore,
	}
	var putResp reviewpb.CreateResponse
	err := request(reviewPutSubj, &putReq, &putResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(putResp.StatusCode))
	//////////////

	// Check updated review
	getReq := reviewpb.GetRequest{
		ReviewId: putResp.ReviewId,
	}
	var getResp reviewpb.GetResponse
	err = request(reviewGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(getResp.StatusCode))
	assert.Equal(t, testUserID, int(getResp.UserId))
	assert.Equal(t, testBookID, int(getResp.BookId))
	assert.Equal(t, testTitle, getResp.Title)
	assert.Equal(t, testText, getResp.Text)
	assert.Equal(t, testScore, int(getResp.Score))
	//////////////

	// Delete review
	deleteReq := reviewpb.DeleteRequest{
		ReviewId: putResp.ReviewId,
	}
	var deleteResp reviewpb.EmptyResponse
	err = request(reviewDeleteSubj, &deleteReq, &deleteResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, int(deleteResp.StatusCode))
	//////////////

	// Check deleted review not found
	getReq = reviewpb.GetRequest{
		ReviewId: putResp.ReviewId,
	}
	getResp = reviewpb.GetResponse{}
	err = request(reviewGetSubj, &getReq, &getResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, int(getResp.StatusCode))
}

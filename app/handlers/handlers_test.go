package handlers_test

import (
	"app/database"
	"app/database/mocks"
	pb "app/explore_service_protos"
	"app/handlers"
	"context"
	"fmt"
	"testing"

	"github.com/zeebo/assert"
)

func TestListLikedYou(t *testing.T) {
	ctx := context.Background()
	page := 1
	recipientId := "1"

	dbResponse := []database.DecisionModel{
		{
			Id:          "1",
			ActorId:     2,
			RecipientId: 1,
			Liked:       true,
			Created_at:  20090520145024798,
			Updated_at:  20090520145024798,
		},
	}
	emptyDbResponse := []database.DecisionModel{}

	tests := []struct {
		name         string
		reader       func(t *testing.T) database.Reader
		writer       func(t *testing.T) database.Writer
		expectations func(t *testing.T, output *pb.ListLikedYouResponse, err error)
	}{
		{
			name: "successful_request",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("FindLikesByRecipientIdPaginated", ctx, recipientId, page).Return(dbResponse, nil)
				mockReader.Mock.On("GetLimit").Return(10)
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				mockWriter := mocks.NewWriter(t)
				mockWriter.Mock.On("UpdateLikesAsViewed", ctx, recipientId, dbResponse).Return(nil)
				return mockWriter
			},
			expectations: func(t *testing.T, output *pb.ListLikedYouResponse, err error) {
				likers := []*pb.ListLikedYouResponse_Liker{}
				likers = append(likers, &pb.ListLikedYouResponse_Liker{
					ActorId:       "2",
					UnixTimestamp: 20090520145024798,
				})
				expectedResponse := &pb.ListLikedYouResponse{
					Likers: likers,
				}

				assert.NoError(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
		{
			name: "no_likes_for_given_ricipient",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("FindLikesByRecipientIdPaginated", ctx, recipientId, page).Return(emptyDbResponse, nil)
				mockReader.Mock.On("GetLimit").Return(10)
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				return mocks.NewWriter(t)
			},
			expectations: func(t *testing.T, output *pb.ListLikedYouResponse, err error) {
				likers := []*pb.ListLikedYouResponse_Liker{}
				expectedResponse := &pb.ListLikedYouResponse{
					Likers: likers,
				}

				assert.NoError(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
		{
			name: "error_while_getting_likes",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("FindLikesByRecipientIdPaginated", ctx, recipientId, page).Return(emptyDbResponse, fmt.Errorf("generic error"))
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				return mocks.NewWriter(t)
			},
			expectations: func(t *testing.T, output *pb.ListLikedYouResponse, err error) {
				expectedResponse := &pb.ListLikedYouResponse{}
				assert.Error(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
		{
			name: "error_while_updating_likes_as_viewed",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("FindLikesByRecipientIdPaginated", ctx, recipientId, page).Return(dbResponse, nil)
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				mockWriter := mocks.NewWriter(t)
				mockWriter.Mock.On("UpdateLikesAsViewed", ctx, recipientId, dbResponse).Return(fmt.Errorf("some error"))
				return mockWriter
			},
			expectations: func(t *testing.T, output *pb.ListLikedYouResponse, err error) {
				expectedResponse := &pb.ListLikedYouResponse{}
				assert.Error(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
	}

	requestPage := fmt.Sprintf("%d", page)
	request := &pb.ListLikedYouRequest{
		RecipientUserId: recipientId,
		PaginationToken: &requestPage,
	}

	for _, test := range tests {
		server := handlers.Server{
			DatabaseReader: test.reader(t),
			DatabaseWriter: test.writer(t),
		}
		t.Run(test.name, func(t *testing.T) {
			output, err := server.ListLikedYou(ctx, request)
			test.expectations(t, output, err)
		})
	}
}

func TestListNewLikedYou(t *testing.T) {
	ctx := context.Background()
	page := 1
	recipientId := "1"

	dbResponse := []database.DecisionModel{
		{
			Id:          "1",
			ActorId:     2,
			RecipientId: 1,
			Liked:       true,
			Created_at:  20090520145024798,
			Updated_at:  20090520145024798,
		},
	}
	emptyDbResponse := []database.DecisionModel{}

	tests := []struct {
		name         string
		reader       func(t *testing.T) database.Reader
		writer       func(t *testing.T) database.Writer
		expectations func(t *testing.T, output *pb.ListLikedYouResponse, err error)
	}{
		{
			name: "successful_request",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("FindNewLikesByRecipientIdPaginated", ctx, recipientId, page).Return(dbResponse, nil)
				mockReader.Mock.On("GetLimit").Return(10)
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				mockWriter := mocks.NewWriter(t)
				mockWriter.Mock.On("UpdateLikesAsViewed", ctx, recipientId, dbResponse).Return(nil)
				return mockWriter
			},
			expectations: func(t *testing.T, output *pb.ListLikedYouResponse, err error) {
				likers := []*pb.ListLikedYouResponse_Liker{}
				likers = append(likers, &pb.ListLikedYouResponse_Liker{
					ActorId:       "2",
					UnixTimestamp: 20090520145024798,
				})
				expectedResponse := &pb.ListLikedYouResponse{
					Likers: likers,
				}

				assert.NoError(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
		{
			name: "no_likes_for_given_ricipient",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("FindNewLikesByRecipientIdPaginated", ctx, recipientId, page).Return(emptyDbResponse, nil)
				mockReader.Mock.On("GetLimit").Return(10)
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				return mocks.NewWriter(t)
			},
			expectations: func(t *testing.T, output *pb.ListLikedYouResponse, err error) {
				likers := []*pb.ListLikedYouResponse_Liker{}
				expectedResponse := &pb.ListLikedYouResponse{
					Likers: likers,
				}

				assert.NoError(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
		{
			name: "error_while_getting_likes",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("FindNewLikesByRecipientIdPaginated", ctx, recipientId, page).Return(emptyDbResponse, fmt.Errorf("generic error"))
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				return mocks.NewWriter(t)
			},
			expectations: func(t *testing.T, output *pb.ListLikedYouResponse, err error) {
				expectedResponse := &pb.ListLikedYouResponse{}
				assert.Error(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
		{
			name: "error_while_updating_likes_as_viewed",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("FindNewLikesByRecipientIdPaginated", ctx, recipientId, page).Return(dbResponse, nil)
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				mockWriter := mocks.NewWriter(t)
				mockWriter.Mock.On("UpdateLikesAsViewed", ctx, recipientId, dbResponse).Return(fmt.Errorf("some error"))
				return mockWriter
			},
			expectations: func(t *testing.T, output *pb.ListLikedYouResponse, err error) {
				expectedResponse := &pb.ListLikedYouResponse{}
				assert.Error(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
	}

	requestPage := fmt.Sprintf("%d", page)
	request := &pb.ListLikedYouRequest{
		RecipientUserId: recipientId,
		PaginationToken: &requestPage,
	}

	for _, test := range tests {
		server := handlers.Server{
			DatabaseReader: test.reader(t),
			DatabaseWriter: test.writer(t),
		}
		t.Run(test.name, func(t *testing.T) {
			output, err := server.ListNewLikedYou(ctx, request)
			test.expectations(t, output, err)
		})
	}
}

func TestCountLikedYou(t *testing.T) {
	ctx := context.Background()
	likes := 20
	recipientId := "1"

	dbResponse := database.UserModel{
		Id:        "1",
		Name:      "John Doe",
		Likes:     uint(likes),
		Gender:    "M",
		CreatedAt: 20090520145024798,
		UpdatedAt: 20090520145024798,
		IsAactive: true,
	}
	emptyDbResponse := database.UserModel{}

	tests := []struct {
		name         string
		reader       func(t *testing.T) database.Reader
		writer       func(t *testing.T) database.Writer
		expectations func(t *testing.T, output *pb.CountLikedYouResponse, err error)
	}{
		{
			name: "successful_request",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("GetUserById", ctx, recipientId).Return(dbResponse, nil)
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				return mocks.NewWriter(t)
			},
			expectations: func(t *testing.T, output *pb.CountLikedYouResponse, err error) {
				expectedResponse := &pb.CountLikedYouResponse{Count: uint64(likes)}

				assert.NoError(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
		{
			name: "error_getting_total_likes",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("GetUserById", ctx, recipientId).Return(emptyDbResponse, fmt.Errorf("generic error"))
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				return mocks.NewWriter(t)
			},
			expectations: func(t *testing.T, output *pb.CountLikedYouResponse, err error) {
				expectedResponse := &pb.CountLikedYouResponse{}

				assert.Error(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
	}

	request := &pb.CountLikedYouRequest{
		RecipientUserId: recipientId,
	}

	for _, test := range tests {
		server := handlers.Server{
			DatabaseReader: test.reader(t),
			DatabaseWriter: test.writer(t),
		}
		t.Run(test.name, func(t *testing.T) {
			output, err := server.CountLikedYou(ctx, request)
			test.expectations(t, output, err)
		})
	}
}

func TestPutDecision(t *testing.T) {
	ctx := context.Background()
	recipientId := "1"
	actorId := "2"

	tests := []struct {
		name         string
		reader       func(t *testing.T) database.Reader
		writer       func(t *testing.T) database.Writer
		request      *pb.PutDecisionRequest
		expectations func(t *testing.T, output *pb.PutDecisionResponse, err error)
	}{
		{
			name: "successful_request_with_match",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("GetIsMatch", ctx, actorId, recipientId).Return(true, nil)
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				mockWriter := mocks.NewWriter(t)
				entry := database.PutDecisionEntry{
					ActorId:     actorId,
					RecipientId: recipientId,
					Like:        true,
				}
				mockWriter.Mock.On("InsertOrUpdateDecision", ctx, entry).Return(nil)
				mockWriter.Mock.On("UpdateUserTotalLikes", ctx, recipientId).Return(nil)

				return mockWriter
			},
			request: &pb.PutDecisionRequest{
				ActorUserId:     actorId,
				RecipientUserId: recipientId,
				LikedRecipient:  true,
			},
			expectations: func(t *testing.T, output *pb.PutDecisionResponse, err error) {
				expectedResponse := &pb.PutDecisionResponse{MutualLikes: true}

				assert.NoError(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
		{
			name: "successful_request_with_match",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("GetIsMatch", ctx, actorId, recipientId).Return(false, nil)
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				mockWriter := mocks.NewWriter(t)
				entry := database.PutDecisionEntry{
					ActorId:     actorId,
					RecipientId: recipientId,
					Like:        false,
				}
				mockWriter.Mock.On("InsertOrUpdateDecision", ctx, entry).Return(nil)
				mockWriter.Mock.On("UpdateUserTotalLikes", ctx, recipientId).Return(nil)

				return mockWriter
			},
			request: &pb.PutDecisionRequest{
				ActorUserId:     actorId,
				RecipientUserId: recipientId,
				LikedRecipient:  false,
			},
			expectations: func(t *testing.T, output *pb.PutDecisionResponse, err error) {
				expectedResponse := &pb.PutDecisionResponse{MutualLikes: false}

				assert.NoError(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
		{
			name: "error_update_decision",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				mockWriter := mocks.NewWriter(t)
				entry := database.PutDecisionEntry{
					ActorId:     actorId,
					RecipientId: recipientId,
					Like:        false,
				}
				mockWriter.Mock.On("InsertOrUpdateDecision", ctx, entry).Return(fmt.Errorf("generic error"))
				return mockWriter
			},
			request: &pb.PutDecisionRequest{
				ActorUserId:     actorId,
				RecipientUserId: recipientId,
				LikedRecipient:  false,
			},
			expectations: func(t *testing.T, output *pb.PutDecisionResponse, err error) {
				expectedResponse := &pb.PutDecisionResponse{MutualLikes: false}

				fmt.Printf("Expected: %+v|\n", expectedResponse)
				fmt.Printf("output: %+v|\n", output)

				assert.Error(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
		{
			name: "error_update_total_like",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				mockWriter := mocks.NewWriter(t)
				entry := database.PutDecisionEntry{
					ActorId:     actorId,
					RecipientId: recipientId,
					Like:        false,
				}
				mockWriter.Mock.On("InsertOrUpdateDecision", ctx, entry).Return(nil)
				mockWriter.Mock.On("UpdateUserTotalLikes", ctx, recipientId).Return(fmt.Errorf("generic error"))

				return mockWriter
			},
			request: &pb.PutDecisionRequest{
				ActorUserId:     actorId,
				RecipientUserId: recipientId,
				LikedRecipient:  false,
			},
			expectations: func(t *testing.T, output *pb.PutDecisionResponse, err error) {
				expectedResponse := &pb.PutDecisionResponse{MutualLikes: false}

				assert.Error(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
		{
			name: "error_check_if_is_match",
			reader: func(t *testing.T) database.Reader {
				mockReader := mocks.NewReader(t)
				mockReader.Mock.On("GetIsMatch", ctx, actorId, recipientId).Return(false, fmt.Errorf("generic error"))
				return mockReader
			},
			writer: func(t *testing.T) database.Writer {
				mockWriter := mocks.NewWriter(t)
				entry := database.PutDecisionEntry{
					ActorId:     actorId,
					RecipientId: recipientId,
					Like:        false,
				}
				mockWriter.Mock.On("InsertOrUpdateDecision", ctx, entry).Return(nil)
				mockWriter.Mock.On("UpdateUserTotalLikes", ctx, recipientId).Return(nil)

				return mockWriter
			},
			request: &pb.PutDecisionRequest{
				ActorUserId:     actorId,
				RecipientUserId: recipientId,
				LikedRecipient:  false,
			},
			expectations: func(t *testing.T, output *pb.PutDecisionResponse, err error) {
				expectedResponse := &pb.PutDecisionResponse{MutualLikes: false}

				assert.Error(t, err)
				assert.Equal(t, expectedResponse, output)
			},
		},
	}

	for _, test := range tests {
		server := handlers.Server{
			DatabaseReader: test.reader(t),
			DatabaseWriter: test.writer(t),
		}
		t.Run(test.name, func(t *testing.T) {
			output, err := server.PutDecision(ctx, test.request)
			test.expectations(t, output, err)
		})
	}
}

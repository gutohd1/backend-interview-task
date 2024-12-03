package handlers

import (
	"app/database"
	"context"
	"fmt"
	"log"
	"strconv"

	pb "app/explore_service_protos"
)

type Server struct {
	DatabaseReader database.Reader
	DatabaseWriter database.Writer
	pb.UnimplementedExploreServiceServer
}

func (s *Server) SortLikers(likes []database.DecisionModel) []*pb.ListLikedYouResponse_Liker {

	likers := []*pb.ListLikedYouResponse_Liker{}

	for _, like := range likes {
		liker := &pb.ListLikedYouResponse_Liker{
			ActorId:       fmt.Sprintf("%d", like.ActorId),
			UnixTimestamp: like.Updated_at,
		}
		likers = append(likers, liker)
	}

	return likers
}

func (s *Server) GetNextPage(currPage int, limit int, pageSize int) *string {
	next := fmt.Sprintf("%d", currPage+1)

	if pageSize < limit {
		return nil
	}

	return &next
}

func (s *Server) ListLikedYou(ctx context.Context, request *pb.ListLikedYouRequest) (*pb.ListLikedYouResponse, error) {
	page := 1
	if request.PaginationToken != nil {
		var err error
		page, err = strconv.Atoi(*request.PaginationToken)
		if err != nil {
			log.Printf("Error converting page number: %s", err)
		}
	}

	likes, err := s.DatabaseReader.FindLikesByRecipientIdPaginated(ctx, request.RecipientUserId, page)
	if err != nil {
		log.Printf("Error on FindLikesByRecipientIdPaginated: %s", err)
		return &pb.ListLikedYouResponse{}, fmt.Errorf("unable to find likes for ListLikedYou")
	}

	if len(likes) > 0 {
		err = s.DatabaseWriter.UpdateLikesAsViewed(ctx, request.RecipientUserId, likes)
		if err != nil {
			log.Printf("Error on UpdateLikesAsViewed: %s", err)
			return &pb.ListLikedYouResponse{}, fmt.Errorf("unable to update likes")
		}
	}

	nextPage := s.GetNextPage(page, s.DatabaseReader.GetLimit(), len(likes))
	likers := s.SortLikers(likes)

	return &pb.ListLikedYouResponse{
		Likers:              likers,
		NextPaginationToken: nextPage,
	}, nil
}

func (s *Server) ListNewLikedYou(ctx context.Context, request *pb.ListLikedYouRequest) (*pb.ListLikedYouResponse, error) {
	page := 1
	if request.PaginationToken != nil {
		var err error
		page, err = strconv.Atoi(*request.PaginationToken)
		if err != nil {
			log.Printf("Error converting page number: %s", err)
		}
	}

	likes, err := s.DatabaseReader.FindNewLikesByRecipientIdPaginated(ctx, request.RecipientUserId, page)
	if err != nil {
		log.Printf("Error on FindNewLikesByRecipientIdPaginated: %s", err)
		return &pb.ListLikedYouResponse{}, fmt.Errorf("unable to find likes for ListNewLikedYou")
	}

	if len(likes) > 0 {
		err = s.DatabaseWriter.UpdateLikesAsViewed(ctx, request.RecipientUserId, likes)
		if err != nil {
			log.Printf("Error on UpdateLikesAsViewed: %s", err)
			return &pb.ListLikedYouResponse{}, fmt.Errorf("unable to update likes")
		}
	}

	nextPage := s.GetNextPage(page, s.DatabaseReader.GetLimit(), len(likes))
	likers := s.SortLikers(likes)

	return &pb.ListLikedYouResponse{
		Likers:              likers,
		NextPaginationToken: nextPage,
	}, nil
}
func (s *Server) CountLikedYou(ctx context.Context, request *pb.CountLikedYouRequest) (*pb.CountLikedYouResponse, error) {
	user, err := s.DatabaseReader.GetUserById(ctx, request.RecipientUserId)

	if err != nil {
		log.Printf("Error on GetUserById: %s", err)
		return &pb.CountLikedYouResponse{}, err
	}

	return &pb.CountLikedYouResponse{Count: uint64(user.Likes)}, nil
}
func (s *Server) PutDecision(ctx context.Context, request *pb.PutDecisionRequest) (*pb.PutDecisionResponse, error) {
	response := &pb.PutDecisionResponse{MutualLikes: false}
	err := s.DatabaseWriter.InsertOrUpdateDecision(ctx, database.PutDecisionEntry{
		ActorId:     request.ActorUserId,
		RecipientId: request.RecipientUserId,
		Like:        request.LikedRecipient,
	})
	if err != nil {
		log.Printf("Error on GetUserById while InsertOrUpdateDecision: %s", err)
		return response, fmt.Errorf("unable to create or update decision")
	}

	err = s.DatabaseWriter.UpdateUserTotalLikes(ctx, request.RecipientUserId)
	if err != nil {
		log.Printf("Error on GetUserById while UpdateUserTotalLikes: %s", err)
		return response, fmt.Errorf("unable to update user total likes")
	}

	var isMatch bool
	isMatch, err = s.DatabaseReader.GetIsMatch(ctx, request.ActorUserId, request.RecipientUserId)
	if err != nil {
		log.Printf("Error on GetUserById while GetIsMatch: %s", err)
		return response, err
	}

	response.MutualLikes = isMatch
	return response, nil
}

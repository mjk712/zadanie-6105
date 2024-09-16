package postgresql

import (
	"api/internal/http-server/handlers/bids/bidReviews"
	"api/internal/http-server/handlers/bids/redactBid"
	"api/internal/http-server/handlers/tenders/redactTender"
	"api/internal/models"
	"api/internal/storage/query"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"strconv"
	//_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func New(connectionString string) (*Storage, error) {
	const op = "storage.postgres.New"
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	m, err := migrate.New("file:///app/internal/storage/migrations", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = m.Up()
	if err != nil {
		fmt.Println(fmt.Errorf("%s: %w", op, err))
	}

	return &Storage{db}, nil
}

func (s *Storage) GetTenders(limit int, offset int, serviceTypes []string) ([]*models.Tender, error) {
	const op = "storage.postgres.tender.getAll"

	query := `SELECT * FROM tender`
	args := []interface{}{}
	argIndex := 1

	if len(serviceTypes) > 0 {
		query += "WHERE service_type IN("
		for i, serviceType := range serviceTypes {
			if i > 0 {
				query += ", "
			}
			query += fmt.Sprintf("$%d", argIndex)
			args = append(args, serviceType)
			argIndex++
		}
		query += ")"
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	var tenders = make([]*models.Tender, 0)
	rows, err := s.db.Queryx(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		var t models.Tender
		if err := rows.StructScan(&t); err != nil {
			return nil, fmt.Errorf("scan rows %s: %w", op, err)
		}
		tenders = append(tenders, &t)
	}
	rows.Close()
	return tenders, nil
}

func (s *Storage) NewTender(tender *models.Tender) (*models.Tender, error) {
	const (
		op     = "storage.postgres.tender.addTender"
		status = "Created"
	)

	var userId uuid.UUID
	err := s.db.QueryRowx(query.GetIdByUsername, tender.CreatorUsername).Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var isResponsible bool
	err = s.db.QueryRowx(query.CheckIfUserIsResponsible, tender.OrganizationId, userId).Scan(&isResponsible)
	if err != nil {
		return nil, fmt.Errorf("check responsible %s: %w", op, err)
	}
	if !isResponsible {
		return nil, fmt.Errorf("%s: user is not responsible")
	}
	var t models.Tender
	err = s.db.QueryRowx(query.InsertTender, tender.Name, tender.Description, tender.ServiceType, status, tender.OrganizationId, tender.CreatorUsername).StructScan(&t)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &t, nil
}

func (s *Storage) GetMyTender(username string, limit int, offset int) (*models.Tender, error) {
	const (
		op = "storage.postgres.tender.myTender"
	)
	var t models.Tender
	err := s.db.QueryRowx(query.GetMyTender, username, limit, offset).StructScan(&t)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &t, nil
}

func (s *Storage) GetTenderStatus(username string, id string) (string, error) {
	const (
		op = "storage.postgres.tender.tenderStatus"
	)
	var t models.Tender
	err := s.db.QueryRowx(query.GetTenderStatus, username, id).StructScan(&t)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return t.Status, nil
}

func (s *Storage) PutTenderStatus(username string, id string, status string) (*models.Tender, error) {
	const (
		op = "storage.postgres.tender.tenderStatus"
	)
	//todo переделать под один запрос как в бид
	var t models.Tender
	_, err := s.db.Query(query.PutTenderStatus, status, username, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	er := s.db.QueryRowx(query.GetTender, id).StructScan(&t)
	if er != nil {
		return nil, fmt.Errorf("%s getTenders tender: %w", op, err)
	}

	return &t, nil
}

func (s *Storage) UpdateTender(username string, id string, updData *redactTender.Request) (*models.Tender, error) {
	const (
		op = "storage.postgres.tender.update"
	)

	query := "UPDATE tender SET"
	params := []interface{}{}
	idx := 1

	if updData.Name != nil {
		if idx > 1 {
			query += ", "
		}
		query += " name = $" + strconv.Itoa(idx)
		params = append(params, *updData.Name)
		idx++
	}
	if updData.Description != nil {
		if idx > 1 {
			query += ", "
		}
		query += " description = $" + strconv.Itoa(idx)
		params = append(params, *updData.Description)
		idx++
	}
	if updData.ServiceType != nil {
		if idx > 1 {
			query += ", "
		}
		query += " service_type = $" + strconv.Itoa(idx)
		params = append(params, *updData.ServiceType)
		idx++
	}
	if updData.Status != nil {
		if idx > 1 {
			query += ", "
		}
		query += " tenderStatus = $" + strconv.Itoa(idx) + "::tenderStatus "
		params = append(params, *updData.Status)
		idx++
	}

	query += " WHERE id = $" + strconv.Itoa(idx)
	idx++
	params = append(params, id)
	query += " AND creator_username = $" + strconv.Itoa(idx)
	params = append(params, username)
	query += " RETURNING *;"

	fmt.Println(query, params)

	var updatedTender models.Tender
	err := s.db.QueryRowx(query, params...).StructScan(&updatedTender)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &updatedTender, nil
}

func (s *Storage) RollbackTender(username string, version int, id string) (*models.Tender, error) {
	const (
		op = "storage.postgres.tender.rollback"
	)

	var t models.Tender
	err := s.db.QueryRowx(query.TenderRollback, id, username, version).StructScan(&t)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &t, nil
}

func (s *Storage) NewBid(bid *models.Bid) (*models.Bid, error) {
	const (
		op     = "storage.postgres.bids.addBid"
		status = "Created"
	)

	var organizationId uuid.UUID
	err := s.db.QueryRowx(query.GetOrganizationByTenderId, bid.TenderId).Scan(&organizationId)
	if err != nil {
		return nil, fmt.Errorf("get tender organization %s: %w", op, err)
	}
	var isResponsible bool

	err = s.db.QueryRowx(query.CheckIfUserIsResponsible, organizationId, bid.AuthorId).Scan(&isResponsible)
	if err != nil {
		return nil, fmt.Errorf("check responsible %s: %w", op, err)
	}

	if !isResponsible {
		return nil, fmt.Errorf("%s: user is not responsible")
	}

	var b models.Bid
	err = s.db.QueryRowx(query.InsertBid, bid.Name, bid.Description, status, bid.TenderId, bid.AuthorType, bid.AuthorId).StructScan(&b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &b, nil
}

func (s *Storage) GetMyBids(username string, limit int, offset int) ([]*models.Bid, error) {
	const (
		op = "storage.postgres.bids.myBids"
	)
	var uid uuid.UUID
	err := s.db.QueryRowx(query.GetIdByUsername, username).Scan(&uid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var bids = make([]*models.Bid, 0)
	rows, err := s.db.Queryx(query.GetMyBids, uid, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		var b models.Bid
		if err := rows.StructScan(&b); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		bids = append(bids, &b)
	}
	rows.Close()
	return bids, nil

}

func (s *Storage) GetBidTenderList(tenderId string, limit int, offset int) ([]*models.Bid, error) {
	const (
		op = "storage.postgres.bids.getBidTenderList"
	)

	var bids = make([]*models.Bid, 0)
	rows, err := s.db.Queryx(query.GetBidsByTenderId, tenderId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		var b models.Bid
		if err := rows.StructScan(&b); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		bids = append(bids, &b)
	}
	rows.Close()
	return bids, nil
}

func (s *Storage) GetBidStatus(username string, id string) (string, error) {
	const (
		op = "storage.postgres.bids.bidStatus"
	)
	var userId uuid.UUID
	err := s.db.QueryRowx(query.GetIdByUsername, username).Scan(&userId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var b models.Bid
	err = s.db.QueryRowx(query.GetBidStatus, userId, id).StructScan(&b)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return b.Status, nil
}

func (s *Storage) PutBidStatus(username string, id string, status string) (*models.Bid, error) {
	const (
		op = "storage.postgres.bids.changeStatus"
	)
	var userId uuid.UUID
	err := s.db.QueryRowx(query.GetIdByUsername, username).Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var b models.Bid
	err = s.db.QueryRowx(query.PutBidStatus, status, userId, id).StructScan(&b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &b, nil
}

func (s *Storage) UpdateBid(username string, id string, updData *redactBid.Request) (*models.Bid, error) {
	const (
		op = "storage.postgres.bids.update"
	)

	updQuery := "UPDATE bid SET"
	params := []interface{}{}
	idx := 1

	var userId uuid.UUID
	err := s.db.QueryRowx(query.GetIdByUsername, username).Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if updData.Name != nil {
		if idx > 1 {
			updQuery += ", "
		}
		updQuery += " name = $" + strconv.Itoa(idx)
		params = append(params, *updData.Name)
		idx++
	}
	if updData.Description != nil {
		if idx > 1 {
			updQuery += ", "
		}
		updQuery += " description = $" + strconv.Itoa(idx)
		params = append(params, *updData.Description)
		idx++
	}

	updQuery += " WHERE id = $" + strconv.Itoa(idx)
	idx++
	params = append(params, id)
	updQuery += " AND author_id = $" + strconv.Itoa(idx)
	params = append(params, userId)
	updQuery += " RETURNING id,name,description,status,author_type,author_id,version,created_at;"

	var updatedBid models.Bid
	err = s.db.QueryRowx(updQuery, params...).StructScan(&updatedBid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &updatedBid, nil
}

func (s *Storage) PutBidDecision(username string, id string, decision string) (*models.Bid, error) {
	const (
		op = "storage.postgres.bids.submitDecision"
	)
	var userId uuid.UUID
	err := s.db.QueryRowx(query.GetIdByUsername, username).Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var b models.Bid
	err = s.db.QueryRowx(query.PutBidDecision, decision, userId, id).StructScan(&b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &b, nil
}

func (s *Storage) PutBidFeedback(username string, id string, feedback string) (*models.Bid, error) {
	const (
		op = "storage.postgres.bids.putFeedback"
	)
	var userId uuid.UUID
	err := s.db.QueryRowx(query.GetIdByUsername, username).Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf(" GetIdByUsername %s: %w", op, err)
	}

	var feedbackId uuid.UUID
	err = s.db.QueryRowx(query.PutBidFeedback, id, userId, feedback).Scan(&feedbackId)
	if err != nil {
		return nil, fmt.Errorf(" PutBidFeedback %s: %w", op, err)
	}

	var b models.Bid
	err = s.db.QueryRowx(query.GetBidByFeedbackId, feedbackId).StructScan(&b)
	if err != nil {
		return nil, fmt.Errorf(" GetBidByFeedbackId %s: %w", op, err)
	}

	return &b, nil
}

func (s *Storage) RollbackBid(username string, version int, id string) (*models.Bid, error) {
	const (
		op = "storage.postgres.bids.rollback"
	)
	var userId uuid.UUID
	err := s.db.QueryRowx(query.GetIdByUsername, username).Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var b models.Bid
	err = s.db.QueryRowx(query.BidRollback, id, userId, version).StructScan(&b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &b, nil
}

func (s *Storage) GetBidReviews(tenderId string, authorUsername string, requesterUsername string, limit int, offset int) ([]*bidReviews.Response, error) {
	const (
		op = "storage.postgres.bids.getBidReviews"
	)

	var authorId uuid.UUID
	err := s.db.QueryRowx(query.GetIdByUsername, authorUsername).Scan(&authorId)
	if err != nil {
		return nil, fmt.Errorf("get author id %s: %w", op, err)
	}

	var requesterId uuid.UUID
	err = s.db.QueryRowx(query.GetIdByUsername, requesterUsername).Scan(&requesterId)
	if err != nil {
		return nil, fmt.Errorf("get requestor id %s: %w", op, err)
	}

	var organizationId uuid.UUID
	err = s.db.QueryRowx(query.GetOrganizationByTenderId, tenderId).Scan(&organizationId)
	if err != nil {
		return nil, fmt.Errorf("get tender organization %s: %w", op, err)
	}

	var isResponsible bool
	err = s.db.QueryRowx(query.CheckIfUserIsResponsible, organizationId, requesterId).Scan(&isResponsible)
	if err != nil {
		return nil, fmt.Errorf("check responsible %s: %w", op, err)
	}
	if !isResponsible {
		return nil, fmt.Errorf("%s: user is not responsible")
	}

	var feedbacks = make([]*bidReviews.Response, 0)
	rows, err := s.db.Queryx(query.GetBidsReviews, tenderId, authorId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get reviews %s: %w", op, err)
	}
	for rows.Next() {
		var f bidReviews.Response
		if err := rows.StructScan(&f); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		feedbacks = append(feedbacks, &f)
	}
	rows.Close()
	return feedbacks, nil
}

func (s *Storage) IsUserAttachedToOrganization(username string, tenderId string) (bool, error) {
	const (
		op = "storage.postgres.IsUserAttachedToOrganization"
	)
	var userId uuid.UUID
	err := s.db.QueryRowx(query.GetIdByUsername, username).Scan(&userId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	fmt.Println(userId)
	var orgId uuid.UUID
	err = s.db.QueryRowx(query.GetOrganizationByTenderId, tenderId).Scan(&orgId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var isResponsible bool

	err = s.db.QueryRowx(query.CheckIfUserIsResponsible, orgId, userId).Scan(&isResponsible)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isResponsible, nil
}

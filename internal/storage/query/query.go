package query

var InsertTender = `INSERT INTO tender(name,description,service_type,status,organization_id,creator_username) VALUES($1,$2,$3,$4,$5,$6) RETURNING id,name,description,service_type,status,version,created_at;`

var GetMyTender = `SELECT * FROM tender WHERE creator_username = $1 LIMIT $2 OFFSET $3;`

var GetTenderStatus = `SELECT status FROM tender WHERE creator_username = $1 AND id = $2;`

var PutTenderStatus = `UPDATE tender SET status = $1::status WHERE creator_username = $2 AND id = $3;`

var GetTender = `SELECT * FROM tender WHERE id = $1;`

var TenderRollback = `
UPDATE tender t
SET name = th.name,
description = th.description,
service_type = th.service_type,
status = th.status,
version = t.version+1,
updated_at = CURRENT_TIMESTAMP 
FROM tender_history th
WHERE t.id = th.tender_id
AND t.id = $1
AND t.creator_username = $2
AND th.version = $3
RETURNING t.*;`

var InsertBid = `INSERT INTO bid(name,description,status,tender_id,author_type,author_id) VALUES($1,$2,$3,$4,$5,$6) RETURNING id,name,description,status,author_type,author_id,version,created_at;`

var GetMyBids = `SELECT id,name,description,status,author_type,author_id,version,created_at FROM bid WHERE author_id =$1 LIMIT $2 OFFSET $3;`

var GetBidByFeedbackId = `
SELECT b.id,b.name,b.description,b.status,b.author_type,b.author_id,b.version,b.created_at
FROM bid b
JOIN feedback f ON b.id = f.bid_id
WHERE f.id = $1;`

var GetIdByUsername = `SELECT id FROM employee WHERE username = $1;`

var GetBidsByTenderId = `SELECT id,name,description,status,author_type,author_id,version,created_at FROM bid WHERE tender_id =$1 LIMIT $2 OFFSET $3;`

var GetBidStatus = `SELECT status FROM bid WHERE author_id = $1 AND id = $2;`

var PutBidStatus = `UPDATE bid SET status = $1::status WHERE author_id = $2 AND id = $3 RETURNING id,name,description,status,author_type,author_id,version,created_at;`

var PutBidDecision = `UPDATE bid SET decision = $1::decision WHERE author_id = $2 AND id = $3 RETURNING id,name,description,status,author_type,author_id,version,created_at;`

var PutBidFeedback = `
INSERT INTO feedback (bid_id, reviewer_id, feedback_text, created_at, updated_at)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id;`

var BidRollback = `
UPDATE bid b
SET name = bh.name,
    description = bh.description,
    status = bh.status,
    tender_id = bh.tender_id,
    author_type = bh.author_type,
    author_id = bh.author_id,
    decision = bh.decision,
    feedback = bh.feedback,
    version = b.version + 1, -- Увеличиваем версию
    updated_at = CURRENT_TIMESTAMP
FROM bid_history bh
WHERE b.id = bh.bid_id
  AND b.id = $1        
  AND b.author_id = $2  
  AND bh.version = $3   
RETURNING b.id,b.name,b.description,b.status,b.author_type,b.author_id,b.version,b.created_at;
`
var GetBidsReviews = `
SELECT b.tender_id, f.feedback_text, f.created_at
FROM bid b
LEFT JOIN feedback f ON b.id = f.bid_id
WHERE b.tender_id = $1 AND b.author_id = $2
LIMIT $3 OFFSET $4;`

var GetOrganizationByTenderId = `SELECT organization_id FROM tender WHERE id = $1;`

var CheckIfUserIsResponsible = `SELECT EXISTS (SELECT 1 FROM organization_responsible WHERE organization_id = $1 AND user_id = $2);`

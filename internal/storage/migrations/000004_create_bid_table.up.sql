CREATE TYPE decision AS ENUM (
    'Approved',
    'Rejected'
    );
CREATE TYPE author_type AS ENUM (
    'User',
    'Organization'
    );

CREATE TABLE IF NOT EXISTS bid (
                                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                      name VARCHAR(100) NOT NULL,
                                      description TEXT,
                                      status status,
                                      tender_id UUID,
                                      author_type author_type,
                                      author_id UUID,
                                      version INTEGER DEFAULT 1,
                                      decision decision,
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      FOREIGN KEY (tender_id) REFERENCES tender(id),
                                      FOREIGN KEY (author_id) REFERENCES employee(id)
);

CREATE TABLE IF NOT EXISTS bid_history (
                                   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                   bid_id UUID REFERENCES bid(id),
                                   name VARCHAR(100) NOT NULL,
                                   description TEXT,
                                   status status,
                                   tender_id UUID,
                                   author_type author_type,
                                   author_id UUID,
                                   version INTEGER DEFAULT 1,
                                   decision decision,
                                   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                   FOREIGN KEY (tender_id) REFERENCES tender(id),
                                   FOREIGN KEY (author_id) REFERENCES employee(id)
);

CREATE TABLE IF NOT EXISTS feedback(
                                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                       bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
                                       reviewer_id UUID REFERENCES employee(id) ON DELETE SET NULL,
                                       feedback_text TEXT NOT NULL,
                                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

CREATE OR REPLACE FUNCTION save_bid_version()
    RETURNS TRIGGER AS $$
BEGIN
    -- Сохраняем текущую версию заявки в таблицу истории перед обновлением
    INSERT INTO bid_history (bid_id, name, description, status, tender_id, author_type, author_id, version, decision, created_at, updated_at)
    VALUES (OLD.id, OLD.name, OLD.description, OLD.status, OLD.tender_id, OLD.author_type, OLD.author_id, OLD.version, OLD.decision,  OLD.created_at, OLD.updated_at);

    -- Увеличиваем версию и обновляем время обновления
    NEW.version = OLD.version + 1;
    NEW.updated_at = CURRENT_TIMESTAMP;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_bid_version
    BEFORE UPDATE ON bid
    FOR EACH ROW
EXECUTE FUNCTION save_bid_version();

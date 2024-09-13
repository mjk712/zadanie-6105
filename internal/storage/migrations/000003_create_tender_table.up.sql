CREATE TYPE service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
    );
CREATE TYPE status AS ENUM (
    'Created',
    'Published',
    'Closed'
    );

CREATE TABLE IF NOT EXISTS tender (
                                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                            name VARCHAR(100) NOT NULL,
                                            description TEXT,
                                            service_type service_type,
                                            status status,
                                            organization_id UUID,
                                            creator_username TEXT,
                                            version INTEGER DEFAULT 1,
                                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                            FOREIGN KEY (organization_id) REFERENCES organization(id)
);
CREATE TABLE IF NOT EXISTS tender_history(
                                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                             tender_id UUID REFERENCES tender(id),
                                             name VARCHAR(100) NOT NULL,
                                             description TEXT,
                                             service_type service_type,
                                             status status,
                                             organization_id UUID,
                                             creator_username TEXT,
                                             version INTEGER NOT NULL ,
                                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                             FOREIGN KEY (organization_id) REFERENCES organization(id)
);

CREATE OR REPLACE FUNCTION save_tender_version()
    RETURNS TRIGGER AS $$
BEGIN
    -- Сохраняем текущую версию из таблицы tender в таблицу tender_history перед обновлением
    INSERT INTO tender_history (tender_id, name, description, service_type, status, organization_id, creator_username, version, created_at, updated_at)
    VALUES (OLD.id, OLD.name, OLD.description, OLD.service_type, OLD.status::status, OLD.organization_id, OLD.creator_username, OLD.version, OLD.created_at, OLD.updated_at);

    -- Увеличиваем версию в таблице tender и обновляем время обновления
    NEW.version = OLD.version + 1;
    NEW.updated_at = CURRENT_TIMESTAMP;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_tender_version
    BEFORE UPDATE ON tender
    FOR EACH ROW
EXECUTE FUNCTION save_tender_version();

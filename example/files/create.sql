CREATE TABLE IF NOT EXISTS metadata (
    id SERIAL PRIMARY KEY,
    filepath TEXT NOT NULL,
    executed_at TIMESTAMPTZ DEFAULT timezone('Europe/Berlin', now()) NOT NULL
);


CREATE TABLE IF NOT EXISTS identities (
    id TEXT PRIMARY KEY,
    name TEXT
);

CREATE TABLE IF NOT EXISTS entities (
    id INTEGER NOT NULL,
    identity TEXT NOT NULL,
    PRIMARY KEY(id, identity),
    CONSTRAINT fk_identity FOREIGN KEY(identity) REFERENCES identities(id) ON DELETE CASCADE,
    entity_id TEXT NOT NULL,
    name TEXT NOT NULL,
    is_device BOOLEAN NOT NULL,
    allow_rules BOOLEAN NOT NULL,
    has_attribute BOOLEAN NOT NULL,
    attribute TEXT,
    is_victron_sensor BOOLEAN NOT NULL,
    sensor_type INTEGER NOT NULL DEFAULT 0,
    has_numeric_state BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS states (
    id SERIAL PRIMARY KEY,
    entity_id INTEGER NOT NULL,
    identity TEXT NOT NULL,
    CONSTRAINT fk_entity_id FOREIGN KEY(identity, entity_id) REFERENCES entities(identity, id) ON DELETE CASCADE,
    state TEXT NOT NULL,
    record_time TIMESTAMPTZ DEFAULT timezone('Europe/Berlin', now()) NOT NULL
);

CREATE TABLE IF NOT EXISTS honua_services (
    id INTEGER NOT NULL,
    identity TEXT NOT NULL,
    PRIMARY KEY(id, identity),
    CONSTRAINT fk_identity FOREIGN KEY(identity) REFERENCES identities(id) ON DELETE CASCADE,
    domain TEXT UNIQUE,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS allowed_sensors (
    id SERIAL PRIMARY KEY,
    identity TEXT NOT NULL,
    device_id INTEGER NOT NULL,
    sensor_id INTEGER NOT NULL,
    CONSTRAINT fk_device_id FOREIGN KEY(identity, device_id) REFERENCES entities(identity, id) ON DELETE CASCADE,
    CONSTRAINT fk_sensor_id FOREIGN KEY(identity, sensor_id) REFERENCES entities(identity, id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS allowed_service (
    id SERIAL PRIMARY KEY,
    identity TEXT NOT NULL,
    device_id INTEGER NOT NULL,
    service_id INTEGER NOT NULL,
    CONSTRAINT fk_device_id FOREIGN KEY(identity, device_id) REFERENCES entities(identity, id) ON DELETE CASCADE,
    CONSTRAINT fk_service_id FOREIGN KEY(identity, service_id) REFERENCES honua_services(identity, id) ON DELETE CASCADE
);
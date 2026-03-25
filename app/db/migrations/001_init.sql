-- migrate:up

-- ============================================================
-- EXTENSIONS
-- ============================================================
CREATE EXTENSION IF NOT EXISTS "btree_gist";

-- ============================================================
-- DEALERSHIPS
-- ============================================================
CREATE TABLE dealerships (
    id              BIGINT          PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uuid            CHAR(36)        NOT NULL UNIQUE,
    name            VARCHAR(255)    NOT NULL,
    timezone        VARCHAR(64)    NOT NULL DEFAULT 'Asia/Ho_Chi_Minh',
    open_time       TIME            NOT NULL DEFAULT '08:00:00',
    close_time      TIME            NOT NULL DEFAULT '18:00:00',
    is_weekend_open BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- ============================================================
-- SERVICE BAYS
-- ============================================================
CREATE TABLE service_bays (
    id              BIGINT          PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uuid            CHAR(36)        NOT NULL UNIQUE,
    dealership_id   BIGINT          NOT NULL REFERENCES dealerships(id) ON DELETE CASCADE,
    name            VARCHAR(100)    NOT NULL,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    UNIQUE (dealership_id, name)
);

-- ============================================================
-- SERVICE TYPES
-- ============================================================
CREATE TABLE service_types (
    id                BIGINT          PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uuid              CHAR(36)        NOT NULL UNIQUE,
    name              VARCHAR(255)    NOT NULL,
    duration_minutes  INT             NOT NULL CHECK (duration_minutes > 0),
    created_at        TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- ============================================================
-- TECHNICIANS
-- ============================================================
CREATE TABLE technicians (
    id              BIGINT          PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uuid            CHAR(36)        NOT NULL UNIQUE,
    dealership_id   BIGINT          NOT NULL REFERENCES dealerships(id) ON DELETE CASCADE,
    name            VARCHAR(255)    NOT NULL,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- ============================================================
-- TECHNICIANS <-> SERVICE TYPES  (many-to-many)
-- ============================================================
CREATE TABLE technician_service_types (
    technician_id    BIGINT      NOT NULL REFERENCES technicians(id) ON DELETE CASCADE,
    service_type_id  BIGINT      NOT NULL REFERENCES service_types(id) ON DELETE CASCADE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (technician_id, service_type_id)
);

-- ============================================================
-- CUSTOMERS
-- ============================================================
CREATE TABLE customers (
    id          BIGINT          PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uuid        CHAR(36)        NOT NULL UNIQUE,
    name        VARCHAR(255)    NOT NULL,
    email       VARCHAR(255)    NOT NULL UNIQUE,
    phone       VARCHAR(50)     NOT NULL,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- ============================================================
-- VEHICLES
-- ============================================================
CREATE TABLE vehicles (
    id           BIGINT          PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uuid         CHAR(36)        NOT NULL UNIQUE,
    customer_id  BIGINT          REFERENCES customers(id) ON DELETE SET NULL, -- nullable
    name         VARCHAR(255)    NOT NULL,
    created_at   TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- ============================================================
-- APPOINTMENTS
-- ============================================================
CREATE TABLE appointments (
    id               BIGINT          PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uuid             CHAR(36)        NOT NULL UNIQUE,
    customer_id      BIGINT          NOT NULL REFERENCES customers(id),
    vehicle_id       BIGINT          REFERENCES vehicles(id) ON DELETE SET NULL, -- nullable
    dealership_id    BIGINT          NOT NULL REFERENCES dealerships(id),
    technician_id    BIGINT          NOT NULL REFERENCES technicians(id),
    service_bay_id   BIGINT          NOT NULL REFERENCES service_bays(id),
    service_type_id  BIGINT          NOT NULL REFERENCES service_types(id),
    status           VARCHAR(200)    NOT NULL,
    start_at         TIMESTAMPTZ     NOT NULL,
    end_at           TIMESTAMPTZ     NOT NULL,
    description      TEXT,
    created_at       TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    CHECK (end_at > start_at),

    CONSTRAINT no_overlapping_technician_appointments
    EXCLUDE USING GIST (
        technician_id WITH =, 
        tstzrange(start_at, end_at, '[)') WITH &&
    ) WHERE (status != 'cancelled'),

    CONSTRAINT no_overlapping_service_bay_appointments
    EXCLUDE USING GIST (
        service_bay_id WITH =, 
        tstzrange(start_at, end_at, '[)') WITH &&
    ) WHERE (status != 'cancelled')
);

-- Indexes
CREATE INDEX idx_appointments_technician  ON appointments (technician_id, start_at, end_at);
CREATE INDEX idx_appointments_service_bay ON appointments (service_bay_id, start_at, end_at);

-- migrate:down

-- ============================================================
-- APPOINTMENTS
-- ============================================================
DROP INDEX IF EXISTS idx_appointments_technician;
DROP INDEX IF EXISTS idx_appointments_service_bay;
DROP TABLE IF EXISTS appointments;

-- ============================================================
-- VEHICLES
-- ============================================================
DROP TABLE IF EXISTS vehicles;

-- ============================================================
-- TECHNICIANS <-> SERVICE TYPES
-- ============================================================
DROP TABLE IF EXISTS technician_service_types;

-- ============================================================
-- TECHNICIANS
-- ============================================================
DROP TABLE IF EXISTS technicians;

-- ============================================================
-- SERVICE BAYS
-- ============================================================
DROP TABLE IF EXISTS service_bays;

-- ============================================================
-- CUSTOMERS
-- ============================================================
DROP TABLE IF EXISTS customers;

-- ============================================================
-- SERVICE TYPES
-- ============================================================
DROP TABLE IF EXISTS service_types;

-- ============================================================
-- DEALERSHIPS
-- ============================================================
DROP TABLE IF EXISTS dealerships;

-- ============================================================
-- EXTENSIONS
-- ============================================================
DROP EXTENSION IF EXISTS "btree_gist";
-- 0. Users Table
CREATE TABLE users (
    user_id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'OPERATOR' CHECK (role IN ('ADMIN', 'OPERATOR')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 1. Machines Table
CREATE TABLE machines (
    machine_id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'OPERATIONAL',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 2. Sensors Table
CREATE TABLE sensors (
    sensor_id VARCHAR(50) PRIMARY KEY,
    machine_id VARCHAR(50) NOT NULL REFERENCES machines(machine_id) ON DELETE CASCADE,
    metric_name VARCHAR(50) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sensors_machine_id ON sensors(machine_id);

-- 3. Telemetry Logs Table
CREATE TABLE telemetry_logs (
    log_id BIGSERIAL PRIMARY KEY,
    sensor_id VARCHAR(50) NOT NULL REFERENCES sensors(sensor_id) ON DELETE CASCADE,
    metric_value DOUBLE PRECISION NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_telemetry_sensor_recorded ON telemetry_logs(sensor_id, recorded_at DESC);

-- 4. AI Evaluations Table
CREATE TABLE ai_evaluations (
    eval_id BIGSERIAL PRIMARY KEY,
    machine_id VARCHAR(50) NOT NULL REFERENCES machines(machine_id) ON DELETE CASCADE,
    health_score DOUBLE PRECISION NOT NULL,
    is_anomaly BOOLEAN NOT NULL DEFAULT FALSE,
    rul_hours DOUBLE PRECISION NOT NULL,
    evaluated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_ai_evaluations_machine_time ON ai_evaluations(machine_id, evaluated_at DESC);

-- 5. Alerts Table
CREATE TABLE alerts (
    alert_id BIGSERIAL PRIMARY KEY,
    machine_id VARCHAR(50) NOT NULL REFERENCES machines(machine_id) ON DELETE CASCADE,
    severity VARCHAR(20) NOT NULL,
    message TEXT NOT NULL,
    is_resolved BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_alerts_machine_resolved ON alerts(machine_id, is_resolved);

-- 6. Command Logs Table
CREATE TABLE command_logs (
    command_id BIGSERIAL PRIMARY KEY,
    machine_id VARCHAR(50) NOT NULL REFERENCES machines(machine_id) ON DELETE CASCADE,
    command_type VARCHAR(50) NOT NULL,
    issued_by VARCHAR(50) NOT NULL,
    reason TEXT NOT NULL,
    executed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_command_logs_machine ON command_logs(machine_id);
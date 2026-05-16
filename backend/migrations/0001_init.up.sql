CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- USERS & PROFILES
-- ============================================================================

CREATE TYPE user_role AS ENUM ('admin', 'coach', 'student', 'parent', 'coordinator');

CREATE TABLE users (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    role            user_role NOT NULL,
    full_name       text NOT NULL,
    phone           text NOT NULL,
    email           text,
    password_hash   text NOT NULL,
    is_active       boolean NOT NULL DEFAULT true,
    city            text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT users_phone_unique UNIQUE (phone),
    CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);

-- ----------------------------------------------------------------------------
-- Pre-registration applications (öğrenci & koç başvuru formları)
-- ----------------------------------------------------------------------------

CREATE TYPE registration_status AS ENUM ('pending', 'approved', 'rejected');
CREATE TYPE registration_kind   AS ENUM ('student', 'coach');

CREATE TABLE registrations (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    kind            registration_kind NOT NULL,
    status          registration_status NOT NULL DEFAULT 'pending',
    payload         jsonb NOT NULL,
    full_name       text NOT NULL,
    phone           text NOT NULL,
    email           text,
    city            text,
    reviewed_by     uuid REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at     timestamptz,
    review_note     text,
    created_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_registrations_status ON registrations(status);
CREATE INDEX idx_registrations_kind ON registrations(kind);

-- ----------------------------------------------------------------------------
-- Profile tables (extend users by role)
-- ----------------------------------------------------------------------------

CREATE TABLE students (
    user_id         uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    school          text,
    grade           text,
    parent_id       uuid REFERENCES users(id) ON DELETE SET NULL,
    coordinator_id  uuid REFERENCES users(id) ON DELETE SET NULL,
    notes           text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_students_parent ON students(parent_id);
CREATE INDEX idx_students_coordinator ON students(coordinator_id);

CREATE TYPE coach_status AS ENUM ('pending', 'approved', 'rejected', 'suspended');

CREATE TABLE coaches (
    user_id         uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    status          coach_status NOT NULL DEFAULT 'approved',
    bio             text,
    experience      text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE parents (
    user_id         uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    notes           text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE coordinators (
    user_id         uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    foundation_name text,
    notes           text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

-- ============================================================================
-- COACH-STUDENT ASSIGNMENTS
-- ============================================================================

CREATE TABLE coach_student_assignments (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    coach_id        uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    student_id      uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    started_at      timestamptz NOT NULL DEFAULT now(),
    ended_at        timestamptz,
    is_active       boolean NOT NULL DEFAULT true,
    assigned_by     uuid REFERENCES users(id) ON DELETE SET NULL,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_assignments_active_unique
    ON coach_student_assignments(student_id)
    WHERE is_active = true;

CREATE INDEX idx_assignments_coach ON coach_student_assignments(coach_id);
CREATE INDEX idx_assignments_student ON coach_student_assignments(student_id);

-- ============================================================================
-- EVALUATION WEEKS (admin tanımlı haftalar)
-- ============================================================================

CREATE TABLE evaluation_weeks (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    week_no         int NOT NULL,
    label           text NOT NULL,
    start_date      date NOT NULL,
    end_date        date NOT NULL,
    is_open         boolean NOT NULL DEFAULT true,
    opened_by_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
    opened_at       timestamptz,
    closed_at       timestamptz,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT evaluation_weeks_week_no_unique UNIQUE (week_no),
    CONSTRAINT evaluation_weeks_dates_check CHECK (start_date <= end_date)
);

CREATE INDEX idx_eval_weeks_dates ON evaluation_weeks(start_date, end_date);
CREATE INDEX idx_eval_weeks_is_open ON evaluation_weeks(is_open);

-- ============================================================================
-- EVALUATIONS
-- ============================================================================

CREATE TYPE evaluation_status AS ENUM ('pending', 'approved', 'edited_by_admin');

CREATE TABLE evaluations (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id       uuid NOT NULL REFERENCES coach_student_assignments(id) ON DELETE CASCADE,
    evaluation_week_id  uuid NOT NULL REFERENCES evaluation_weeks(id) ON DELETE RESTRICT,

    course_status       text,
    homework_done       boolean,
    motivation          smallint,
    behavior            smallint,
    general_note        text,
    admin_only_note     text,

    status              evaluation_status NOT NULL DEFAULT 'pending',
    submitted_by        uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    submitted_at        timestamptz NOT NULL DEFAULT now(),
    last_edited_by      uuid REFERENCES users(id) ON DELETE SET NULL,
    last_edited_at      timestamptz,
    approved_by         uuid REFERENCES users(id) ON DELETE SET NULL,
    approved_at         timestamptz,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT evaluations_assignment_week_unique UNIQUE (assignment_id, evaluation_week_id),
    CONSTRAINT evaluations_motivation_range CHECK (motivation IS NULL OR (motivation BETWEEN 1 AND 5)),
    CONSTRAINT evaluations_behavior_range CHECK (behavior IS NULL OR (behavior BETWEEN 1 AND 5))
);

CREATE INDEX idx_evaluations_status ON evaluations(status);
CREATE INDEX idx_evaluations_week ON evaluations(evaluation_week_id);
CREATE INDEX idx_evaluations_assignment ON evaluations(assignment_id);
CREATE INDEX idx_evaluations_submitted_by ON evaluations(submitted_by);

-- Tarihsel sürümler — admin değiştirdikçe önceki snapshot buraya yazılır
CREATE TABLE evaluation_versions (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    evaluation_id   uuid NOT NULL REFERENCES evaluations(id) ON DELETE CASCADE,
    version_no      int NOT NULL,
    snapshot        jsonb NOT NULL,
    edited_by       uuid REFERENCES users(id) ON DELETE SET NULL,
    edited_at       timestamptz NOT NULL DEFAULT now(),
    change_reason   text,
    CONSTRAINT evaluation_versions_unique UNIQUE (evaluation_id, version_no)
);

CREATE INDEX idx_eval_versions_evaluation ON evaluation_versions(evaluation_id);

-- ============================================================================
-- MESSAGING & FEEDBACK
-- ============================================================================

CREATE TYPE thread_kind AS ENUM ('general', 'feedback', 'complaint');
CREATE TYPE thread_status AS ENUM ('open', 'closed');

CREATE TABLE message_threads (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    kind            thread_kind NOT NULL DEFAULT 'general',
    subject         text NOT NULL,
    opened_by_user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_to_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
    status          thread_status NOT NULL DEFAULT 'open',
    last_message_at timestamptz,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_threads_opened_by ON message_threads(opened_by_user_id);
CREATE INDEX idx_threads_status ON message_threads(status);
CREATE INDEX idx_threads_kind ON message_threads(kind);

CREATE TABLE messages (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    thread_id       uuid NOT NULL REFERENCES message_threads(id) ON DELETE CASCADE,
    from_user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body            text NOT NULL,
    read_at         timestamptz,
    created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_messages_thread ON messages(thread_id);
CREATE INDEX idx_messages_from_user ON messages(from_user_id);

-- ============================================================================
-- SMS LOGS
-- ============================================================================

CREATE TYPE sms_status AS ENUM ('mock_sent', 'sent', 'failed');

CREATE TABLE sms_logs (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    to_phone        text NOT NULL,
    to_user_id      uuid REFERENCES users(id) ON DELETE SET NULL,
    body            text NOT NULL,
    sent_by_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
    status          sms_status NOT NULL DEFAULT 'mock_sent',
    provider_name   text NOT NULL DEFAULT 'mock',
    provider_response jsonb,
    template_key    text,
    sent_at         timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_sms_logs_to_user ON sms_logs(to_user_id);
CREATE INDEX idx_sms_logs_sent_at ON sms_logs(sent_at DESC);

-- ============================================================================
-- AUDIT LOGS
-- ============================================================================

CREATE TABLE audit_logs (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_user_id   uuid REFERENCES users(id) ON DELETE SET NULL,
    entity          text NOT NULL,
    entity_id       text NOT NULL,
    action          text NOT NULL,
    before          jsonb,
    after           jsonb,
    metadata        jsonb,
    at              timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_entity ON audit_logs(entity, entity_id);
CREATE INDEX idx_audit_logs_actor ON audit_logs(actor_user_id);
CREATE INDEX idx_audit_logs_at ON audit_logs(at DESC);

-- ============================================================================
-- updated_at trigger
-- ============================================================================

CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
    t text;
BEGIN
    FOR t IN
        SELECT unnest(ARRAY[
            'users','registrations','students','coaches','parents','coordinators',
            'coach_student_assignments','evaluation_weeks','evaluations',
            'message_threads'
        ])
    LOOP
        EXECUTE format(
            'CREATE TRIGGER trg_set_updated_at_%I BEFORE UPDATE ON %I FOR EACH ROW EXECUTE FUNCTION set_updated_at()',
            t, t
        );
    END LOOP;
END $$;

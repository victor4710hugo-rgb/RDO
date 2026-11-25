-- RDO | Etapa D2 | Especificacao SQLite v1.0
-- PROPOSTA PARA REVISAO. NAO EXECUTAR NO BANCO EXISTENTE SEM PLANO DE MIGRACAO.
-- SQLite recomendado: PRAGMA foreign_keys = ON em TODA conexao.
-- Datas: TEXT ISO-8601; datas operacionais YYYY-MM-DD; horas HH:MM.
-- Dinheiro: INTEGER em centavos. IDs: UUID armazenado em TEXT.

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    checksum TEXT,
    applied_at TEXT NOT NULL
);

CREATE TABLE clients (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    document TEXT,
    email TEXT,
    phone TEXT,
    address TEXT,
    city TEXT,
    state TEXT,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    deleted_at TEXT
);

CREATE TABLE units (
    id TEXT PRIMARY KEY,
    client_id TEXT NOT NULL,
    name TEXT NOT NULL,
    code TEXT,
    address TEXT,
    city TEXT,
    state TEXT,
    latitude REAL,
    longitude REAL,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    deleted_at TEXT,
    UNIQUE (id, client_id),
    FOREIGN KEY (client_id) REFERENCES clients(id)
);

CREATE TABLE organizations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    document TEXT,
    organization_type TEXT NOT NULL CHECK (organization_type IN ('internal','contractor')),
    email TEXT,
    phone TEXT,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    organization_id TEXT,
    name TEXT NOT NULL,
    email TEXT NOT NULL COLLATE NOCASE UNIQUE,
    phone TEXT,
    job_title TEXT,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive','blocked')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    deleted_at TEXT,
    FOREIGN KEY (organization_id) REFERENCES organizations(id)
);

CREATE TABLE workers (
    id TEXT PRIMARY KEY,
    organization_id TEXT,
    user_id TEXT UNIQUE,
    name TEXT NOT NULL,
    document TEXT,
    job_title TEXT,
    worker_type TEXT NOT NULL CHECK (worker_type IN ('internal','contractor')),
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE teams (
    id TEXT PRIMARY KEY,
    organization_id TEXT,
    name TEXT NOT NULL,
    team_type TEXT NOT NULL CHECK (team_type IN ('internal','contractor')),
    description TEXT,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations(id)
);

CREATE TABLE team_members (
    id TEXT PRIMARY KEY,
    team_id TEXT NOT NULL,
    worker_id TEXT NOT NULL,
    member_role TEXT,
    start_date TEXT NOT NULL,
    end_date TEXT,
    created_at TEXT NOT NULL,
    CHECK (end_date IS NULL OR end_date >= start_date),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (worker_id) REFERENCES workers(id)
);

CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    client_id TEXT NOT NULL,
    name TEXT NOT NULL,
    code TEXT,
    description TEXT,
    manager_user_id TEXT,
    start_date TEXT,
    end_date TEXT,
    actual_start_date TEXT,
    actual_end_date TEXT,
    contract_value_cents INTEGER CHECK (contract_value_cents IS NULL OR contract_value_cents >= 0),
    workforce_control_mode TEXT NOT NULL DEFAULT 'aggregate'
        CHECK (workforce_control_mode IN ('aggregate','individual','both')),
    status TEXT NOT NULL DEFAULT 'planned'
        CHECK (status IN ('planned','in_progress','on_hold','completed','cancelled')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    deleted_at TEXT,
    UNIQUE (id, client_id),
    FOREIGN KEY (client_id) REFERENCES clients(id),
    FOREIGN KEY (manager_user_id) REFERENCES users(id),
    CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date)
);

CREATE TABLE project_units (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    unit_id TEXT NOT NULL,
    client_id TEXT NOT NULL,
    is_primary INTEGER NOT NULL DEFAULT 0 CHECK (is_primary IN (0,1)),
    created_at TEXT NOT NULL,
    UNIQUE (project_id, unit_id),
    UNIQUE (project_id, unit_id, client_id),
    FOREIGN KEY (project_id, client_id) REFERENCES projects(id, client_id),
    FOREIGN KEY (unit_id, client_id) REFERENCES units(id, client_id)
);

CREATE TABLE work_orders (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    unit_id TEXT,
    number TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    responsible_user_id TEXT,
    priority TEXT NOT NULL DEFAULT 'medium'
        CHECK (priority IN ('low','medium','high','critical')),
    status TEXT NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft','planned','in_progress','blocked','completed','cancelled')),
    planned_start_date TEXT,
    planned_end_date TEXT,
    actual_start_date TEXT,
    actual_end_date TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    deleted_at TEXT,
    UNIQUE (project_id, number),
    UNIQUE (id, project_id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (project_id, unit_id) REFERENCES project_units(project_id, unit_id),
    FOREIGN KEY (responsible_user_id) REFERENCES users(id),
    CHECK (planned_end_date IS NULL OR planned_start_date IS NULL OR planned_end_date >= planned_start_date),
    CHECK (actual_end_date IS NULL OR actual_start_date IS NULL OR actual_end_date >= actual_start_date)
);

-- A FK composta em work_orders garante que a unidade opcional esteja vinculada ao projeto.

CREATE TABLE work_order_assignments (
    id TEXT PRIMARY KEY,
    work_order_id TEXT NOT NULL,
    team_id TEXT,
    user_id TEXT,
    assignment_role TEXT,
    start_date TEXT,
    end_date TEXT,
    created_at TEXT NOT NULL,
    CHECK ((team_id IS NOT NULL AND user_id IS NULL) OR (team_id IS NULL AND user_id IS NOT NULL)),
    CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date),
    FOREIGN KEY (work_order_id) REFERENCES work_orders(id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE activities (
    id TEXT PRIMARY KEY,
    work_order_id TEXT NOT NULL,
    code TEXT,
    title TEXT NOT NULL,
    description TEXT,
    sequence INTEGER,
    priority TEXT NOT NULL DEFAULT 'medium'
        CHECK (priority IN ('low','medium','high','critical')),
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending','in_progress','blocked','done','cancelled')),
    planned_progress REAL NOT NULL DEFAULT 0 CHECK (planned_progress BETWEEN 0 AND 100),
    actual_progress REAL NOT NULL DEFAULT 0 CHECK (actual_progress BETWEEN 0 AND 100),
    planned_start_date TEXT,
    planned_end_date TEXT,
    actual_start_date TEXT,
    actual_end_date TEXT,
    weight REAL CHECK (weight IS NULL OR weight > 0),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    deleted_at TEXT,
    UNIQUE (id, work_order_id),
    FOREIGN KEY (work_order_id) REFERENCES work_orders(id)
);

CREATE TABLE activity_assignments (
    id TEXT PRIMARY KEY,
    activity_id TEXT NOT NULL,
    team_id TEXT,
    user_id TEXT,
    assignment_role TEXT,
    start_date TEXT,
    end_date TEXT,
    created_at TEXT NOT NULL,
    CHECK ((team_id IS NOT NULL AND user_id IS NULL) OR (team_id IS NULL AND user_id IS NOT NULL)),
    CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date),
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE activity_dependencies (
    id TEXT PRIMARY KEY,
    activity_id TEXT NOT NULL,
    depends_on_activity_id TEXT NOT NULL,
    dependency_type TEXT NOT NULL
        CHECK (dependency_type IN ('finish_to_start','start_to_start','finish_to_finish','start_to_finish')),
    created_at TEXT NOT NULL,
    CHECK (activity_id <> depends_on_activity_id),
    UNIQUE (activity_id, depends_on_activity_id, dependency_type),
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    FOREIGN KEY (depends_on_activity_id) REFERENCES activities(id)
);

CREATE TABLE plans (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    plan_type TEXT NOT NULL CHECK (plan_type IN ('daily','weekly','monthly','yearly')),
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','active','completed','cancelled')),
    created_by_user_id TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    CHECK (end_date >= start_date),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (created_by_user_id) REFERENCES users(id)
);

CREATE TABLE plan_items (
    id TEXT PRIMARY KEY,
    plan_id TEXT NOT NULL,
    activity_id TEXT NOT NULL,
    team_id TEXT,
    planned_date TEXT NOT NULL,
    planned_hours REAL CHECK (planned_hours IS NULL OR planned_hours >= 0),
    planned_quantity REAL CHECK (planned_quantity IS NULL OR planned_quantity >= 0),
    notes TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (plan_id) REFERENCES plans(id),
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    FOREIGN KEY (team_id) REFERENCES teams(id)
);

CREATE TABLE field_reports (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    report_number TEXT NOT NULL,
    report_date TEXT NOT NULL,
    area TEXT,
    weather_condition TEXT,
    weather_notes TEXT,
    shift_start TEXT,
    shift_end TEXT,
    general_notes TEXT,
    next_shift_instructions TEXT,
    machine_condition TEXT CHECK (machine_condition IS NULL OR machine_condition IN ('operational','partially_operational','blocked')),
    status TEXT NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft','submitted','approved','rejected','closed')),
    version INTEGER NOT NULL DEFAULT 1 CHECK (version >= 1),
    payload TEXT,
    created_by_user_id TEXT,
    submitted_at TEXT,
    approved_at TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE (project_id, report_number),
    UNIQUE (id, project_id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (created_by_user_id) REFERENCES users(id)
);

CREATE TABLE field_report_units (
    id TEXT PRIMARY KEY,
    field_report_id TEXT NOT NULL,
    project_id TEXT NOT NULL,
    unit_id TEXT NOT NULL,
    area TEXT,
    created_at TEXT NOT NULL,
    UNIQUE (field_report_id, unit_id),
    FOREIGN KEY (field_report_id, project_id) REFERENCES field_reports(id, project_id),
    FOREIGN KEY (project_id, unit_id) REFERENCES project_units(project_id, unit_id)
);

CREATE TABLE field_report_work_orders (
    id TEXT PRIMARY KEY,
    field_report_id TEXT NOT NULL,
    project_id TEXT NOT NULL,
    work_order_id TEXT NOT NULL,
    notes TEXT,
    created_at TEXT NOT NULL,
    UNIQUE (field_report_id, work_order_id),
    FOREIGN KEY (field_report_id, project_id) REFERENCES field_reports(id, project_id),
    FOREIGN KEY (work_order_id, project_id) REFERENCES work_orders(id, project_id)
);

CREATE TABLE field_report_teams (
    id TEXT PRIMARY KEY,
    field_report_id TEXT NOT NULL,
    team_id TEXT NOT NULL,
    workforce_mode TEXT NOT NULL CHECK (workforce_mode IN ('aggregate','individual')),
    workers_count INTEGER CHECK (workers_count IS NULL OR workers_count >= 0),
    shift_code TEXT,
    shift_start TEXT,
    shift_end TEXT,
    break_minutes INTEGER CHECK (break_minutes IS NULL OR break_minutes >= 0),
    hours_worked REAL CHECK (hours_worked IS NULL OR hours_worked >= 0),
    notes TEXT,
    created_at TEXT NOT NULL,
    CHECK (
      (workforce_mode = 'aggregate' AND workers_count IS NOT NULL) OR
      (workforce_mode = 'individual' AND workers_count IS NULL)
    ),
    FOREIGN KEY (field_report_id) REFERENCES field_reports(id),
    FOREIGN KEY (team_id) REFERENCES teams(id)
);

CREATE TABLE field_report_workers (
    id TEXT PRIMARY KEY,
    field_report_id TEXT NOT NULL,
    field_report_team_id TEXT NOT NULL,
    worker_id TEXT NOT NULL,
    shift_start TEXT,
    shift_end TEXT,
    break_minutes INTEGER CHECK (break_minutes IS NULL OR break_minutes >= 0),
    hours_worked REAL CHECK (hours_worked IS NULL OR hours_worked >= 0),
    attendance_status TEXT NOT NULL DEFAULT 'present'
        CHECK (attendance_status IN ('present','absent','partial')),
    notes TEXT,
    created_at TEXT NOT NULL,
    UNIQUE (field_report_team_id, worker_id),
    FOREIGN KEY (field_report_id) REFERENCES field_reports(id),
    FOREIGN KEY (field_report_team_id) REFERENCES field_report_teams(id),
    FOREIGN KEY (worker_id) REFERENCES workers(id)
);

CREATE TABLE field_report_activities (
    id TEXT PRIMARY KEY,
    field_report_id TEXT NOT NULL,
    activity_id TEXT NOT NULL,
    team_id TEXT,
    execution_sequence INTEGER NOT NULL DEFAULT 1,
    planned_progress REAL CHECK (planned_progress IS NULL OR planned_progress BETWEEN 0 AND 100),
    actual_progress REAL CHECK (actual_progress IS NULL OR actual_progress BETWEEN 0 AND 100),
    hours_worked REAL CHECK (hours_worked IS NULL OR hours_worked >= 0),
    quantity_executed REAL CHECK (quantity_executed IS NULL OR quantity_executed >= 0),
    unit TEXT,
    status TEXT NOT NULL CHECK (status IN ('not_started','in_progress','blocked','completed')),
    notes TEXT,
    non_completion_reason TEXT,
    created_at TEXT NOT NULL,
    FOREIGN KEY (field_report_id) REFERENCES field_reports(id),
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    FOREIGN KEY (team_id) REFERENCES teams(id)
);

CREATE UNIQUE INDEX uq_report_activity_team_sequence
ON field_report_activities(field_report_id, activity_id, team_id, execution_sequence)
WHERE team_id IS NOT NULL;
CREATE UNIQUE INDEX uq_report_activity_no_team_sequence
ON field_report_activities(field_report_id, activity_id, execution_sequence)
WHERE team_id IS NULL;

CREATE TABLE occurrences (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    field_report_id TEXT,
    work_order_id TEXT,
    activity_id TEXT,
    reported_by_user_id TEXT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    occurrence_type TEXT NOT NULL,
    severity TEXT NOT NULL CHECK (severity IN ('low','medium','high','critical')),
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open','in_progress','resolved','cancelled')),
    occurred_at TEXT NOT NULL,
    resolved_at TEXT,
    resolution_notes TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (field_report_id, project_id) REFERENCES field_reports(id, project_id),
    FOREIGN KEY (work_order_id, project_id) REFERENCES work_orders(id, project_id),
    FOREIGN KEY (activity_id, work_order_id) REFERENCES activities(id, work_order_id),
    CHECK (activity_id IS NULL OR work_order_id IS NOT NULL),
    FOREIGN KEY (reported_by_user_id) REFERENCES users(id)
);

CREATE TABLE equipment (
    id TEXT PRIMARY KEY,
    organization_id TEXT,
    name TEXT NOT NULL,
    code TEXT,
    equipment_type TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive','maintenance')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations(id)
);

CREATE TABLE field_report_equipment (
    id TEXT PRIMARY KEY,
    field_report_id TEXT NOT NULL,
    equipment_id TEXT NOT NULL,
    quantity REAL NOT NULL DEFAULT 1 CHECK (quantity > 0),
    start_time TEXT,
    end_time TEXT,
    condition_notes TEXT,
    created_at TEXT NOT NULL,
    FOREIGN KEY (field_report_id) REFERENCES field_reports(id),
    FOREIGN KEY (equipment_id) REFERENCES equipment(id)
);

CREATE TABLE media_files (
    id TEXT PRIMARY KEY,
    file_name TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    file_size INTEGER NOT NULL CHECK (file_size >= 0),
    storage_type TEXT NOT NULL CHECK (storage_type IN ('local','object_storage')),
    storage_path TEXT NOT NULL,
    checksum TEXT,
    uploaded_at TEXT,
    created_by_user_id TEXT,
    created_at TEXT NOT NULL,
    FOREIGN KEY (created_by_user_id) REFERENCES users(id)
);

-- Vínculo explícito para permitir anexos em respostas de checklist e outras entidades.
-- entity_type/entity_id é polimórfico: a existência do pai é validada pela camada de serviço.
CREATE TABLE media_attachments (
    id TEXT PRIMARY KEY,
    media_file_id TEXT NOT NULL,
    entity_type TEXT NOT NULL CHECK (entity_type IN ('field_report','occurrence','work_order','activity','checklist_answer','signature')),
    entity_id TEXT NOT NULL,
    created_at TEXT NOT NULL,
    UNIQUE (media_file_id, entity_type, entity_id),
    FOREIGN KEY (media_file_id) REFERENCES media_files(id)
);

CREATE TABLE signatures (
    id TEXT PRIMARY KEY,
    field_report_id TEXT NOT NULL,
    signer_user_id TEXT,
    signer_name TEXT NOT NULL,
    signer_role TEXT NOT NULL,
    signature_type TEXT NOT NULL CHECK (signature_type IN ('internal','client','supervisor','other')),
    signature_method TEXT NOT NULL CHECK (signature_method IN ('drawn','typed','certified')),
    signature_file_id TEXT,
    provider TEXT,
    external_reference TEXT,
    signed_at TEXT,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','signed','rejected')),
    verification_data TEXT,
    notes TEXT,
    created_at TEXT NOT NULL,
    FOREIGN KEY (field_report_id) REFERENCES field_reports(id),
    FOREIGN KEY (signer_user_id) REFERENCES users(id),
    FOREIGN KEY (signature_file_id) REFERENCES media_files(id)
);

CREATE TABLE checklist_templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    category TEXT NOT NULL CHECK (category IN ('safety','quality','operational')),
    version INTEGER NOT NULL DEFAULT 1 CHECK (version >= 1),
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE (name, version)
);

CREATE TABLE checklist_items (
    id TEXT PRIMARY KEY,
    template_id TEXT NOT NULL,
    description TEXT NOT NULL,
    sequence INTEGER NOT NULL,
    is_critical INTEGER NOT NULL DEFAULT 0 CHECK (is_critical IN (0,1)),
    required INTEGER NOT NULL DEFAULT 1 CHECK (required IN (0,1)),
    created_at TEXT NOT NULL,
    UNIQUE (template_id, sequence),
    FOREIGN KEY (template_id) REFERENCES checklist_templates(id)
);

CREATE TABLE field_report_checklists (
    id TEXT PRIMARY KEY,
    field_report_id TEXT NOT NULL,
    template_id TEXT NOT NULL,
    template_version INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','completed')),
    completed_by_user_id TEXT,
    completed_at TEXT,
    created_at TEXT NOT NULL,
    FOREIGN KEY (field_report_id) REFERENCES field_reports(id),
    FOREIGN KEY (template_id) REFERENCES checklist_templates(id),
    FOREIGN KEY (completed_by_user_id) REFERENCES users(id)
);

CREATE TABLE field_report_checklist_answers (
    id TEXT PRIMARY KEY,
    field_report_checklist_id TEXT NOT NULL,
    checklist_item_id TEXT NOT NULL,
    answer TEXT NOT NULL CHECK (answer IN ('compliant','non_compliant','not_applicable','not_checked')),
    notes TEXT,
    answered_by_user_id TEXT,
    answered_at TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE (field_report_checklist_id, checklist_item_id),
    FOREIGN KEY (field_report_checklist_id) REFERENCES field_report_checklists(id),
    FOREIGN KEY (checklist_item_id) REFERENCES checklist_items(id),
    FOREIGN KEY (answered_by_user_id) REFERENCES users(id)
);

CREATE TABLE roles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TEXT NOT NULL
);

CREATE TABLE permissions (
    id TEXT PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE role_permissions (
    role_id TEXT NOT NULL,
    permission_id TEXT NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id)
);

CREATE TABLE user_roles (
    user_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    created_at TEXT NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE project_memberships (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    access_level TEXT NOT NULL CHECK (access_level IN ('read','contribute','manage')),
    start_date TEXT,
    end_date TEXT,
    created_at TEXT NOT NULL,
    CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE devices (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    device_name TEXT,
    platform TEXT NOT NULL CHECK (platform IN ('android','ios','desktop')),
    app_version TEXT,
    last_seen_at TEXT,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive')),
    created_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE sync_records (
    id TEXT PRIMARY KEY,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    device_id TEXT,
    operation TEXT NOT NULL CHECK (operation IN ('create','update','delete')),
    status TEXT NOT NULL CHECK (status IN ('pending','syncing','synced','error','conflict')),
    record_version INTEGER NOT NULL DEFAULT 1 CHECK (record_version >= 1),
    attempts INTEGER NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    last_error TEXT,
    synced_at TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (device_id) REFERENCES devices(id)
);

CREATE TABLE sync_conflicts (
    id TEXT PRIMARY KEY,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    device_id TEXT,
    server_version INTEGER NOT NULL,
    client_version INTEGER NOT NULL,
    client_payload TEXT NOT NULL,
    server_payload TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','resolved')),
    resolution TEXT,
    resolved_by_user_id TEXT,
    resolved_at TEXT,
    created_at TEXT NOT NULL,
    FOREIGN KEY (device_id) REFERENCES devices(id),
    FOREIGN KEY (resolved_by_user_id) REFERENCES users(id)
);

CREATE TABLE audit_logs (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    device_id TEXT,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    action TEXT NOT NULL CHECK (action IN ('create','update','delete','approve','close','sign')),
    old_values TEXT,
    new_values TEXT,
    created_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (device_id) REFERENCES devices(id)
);

CREATE TABLE transcriptions (
    id TEXT PRIMARY KEY,
    media_file_id TEXT NOT NULL,
    field_report_id TEXT,
    transcription_text TEXT,
    structured_data TEXT,
    processing_status TEXT NOT NULL DEFAULT 'pending'
        CHECK (processing_status IN ('pending','processing','completed','error')),
    review_status TEXT NOT NULL DEFAULT 'pending_review'
        CHECK (review_status IN ('pending_review','approved','rejected')),
    reviewed_by_user_id TEXT,
    reviewed_at TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (media_file_id) REFERENCES media_files(id),
    FOREIGN KEY (field_report_id) REFERENCES field_reports(id),
    FOREIGN KEY (reviewed_by_user_id) REFERENCES users(id)
);

CREATE TABLE notifications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    notification_type TEXT NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    entity_type TEXT,
    entity_id TEXT,
    read_at TEXT,
    created_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Indices para FKs, filtros e consultas operacionais frequentes.
CREATE INDEX idx_units_client ON units(client_id);
CREATE INDEX idx_projects_client_status ON projects(client_id, status);
CREATE INDEX idx_project_units_unit ON project_units(unit_id);
CREATE INDEX idx_work_orders_project_status ON work_orders(project_id, status);
CREATE INDEX idx_work_orders_unit ON work_orders(unit_id);
CREATE INDEX idx_activities_work_order_status ON activities(work_order_id, status);
CREATE INDEX idx_activity_dependencies_depends ON activity_dependencies(depends_on_activity_id);
CREATE INDEX idx_plans_project_period ON plans(project_id, start_date, end_date);
CREATE INDEX idx_plan_items_date ON plan_items(planned_date);
CREATE INDEX idx_field_reports_project_date ON field_reports(project_id, report_date);
CREATE INDEX idx_field_report_work_orders_os ON field_report_work_orders(work_order_id);
CREATE INDEX idx_field_report_teams_report ON field_report_teams(field_report_id);
CREATE INDEX idx_field_report_workers_worker ON field_report_workers(worker_id);
CREATE INDEX idx_field_report_activities_activity ON field_report_activities(activity_id);
CREATE INDEX idx_occurrences_project_status ON occurrences(project_id, status);
CREATE INDEX idx_media_attachments_entity ON media_attachments(entity_type, entity_id);
CREATE INDEX idx_signatures_report_status ON signatures(field_report_id, status);
CREATE INDEX idx_checklist_answers_checklist ON field_report_checklist_answers(field_report_checklist_id);
CREATE INDEX idx_project_memberships_user ON project_memberships(user_id);
CREATE INDEX idx_sync_records_status_updated ON sync_records(status, updated_at);
CREATE INDEX idx_sync_records_entity ON sync_records(entity_type, entity_id);
CREATE INDEX idx_sync_conflicts_status ON sync_conflicts(status);
CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id, created_at);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, read_at, created_at);
CREATE INDEX idx_team_members_worker_dates ON team_members(worker_id, start_date, end_date);

-- Índices únicos parciais: um único vínculo ativo da mesma pessoa na mesma equipe.
CREATE UNIQUE INDEX uq_team_members_active
ON team_members(team_id, worker_id)
WHERE end_date IS NULL;

-- No máximo uma unidade principal por projeto.
CREATE UNIQUE INDEX uq_project_primary_unit
ON project_units(project_id)
WHERE is_primary = 1;

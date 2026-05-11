-- Fix education business rows accidentally created with tenant_id=0.
-- Run after verifying these ids belong to tenant_id=1.

BEGIN;

-- Sanity: inspect target rows before running.
-- SELECT id, name, code, tenant_id FROM edu_room WHERE id IN (1, 2);
-- SELECT id, name, code, tenant_id FROM edu_course WHERE id IN (1, 2);
-- SELECT id, name, phone, tenant_id FROM edu_student WHERE id IN (1, 2);

-- Move current tenant_id=1 room code out of the way, then migrate the intended record.
-- This avoids uk_edu_room_tenant_code conflicts on (tenant_id, code).
UPDATE edu_room
SET code = code || '-OLD-' || id
WHERE id = 1 AND tenant_id = 1 AND code = 'JS-001';

UPDATE edu_room
SET tenant_id = 1
WHERE id = 2 AND tenant_id = 0;

UPDATE edu_course
SET tenant_id = 1
WHERE id = 2 AND tenant_id = 0;

UPDATE edu_student
SET tenant_id = 1
WHERE id = 2 AND tenant_id = 0;

UPDATE edu_room_weekly_rule
SET tenant_id = 1
WHERE room_id = 2 AND tenant_id = 0;

UPDATE edu_room_exception
SET tenant_id = 1
WHERE room_id = 2 AND tenant_id = 0;

UPDATE edu_student_contact
SET tenant_id = 1
WHERE student_id = 2 AND tenant_id = 0;

-- Sanity: these should return 0 rows.
-- SELECT id, name, code, tenant_id FROM edu_room WHERE tenant_id = 0;
-- SELECT id, name, code, tenant_id FROM edu_course WHERE tenant_id = 0;
-- SELECT id, name, phone, tenant_id FROM edu_student WHERE tenant_id = 0;

COMMIT;

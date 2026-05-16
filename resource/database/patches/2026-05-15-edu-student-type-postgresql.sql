-- 学生表增加 student_type 字段
ALTER TABLE edu_student ADD COLUMN IF NOT EXISTS student_type VARCHAR(32) DEFAULT 'formal' NOT NULL;
COMMENT ON COLUMN edu_student.student_type IS '学生类型: formal-正式, trial-体验, vip-VIP';

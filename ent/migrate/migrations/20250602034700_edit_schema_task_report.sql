-- Modify "task_reports" table
ALTER TABLE "public"."task_reports" DROP COLUMN "title", DROP COLUMN "status", DROP COLUMN "progress_percentage", DROP COLUMN "reported_at", DROP COLUMN "issues_encountered", DROP COLUMN "next_steps", DROP COLUMN "estimated_completion";

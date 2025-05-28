-- Modify "projects" table
ALTER TABLE "public"."projects" ALTER COLUMN "start_at" DROP NOT NULL;
-- Modify "tasks" table
ALTER TABLE "public"."tasks" ALTER COLUMN "process" SET DEFAULT 0, ADD COLUMN "due_date" timestamptz NULL;
-- Create index "tasks_code_key" to table: "tasks"
CREATE UNIQUE INDEX "tasks_code_key" ON "public"."tasks" ("code");
-- Create "task_assignees" table
CREATE TABLE "public"."task_assignees" ("task_id" bigint NOT NULL, "employee_id" bigint NOT NULL, PRIMARY KEY ("task_id", "employee_id"), CONSTRAINT "task_assignees_employee_id" FOREIGN KEY ("employee_id") REFERENCES "public"."employees" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "task_assignees_task_id" FOREIGN KEY ("task_id") REFERENCES "public"."tasks" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);

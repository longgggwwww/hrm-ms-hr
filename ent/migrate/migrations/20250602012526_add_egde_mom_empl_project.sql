-- Modify "employees" table
ALTER TABLE "public"."employees" DROP COLUMN "project_members";
-- Create "project_members" table
CREATE TABLE "public"."project_members" ("project_id" bigint NOT NULL, "employee_id" bigint NOT NULL, PRIMARY KEY ("project_id", "employee_id"), CONSTRAINT "project_members_employee_id" FOREIGN KEY ("employee_id") REFERENCES "public"."employees" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "project_members_project_id" FOREIGN KEY ("project_id") REFERENCES "public"."projects" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);

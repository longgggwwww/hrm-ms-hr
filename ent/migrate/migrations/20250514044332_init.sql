-- Create "companies" table
CREATE TABLE "public"."companies" ("id" uuid NOT NULL, "name" character varying NOT NULL, "code" character varying NOT NULL, "address" character varying NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create index "companies_code_key" to table: "companies"
CREATE UNIQUE INDEX "companies_code_key" ON "public"."companies" ("code");
-- Create "branches" table
CREATE TABLE "public"."branches" ("id" uuid NOT NULL, "name" character varying NOT NULL, "code" character varying NOT NULL, "address" character varying NULL, "contact_info" character varying NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "company_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "branches_companies_branches" FOREIGN KEY ("company_id") REFERENCES "public"."companies" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "branches_code_key" to table: "branches"
CREATE UNIQUE INDEX "branches_code_key" ON "public"."branches" ("code");
-- Create "departments" table
CREATE TABLE "public"."departments" ("id" uuid NOT NULL, "name" character varying NOT NULL, "code" character varying NOT NULL, "branch_id" uuid NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create index "departments_code_key" to table: "departments"
CREATE UNIQUE INDEX "departments_code_key" ON "public"."departments" ("code");
-- Create "positions" table
CREATE TABLE "public"."positions" ("id" uuid NOT NULL, "name" character varying NOT NULL, "code" character varying NOT NULL, "parent_id" uuid NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "department_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "positions_departments_positions" FOREIGN KEY ("department_id") REFERENCES "public"."departments" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "positions_code_key" to table: "positions"
CREATE UNIQUE INDEX "positions_code_key" ON "public"."positions" ("code");
-- Create "employees" table
CREATE TABLE "public"."employees" ("id" uuid NOT NULL, "user_id" character varying NOT NULL, "code" character varying NOT NULL, "status" boolean NOT NULL, "joining_at" timestamptz NOT NULL, "branch_id" uuid NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "department_id" uuid NOT NULL, "position_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "employees_positions_employees" FOREIGN KEY ("position_id") REFERENCES "public"."positions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "employees_code_key" to table: "employees"
CREATE UNIQUE INDEX "employees_code_key" ON "public"."employees" ("code");
-- Create index "employees_user_id_key" to table: "employees"
CREATE UNIQUE INDEX "employees_user_id_key" ON "public"."employees" ("user_id");

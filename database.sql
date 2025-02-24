-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.
-- This is test table. Remove this table and replace with your own tables. 
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
	"estates" (
		"id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4 ()),
		"length" integer NOT NULL,
		"width" integer NOT NULL,
		"created_at" timestamp NOT NULL DEFAULT (now ()),
		"updated_at" timestamp NOT NULL DEFAULT (now ()),
		"deleted_at" timestamp
	);

CREATE TABLE
	"trees" (
		"id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4 ()),
		"estate_id" uuid NOT NULL,
		"x" integer NOT NULL,
		"y" integer NOT NULL,
		"height" integer NOT NULL,
		"created_at" timestamp NOT NULL DEFAULT (now ()),
		"updated_at" timestamp NOT NULL DEFAULT (now ()),
		"deleted_at" timestamp
	);

CREATE UNIQUE INDEX ON "trees" ("estate_id", "x", "y");

ALTER TABLE "trees" ADD FOREIGN KEY ("estate_id") REFERENCES "estates" ("id");
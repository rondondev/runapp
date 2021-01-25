CREATE TYPE "user_type" AS ENUM (
  'admin',
  'coach',
  'athlete'
);

CREATE TYPE "training_sport" AS ENUM (
  'running',
  'cycling',
  'swimming',
  'weight'
);

CREATE TYPE "training_status" AS ENUM (
  'new',
  'notified',
  'overdue',
  'done',
  'done_feedback'
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "type" user_type NOT NULL,
  "name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "password_hash" varchar NOT NULL,
  "phone" varchar,
  "birth" date,
  "active" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz
);

CREATE TABLE "training" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "date" date NOT NULL,
  "sport" training_sport NOT NULL,
  "type" varchar,
  "intensity" varchar,
  "details" varchar NOT NULL,
  "status" training_status NOT NULL DEFAULT 'new',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz
);

CREATE TABLE "training_feedback" (
  "id" bigserial PRIMARY KEY,
  "training_id" bigint NOT NULL,
  "borg_scale" int NOT NULL
);

ALTER TABLE "training" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "training_feedback" ADD FOREIGN KEY ("training_id") REFERENCES "training" ("id");

CREATE INDEX ON "training" ("user_id");

-- Create "todos" table
CREATE TABLE "public"."todos" (
  "id" serial NOT NULL,
  "deadline" timestamp NULL,
  "task" text NOT NULL,
  "is_completed" boolean NOT NULL DEFAULT false,
  "todo_parent" integer NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "todos_todos_parent" FOREIGN KEY ("todo_parent") REFERENCES "public"."todos" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

{{ define "user/schema" -}}

create extension if not exists "uuid-ossp";

create schema if not exists iam;

create table if not exists "iam"."user" (
  id uuid primary key default uuid_generate_v4() not null,
  email text unique not null,
  family_name text,
  given_name text,
  city text,
  age integer,
  create_author uuid not null,
  create_time timestamptz,
  update_author uuid,
  update_time timestamptz,
  delete_author uuid,
  delete_time timestamptz
);

create unique index if not exists unique_email ON "iam"."user" using btree(email) where delete_time is not null;

{{- end }}

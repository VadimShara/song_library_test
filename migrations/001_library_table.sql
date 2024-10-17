-- +goose Up
CREATE TABLE library(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "group" CHARACTER VARYING(30),
    song CHARACTER VARYING(30),
    releaseDate CHARACTER VARYING(15),
    "text" TEXT,
    link CHARACTER VARYING(50)
);
CREATE TABLE IF NOT EXISTS secret
(
    id uuid NOT NULL,
    meta_data jsonb,
    context bytea,
    user_id bigint,
    version bigint,
    added timestamp with time zone,
    updated timestamp with time zone,
    CONSTRAINT secret_pkey PRIMARY KEY (id)
)
CREATE TABLE IF NOT EXISTS public.links3
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    link text COLLATE pg_catalog."default",
    short text COLLATE pg_catalog."default",
    CONSTRAINT links3_pkey PRIMARY KEY (id),
    CONSTRAINT links3_link_key UNIQUE (link),
    CONSTRAINT links3_short_key UNIQUE (short)
    )

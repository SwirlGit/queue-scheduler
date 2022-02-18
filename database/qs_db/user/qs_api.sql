CREATE USER qs_api WITH
    LOGIN
    NOSUPERUSER
    NOCREATEDB
    NOINHERIT
    NOREPLICATION
    PASSWORD 'qs_api';

COMMENT ON ROLE qs_api
    IS 'api user';

GRANT CONNECT ON DATABASE qs_db TO qs_api;

GRANT USAGE ON SEQUENCE public.queues_id_seq TO qs_api;
GRANT INSERT, SELECT ON TABLE public.queues TO qs_api;
GRANT USAGE ON TYPE public.STATE TO qs_api;
GRANT USAGE ON SEQUENCE public.jobs_id_seq TO qs_api;
GRANT INSERT ON TABLE public.jobs TO qs_api;
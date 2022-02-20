CREATE USER qs_checker WITH
    LOGIN
    NOSUPERUSER
    NOCREATEDB
    NOINHERIT
    NOREPLICATION
    PASSWORD 'qs_checker';

COMMENT ON ROLE qs_checker
    IS 'checker user';

GRANT CONNECT ON DATABASE qs_db TO qs_checker;

GRANT USAGE ON TYPE public.QUEUE_STATE TO qs_checker;
GRANT USAGE ON SEQUENCE public.queues_id_seq TO qs_checker;
GRANT SELECT, UPDATE ON TABLE public.queues TO qs_checker;
GRANT USAGE ON TYPE public.JOB_STATE TO qs_checker;
GRANT USAGE ON SEQUENCE public.jobs_id_seq TO qs_checker;
GRANT SELECT, UPDATE ON TABLE public.jobs TO qs_checker;
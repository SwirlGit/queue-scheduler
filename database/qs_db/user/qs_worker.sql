CREATE USER qs_worker WITH
    LOGIN
    NOSUPERUSER
    NOCREATEDB
    NOINHERIT
    NOREPLICATION
    PASSWORD 'qs_worker';

COMMENT ON ROLE qs_worker
    IS 'worker user';

GRANT CONNECT ON DATABASE qs_db TO qs_worker;

GRANT USAGE ON TYPE public.QUEUE_STATE TO qs_worker;
GRANT USAGE ON SEQUENCE public.queues_id_seq TO qs_worker;
GRANT SELECT, UPDATE ON TABLE public.queues TO qs_worker;
GRANT USAGE ON TYPE public.JOB_STATE TO qs_worker;
GRANT USAGE ON SEQUENCE public.jobs_id_seq TO qs_worker;
GRANT SELECT, UPDATE ON TABLE public.jobs TO qs_worker;
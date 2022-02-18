CREATE TABLE public.jobs
(
    id              BIGSERIAL PRIMARY KEY,
    ref_queue_id    BIGINT      NOT NULL,
    date_time       TIMESTAMPTZ NOT NULL,
    action          VARCHAR     NOT NULL,
    state           STATE       NOT NULL DEFAULT 'new':: STATE,
    last_heart_beat TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_query_id FOREIGN KEY (ref_query_id) REFERENCES queries (id)
);

CREATE INDEX idx_date_time_running ON public.jobs (date_time) WHERE state = 'new'::STATE;
CREATE TABLE public.jobs
(
    id              BIGSERIAL PRIMARY KEY,
    ref_queue_id    BIGINT      NOT NULL,
    date_time       TIMESTAMPTZ NOT NULL,
    action          VARCHAR     NOT NULL,
    state           JOB_STATE   NOT NULL DEFAULT 'new':: JOB_STATE,
    last_heart_beat TIMESTAMPTZ,
    CONSTRAINT fk_queue_id FOREIGN KEY (ref_queue_id) REFERENCES queues (id)
);

CREATE INDEX idx_date_time_running ON public.jobs (date_time) WHERE state = 'new'::JOB_STATE;
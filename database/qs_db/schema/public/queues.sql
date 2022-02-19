CREATE TABLE public.queues
(
    id       SERIAL PRIMARY KEY,
    queue_id VARCHAR     NOT NULL,
    state    QUEUE_STATE NOT NULL DEFAULT 'ready'::QUEUE_STATE
);

CREATE UNIQUE INDEX idx_queue_id ON public.queues (queue_id);
-- name: CreatePresentation :one
INSERT INTO presentations
  (currentpollindex)
VALUES
  ($1)
RETURNING id;

-- name: GetPresentation :one
SELECT *
FROM presentations
WHERE id = $1;


-- name: CreatePresentationAndPolls :one
WITH presentations_cte AS (
INSERT INTO presentations
  (id, currentpollindex)
VALUES
  (uuid_generate_v4(), 0)
RETURNING id
),
data_cte AS
(
  SELECT
  to_jsonb($1
::jsonb) AS data
),
polls_cte AS
(
INSERT INTO polls
  (question, presentationid, pollindex)
SELECT
  (arr.elem ->> 'question')
::TEXT AS question,
    pc.id AS presentationid,
    arr.idx AS pollindex
  FROM presentations_cte pc, data_cte dc,
    jsonb_array_elements
(dc.data)
WITH ORDINALITY AS arr
(elem, idx)
  RETURNING id, pollindex
),
options_cte AS
(
INSERT INTO options
  (pollid, optionkey, optionvalue)
SELECT
  pc.id,
  (o ->> 'key')
::TEXT AS optionkey,
(o ->> 'value')::TEXT AS optionvalue
  FROM polls_cte pc, data_cte dc,
    LATERAL jsonb_array_elements
(dc.data -> pc.pollindex - 1 -> 'options') AS o
  RETURNING id,pollid
)

SELECT id
FROM presentations_cte;

-- name: GetPollVotes :many
SELECT
  votes.optionkey,
  votes.clientid
FROM votes
  JOIN polls ON votes.pollid = polls.id
  JOIN presentations ON polls.presentationID = presentations.id
WHERE presentations.id = $1
  AND polls.id = $2
ORDER BY votes.optionkey;

-- name: GetPresentationCurrentPoll :one
SELECT
  p.id AS id,
  p.question AS question,
  (
    SELECT jsonb_agg(jsonb_build_object(
      'optionkey', o.optionkey,
      'optionvalue', o.optionvalue
    )) AS options
  FROM options o
  WHERE o.pollid = p.id
  ) AS options
FROM polls p
WHERE p.presentationid = $1 and p.pollindex=(SELECT currentpollindex
  FROM presentations
  WHERE id=$1);


-- name: MoveForwardToNextPoll :one
WITH updated_polls_cte AS(
   UPDATE presentations
	SET currentpollindex = LEAST(currentpollindex + 1, (
	  SELECT max(pollindex)
	  FROM polls
	  WHERE polls.presentationid = $1
	))
	WHERE id = $1
   RETURNING id, currentpollindex
)
SELECT
  p.id AS id,
  p.question AS question,
  (
    SELECT jsonb_agg(jsonb_build_object(
      'optionkey', o.optionkey,
      'optionvalue', o.optionvalue
    )) AS options
  FROM options o
  INNER JOIN polls p ON o.pollid = p.id
  WHERE p.presentationid = upc.id AND p.pollindex=upc.currentpollindex
  ) AS options
FROM polls p
INNER JOIN updated_polls_cte upc ON p.presentationid = upc.id AND p.pollindex=upc.currentpollindex;


-- name: MoveBackwardToPreviousPoll :one
WITH updated_polls_cte AS(
   UPDATE presentations
	SET currentpollindex = GREATEST(currentpollindex - 1, (
	  SELECT min(pollindex)
	  FROM polls
	  WHERE polls.presentationid = $1
	))
	WHERE id = $1
   RETURNING id, currentpollindex
)
SELECT
  p.id AS id,
  p.question AS question,
  (
    SELECT jsonb_agg(jsonb_build_object(
      'optionkey', o.optionkey,
      'optionvalue', o.optionvalue
    )) AS options
  FROM options o
  INNER JOIN polls p ON o.pollid = p.id
  WHERE p.presentationid = upc.id AND p.pollindex=upc.currentpollindex
  ) AS options
FROM polls p
INNER JOIN updated_polls_cte upc ON p.presentationid = upc.id AND p.pollindex=upc.currentpollindex;

-- name: GetPollByPID :one
SELECT *
FROM polls
WHERE id = $1 and presentationid = $2;

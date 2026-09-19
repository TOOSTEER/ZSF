-- 001_init.sql
-- Базовая схема согласно разделу 8 ТЗ.

CREATE TABLE IF NOT EXISTS teams (
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(50) NOT NULL,
    color    VARCHAR(20) NOT NULL,
    host_id  INTEGER -- REFERENCES users(id), добавляется ниже через ALTER,
                      -- т.к. users ссылается на teams (взаимная связь)
);

CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    group_name    VARCHAR(50) NOT NULL,
    role          VARCHAR(20) NOT NULL CHECK (role IN ('host', 'member')),
    team_id       INTEGER REFERENCES teams(id),
    login         VARCHAR(50) UNIQUE,
    password_hash VARCHAR(255), -- только для host
    created_at    TIMESTAMP DEFAULT NOW()
);

ALTER TABLE teams
    ADD CONSTRAINT fk_teams_host FOREIGN KEY (host_id) REFERENCES users(id);

CREATE TABLE IF NOT EXISTS tasks (
    id                 SERIAL PRIMARY KEY,
    title              VARCHAR(200) NOT NULL,
    type               VARCHAR(20) NOT NULL CHECK (type IN ('point', 'extra', 'crossword')),
    points_percent     INTEGER NOT NULL,           -- 12 или 3
    order_index        INTEGER NOT NULL,
    crossword_subtype  VARCHAR(20)                 -- 'classic' | 'circles' | 'wordsearch' | 'text' | NULL
        CHECK (crossword_subtype IS NULL OR crossword_subtype IN ('classic', 'circles', 'wordsearch', 'text')),
    created_at         TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS team_progress (
    id             SERIAL PRIMARY KEY,
    team_id        INTEGER REFERENCES teams(id) ON DELETE CASCADE,
    task_id        INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    is_completed   BOOLEAN DEFAULT FALSE,
    points_awarded INTEGER DEFAULT 0,
    completed_by   INTEGER REFERENCES users(id),
    completed_at   TIMESTAMP,
    UNIQUE (team_id, task_id)
);

CREATE TABLE IF NOT EXISTS crossword_answers (
    id              SERIAL PRIMARY KEY,
    task_id         INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    question_number INTEGER NOT NULL,
    correct_answer  VARCHAR(200) NOT NULL,
    UNIQUE (task_id, question_number)
);

-- Дополнение сверх ТЗ: таблица с текстами вопросов/подсказок кроссворда.
-- Нужна, чтобы фронтенд мог отрисовать список вопросов ДО того, как
-- будут внесены реальные ответы (crossword_answers может быть пустой).
-- TODO: уточнить окончательную структуру, когда придут реальные кроссворды
CREATE TABLE IF NOT EXISTS crossword_questions (
    id              SERIAL PRIMARY KEY,
    task_id         INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    question_number INTEGER NOT NULL,
    question_text   VARCHAR(500) NOT NULL,
    direction       VARCHAR(10) CHECK (direction IS NULL OR direction IN ('across', 'down')),
    answer_length   INTEGER,
    UNIQUE (task_id, question_number)
);

CREATE INDEX IF NOT EXISTS idx_team_progress_team ON team_progress(team_id);
CREATE INDEX IF NOT EXISTS idx_team_progress_task ON team_progress(task_id);
CREATE INDEX IF NOT EXISTS idx_crossword_answers_task ON crossword_answers(task_id);
CREATE INDEX IF NOT EXISTS idx_crossword_questions_task ON crossword_questions(task_id);

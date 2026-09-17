-- 002_seed.sql
-- Тестовые данные согласно разделу 8 ТЗ.

-- 5 команд с цветами (host_id проставляется автоматически при
-- регистрации первого водящего для каждой команды — см. internal/auth).
INSERT INTO teams (name, color, host_id) VALUES
    ('Team 1', '#FF0000', NULL), -- Красный
    ('Team 2', '#008000', NULL), -- Зелёный
    ('Team 3', '#0000FF', NULL), -- Синий
    ('Team 4', '#FFFF00', NULL), -- Жёлтый
    ('Team 5', '#FFC0CB', NULL)  -- Розовый
ON CONFLICT DO NOTHING;

-- 7 точек (type='point', 12% каждая).
-- В разделе 4 ТЗ указано "6 точек", но перечислено 7 названий —
-- используем все перечисленные, как явно указано в примечании к ТЗ.
INSERT INTO tasks (title, type, points_percent, order_index, crossword_subtype) VALUES
    ('RK3',     'point', 12, 1, 'classic'),
    ('RK4',     'point', 12, 2, 'circles'),
    ('RK5',     'point', 12, 3, 'classic'),
    ('RK6',     'point', 12, 4, 'classic'),
    ('RK9',     'point', 12, 5, 'wordsearch'),
    ('RK1/RK2', 'point', 12, 6, 'text'),
    ('SSFRK',   'point', 12, 7, 'circles')
ON CONFLICT DO NOTHING;

-- 6 доп. заданий (type='extra', 3% каждое).
INSERT INTO tasks (title, type, points_percent, order_index, crossword_subtype) VALUES
    ('Extra 1', 'extra', 3, 8,  NULL),
    ('Extra 2', 'extra', 3, 9,  NULL),
    ('Extra 3', 'extra', 3, 10, NULL),
    ('Extra 4', 'extra', 3, 11, NULL),
    ('Extra 5', 'extra', 3, 12, NULL),
    ('Extra 6', 'extra', 3, 13, NULL)
ON CONFLICT DO NOTHING;

-- 7 кроссвордов на общие темы (type='crossword', 3% каждый).
-- Названия — общие темы-заготовки; конкретные вопросы/ответы
-- организатор внесёт позже через crossword_questions/crossword_answers.
-- crossword_subtype назначен по кругу между 4 доступными шаблонами
-- вёрстки (classic/circles/wordsearch/text).
-- TODO: уточнить темы и распределение по шаблонам, когда появится контент
INSERT INTO tasks (title, type, points_percent, order_index, crossword_subtype) VALUES
    ('Crossword 1',        'crossword', 3, 14, 'classic'),
    ('Crossword 2',                'crossword', 3, 15, 'circles'),
    ('Crossword 3',              'crossword', 3, 16, 'wordsearch'),
    ('Crossword 4',          'crossword', 3, 17, 'text'),
    ('Crossword 5',                  'crossword', 3, 18, 'classic'),
    ('Crossword 6',             'crossword', 3, 19, 'circles'),
    ('Crossword 7',        'crossword', 3, 20, 'wordsearch')
ON CONFLICT DO NOTHING;

-- ВАЖНО: правильные ответы для точек и кроссвордов вносятся вручную позже.
-- Таблица crossword_answers намеренно оставлена пустой.
-- Пример вставки одного ответа (закомментировано):
--   INSERT INTO crossword_answers (task_id, question_number, correct_answer)
--   VALUES (1, 1, 'ПРИМЕР');

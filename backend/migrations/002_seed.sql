-- 002_seed.sql — команды, водящие, задания.
-- Имена команд и людей записаны через U&'...' (Unicode escapes),
-- чтобы данные в БД гарантированно оказались в UTF-8 независимо от
-- кодировки клиента (psql на Windows).

SET client_encoding = 'UTF8';
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ============================================================
-- 1. Команды
-- ============================================================
INSERT INTO teams (name, color) VALUES
    (U&'\0421\0438\043D\044F\044F \043A\043E\043C\0430\043D\0434\0430',    '#2563EB'), -- Синяя команда
    (U&'\041A\0440\0430\0441\043D\0430\044F \043A\043E\043C\0430\043D\0434\0430',  '#DC2626'), -- Красная команда
    (U&'\041E\0440\0430\043D\0436\0435\0432\0430\044F \043A\043E\043C\0430\043D\0434\0430','#F97316'), -- Оранжевая команда
    (U&'\0417\0435\043B\0451\043D\0430\044F \043A\043E\043C\0430\043D\0434\0430',  '#16A34A'), -- Зелёная команда
    (U&'\0416\0451\043B\0442\0430\044F \043A\043E\043C\0430\043D\0434\0430',       '#EAB308'); -- Жёлтая команда

-- ============================================================
-- 2. Водящие
--    Вход по ЛОГИНУ (латиница) + паролю.
--    Логины: host_blue, host_red, host_orange, host_green, host_yellow.
-- ============================================================
INSERT INTO users (first_name, last_name, group_name, role, team_id, login, password_hash) VALUES
    (U&'\0421\0430\0448\0430',    U&'\041C\043E\0438\0441\0435\0435\0432',   '—', 'host', 1, 'host_blue',   crypt('SashaMoiseev',     gen_salt('bf'))),
    (U&'\041F\043E\043B\0438\043D\0430', U&'\0428\0435\0432\0447\0435\043D\043A\043E','—', 'host', 2, 'host_red',    crypt('PolinaShevchenko', gen_salt('bf'))),
    (U&'\0413\0435\0440\0430',    U&'\0422\0430\0432\0430\0434\044F\043D',   '—', 'host', 3, 'host_orange', crypt('Armyanchick',      gen_salt('bf'))),
    (U&'\041C\0430\0448\0430',    U&'\0421\0442\0430\0440\043E\0432\0438\043A','—', 'host', 4, 'host_green',  crypt('MashaStarovick',   gen_salt('bf'))),
    (U&'\0414\0438\043C\0430',    U&'\0422\0443\0433\0443\0448\0435\0432',   '—', 'host', 5, 'host_yellow', crypt('DimaTugushev',     gen_salt('bf')));

UPDATE teams
SET host_id = sub.id
FROM (SELECT id, team_id FROM users WHERE role = 'host') AS sub
WHERE teams.id = sub.team_id;

-- ============================================================
-- 3. Задания
-- ============================================================
INSERT INTO tasks (title, type, points_percent, order_index, crossword_subtype) VALUES
    ('RK3',     'point', 12, 1, 'classic'),
    ('RK4',     'point', 12, 2, 'circles'),
    ('RK5',     'point', 12, 3, 'classic'),
    ('RK6',     'point', 12, 4, 'classic'),
    ('RK9',     'point', 12, 5, 'wordsearch'),
    ('RK1/RK2', 'point', 12, 6, 'text'),
    ('SSFRK',   'point', 12, 7, 'circles'),

    ('Extra 1', 'extra', 3, 8,  NULL),
    ('Extra 2', 'extra', 3, 9,  NULL),
    ('Extra 3', 'extra', 3, 10, NULL),
    ('Extra 4', 'extra', 3, 11, NULL),
    ('Extra 5', 'extra', 3, 12, NULL),
    ('Extra 6', 'extra', 3, 13, NULL),

    ('Crossword 1', 'crossword', 3, 14, 'classic'),
    ('Crossword 2', 'crossword', 3, 15, 'circles'),
    ('Crossword 3', 'crossword', 3, 16, 'wordsearch'),
    ('Crossword 4', 'crossword', 3, 17, 'text'),
    ('Crossword 5', 'crossword', 3, 18, 'classic'),
    ('Crossword 6', 'crossword', 3, 19, 'circles'),
    ('Crossword 7', 'crossword', 3, 20, 'wordsearch');
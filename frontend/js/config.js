// config.js — глобальный конфиг фронтенда квеста.

const CONFIG = {
    API_URL: 'http://localhost:8080/api',
    TEAM_COLORS: {
        1: '#FF0000',
        2: '#008000',
        3: '#0000FF',
        4: '#FFFF00',
        5: '#FFC0CB',
    },
    MAX_PERCENT: 100,
    TOTAL_PERCENT: 123,

    // Ключи — слаги (как в data-task), значения — имя HTML-шаблона кроссворда
    CROSSWORD_TYPES: {
        'rk3':    'classic',
        'rk4':    'circles',
        'rk5':    'classic',
        'rk6':    'classic',
        'rk9':    'wordsearch',
        'rk1rk2': 'text',
        'ssfrk':  'circles',
    },

    // Резервный маппинг task_id → слаг, если API не вернул ответ
    // (совпадает с порядком сидов в БД — 002_seed.sql).
    TASK_SLUGS_BY_ID: {
        1: 'rk3',
        2: 'rk4',
        3: 'rk5',
        4: 'rk6',
        5: 'rk9',
        6: 'rk1rk2',
        7: 'ssfrk',
    },

    DESIGN_WIDTH: 1440,
};

function initStageScaling() {
    const wrapper = document.querySelector('.stage-wrapper');
    const stage = document.querySelector('.stage');
    if (!wrapper || !stage) return;

    function applyScale() {
        const wrapperWidth = wrapper.clientWidth || window.innerWidth;
        const scale = wrapperWidth / CONFIG.DESIGN_WIDTH;
        stage.style.transform = `scale(${scale})`;
        const naturalHeight = stage.scrollHeight || stage.offsetHeight;
        wrapper.style.height = `${naturalHeight * scale}px`;
    }

    applyScale();
    window.addEventListener('resize', applyScale);
    window.addEventListener('load', applyScale);
    setTimeout(applyScale, 300);
}

/** Общая функция установки экрана и data-атрибутов на .stage. */
function setCrosswordScreen(screenNumber) {
    const stage = document.querySelector('.stage');
    if (!stage) return;
    stage.dataset.screen = String(screenNumber);
}

/** Приводит название из БД к слагу: "RK3" → "rk3", "РК1/РК2" → "rk1rk2". */
function slugify(str) {
    return String(str || '')
        .toLowerCase()
        .replace(/[\s\/\\_\-\.]+/g, '');
}

/**
 * Устанавливает data-task для .stage ДО запроса к API, чтобы фон применился
 * даже если бэкенд недоступен или вернул ошибку.
 */
function setTaskSlugFromUrl() {
    const stage = document.querySelector('.stage');
    if (!stage) return;

    const params = new URLSearchParams(window.location.search);
    const urlSlug = params.get('slug');
    const taskId = params.get('id');

    if (urlSlug) {
        stage.dataset.task = slugify(urlSlug);
        return;
    }
    if (taskId && CONFIG.TASK_SLUGS_BY_ID[taskId]) {
        stage.dataset.task = CONFIG.TASK_SLUGS_BY_ID[taskId];
        return;
    }
    // Оставим пустым — CSS просто не подставит фон.
}
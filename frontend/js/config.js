// config.js — глобальный конфиг фронтенда квеста.

const CONFIG = {
    // Автоопределение адреса API:
    // — локально (localhost) → http://localhost:8080/api
    // — на сервере (155.212.129.123 или домен) → http://<host>/api
    API_URL: (function () {
        const host = location.hostname;
        if (host === 'localhost' || host === '127.0.0.1') {
            return 'http://localhost:8080/api';
        }
        return location.protocol + '//' + host + '/api';
    })(),

    TEAM_COLORS: {
        1: '#2563EB',   // Синяя команда
        2: '#DC2626',   // Красная команда
        3: '#F97316',   // Оранжевая команда
        4: '#16A34A',   // Зелёная команда
        5: '#EAB308',   // Жёлтая команда
    },
    MAX_PERCENT: 100,
    TOTAL_PERCENT: 123,
    PROGRESS_MAX_FILL_PX: 1221,

    CROSSWORD_TYPES: {
        'rk3':    'classic',
        'rk4':    'circles',
        'rk5':    'classic',
        'rk6':    'classic',
        'rk9':    'wordsearch',
        'rk1rk2': 'text',
        'ssfrk':  'circles',
    },

    TASK_SLUGS_BY_ID: {
        1: 'rk3', 2: 'rk4', 3: 'rk5',
        4: 'rk6', 5: 'rk9', 6: 'rk1rk2', 7: 'ssfrk',
    },

    DESIGN_WIDTH: 1440,
    DESIGN_HEIGHT: 2560,
};

const ANSWER_LAYOUT = {
    rk3: {
        screen2: [
            { n: 1, top: 562,  left: 251, width: 1065 },
            { n: 2, top: 722,  left: 251, width: 1065 },
            { n: 3, top: 882,  left: 251, width: 1065 },
            { n: 4, top: 1042, left: 251, width: 1065 },
            { n: 5, top: 1202, left: 251, width: 1065 },
            { n: 6, top: 1362, left: 251, width: 1065 },
            { n: 7, top: 1522, left: 251, width: 1065 },
        ],
    },
    rk5: {
        screen2: [
            { n: 1, top: 534,  left: 251, width: 1065 },
            { n: 2, top: 694,  left: 251, width: 1065 },
            { n: 3, top: 854,  left: 251, width: 1065 },
            { n: 4, top: 1014, left: 251, width: 1065 },
            { n: 5, top: 1174, left: 251, width: 1065 },
            { n: 6, top: 1334, left: 251, width: 1065 },
            { n: 7, top: 1494, left: 251, width: 1065 },
            { n: 8, top: 1654, left: 251, width: 1065 },
            { n: 9, top: 1814, left: 251, width: 1065 },
        ],
    },
    rk6: {
        screen2: [
            { n: 1, top: 562,  left: 251, width: 1065 },
            { n: 2, top: 722,  left: 251, width: 1065 },
            { n: 3, top: 882,  left: 251, width: 1065 },
            { n: 4, top: 1042, left: 251, width: 1065 },
            { n: 5, top: 1202, left: 251, width: 1065 },
            { n: 6, top: 1362, left: 251, width: 1065 },
            { n: 7, top: 1522, left: 251, width: 1065 },
        ],
    },
    rk9: {
        screen2: [
            { n: 1, top: 562,  left: 251, width: 1065 },
            { n: 2, top: 722,  left: 251, width: 1065 },
            { n: 3, top: 882,  left: 251, width: 1065 },
            { n: 4, top: 1042, left: 251, width: 1065 },
            { n: 5, top: 1202, left: 251, width: 1065 },
            { n: 6, top: 1362, left: 251, width: 1065 },
            { n: 7, top: 1522, left: 251, width: 1065 },
        ],
    },
    rk1rk2: {
        screen2: [
            { n: 1, top: 534,  left: 251, width: 1065 },
            { n: 2, top: 694,  left: 251, width: 1065 },
            { n: 3, top: 854,  left: 251, width: 1065 },
            { n: 4, top: 1014, left: 251, width: 1065 },
            { n: 5, top: 1174, left: 251, width: 1065 },
            { n: 6, top: 1334, left: 251, width: 1065 },
            { n: 7, top: 1494, left: 251, width: 1065 },
            { n: 8, top: 1654, left: 251, width: 1065 },
            { n: 9, top: 1814, left: 251, width: 1065 },
        ],
    },
    rk4: {
        screen1: [
            { n: 1, top: 887,  left: 735, width: 560 },
            { n: 2, top: 1589, left: 150, width: 635 },
            { n: 3, top: 1995, left: 864, width: 520 },
            { n: 4, top: 2447, left: 145, width: 554 },
        ],
        screen2: [
            { n: 5, top: 1009, left: 137, width: 850 },
            { n: 6, top: 1686, left: 866, width: 492 },
            { n: 7, top: 2087, left: 126, width: 622 },
        ],
    },
    ssfrk: {
        screen1: [
            { n: 1, top: 1147, left: 843, width: 508 },
            { n: 2, top: 1540, left: 199, width: 421 },
            { n: 3, top: 2223, left: 651, width: 467 },
        ],
        screen2: [
            { n: 4, top: 654,  left: 118, width: 421 },
            { n: 5, top: 1149, left: 571, width: 805 },
            { n: 6, top: 1744, left: 123, width: 492 },
            { n: 7, top: 2117, left: 754, width: 622 },
        ],
    },
};

function initStageScaling() {
    const wrapper = document.querySelector('.stage-wrapper');
    const stage = document.querySelector('.stage');
    if (!wrapper || !stage) return;

    function applyScale() {
        const wrapperWidth = wrapper.clientWidth || window.innerWidth;
        const scale = wrapperWidth / CONFIG.DESIGN_WIDTH;
        stage.style.transform = `scale(${scale})`;
        wrapper.style.height = `${CONFIG.DESIGN_HEIGHT * scale}px`;
    }

    applyScale();
    window.addEventListener('resize', applyScale);
    window.addEventListener('load', applyScale);
    setTimeout(applyScale, 300);
}

function setCrosswordScreen(screenNumber) {
    const stage = document.querySelector('.stage');
    if (!stage) return;
    stage.dataset.screen = String(screenNumber);
}

function slugify(str) {
    return String(str || '').toLowerCase().replace(/[\s\/\\_\-\.]+/g, '');
}

function setTaskSlugFromUrl() {
    const stage = document.querySelector('.stage');
    if (!stage) return;
    const params = new URLSearchParams(window.location.search);
    const urlSlug = params.get('slug');
    const taskId = params.get('id');
    if (urlSlug) { stage.dataset.task = slugify(urlSlug); return; }
    if (taskId && CONFIG.TASK_SLUGS_BY_ID[taskId]) {
        stage.dataset.task = CONFIG.TASK_SLUGS_BY_ID[taskId];
    }
}

function getCurrentSlug() {
    const stage = document.querySelector('.stage');
    return stage ? stage.dataset.task : '';
}

function renderAnswerFields(slug, screenNumber) {
    const container = document.getElementById('answer-fields-' + screenNumber);
    if (!container) return;

    const layout = ANSWER_LAYOUT[slug];
    const items = layout && layout['screen' + screenNumber];
    container.innerHTML = '';
    if (!items || !items.length) return;

    items.forEach((item) => {
        const input = document.createElement('input');
        input.type = 'text';
        input.className = 'answer-line';
        input.dataset.n = String(item.n);
        input.autocomplete = 'off';
        input.style.top  = (item.top - 60) + 'px';
        input.style.left = item.left + 'px';
        input.style.width = item.width + 'px';
        container.appendChild(input);
    });
}

function collectAnswers() {
    const inputs = Array.from(document.querySelectorAll('.answer-line'));
    inputs.sort((a, b) => parseInt(a.dataset.n, 10) - parseInt(b.dataset.n, 10));
    return inputs.map((inp) => ({
        question_number: parseInt(inp.dataset.n, 10),
        answer: inp.value,
    }));
}

function markWrongAnswers(wrongNumbers) {
    document.querySelectorAll('.answer-line').forEach((inp) => inp.classList.remove('field-error'));
    (wrongNumbers || []).forEach((n) => {
        const inp = document.querySelector(`.answer-line[data-n="${n}"]`);
        if (inp) inp.classList.add('field-error');
    });
}
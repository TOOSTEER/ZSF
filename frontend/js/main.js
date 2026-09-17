// main.js — логика главной страницы (main.html)

document.addEventListener('DOMContentLoaded', async () => {
    initStageScaling();

    if (!requireAuth()) return;

    const user = Storage.getUser();
    if (user && user.role === 'host') {
        window.location.href = 'host.html';
        return;
    }

    const els = {
        teamName: document.getElementById('team-name'),
        teamColorDot: document.getElementById('team-color-dot'),
        progressFill: document.getElementById('progress-fill'),
        progressPercent: document.getElementById('progress-percent'),
        buttons: {
            'rk3':    document.getElementById('btn-rk3'),
            'rk4':    document.getElementById('btn-rk4'),
            'rk5':    document.getElementById('btn-rk5'),
            'rk6':    document.getElementById('btn-rk6'),
            'rk9':    document.getElementById('btn-rk9'),
            'rk1rk2': document.getElementById('btn-rk1rk2'),
            'ssfrk':  document.getElementById('btn-ssfrk'),
        },
    };

    // Навешиваем обработчики кликов сразу — до запроса к API,
    // чтобы кнопки работали даже при пустом прогрессе.
    Object.entries(els.buttons).forEach(([slug, btn]) => {
        if (!btn) return;
        btn.addEventListener('click', () => {
            const taskId = btn.dataset.taskId || CONFIG.TASK_SLUGS_BY_ID[Object.keys(CONFIG.TASK_SLUGS_BY_ID).find(k => CONFIG.TASK_SLUGS_BY_ID[k] === slug)];
            const subtype = CONFIG.CROSSWORD_TYPES[slug];
            if (!subtype) {
                console.warn('Не найден шаблон кроссворда для', slug);
                return;
            }
            const query = taskId
                ? `?id=${taskId}&slug=${slug}`
                : `?slug=${slug}`;
            window.location.href = `crosswords/${subtype}.html${query}`;
        });
    });

    try {
        const me = await Api.me();
        Storage.setUser(me.user);
        Storage.setTeam(me.team);

        const team = me.team;
        if (!team) throw new Error('Команда не найдена');

        els.teamName.textContent = `${me.user.first_name} ${me.user.last_name}\nКоманда: ${team.name}`;
        els.teamName.style.whiteSpace = 'pre-line';
        els.teamColorDot.style.background = team.color;

        const progress = await Api.teamProgress(team.id);
        renderProgressBar(progress);
        renderTaskButtons(progress.items, els.buttons);
    } catch (err) {
        console.error(err);
        if (err.status === 401) {
            Storage.clearAll();
            window.location.href = 'index.html';
        }
    }

    function renderProgressBar(progress) {
        const displayPercent = Math.min(progress.display_percent, CONFIG.MAX_PERCENT);
        els.progressFill.style.width = `${displayPercent}%`;
        const suffix = progress.total_percent > CONFIG.MAX_PERCENT
            ? `${progress.total_percent}% / ${CONFIG.MAX_PERCENT}%`
            : `${progress.total_percent}%`;
        els.progressPercent.textContent = suffix;
    }

    function renderTaskButtons(items, buttons) {
        Object.entries(buttons).forEach(([slug, btn]) => {
            if (!btn) return;
            const item = items.find((i) => slugify(i.title) === slug);
            if (!item) return;

            btn.dataset.taskId = item.task_id;
            btn.classList.toggle('completed', item.is_completed);
        });
    }
});
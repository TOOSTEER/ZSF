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
            'РК3': document.getElementById('btn-rk3'),
            'РК4': document.getElementById('btn-rk4'),
            'РК5': document.getElementById('btn-rk5'),
            'РК6': document.getElementById('btn-rk6'),
            'РК9': document.getElementById('btn-rk9'),
            'РК1/РК2': document.getElementById('btn-rk1rk2'),
            'ССФРК': document.getElementById('btn-ssfrk'),
        },
    };

    try {
        const me = await Api.me();
        Storage.setUser(me.user);
        Storage.setTeam(me.team);

        const team = me.team;
        if (!team) throw new Error('Команда не найдена');

        els.teamName.textContent = team.name;
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
        Object.entries(buttons).forEach(([title, btn]) => {
            if (!btn) return;
            const item = items.find((i) => i.title === title);
            if (!item) {
                btn.textContent = title;
                btn.classList.add('disabled');
                return;
            }

            btn.textContent = title;
            btn.classList.toggle('completed', item.is_completed);

            const subtype = CONFIG.CROSSWORD_TYPES[title] || item.crossword_subtype;
            btn.addEventListener('click', () => {
                window.location.href = `crosswords/${subtype}.html?id=${item.task_id}`;
            });
        });
    }
});

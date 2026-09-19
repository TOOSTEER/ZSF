// host.js — логика панели водящего (host.html)

document.addEventListener('DOMContentLoaded', async () => {
    initStageScaling();

    if (!requireAuth()) return;
    if (!requireRole('host')) return;

    const els = {
        teamColor: document.getElementById('host-team-color'),
        teamTitle: document.getElementById('host-team-title'),
        progressFill: document.getElementById('host-progress-fill'),
        progressPercent: document.getElementById('host-progress-percent'),
        pointsGrid: document.getElementById('grid-points'),
        extraGrid: document.getElementById('grid-extra'),
        crosswordGrid: document.getElementById('grid-crossword'),
        modalOverlay: document.getElementById('award-modal'),
        modalTaskTitle: document.getElementById('award-task-title'),
        modalPointsInput: document.getElementById('award-points-input'),
        modalConfirmBtn: document.getElementById('award-confirm'),
        modalCancelBtn: document.getElementById('award-cancel'),
        modalError: document.getElementById('award-error'),
    };

    let team = null;
    let activeItem = null;

    try {
        const me = await Api.me();
        Storage.setUser(me.user);
        Storage.setTeam(me.team);
        team = me.team;
        if (!team) throw new Error('Команда не найдена');

        els.teamTitle.textContent = team.name;
        els.teamColor.style.background = team.color;

        await reloadProgress();
    } catch (err) {
        console.error(err);
        if (err.status === 401) {
            Storage.clearAll();
            window.location.href = 'index.html';
        }
    }

    async function reloadProgress() {
        const progress = await Api.teamProgress(team.id);
        const pct = Math.min(progress.display_percent, CONFIG.MAX_PERCENT);
        // Внутри host.css заливка сама ограничена max-width: calc(100% - 22px),
        // поэтому здесь достаточно процента.
        els.progressFill.style.width = `${pct}%`;

        els.progressPercent.textContent = progress.total_percent > CONFIG.MAX_PERCENT
            ? `${progress.total_percent}% / ${CONFIG.MAX_PERCENT}%`
            : `${progress.total_percent}%`;

        renderGroup(els.pointsGrid, progress.items.filter((i) => i.type === 'point'));
        renderGroup(els.extraGrid, progress.items.filter((i) => i.type === 'extra'));
        renderGroup(els.crosswordGrid, progress.items.filter((i) => i.type === 'crossword'));
    }

    function renderGroup(container, items) {
        container.innerHTML = '';
        items.forEach((item) => {
            const btn = document.createElement('button');
            btn.className = 'host-task-btn' + (item.is_completed ? ' completed' : '');
            btn.innerHTML = `<span>${item.title}</span><span class="checkmark">✓</span>`;

            if (item.type === 'crossword') {
                btn.disabled = true;
                btn.title = 'Начисляется автоматически';
            } else {
                btn.addEventListener('click', () => openAwardModal(item));
            }
            container.appendChild(btn);
        });
    }

    function openAwardModal(item) {
        if (item.is_completed) return;
        activeItem = item;
        els.modalTaskTitle.textContent = item.title;
        els.modalPointsInput.value = item.points_percent;
        els.modalError.textContent = '';
        els.modalOverlay.classList.remove('hidden');
    }

    function closeAwardModal() {
        els.modalOverlay.classList.add('hidden');
        activeItem = null;
    }

    els.modalCancelBtn.addEventListener('click', closeAwardModal);

    els.modalConfirmBtn.addEventListener('click', async () => {
        if (!activeItem) return;
        const points = parseInt(els.modalPointsInput.value, 10);
        if (!points || points <= 0) {
            els.modalError.textContent = 'Введите корректное количество баллов';
            return;
        }
        try {
            await Api.awardProgress({ team_id: team.id, task_id: activeItem.task_id, points });
            closeAwardModal();
            await reloadProgress();
        } catch (err) {
            els.modalError.textContent = err.status === 409
                ? 'Баллы за это задание уже начислены'
                : (err.message || 'Не удалось начислить баллы');
        }
    });
});
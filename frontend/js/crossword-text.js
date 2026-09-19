// crossword-text.js — тип D: text.html (РК1/РК2)

document.addEventListener('DOMContentLoaded', async () => {
    if (typeof initStageScaling === 'function') initStageScaling();
    if (!requireAuth()) return;
    if (!requireRole('member')) return;

    const params = new URLSearchParams(window.location.search);
    const taskId = params.get('id');
    if (!taskId) {
        document.body.innerHTML = '<p class="error-text">Не указан id задания (?id=...)</p>';
        return;
    }

    if (typeof setTaskSlugFromUrl === 'function') setTaskSlugFromUrl();
    const slug = typeof getCurrentSlug === 'function' ? getCurrentSlug() : '';

    // Явная установка фона для rk1rk2 — на случай если data-атрибут не сработал.
    const stage = document.querySelector('.stage');
    if (stage && slug === 'rk1rk2') {
        const setBg = (n) => {
            const path = `/images/backgrounds/rk1rk2-${n}.png`;
            stage.style.backgroundImage = `url('${path}')`;
        };
        setBg(1);
        stage.dataset.task = slug;
        stage._setBg = setBg;
    }

    // Рендер полей (все на экране 2)
    try {
        if (typeof renderAnswerFields === 'function') {
            renderAnswerFields(slug, 2);
        }
    } catch (e) { console.error('renderAnswerFields failed:', e); }

    const els = {
        stage: stage,
        screen1: document.getElementById('screen-1'),
        screen2: document.getElementById('screen-2'),
        arrowNext: document.getElementById('arrow-next'),
        arrowPrev: document.getElementById('arrow-prev'),
        checkBtn: document.getElementById('check-btn'),
        resultModal: document.getElementById('result-modal'),
        resultText: document.getElementById('result-text'),
        resultClose: document.getElementById('result-close'),
    };

    // Уточняем слаг из API
    try {
        const structure = await Api.crosswordStructure(taskId);
        const apiSlug = typeof slugify === 'function' ? slugify(structure.title) : '';
        if (apiSlug && CONFIG.CROSSWORD_TYPES[apiSlug] && els.stage) {
            els.stage.dataset.task = apiSlug;
            if (els.stage._setBg) els.stage._setBg(1);
        }
    } catch (err) {
        console.error('crosswordStructure error:', err);
    }

    function showScreen(n) {
        if (els.screen1) els.screen1.classList.toggle('active', n === 1);
        if (els.screen2) els.screen2.classList.toggle('active', n === 2);
        if (typeof setCrosswordScreen === 'function') setCrosswordScreen(n);
        if (els.stage && els.stage._setBg) els.stage._setBg(n);
    }
    showScreen(1);

    if (els.arrowNext) els.arrowNext.addEventListener('click', () => showScreen(2));
    if (els.arrowPrev) els.arrowPrev.addEventListener('click', () => showScreen(1));

    // Кнопка «На главную» — это <a href="../main.html">. Никаких preventDefault.

    if (els.checkBtn) els.checkBtn.addEventListener('click', async () => {
        const answers = typeof collectAnswers === 'function' ? collectAnswers() : [];
        const filled = answers.filter((a) => a.answer && a.answer.trim().length > 0);
        if (!filled.length) {
            showResultModal(false, 'Введите хотя бы один ответ');
            return;
        }
        try {
            const result = await Api.checkCrossword(taskId, answers);
            if (result.awarded) {
                showResultModal(true, 'Верно! +3% команде');
            } else {
                if (typeof markWrongAnswers === 'function') markWrongAnswers(result.wrong_numbers || []);
                showResultModal(false, 'Есть ошибки — попробуйте ещё раз');
            }
        } catch (err) {
            if (err.status === 409) {
                showResultModal(true, 'Задание уже выполнено вашей командой');
            } else {
                showResultModal(false, err.message || 'Не удалось проверить ответы');
            }
        }
    });

    function showResultModal(success, text) {
        if (!els.resultModal) return;
        if (els.resultText) {
            els.resultText.textContent = text;
            els.resultText.style.color = success ? 'var(--color-success)' : 'var(--color-danger)';
        }
        els.resultModal.classList.remove('hidden');
    }
    if (els.resultClose) {
        els.resultClose.addEventListener('click', () => els.resultModal && els.resultModal.classList.add('hidden'));
    }
});
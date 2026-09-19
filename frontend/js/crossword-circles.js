// crossword-circles.js — тип B: circles.html (РК4, ССФРК)

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

    // Рендер полей на обоих экранах
    try {
        if (typeof renderAnswerFields === 'function') {
            renderAnswerFields(slug, 1);
            renderAnswerFields(slug, 2);
        }
    } catch (e) { console.error('renderAnswerFields failed:', e); }

    const els = {
        stage: document.querySelector('.stage'),
        screen1: document.getElementById('screen-1'),
        screen2: document.getElementById('screen-2'),
        arrowNext: document.getElementById('arrow-next'),
        arrowPrev: document.getElementById('arrow-prev'),
        checkBtn: document.getElementById('check-btn'),
        resultModal: document.getElementById('result-modal'),
        resultText: document.getElementById('result-text'),
        resultClose: document.getElementById('result-close'),
    };

    // Уточняем слаг из API (на случай если в URL не было slug)
    try {
        const structure = await Api.crosswordStructure(taskId);
        const apiSlug = typeof slugify === 'function' ? slugify(structure.title) : '';
        if (apiSlug && CONFIG.CROSSWORD_TYPES[apiSlug] && els.stage) {
            els.stage.dataset.task = apiSlug;
        }
    } catch (err) { console.error(err); }

    // Показ экрана — фон переключается CSS-селектором [data-task][data-screen].
    // Никаких inline-style, иначе CSS не сработает.
    function showScreen(n) {
        if (els.screen1) els.screen1.classList.toggle('active', n === 1);
        if (els.screen2) els.screen2.classList.toggle('active', n === 2);
        if (typeof setCrosswordScreen === 'function') setCrosswordScreen(n);
    }
    showScreen(1);

    if (els.arrowNext) els.arrowNext.addEventListener('click', () => showScreen(2));
    if (els.arrowPrev) els.arrowPrev.addEventListener('click', () => showScreen(1));

    if (els.checkBtn) els.checkBtn.addEventListener('click', async () => {
        const answers = typeof collectAnswers === 'function' ? collectAnswers() : [];
        const filled = answers.filter((a) => a.answer && a.answer.trim().length > 0);
        if (!filled.length) { showResultModal(false, 'Введите хотя бы один ответ'); return; }
        try {
            const result = await Api.checkCrossword(taskId, answers);
            if (result.awarded) {
                showResultModal(true, 'Верно! +3% команде');
            } else {
                if (typeof markWrongAnswers === 'function') markWrongAnswers(result.wrong_numbers || []);
                showResultModal(false, 'Есть ошибки — попробуйте ещё раз');
            }
        } catch (err) {
            if (err.status === 409) showResultModal(true, 'Задание уже выполнено вашей командой');
            else showResultModal(false, err.message || 'Не удалось проверить ответы');
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
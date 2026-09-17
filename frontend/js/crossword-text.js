// crossword-text.js — тип D: text.html (РК1/РК2)

document.addEventListener('DOMContentLoaded', async () => {
    initStageScaling();
    if (!requireAuth()) return;
    if (!requireRole('member')) return;

    const params = new URLSearchParams(window.location.search);
    const taskId = params.get('id');
    if (!taskId) {
        document.body.innerHTML = '<p class="error-text">Не указан id задания (?id=...)</p>';
        return;
    }

    const els = {
        title: document.getElementById('crossword-title'),
        infoBody: document.getElementById('info-body'),
        answerFields: document.getElementById('answer-fields'),
        screen1: document.getElementById('screen-1'),
        screen2: document.getElementById('screen-2'),
        arrowNext: document.getElementById('arrow-next'),
        arrowPrev: document.getElementById('arrow-prev'),
        checkBtn: document.getElementById('check-btn'),
        backBtn: document.getElementById('back-btn'),
        resultModal: document.getElementById('result-modal'),
        resultText: document.getElementById('result-text'),
        resultClose: document.getElementById('result-close'),
    };

    let questions = [];
    let solved = false;

    try {
        const structure = await Api.crosswordStructure(taskId);
        questions = structure.questions || [];
        els.title.textContent = structure.title || 'Невыпускающие кафедры';
        renderInfo(questions);
        renderAnswerFields(questions);
    } catch (err) {
        console.error(err);
    }

    function renderInfo(qs) {
        if (!qs.length) {
            els.infoBody.innerHTML = '<p>Информация появится позже.</p>';
            return;
        }
        els.infoBody.innerHTML = qs.map((q) => `<p>${q.question_text}</p>`).join('');
    }

    function renderAnswerFields(qs) {
        const list = qs.length ? qs : [{ question_number: 1 }, { question_number: 2 }];
        els.answerFields.innerHTML = list
            .map(
                (q) => `
            <div class="answer-field-row">
                <div class="answer-field-number">${q.question_number}</div>
                <input class="answer-field" data-number="${q.question_number}" type="text" placeholder="Ответ" autocomplete="off" />
            </div>`
            )
            .join('');
    }

    function showScreen(n) {
        els.screen1.classList.toggle('active', n === 1);
        els.screen2.classList.toggle('active', n === 2);
    }
    els.arrowNext.addEventListener('click', () => showScreen(2));
    els.arrowPrev.addEventListener('click', () => showScreen(1));

    els.checkBtn.addEventListener('click', async () => {
        const inputs = els.answerFields.querySelectorAll('.answer-field');
        inputs.forEach((i) => i.classList.remove('field-error'));

        const answers = Array.from(inputs).map((i) => ({
            question_number: parseInt(i.dataset.number, 10),
            answer: i.value,
        }));

        try {
            const result = await Api.checkCrossword(taskId, answers);
            if (result.awarded) {
                solved = true;
                els.backBtn.classList.remove('disabled');
                els.backBtn.disabled = false;
                showResultModal(true, 'Верно! +3% команде');
            } else {
                (result.wrong_numbers || []).forEach((num) => {
                    const field = els.answerFields.querySelector(`.answer-field[data-number="${num}"]`);
                    if (field) field.classList.add('field-error');
                });
                showResultModal(false, 'Есть ошибки — попробуйте ещё раз');
            }
        } catch (err) {
            if (err.status === 409) {
                solved = true;
                els.backBtn.classList.remove('disabled');
                els.backBtn.disabled = false;
                showResultModal(true, 'Задание уже выполнено вашей командой');
            } else {
                showResultModal(false, err.message || 'Не удалось проверить ответы');
            }
        }
    });

    function showResultModal(success, text) {
        els.resultText.textContent = text;
        els.resultText.style.color = success ? 'var(--color-success)' : 'var(--color-danger)';
        els.resultModal.classList.remove('hidden');
    }
    els.resultClose.addEventListener('click', () => els.resultModal.classList.add('hidden'));

    els.backBtn.addEventListener('click', () => {
        if (!solved) return;
        window.location.href = '../main.html';
    });
});

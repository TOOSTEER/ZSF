// crossword-circles.js — тип B: circles.html (РК4, ССФРК)

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

    // 7 кругов всего: 4 на экране 1, 3 на экране 2 (нумерация вопросов 1..7)
    const circleNumbers = [1, 2, 3, 4, 5, 6, 7];
    let solved = false;

    try {
        const structure = await Api.crosswordStructure(taskId);
        els.title.textContent = structure.title;
        fillCircleLabels(structure.questions || []);
    } catch (err) {
        console.error(err);
    }

    function fillCircleLabels(questions) {
        circleNumbers.forEach((num) => {
            const circleEl = document.getElementById(circleIdFor(num));
            const q = questions.find((x) => x.question_number === num);
            if (circleEl) {
                circleEl.textContent = q ? `№${num}` : `№${num} (уточняется)`;
            }
        });
    }

    function circleIdFor(num) {
        return num <= 4 ? `circle-1-${num}` : `circle-2-${num - 4}`;
    }
    function inputIdFor(num) {
        return num <= 4 ? `input-1-${num}` : `input-2-${num - 4}`;
    }

    function showScreen(n) {
        els.screen1.classList.toggle('active', n === 1);
        els.screen2.classList.toggle('active', n === 2);
    }
    els.arrowNext.addEventListener('click', () => showScreen(2));
    els.arrowPrev.addEventListener('click', () => showScreen(1));

    els.checkBtn.addEventListener('click', async () => {
        circleNumbers.forEach((num) => {
            const input = document.getElementById(inputIdFor(num));
            if (input) input.classList.remove('field-error');
        });

        const answers = circleNumbers.map((num) => ({
            question_number: num,
            answer: (document.getElementById(inputIdFor(num)) || {}).value || '',
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
                    const input = document.getElementById(inputIdFor(num));
                    if (input) input.classList.add('field-error');
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

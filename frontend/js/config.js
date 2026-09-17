// config.js — глобальный конфиг фронтенда квеста.

const CONFIG = {
    API_URL: 'http://localhost:8080/api', // TODO: уточнить прод-адрес бэкенда
    TEAM_COLORS: {
        1: '#FF0000', // Красный
        2: '#008000', // Зелёный
        3: '#0000FF', // Синий
        4: '#FFFF00', // Жёлтый
        5: '#FFC0CB', // Розовый
    },
    MAX_PERCENT: 100,
    TOTAL_PERCENT: 123,
    CROSSWORD_TYPES: {
        'РК3': 'classic',
        'РК4': 'circles',
        'РК5': 'classic',
        'РК6': 'classic',
        'РК9': 'wordsearch',
        'РК1/РК2': 'text',
        'ССФРК': 'circles',
    },
    DESIGN_WIDTH: 1440, // ширина макета в условных px, см. css/common.css
};

/**
 * Масштабирует .stage (макет фиксированной ширины CONFIG.DESIGN_WIDTH)
 * под реальную ширину экрана 320–480px (или колонку 480px на широких
 * экранах — см. .stage-wrapper в common.css).
 * Вызывается один раз при загрузке страницы и повторно при resize.
 */
function initStageScaling() {
    const wrapper = document.querySelector('.stage-wrapper');
    const stage = document.querySelector('.stage');
    if (!wrapper || !stage) return;

    function applyScale() {
        const wrapperWidth = wrapper.clientWidth || window.innerWidth;
        const scale = wrapperWidth / CONFIG.DESIGN_WIDTH;
        stage.style.transform = `scale(${scale})`;
        // transform не влияет на поток документа, поэтому высоту обёртки
        // выставляем вручную по фактической (немасштабированной) высоте стейджа.
        const naturalHeight = stage.scrollHeight || stage.offsetHeight;
        wrapper.style.height = `${naturalHeight * scale}px`;
    }

    applyScale();
    window.addEventListener('resize', applyScale);
    // На случай, если контент стейджа подгружается асинхронно (например,
    // вопросы кроссворда) и высота меняется уже после первого рендера.
    window.addEventListener('load', applyScale);
    setTimeout(applyScale, 300);
}

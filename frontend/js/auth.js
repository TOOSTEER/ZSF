// auth.js — логика страницы index.html
// Участник → /api/auth/register (body: first_name, last_name, group_name, role)
// Водящий  → /api/auth/login    (body: login, password)

document.addEventListener('DOMContentLoaded', () => {
    if (typeof initStageScaling === 'function') initStageScaling();

    const existingUser = Storage.getUser();
    if (Storage.getToken() && existingUser) {
        window.location.href = existingUser.role === 'host' ? 'host.html' : 'main.html';
        return;
    }

    const authPage = document.querySelector('.auth-page');
    const memberFields = document.querySelectorAll('.member-field');
    const hostFields = document.querySelectorAll('.host-field');

    const els = {
        roleTabs: document.querySelectorAll('.js-role-tab'),
        form: document.getElementById('auth-form'),
        errorBox: document.getElementById('auth-error'),
        modalOverlay: document.getElementById('success-modal'),
        modalLogin: document.getElementById('modal-login'),
        modalTeam: document.getElementById('modal-team'),
        modalTeamColor: document.getElementById('modal-team-color'),
        modalContinueBtn: document.getElementById('modal-continue'),
        memberFirst: document.getElementById('member-first-name'),
        memberLast: document.getElementById('member-last-name'),
        memberGroup: document.getElementById('member-group'),
        hostLogin: document.getElementById('host-login'),
        hostPassword: document.getElementById('host-password'),
    };

    let selectedRole = 'member';

    function switchRole(role) {
        selectedRole = role;
        const isHost = role === 'host';

        els.roleTabs.forEach((t) => t.classList.toggle('active', t.dataset.role === role));

        // Скрываем/показываем поля через display:none
        memberFields.forEach((el) => el.classList.toggle('hidden-field', isHost));
        hostFields.forEach((el) => el.classList.toggle('hidden-field', !isHost));

        if (els.memberFirst) els.memberFirst.required = !isHost;
        if (els.memberLast)  els.memberLast.required  = !isHost;
        if (els.memberGroup) els.memberGroup.required = !isHost;
        if (els.hostLogin)   els.hostLogin.required   = isHost;
        if (els.hostPassword) els.hostPassword.required = isHost;

        if (authPage) authPage.classList.toggle('host-mode', isHost);
        showError('');
    }

    els.roleTabs.forEach((tab) => tab.addEventListener('click', () => switchRole(tab.dataset.role)));
    switchRole('member');

    function showError(message) {
        if (!els.errorBox) return;
        els.errorBox.textContent = message || '';
        els.errorBox.classList.toggle('hidden-field', !message);
    }

    els.form.addEventListener('submit', async (e) => {
        e.preventDefault();
        showError('');

        try {
            if (selectedRole === 'member') {
                const fd = new FormData(els.form);
                const payload = {
                    first_name: (fd.get('member_first_name') || '').trim(),
                    last_name:  (fd.get('member_last_name')  || '').trim(),
                    group_name: (fd.get('member_group')      || '').trim(),
                    role: 'member',
                };
                if (!payload.first_name || !payload.last_name || !payload.group_name) {
                    showError('Заполните все поля');
                    return;
                }

                const data = await Api.register(payload);
                Storage.setToken(data.token);
                Storage.setUser(data.user);
                Storage.setTeam(data.team);

                els.modalLogin.textContent = data.user.login;
                els.modalTeam.textContent = data.team ? data.team.name : '—';
                if (data.team) els.modalTeamColor.style.background = data.team.color;
                els.modalOverlay.classList.remove('hidden');

            } else {
                const fd = new FormData(els.form);
                const payload = {
                    login: (fd.get('host_login') || '').trim(),
                    password: fd.get('host_password') || '',
                };
                if (!payload.login || !payload.password) {
                    showError('Введите логин и пароль');
                    return;
                }

                const data = await Api.login(payload);
                Storage.setToken(data.token);
                Storage.setUser(data.user);
                Storage.setTeam(data.team);
                window.location.href = 'host.html';
            }
        } catch (err) {
            console.error('[auth] status=%d message=%s data=%o', err.status, err.message, err.data);
            if (err.status === 401)      showError('Неверный логин или пароль');
            else if (err.status === 409) showError('Все команды уже заняты водящими');
            else                          showError(err.message || 'Не удалось выполнить действие');
        }
    });

    els.modalContinueBtn.addEventListener('click', () => {
        window.location.href = 'main.html';
    });
});
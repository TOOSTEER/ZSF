// auth.js — логика страницы регистрации/входа (index.html)

document.addEventListener('DOMContentLoaded', () => {
    initStageScaling();

    // Если токен уже есть — сразу уводим на нужную страницу.
    const existingUser = Storage.getUser();
    if (Storage.getToken() && existingUser) {
        window.location.href = existingUser.role === 'host' ? 'host.html' : 'main.html';
        return;
    }

    const els = {
        roleTabs: document.querySelectorAll('.js-role-tab'),
        formModeLinks: document.querySelectorAll('.js-form-mode'),
        registerForm: document.getElementById('register-form'),
        loginForm: document.getElementById('login-form'),
        hostPasswordField: document.getElementById('field-password'),
        errorBox: document.getElementById('auth-error'),
        modalOverlay: document.getElementById('success-modal'),
        modalLogin: document.getElementById('modal-login'),
        modalTeam: document.getElementById('modal-team'),
        modalTeamColor: document.getElementById('modal-team-color'),
        modalContinueBtn: document.getElementById('modal-continue'),
    };

    let selectedRole = 'member'; // 'member' | 'host' — активная вкладка регистрации
    let successRedirectTo = 'main.html';

    // ---- Переключение роли (Участник / Водящий) в форме регистрации ----
    els.roleTabs.forEach((tab) => {
        tab.addEventListener('click', () => {
            selectedRole = tab.dataset.role;
            els.roleTabs.forEach((t) => t.classList.toggle('active', t === tab));
            els.hostPasswordField.classList.toggle('visually-hidden', selectedRole !== 'host');
            els.hostPasswordField.querySelector('input').required = selectedRole === 'host';
        });
    });

    // ---- Переключение между "Регистрация" и "Вход" (для водящих) ----
    els.formModeLinks.forEach((link) => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const mode = link.dataset.mode;
            els.registerForm.classList.toggle('visually-hidden', mode !== 'register');
            els.loginForm.classList.toggle('visually-hidden', mode !== 'login');
            showError('');
        });
    });

    function showError(message) {
        els.errorBox.textContent = message || '';
        els.errorBox.classList.toggle('visually-hidden', !message);
    }

    // ---- Регистрация ----
    els.registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        showError('');

        const formData = new FormData(els.registerForm);
        const payload = {
            first_name: formData.get('first_name')?.trim(),
            last_name: formData.get('last_name')?.trim(),
            group_name: formData.get('group_name')?.trim(),
            role: selectedRole,
        };
        if (selectedRole === 'host') {
            payload.password = formData.get('password');
        }

        try {
            const data = await Api.register(payload);
            Storage.setToken(data.token);
            Storage.setUser(data.user);
            Storage.setTeam(data.team);

            if (selectedRole === 'member') {
                openSuccessModal(data);
                successRedirectTo = 'main.html';
            } else {
                window.location.href = 'host.html';
            }
        } catch (err) {
            showError(err.message || 'Не удалось зарегистрироваться');
        }
    });

    // ---- Вход (водящий) ----
    els.loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        showError('');

        const formData = new FormData(els.loginForm);
        const payload = {
            login: formData.get('login')?.trim(),
            password: formData.get('password'),
        };

        try {
            const data = await Api.login(payload);
            Storage.setToken(data.token);
            Storage.setUser(data.user);
            Storage.setTeam(data.team);
            window.location.href = 'host.html';
        } catch (err) {
            showError(err.message || 'Неверный логин или пароль');
        }
    });

    function openSuccessModal(data) {
        els.modalLogin.textContent = data.user.login;
        els.modalTeam.textContent = data.team ? data.team.name : '—';
        if (data.team) {
            els.modalTeamColor.style.background = data.team.color;
        }
        els.modalOverlay.classList.remove('hidden');
    }

    els.modalContinueBtn.addEventListener('click', () => {
        window.location.href = successRedirectTo;
    });
});

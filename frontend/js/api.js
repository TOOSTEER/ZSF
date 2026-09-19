// api.js — обёртка над fetch для похода в бэкенд.

const TOKEN_KEY = 'quest_token';
const USER_KEY  = 'quest_user';
const TEAM_KEY  = 'quest_team';

const Storage = {
    getToken: () => localStorage.getItem(TOKEN_KEY),
    setToken: (token) => localStorage.setItem(TOKEN_KEY, token),
    clearToken: () => localStorage.removeItem(TOKEN_KEY),
    getUser: () => { try { return JSON.parse(localStorage.getItem(USER_KEY) || 'null'); } catch (e) { return null; } },
    setUser: (user) => localStorage.setItem(USER_KEY, JSON.stringify(user)),
    getTeam: () => { try { return JSON.parse(localStorage.getItem(TEAM_KEY) || 'null'); } catch (e) { return null; } },
    setTeam: (team) => localStorage.setItem(TEAM_KEY, JSON.stringify(team)),
    clearAll: () => {
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(USER_KEY);
        localStorage.removeItem(TEAM_KEY);
    },
};

async function apiRequest(path, options = {}) {
    const { method = 'GET', body, auth = true } = options;
    const headers = { 'Content-Type': 'application/json' };
    if (auth) {
        const token = Storage.getToken();
        if (token) headers['Authorization'] = `Bearer ${token}`;
    }

    let response;
    try {
        response = await fetch(`${CONFIG.API_URL}${path}`, {
            method, headers,
            body: body ? JSON.stringify(body) : undefined,
        });
    } catch (networkErr) {
        throw new ApiError(0, 'Нет соединения с сервером');
    }

    let data = null;
    const text = await response.text();
    if (text) {
        try { data = JSON.parse(text); } catch (e) { data = null; }
    }
    if (!response.ok) {
        const message = (data && data.error) || `Ошибка ${response.status}`;
        throw new ApiError(response.status, message, data);
    }
    return data;
}

class ApiError extends Error {
    constructor(status, message, data) {
        super(message);
        this.status = status;
        this.data = data;
    }
}

const Api = {
    register: (payload) => apiRequest('/auth/register', { method: 'POST', body: payload, auth: false }),
    login:    (payload) => apiRequest('/auth/login',    { method: 'POST', body: payload, auth: false }),
    me: () => apiRequest('/me'),
    teamProgress: (teamId) => apiRequest(`/teams/${teamId}/progress`),
    task: (taskId) => apiRequest(`/tasks/${taskId}`),
    crosswordStructure: (taskId) => apiRequest(`/tasks/${taskId}/crossword`),
    awardProgress: (payload) => apiRequest('/progress/award', { method: 'POST', body: payload }),
    checkCrossword: (taskId, answers) =>
        apiRequest(`/crossword/${taskId}/check`, { method: 'POST', body: { answers } }),
};

function requireAuth() {
    if (!Storage.getToken()) { window.location.href = 'index.html'; return false; }
    return true;
}
function requireRole(role) {
    const user = Storage.getUser();
    if (!user || user.role !== role) { window.location.href = 'main.html'; return false; }
    return true;
}
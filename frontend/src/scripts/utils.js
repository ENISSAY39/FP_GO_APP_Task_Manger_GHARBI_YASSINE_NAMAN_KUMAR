export const el = id => document.getElementById(id);

export function escapeHtml(s) {
    if (!s && s !== 0) return "";
    return String(s).replace(/[&<>\"]/g, c => ({
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;'
    }[c]));
}

export function fmt(d) {
    try {
        return new Date(d).toLocaleString();
    } catch (e) {
        return d;
    }
}

export function setStatusRaw(txt) {
    console.log(txt);
}

export function getApiBase() {
    return el("apiBase")?.value.trim().replace(/\/+$/, "") || "";
}

export function getToken() {
    return el("jwt")?.value.trim() || "";
}

export function buildHeaders(additional = {}) {
    const headers = {
        "Content-Type": "application/json",
        ...additional
    };
    const token = getToken();
    if (token) {
        headers["Authorization"] = "Bearer " + token;
    }
    return headers;
}

import { getApiBase, buildHeaders } from "./utils.js";

export async function apiFetch(path, opts = {}) {
    const base = getApiBase();
    const headers = buildHeaders(opts.headers || {});
    const response = await fetch(base + path, {
        ...opts,
        headers
    });
    let json = null;
    try {
        json = await response.json();
    } catch (e) {
        // ignore invalid JSON
    }
    return {
        status: response.status,
        ok: response.ok,
        json
    };
}

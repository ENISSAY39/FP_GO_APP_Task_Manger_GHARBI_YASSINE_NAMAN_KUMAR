// frontend/js/lib/api.js
const API_BASE = (() => {
  if (typeof window === 'undefined') return 'http://localhost:3000/api';
  // Use the same hostname as the page (127.0.0.1 or localhost) and backend port 3000.
  const host = window.location.hostname; // e.g. "127.0.0.1" or "localhost"
  const backendPort = "3000";
  return `http://${host}:${backendPort}/api`;
})();

function getStoredToken() {
  try { return (typeof window !== 'undefined') ? localStorage.getItem('token') : null; }
  catch (e) { return null; }
}
function setStoredToken(token) {
  try {
    if (typeof window !== 'undefined') {
      if (token == null) localStorage.removeItem('token');
      else localStorage.setItem('token', token);
    }
  } catch (e) {}
}

async function request(path, { method = "GET", body = null, headers = {} } = {}) {
  const url = `${API_BASE}${path}`;
  const hdrs = { ...headers };
  const token = getStoredToken();
  if (token && !hdrs["Authorization"] && !hdrs["authorization"]) hdrs["Authorization"] = `Bearer ${token}`;

  const opts = {
    method,
    credentials: "include",
    headers: hdrs,
  };

  if (body && !(body instanceof FormData)) {
    opts.headers["Content-Type"] = "application/json";
    opts.body = JSON.stringify(body);
  } else if (body instanceof FormData) {
    opts.body = body;
  }

  let res;
  try { res = await fetch(url, opts); }
  catch (networkErr) { throw new Error(networkErr?.message || "Network request failed"); }

  if (res.status === 204) return null;
  const text = await res.text();
  let data = null;
  try { data = text ? JSON.parse(text) : null; } catch (e) { data = text; }

  if (!res.ok) {
    const errMsg = data && (data.error || data.message) ? (data.error || data.message) : res.statusText;
    throw new Error(errMsg || `HTTP ${res.status}`);
  }

  return data;
}

// AUTH
export async function signup(payload) { return request("/signup", { method: "POST", body: payload }); }
export async function login(payload)  {
  const data = await request("/login", { method: "POST", body: payload });
  if (data && typeof data === 'object') {
    if (data.token) setStoredToken(data.token);
    else if (data.access_token) setStoredToken(data.access_token);
  }
  return data;
}
export function validate() { return request("/validate", { method: "GET" }); }
export async function logout() { try { await request("/logout", { method: "POST" }); } catch (e) {} setStoredToken(null); }

// PROJECTS
export function createProject(payload) { return request("/projects", { method: "POST", body: payload }); }
export function listProjects()         { return request("/projects", { method: "GET" }); }
export function getProject(id)         { return request(`/projects/${id}`, { method: "GET" }); }
export function updateProject(id, payload) { return request(`/projects/${id}`, { method: "PUT", body: payload }); }
export function deleteProject(id)      { return request(`/projects/${id}`, { method: "DELETE" }); }

// TASKS
export function createTask(projectId, payload) { return request(`/projects/${projectId}/tasks`, { method: "POST", body: payload }); }
export function listTasksForProject(projectId) { return request(`/projects/${projectId}/tasks`, { method: "GET" }); }
export function getTask(id)                    { return request(`/tasks/${id}`, { method: "GET" }); }
export function updateTask(id, payload)        { return request(`/tasks/${id}`, { method: "PUT", body: payload }); }
export function deleteTask(id)                 { return request(`/tasks/${id}`, { method: "DELETE" }); }
export function assignTask(id, payload)        { return request(`/tasks/${id}/assign`, { method: "PUT", body: payload }); }

// frontend/js/user/userDashboard.js
import { validate, logout } from '/js/lib/api.js';

const userInfo = document.getElementById('userInfo');
const tasksList = document.getElementById('tasksList');
const projectsList = document.getElementById('projectsList');
const btnLogout = document.getElementById('btnLogout');

async function onLoad() {
  try {
    const v = await validate();
    const user = v.user || {};
    userInfo.textContent = `${user.email || 'unknown'} — ${user.is_admin ? 'admin' : 'user'}`;
    renderTasks(v.tasks || []);
    renderProjects(v.projects || []);
  } catch (e) {
    // not authenticated -> redirect to login
    window.location.href = '/pages/auth/login.html';
  }
}

function renderTasks(tasks) {
  tasksList.innerHTML = '';
  if (!Array.isArray(tasks) || tasks.length === 0) {
    tasksList.innerHTML = '<div class="small">No tasks for now</div>';
    return;
  }
  tasks.forEach(t => {
    const d = document.createElement('div');
    d.className = 'task-card';
    d.innerHTML = `
      <div class="task-title">${escapeHtml(t.title || 'Untitled')}</div>
      <div class="task-meta">Project: ${escapeHtml(t.project?.name || '')} • Status: ${escapeHtml(t.status || '')}</div>
      <div style="margin-top:8px" class="small">${escapeHtml(t.description || '')}</div>
    `;
    tasksList.appendChild(d);
  });
}

function renderProjects(projects) {
  projectsList.innerHTML = '';
  if (!Array.isArray(projects) || projects.length === 0) {
    projectsList.innerHTML = '<div class="small">No projects</div>';
    return;
  }
  projects.forEach(p => {
    const el = document.createElement('div');
    el.className = 'project-card';
    el.innerHTML = `<strong>${escapeHtml(p.name || '(no name)')}</strong><div class="small">${escapeHtml(p.description || '')}</div>`;
    projectsList.appendChild(el);
  });
}

btnLogout.addEventListener('click', async () => {
  try {
    await logout();
  } catch (e) {
    // ignore network errors but proceed to redirect
  } finally {
    window.location.href = '/pages/auth/login.html';
  }
});

// safer escaping (works in all browsers)
function escapeHtml(s) {
  if (!s && s !== 0) return '';
  return String(s)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;');
}

onLoad();

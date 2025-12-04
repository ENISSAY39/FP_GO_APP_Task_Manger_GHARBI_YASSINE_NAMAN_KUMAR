// frontend/pages/manager/managerDashboard.js
// Robust manager dashboard module that dynamically resolves the API module path
// so we don't need to move files around on disk.

(function () {
  // Helper: determine root prefix used by your pages ('' or '/frontend')
  const ROOT = (window.__APP_ROOT_PREFIX !== undefined)
    ? window.__APP_ROOT_PREFIX
    : (window.location.pathname.includes('/frontend/') ? '/frontend' : '');

  // Candidate locations to try for api.js (ordered)
  const API_CANDIDATES = [
    `${ROOT}/js/lib/api.js`,            // common if api.js lives in frontend/js/lib
    `${ROOT}/pages/lib/api.js`,        // common if api.js was expected next to pages
    `${ROOT}/frontend/js/lib/api.js`,  // extra fallback
    `${ROOT}/lib/api.js`
  ];

  // Try to import a module by URL and return it, or null if failed.
  async function importModuleIfExists(url) {
    try {
      // fetch first to ensure 200 (import will fail with CORS/404 sometimes)
      const r = await fetch(url, { method: 'GET', credentials: 'include' });
      if (!r.ok) return null;
      // dynamic import expects a full absolute or same-origin relative path string.
      return await import(url);
    } catch (err) {
      return null;
    }
  }

  // Main async init
  (async function init() {
    // Try to find and import api module
    let api = null;
    for (const candidate of API_CANDIDATES) {
      api = await importModuleIfExists(candidate);
      if (api) {
        // module default exports are under the returned namespace (api.createProject etc)
        break;
      }
    }

    if (!api) {
      console.error('Could not load api module. Tried:', API_CANDIDATES);
      showModuleLoadError();
      return;
    }

    // Pick exported functions (support both default and named exports)
    const createProject = api.createProject ?? (api.default && api.default.createProject);
    const listProjects = api.listProjects ?? (api.default && api.default.listProjects);
    const createTask = api.createTask ?? (api.default && api.default.createTask);
    const logoutApi = api.logout ?? (api.default && api.default.logout);

    // Fallback checks
    if (!createProject || !listProjects || !createTask) {
      console.error('api module loaded but missing expected functions', api);
      showModuleLoadError();
      return;
    }

    // small helpers and UI functions
    const escapeHtml = (s) => (s == null) ? '' : String(s).replaceAll('&','&amp;').replaceAll('<','&lt;').replaceAll('>','&gt;');
    const parseIds = (csv) => (!csv) ? [] : csv.split(',').map(s => Number(s.trim())).filter(n => Number.isFinite(n));
    function showMsg(el, text, isError=false){ if(!el) return; el.textContent=text||''; el.style.color=isError? '#9b1c1c':''; if(text) setTimeout(()=>{ if(el.textContent===text) el.textContent=''; },3000); }

    // load tasks for a given project into a container
    async function loadTasksForProjectInto(projectId, tasksContainer) {
      if (!tasksContainer) return;
      if (!projectId) { tasksContainer.innerHTML='<div class="small">Select a project</div>'; return; }
      try {
        const res = await fetch(`${ROOT}/api/projects/${projectId}/tasks`, { method:'GET', credentials:'include', headers:{'Accept':'application/json'} });
        if (!res.ok) { tasksContainer.innerHTML='<div class="small">Failed to load tasks</div>'; return; }
        const data = await res.json().catch(()=>({}));
        const tasks = Array.isArray(data) ? data : (Array.isArray(data.tasks) ? data.tasks : []);
        if (!tasks || tasks.length===0) { tasksContainer.innerHTML='<div class="small">No tasks for this project</div>'; return; }
        tasksContainer.innerHTML=''; tasks.forEach(t => {
          const d=document.createElement('div'); d.className='task-item';
          d.innerHTML = `<strong>${escapeHtml(t.title)}</strong><div class="small">${escapeHtml(String(t.priority||''))} • ${escapeHtml(t.status||'')}</div>`;
          tasksContainer.appendChild(d);
        });
      } catch (err) { console.error('loadTasksForProject error', err); tasksContainer.innerHTML='<div class="small">Error loading tasks</div>'; }
    }

    // load projects into select elements (NORMALIZE API response here)
    async function loadProjectsToSelects({ selectMemberProject, selectTaskProject, tasksContainer, mgrMsg } = {}) {
      try {
        const raw = await listProjects().catch(()=>null);
        const projects = Array.isArray(raw) ? raw : (raw && Array.isArray(raw.projects) ? raw.projects : []);
        if (selectMemberProject) selectMemberProject.innerHTML = '<option value="">-- select --</option>';
        if (selectTaskProject) selectTaskProject.innerHTML = '<option value="">-- select --</option>';
        if (tasksContainer) tasksContainer.innerHTML = '';
        (projects||[]).forEach(p => {
          const opt1 = document.createElement('option'); opt1.value=p.id; opt1.textContent=p.name || `#${p.id}`; if(selectMemberProject) selectMemberProject.appendChild(opt1);
          const opt2 = document.createElement('option'); opt2.value=p.id; opt2.textContent=p.name || `#${p.id}`; if(selectTaskProject) selectTaskProject.appendChild(opt2);
        });
        if (projects && projects.length>0) await loadTasksForProjectInto(projects[0].id, tasksContainer);
        else if (tasksContainer) tasksContainer.innerHTML = '<div class="small">No projects yet</div>';
      } catch (err) { console.error('loadProjectsToSelects error', err); showMsg(mgrMsg, 'Failed to load projects', true); }
    }

    // Wait for DOM ready and wire UI
    document.addEventListener('DOMContentLoaded', async () => {
      const inputProjectName = document.getElementById('projectName');
      const inputProjectDesc = document.getElementById('projectDesc');
      const btnCreate = document.getElementById('btnCreate');
      const mgrMsg = document.getElementById('mgrMsg');

      const selectMemberProject = document.getElementById('memberProjectSelect');
      const inputMemberIds = document.getElementById('memberIds');
      const btnAddMembers = document.getElementById('btnAddMembers');
      const membersMsg = document.getElementById('membersMsg');

      const selectTaskProject = document.getElementById('taskProjectSelect');
      const inputTaskTitle = document.getElementById('taskTitle');
      const inputTaskDesc = document.getElementById('taskDesc');
      const inputTaskDue = document.getElementById('taskDue');
      const selectTaskPriority = document.getElementById('taskPriority');
      const inputTaskAssignees = document.getElementById('taskAssignees');
      const btnCreateTask = document.getElementById('btnCreateTask');
      const taskMsg = document.getElementById('taskMsg');

      const projectTasksDiv = document.getElementById('projectTasks');

      const btnLogout = document.getElementById('btnLogout');
      const userInfoDiv = document.getElementById('userInfo');

      // validate user (redirect if not authenticated)
      try {
        const vRes = await fetch(`${ROOT}/api/validate`, { method: 'GET', credentials: 'include' });
        if (!vRes.ok) {
          window.location.href = `${ROOT}/pages/auth/login.html`;
          return;
        }
        const j = await vRes.json().catch(()=>({}));
        if (j.user && userInfoDiv) userInfoDiv.textContent = `${j.user.email} — ${j.user.is_admin ? 'manager' : 'user'}`;
      } catch (err) {
        console.warn('validate failed', err);
        window.location.href = `${ROOT}/pages/auth/login.html`;
        return;
      }

      // create project
      if (btnCreate) {
        btnCreate.addEventListener('click', async (ev) => {
          ev.preventDefault();
          const name = (inputProjectName && inputProjectName.value || '').trim();
          const description = (inputProjectDesc && inputProjectDesc.value || '').trim();
          if (!name) { showMsg(mgrMsg, 'Project name required', true); return; }
          try {
            btnCreate.disabled = true;
            await createProject({ name, description });
            showMsg(mgrMsg, 'Project created');
            if (inputProjectName) inputProjectName.value = '';
            if (inputProjectDesc) inputProjectDesc.value = '';
            await loadProjectsToSelects({ selectMemberProject, selectTaskProject, tasksContainer: projectTasksDiv, mgrMsg });
          } catch (err) {
            console.error('createProject error', err);
            showMsg(mgrMsg, err.message || 'Failed to create project', true);
          } finally {
            btnCreate.disabled = false;
          }
        });
      }

      // add members
      if (btnAddMembers) {
        btnAddMembers.addEventListener('click', async (ev) => {
          ev.preventDefault();
          const pid = selectMemberProject ? selectMemberProject.value : '';
          const ids = inputMemberIds ? parseIds(inputMemberIds.value) : [];
          if (!pid) { showMsg(membersMsg, 'Select a project', true); return; }
          if (ids.length === 0) { showMsg(membersMsg, 'Enter member IDs', true); return; }
          try {
            btnAddMembers.disabled = true;
            const res = await fetch(`${ROOT}/api/projects/${pid}/members`, {
              method: 'POST',
              credentials: 'include',
              headers: { 'Content-Type': 'application/json', 'Accept': 'application/json' },
              body: JSON.stringify({ member_ids: ids })
            });
            if (!res.ok) {
              const body = await res.json().catch(()=>({}));
              throw new Error(body.error || body.message || `HTTP ${res.status}`);
            }
            showMsg(membersMsg, 'Members added');
            if (inputMemberIds) inputMemberIds.value = '';
          } catch (err) {
            console.error('addMembers error', err);
            showMsg(membersMsg, err.message || 'Failed to add members', true);
          } finally {
            btnAddMembers.disabled = false;
          }
        });
      }

      // project select change -> load tasks
      if (selectTaskProject) {
        selectTaskProject.addEventListener('change', async () => {
          const pid = selectTaskProject.value;
          await loadTasksForProjectInto(pid, projectTasksDiv);
        });
      }

      // create task
      if (btnCreateTask) {
        btnCreateTask.addEventListener('click', async (ev) => {
          ev.preventDefault();
          const pid = selectTaskProject ? selectTaskProject.value : '';
          const title = inputTaskTitle ? inputTaskTitle.value.trim() : '';
          const description = inputTaskDesc ? inputTaskDesc.value.trim() : '';
          const due_date = inputTaskDue ? inputTaskDue.value || null : null;
          const priority = selectTaskPriority ? Number(selectTaskPriority.value || 1) : 1;
          const assignees = inputTaskAssignees ? parseIds(inputTaskAssignees.value) : [];

          if (!pid) { showMsg(taskMsg, 'Select a project for the task', true); return; }
          if (!title) { showMsg(taskMsg, 'Task title required', true); return; }

          try {
            btnCreateTask.disabled = true;
            await createTask(pid, { title, description, due_date, priority, assignee_ids: assignees });
            showMsg(taskMsg, 'Task created');
            if (inputTaskTitle) inputTaskTitle.value = '';
            if (inputTaskDesc) inputTaskDesc.value = '';
            if (inputTaskDue) inputTaskDue.value = '';
            if (inputTaskAssignees) inputTaskAssignees.value = '';
            await loadTasksForProjectInto(pid, projectTasksDiv);
          } catch (err) {
            console.error('createTask error', err);
            showMsg(taskMsg, err.message || 'Failed to create task', true);
          } finally {
            btnCreateTask.disabled = false;
          }
        });
      }

      // logout
      if (btnLogout) {
        btnLogout.addEventListener('click', async (ev) => {
          ev && ev.preventDefault();
          try {
            if (logoutApi) await logoutApi();
            else await fetch(`${ROOT}/api/logout`, { method:'POST', credentials:'include' });
          } catch (err) {
            try { await fetch(`${ROOT}/api/logout`, { method:'POST', credentials:'include' }); } catch(e) {}
          } finally {
            window.location.href = `${ROOT}/pages/auth/login.html`;
          }
        });
      }

      // finally load projects into selects
      await loadProjectsToSelects({ selectMemberProject, selectTaskProject, tasksContainer: projectTasksDiv, mgrMsg });
    });
  })();
  // end init

  // small utility: render visible banner if module can't be loaded
  function showModuleLoadError() {
    try {
      const b = document.createElement('div');
      b.style.padding = '12px';
      b.style.background = '#fee';
      b.style.border = '1px solid #f99';
      b.style.color = '#900';
      b.style.margin = '12px';
      b.textContent = 'Failed to load manager dashboard script. Check console for details.';
      document.body.insertBefore(b, document.body.firstChild);
    } catch (e) { /* ignore */ }
  }

})();

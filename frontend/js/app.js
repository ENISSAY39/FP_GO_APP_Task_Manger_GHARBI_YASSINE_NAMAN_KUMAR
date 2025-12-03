// app.js — main dashboard logic
(async function(){
  const who = document.getElementById('who');
  const logoutBtn = document.getElementById('logoutBtn');
  const projectsList = document.getElementById('projects-list');
  const createProjectBtn = document.getElementById('create-project');
  const newProjectName = document.getElementById('new-project-name');
  const newProjectDesc = document.getElementById('new-project-desc');

  const selectedProjectName = document.getElementById('selected-project-name');
  const tasksList = document.getElementById('tasks-list');
  const createTaskBtn = document.getElementById('create-task');
  const taskTitle = document.getElementById('task-title');
  const taskDesc = document.getElementById('task-desc');
  const taskPriority = document.getElementById('task-priority');
  const membersArea = document.getElementById('members-area');

  let currentUser = null;
  let projects = [];
  let selectedProject = null;

  // logout: clear cookie by calling backend? we'll just clear and redirect to index
  logoutBtn.onclick = async () => {
    // clear cookie by setting expired cookie through backend or client-side trick not reliable.
    // best approach: call a logout endpoint. For now redirect to index and tell backend cookie will expire eventually.
    window.location.href = '/frontend/index.html';
  };

  async function validate() {
    const res = await api.get('/validate');
    if (!res.ok) {
      alert('Not authenticated — redirect to login');
      window.location.href = '/frontend/index.html';
      throw new Error('Not authenticated');
    }
    const json = await res.json();
    currentUser = json.user;
    who.textContent = currentUser.email + (currentUser.is_admin ? ' (GLOBAL ADMIN)' : '');
    return currentUser;
  }

  async function loadProjects() {
    const res = await api.get('/projects');
    if (!res.ok) { projectsList.innerHTML = '<div class="small">Failed to load projects</div>'; return; }
    const json = await res.json();
    projects = json.projects || [];
    renderProjects();
  }

  function hasManagerRole(project, uid) {
    if (!project || !project.members) return false;
    const m = project.members.find(x => x.user && x.user.id === uid);
    return m && m.role === 'manager';
  }

  function isMember(project, uid) {
    if (!project || !project.members) return false;
    return project.members.some(x => x.user && x.user.id === uid);
  }

  function renderProjects() {
    projectsList.innerHTML = '';
    projects.forEach(p => {
      const div = document.createElement('div');
      div.className = 'project';
      div.innerHTML = `
        <div class="meta">
          <strong>${p.name}</strong>
          <div class="small">${p.description || ''}</div>
        </div>
        <div class="actions">
          <span class="small">owner: ${p.owner ? p.owner.email : '—'}</span>
          <button class="btn small open">Open</button>
        </div>
      `;
      const openBtn = div.querySelector('.open');
      openBtn.onclick = () => selectProject(p.id);
      // show edit/delete buttons if allowed (admin or owner)
      const canEdit = currentUser.is_admin || (p.owner && p.owner.id === currentUser.id) || hasManagerRole(p, currentUser.id);
      if (canEdit) {
        const edit = document.createElement('button');
        edit.className = 'btn small';
        edit.textContent = 'Manage';
        edit.onclick = () => selectProject(p.id, true);
        div.querySelector('.actions').appendChild(edit);
      }
      projectsList.appendChild(div);
    });
  }

  async function selectProject(id, openManage=false) {
    selectedProject = projects.find(p => p.id === id);
    if (!selectedProject) return;
    selectedProjectName.textContent = selectedProject.name;
    await loadTasksForProject(id);
    renderMembers();
  }

  async function loadTasksForProject(id) {
    const res = await api.get(`/projects/${id}/tasks`);
    if (!res.ok) {
      tasksList.innerHTML = '<div class="small">Failed to load tasks</div>';
      return;
    }
    const json = await res.json();
    const arr = json.tasks || [];
    tasksList.innerHTML = '';
    arr.forEach(t => {
      const el = document.createElement('div');
      el.className = 'task';
      // show status toggle if user is assignee or admin/manager
      let canToggle = currentUser.is_admin;
      if (!canToggle) {
        // check if currentUser is assignee
        if (t.assignees && t.assignees.some(a => a.id === currentUser.id)) canToggle = true;
        // or project manager / owner
        if (selectedProject.owner && selectedProject.owner.id === currentUser.id) canToggle = true;
        if (hasManagerRole(selectedProject, currentUser.id)) canToggle = true;
      }
      const assignees = (t.assignees || []).map(a => a.email || a.user?.email || 'user').join(', ');

      el.innerHTML = `
        <div>
          <strong>${t.title}</strong>
          <div class="small">${t.description || ''}</div>
          <div class="small">Assignees: ${assignees}</div>
        </div>
        <div>
          <div class="badge">${t.status}</div>
          ${canToggle ? `<select class="status-select">
            <option ${t.status==='todo'?'selected':''} value="todo">todo</option>
            <option ${t.status==='in_progress'?'selected':''} value="in_progress">in_progress</option>
            <option ${t.status==='done'?'selected':''} value="done">done</option>
          </select>` : ''}
        </div>
      `;
      if (canToggle) {
        el.querySelector('.status-select').onchange = async (ev) => {
          const newStatus = ev.target.value;
          await api.put(`/tasks/${t.id}`, { status: newStatus });
          await loadTasksForProject(selectedProject.id);
        };
      }
      tasksList.appendChild(el);
    });
  }

  function renderMembers() {
    membersArea.innerHTML = '';
    if (!selectedProject) return;
    // only show member manager actions if global admin OR owner OR manager in this project
    const canManage = currentUser.is_admin || (selectedProject.owner && selectedProject.owner.id === currentUser.id) || hasManagerRole(selectedProject, currentUser.id);
    const memberList = document.createElement('div');
    memberList.className = 'member-list';
    (selectedProject.members || []).forEach(m => {
      const div = document.createElement('div');
      div.className = 'member';
      div.innerHTML = `<div>${m.user ? m.user.email : 'user'} <span class="small">(${m.role})</span></div>`;
      if (canManage) {
        const select = document.createElement('select');
        select.innerHTML = `<option value="member" ${m.role==='member'?'selected':''}>member</option>
                            <option value="manager" ${m.role==='manager'?'selected':''}>manager</option>`;
        select.onchange = async () => {
          await api.post(`/projects/${selectedProject.id}/members`, { user_id: m.user.id, role: select.value });
          await loadProjects();
          await selectProject(selectedProject.id);
        };
        const del = document.createElement('button');
        del.className = 'btn';
        del.textContent = 'Remove';
        del.onclick = async () => {
          await api.del(`/projects/${selectedProject.id}/members/${m.user.id}`);
          await loadProjects();
          await selectProject(selectedProject.id);
        };
        div.appendChild(select);
        div.appendChild(del);
      }
      memberList.appendChild(div);
    });

    // show area to add member (only manage)
    if (canManage) {
      const addDiv = document.createElement('div');
      addDiv.innerHTML = `
        <div style="margin-top:8px;">
          <input id="add-user-id" placeholder="user id (phpmyadmin)"/>
          <select id="add-role">
            <option value="member">member</option>
            <option value="manager">manager</option>
          </select>
          <button id="add-member-btn" class="btn primary">Add</button>
        </div>
      `;
      addDiv.querySelector('#add-member-btn').onclick = async () => {
        const uid = parseInt(addDiv.querySelector('#add-user-id').value);
        const role = addDiv.querySelector('#add-role').value;
        if (!uid) return alert('user id required');
        await api.post(`/projects/${selectedProject.id}/members`, { user_id: uid, role });
        await loadProjects();
        await selectProject(selectedProject.id);
      };
      membersArea.appendChild(addDiv);
    }

    membersArea.appendChild(memberList);
  }

  // create project
  createProjectBtn.onclick = async () => {
    const name = newProjectName.value.trim();
    const desc = newProjectDesc.value.trim();
    if (!name) return alert('name required');
    const res = await api.post('/projects', { name, description: desc });
    if (!res.ok) return alert('failed to create');
    newProjectName.value = ''; newProjectDesc.value = '';
    await loadProjects();
  };

  // create task
  createTaskBtn.onclick = async () => {
    if (!selectedProject) return alert('select a project first');
    const title = taskTitle.value.trim();
    if (!title) return alert('title required');
    const desc = taskDesc.value.trim();
    const pr = parseInt(taskPriority.value || '0');
    const res = await api.post(`/projects/${selectedProject.id}/tasks`, { title, description: desc, priority: pr });
    if (!res.ok) return alert('failed to create task');
    taskTitle.value = ''; taskDesc.value = ''; taskPriority.value = '';
    await loadTasksForProject(selectedProject.id);
  };

  // initial flow
  try {
    await validate();
    await loadProjects();
  } catch(e){ console.error(e); }

})();

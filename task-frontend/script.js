/*
  All-in-one frontend to manage the APIs you already have:
  - login: POST /api/login
  - projects: GET /api/projects, POST /api/projects
  - project detail: GET /api/projects/:id
  - tasks: GET /api/projects/:id/tasks, POST /api/projects/:id/tasks
  - task update: PUT /api/tasks/:id
  - task delete: DELETE /api/tasks/:id
  - assign: POST /api/projects/:projectId/tasks/:taskId/assign  body { user_id }
  - unassign: PUT /api/tasks/:taskId/unassign  body { user_id }
  - members: POST /api/projects/:projectId/members  body { user_id, role }
             DELETE /api/projects/:projectId/members/:userId (by param)
             DELETE /api/projects/:projectId/members (by body) - if exists
  Usage:
    - login to get token (auto placed in 'Bearer token' field)
    - Load projects (click)
    - Click View to open details; manage tasks/members from detail pane
*/

const el = id => document.getElementById(id);

// small helpers
function safeJson(r){ try { return r.json; } catch(e){ return null; } }
function fmt(d){ try { return new Date(d).toLocaleString(); } catch(e){ return d; } }

// API helper (adds bearer token automatically)
async function apiFetch(path, opts = {}) {
  const base = el("apiBase").value.trim().replace(/\/+$/, "");
  const token = el("jwt").value.trim();
  const headers = Object.assign({ "Content-Type": "application/json" }, opts.headers || {});
  if (token) headers["Authorization"] = "Bearer " + token;
  const res = await fetch(base + path, Object.assign({ headers }, opts));
  let json = null;
  try { json = await res.json(); } catch(e){}
  return { status: res.status, ok: res.ok, json };
}

// UI render functions
async function loadProjects() {
  setStatusRaw("Loading projects...");
  const r = await apiFetch("/api/projects");
  el("raw").textContent = JSON.stringify(r.json, null, 2);
  if (!r.ok) {
    setStatusRaw("Error loading projects: " + (r.json?.error || r.status));
    el("projectsWrap").innerHTML = `<div class="muted">Failed to load projects</div>`;
    return;
  }
  const projects = r.json.projects || [];
  const wrap = el("projectsWrap");
  wrap.innerHTML = "";
  projects.forEach(p => {
    const div = document.createElement("div");
    div.className = "project-card";
    div.innerHTML = `
      <div style="display:flex; justify-content:space-between; align-items:flex-start;">
        <div>
          <b>${escapeHtml(p.name)}</b>
          <div class="small muted">ID: ${p.id} — owner_id: ${p.owner_id}</div>
        </div>
        <div class="controls">
          <button data-id="${p.id}" class="btnView">View</button>
          <button data-id="${p.id}" class="btnTasks ghost">Tasks</button>
        </div>
      </div>
      <div class="small" style="margin-top:8px">${escapeHtml(p.description || "")}</div>
    `;
    wrap.appendChild(div);
  });

  setStatusRaw("Projects loaded (" + projects.length + ")");
}

// Escape helper
function escapeHtml(s){ if(!s && s !== 0) return ""; return String(s).replace(/[&<>"]/g, c=> ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;'}[c])); }

// click handlers (delegated)
document.addEventListener("click", async (ev) => {
  if (ev.target.matches("#btnLoad")) {
    await loadProjects();
  } else if (ev.target.matches("#btnClearToken")) {
    el("jwt").value = "";
    setStatusRaw("Token cleared");
  } else if (ev.target.matches("#btnLogin")) {
    await handleLogin();
  } else if (ev.target.matches("#btnLogout")) {
    await handleLogout();
  } else if (ev.target.matches("#btnCreateProject")) {
    await handleCreateProject();
  } else if (ev.target.matches(".btnView")) {
    const id = ev.target.dataset.id;
    openDetail(id);
  } else if (ev.target.matches(".btnTasks")) {
    const id = ev.target.dataset.id;
    await loadTasksIntoRaw(id);
    alert("Tasks loaded in Raw box.");
  } else if (ev.target.matches("#btnCloseDetail")) {
    closeDetail();
  } else if (ev.target.matches("#btnRefreshDetail")) {
    const pid = el("detailCard").dataset.projectId;
    if (pid) openDetail(pid, true);
  } else if (ev.target.matches("#btnAddMember")) {
    await handleAddMember();
  } else if (ev.target.matches(".remove-member")) {
    const pid = el("detailCard").dataset.projectId;
    const uid = ev.target.dataset.userid;
    if (!confirm("Remove user "+uid+" from project?")) return;
    await removeMemberByParam(pid, uid);
    await openDetail(pid, true);
  } else if (ev.target.matches("#btnCreateTask")) {
    await handleCreateTask();
  } else if (ev.target.matches(".btnAssign")) {
    const pid = el("detailCard").dataset.projectId;
    const taskId = ev.target.dataset.taskid;
    const userId = prompt("user id to assign (number):", "");
    if (!userId) return;
    await assignUserToTask(pid, taskId, userId);
    await openDetail(pid, true);
  } else if (ev.target.matches(".btnUnassign")) {
    const pid = el("detailCard").dataset.projectId;
    const taskId = ev.target.dataset.taskid;
    const userId = prompt("user id to unassign (number):", "");
    if (!userId) return;
    await unassignUserFromTask(taskId, userId);
    await openDetail(pid, true);
  } else if (ev.target.matches(".btnEditTask")) {
    const pid = el("detailCard").dataset.projectId;
    const taskId = ev.target.dataset.taskid;
    await promptEditTask(pid, taskId);
  } else if (ev.target.matches(".btnDelTask")) {
    const pid = el("detailCard").dataset.projectId;
    const taskId = ev.target.dataset.taskid;
    if (!confirm("Delete task "+taskId+" ?")) return;
    await deleteTask(taskId);
    await openDetail(pid, true);
  } else if (ev.target.matches("#btnBulkAssign")) {
    const pid = el("detailCard").dataset.projectId;
    const myId = prompt("Your user id to assign (number):", "");
    if (!myId) return;
    await bulkAssign(pid, myId);
    await openDetail(pid, true);
  }
});

// login
async function handleLogin(){
  const email = el("loginEmail").value.trim();
  const password = el("loginPass").value;
  el("loginMsg").textContent = "Logging in...";
  const r = await apiFetch("/api/login", {
    method: "POST",
    body: JSON.stringify({ email, password })
  });
  if (!r.ok) {
    el("loginMsg").textContent = "Login failed: " + (r.json?.error || r.status);
    return;
  }
  const token = r.json?.token;
  if (token) {
    el("jwt").value = token;
    el("loginMsg").textContent = "Logged in, token set.";
  } else {
    el("loginMsg").textContent = "Login answered but no token returned.";
  }
}

// logout
async function handleLogout(){
  try {
    el("loginMsg").textContent = "Logging out...";
    const r = await apiFetch("/api/logout", { method: "POST" });
    if (!r.ok) {
      el("loginMsg").textContent = "Logout failed: " + (r.json?.error || r.status);
      return;
    }
    // clear token + UI
    el("jwt").value = "";
    el("loginMsg").textContent = "Logged out.";
    setStatusRaw("Logged out");
    // Clear projects display
    el("projectsWrap").innerHTML = "";
    closeDetail();
  } catch (err) {
    el("loginMsg").textContent = "Logout error: " + (err && err.message ? err.message : String(err));
  }
}

// create project
async function handleCreateProject(){
  const name = el("pName").value.trim();
  const desc = el("pDesc").value.trim();
  if (!name) { el("pStatus").textContent = "Missing name"; return; }
  el("pStatus").textContent = "Creating...";
  const r = await apiFetch("/api/projects", { method: "POST", body: JSON.stringify({ name, description: desc }) });
  el("pStatus").textContent = r.ok ? "Created" : ("Error: "+(r.json?.error || r.status));
  await loadProjects();
  el("pName").value=""; el("pDesc").value="";
}

// open detail
async function openDetail(projectId, forceRefresh=false) {
  el("detailCard").style.display = "block";
  el("detailCard").dataset.projectId = projectId;
  setStatusRaw("Loading project " + projectId + "...");
  const r = await apiFetch("/api/projects/" + projectId);
  el("rawDetail").textContent = JSON.stringify(r.json, null, 2);
  if (!r.ok) {
    el("detailTitle").textContent = "Project detail (error)";
    el("detailInfo").textContent = "Error: " + (r.json?.error || r.status);
    return;
  }
  const p = r.json.project;
  el("detailTitle").textContent = p.name + " (ID " + p.id + ")";
  el("detailInfo").textContent = `owner_id: ${p.owner_id} • created: ${fmt(p.created_at)}`;

  // members
  const membersWrap = el("membersList");
  membersWrap.innerHTML = "";
  (p.members || []).forEach(m => {
    const div = document.createElement("div");
    div.className = "member-row";
    div.innerHTML = `
      <div style="display:flex;justify-content:space-between;align-items:center;">
        <div>
          <b>User ${m.user_id}</b> <span class="small muted">role: ${m.role}</span>
          <div class="small muted">member id: ${m.id} • created: ${fmt(m.created_at)}</div>
        </div>
        <div style="display:flex; gap:8px">
          <button class="ghost remove-member" data-userid="${m.user_id}">Remove</button>
        </div>
      </div>
    `;
    membersWrap.appendChild(div);
  });

  // tasks
  const tasksWrap = el("tasksList");
  tasksWrap.innerHTML = "";
  (p.tasks || []).forEach(t => {
    const div = document.createElement("div");
    div.className = "task-row";
    div.innerHTML = `
      <div style="display:flex; justify-content:space-between; align-items:flex-start;">
        <div style="flex:1">
          <b>${escapeHtml(t.title)}</b> <span class="small muted">(#${t.id})</span>
          <div class="small muted">${escapeHtml(t.description || "")}</div>
          <div class="small muted">status: ${t.status} • priority: ${t.priority} • created: ${fmt(t.created_at)}</div>
          <div class="small muted">assignees: ${ (t.assignees && t.assignees.length) ? t.assignees.map(a=>a.user_id).join(", ") : "[]" }</div>
        </div>
        <div style="display:flex; flex-direction:column; gap:6px;">
          <button class="btnAssign" data-taskid="${t.id}" data-projectid="${p.id}">Assign</button>
          <button class="btnUnassign" data-taskid="${t.id}">Unassign</button>
          <button class="btnEditTask" data-taskid="${t.id}">Edit</button>
          <button class="btnDelTask warn" data-taskid="${t.id}">Delete</button>
        </div>
      </div>
    `;
    tasksWrap.appendChild(div);
  });

  setStatusRaw("Project loaded: " + p.name);
}

// helper: load tasks only into raw
async function loadTasksIntoRaw(projectId) {
  const r = await apiFetch(`/api/projects/${projectId}/tasks`);
  el("raw").textContent = JSON.stringify(r.json, null, 2);
  if (!r.ok) setStatusRaw("Failed load tasks: " + (r.json?.error || r.status));
  else setStatusRaw("Tasks loaded for project " + projectId);
}

// create task
async function handleCreateTask() {
  const pid = el("detailCard").dataset.projectId;
  if (!pid) { alert("Open a project detail first"); return; }
  const title = el("tTitle").value.trim();
  const desc = el("tDesc").value.trim();
  const priority = el("tPriority").value;
  if (!title) { alert("Missing title"); return; }
  const r = await apiFetch(`/api/projects/${pid}/tasks`, {
    method: "POST",
    body: JSON.stringify({ title, description: desc, priority })
  });
  if (!r.ok) {
    alert("Error create task: " + (r.json?.error || r.status));
    return;
  }
  el("tTitle").value = ""; el("tDesc").value = "";
  alert("Task created: id " + r.json.task.id);
  await openDetail(pid, true);
}

// assign user to task
async function assignUserToTask(projectId, taskId, userId) {
  const r = await apiFetch(`/api/projects/${projectId}/tasks/${taskId}/assign`, {
    method: "POST",
    body: JSON.stringify({ user_id: Number(userId) })
  });
  if (!r.ok) {
    alert("Assign failed: " + (r.json?.error || r.status));
    return;
  }
  alert("Assigned");
}

// unassign
async function unassignUserFromTask(taskId, userId) {
  const r = await apiFetch(`/api/tasks/${taskId}/unassign`, {
    method: "PUT",
    body: JSON.stringify({ user_id: Number(userId) })
  });
  if (!r.ok) {
    alert("Unassign failed: " + (r.json?.error || r.status));
    return;
  }
  alert("Unassigned");
}

// edit task with prompt-based small form
async function promptEditTask(projectId, taskId) {
  // fetch single task? simplest: ask new title/desc/status/priority
  const title = prompt("New title (leave blank = unchanged):");
  const description = prompt("New description (leave blank = unchanged):");
  const status = prompt("New status (TODO / IN_PROGRESS / DONE), leave blank = unchanged:");
  const priority = prompt("New priority (LOW/MEDIUM/HIGH), leave blank = unchanged:");
  const body = {};
  if (title !== null && title !== "") body.title = title;
  if (description !== null && description !== "") body.description = description;
  if (status !== null && status !== "") body.status = status;
  if (priority !== null && priority !== "") body.priority = priority;
  if (Object.keys(body).length === 0) return;
  const r = await apiFetch(`/api/tasks/${taskId}`, {
    method: "PUT",
    body: JSON.stringify(body)
  });
  if (!r.ok) { alert("Update failed: " + (r.json?.error || r.status)); return; }
  alert("Task updated");
  await openDetail(projectId, true);
}

// delete task
async function deleteTask(taskId) {
  const r = await apiFetch(`/api/tasks/${taskId}`, { method: "DELETE" });
  if (!r.ok) { alert("Delete failed: " + (r.json?.error || r.status)); return; }
  alert(r.json?.message || "Deleted");
}

// add member (body route)
async function handleAddMember() {
  const pid = el("detailCard").dataset.projectId;
  if (!pid) return alert("Open a project detail first");
  const uid = Number(el("newMemberId").value.trim());
  const role = el("newMemberRole").value;
  if (!uid) return alert("User id required (number)");
  const r = await apiFetch(`/api/projects/${pid}/members`, {
    method: "POST",
    body: JSON.stringify({ user_id: uid, role })
  });
  if (!r.ok) {
    alert("Add member failed: " + (r.json?.error || r.status)); return;
  }
  el("newMemberId").value = "";
  await openDetail(pid, true);
}

// remove member by param
async function removeMemberByParam(projectId, userId) {
  // uses DELETE /projects/:projectId/members/:userId
  const r = await apiFetch(`/api/projects/${projectId}/members/${userId}`, { method: "DELETE" });
  if (!r.ok) {
    alert("Remove member failed: " + (r.json?.error || r.status)); return;
  }
  alert("Member removed");
}

// bulk assign: assign provided user to all tasks in project (simple loop)
async function bulkAssign(projectId, userId) {
  const r = await apiFetch(`/api/projects/${projectId}/tasks`);
  if (!r.ok) return alert("Cannot load tasks: " + (r.json?.error || r.status));
  const tasks = r.json.tasks || [];
  for (const t of tasks) {
    await assignUserToTask(projectId, t.id, userId);
  }
  alert("Done assigning to " + tasks.length + " tasks (attempted).");
}

// close detail
function closeDetail() {
  el("detailCard").style.display = "none";
  el("detailCard").dataset.projectId = "";
  el("rawDetail").textContent = "";
}

// small helper set global raw
function setStatusRaw(txt) { el("raw").textContent = txt; }

// convenience: auto-load projects if token present on page open
(async function init(){
  // if token already filled, try loading
  if (el("jwt").value.trim()) {
    await loadProjects();
  }
})();
import { apiFetch } from "./api.js";
import { el, escapeHtml, fmt, setStatusRaw } from "./utils.js";

export async function loadProjects() {
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

export async function openDetail(projectId, forceRefresh = false) {
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
          <div class="small muted">assignees: ${(t.assignees && t.assignees.length) ? t.assignees.map(a => a.user_id).join(", ") : "[]"}</div>
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


export async function handleCreateProject() {
    const name = el("pName").value.trim();
    const desc = el("pDesc").value.trim();

    if (!name) {
        el("pStatus").textContent = "Missing name";
        return;
    }

    el("pStatus").textContent = "Creating...";

    const r = await apiFetch("/api/projects", {
        method: "POST",
        body: JSON.stringify({
            name,
            description: desc
        })
    });

    el("pStatus").textContent =
        r.ok
            ? "Created"
            : ("Error: " + (r.json?.error || r.status));

    await loadProjects();

    el("pName").value = "";
    el("pDesc").value = "";
}

export function closeDetail() {
    el("detailCard").style.display = "none";
    el("detailCard").dataset.projectId = "";
    el("rawDetail").textContent = "";
}
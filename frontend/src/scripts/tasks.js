import { apiFetch } from "./api.js";
import { el } from "./utils.js";
import { openDetail } from "./projects.js";

export async function loadTasksIntoRaw(projectId) {
    const r = await apiFetch(`/api/projects/${projectId}/tasks`);
    el("raw").textContent = JSON.stringify(r.json, null, 2);
    if (!r.ok) {
        el("raw").textContent = "";
        return;
    }
}

export async function handleCreateTask() {
    const pid = el("detailCard").dataset.projectId;
    if (!pid) {
        alert("Open a project detail first");
        return;
    }
    const title = el("tTitle").value.trim();
    const desc = el("tDesc").value.trim();
    const priority = el("tPriority").value;
    if (!title) {
        alert("Missing title");
        return;
    }
    const r = await apiFetch(`/api/projects/${pid}/tasks`, {
        method: "POST",
        body: JSON.stringify({ title, description: desc, priority })
    });
    if (!r.ok) {
        alert("Error create task: " + (r.json?.error || r.status));
        return;
    }
    el("tTitle").value = "";
    el("tDesc").value = "";
    alert("Task created: id " + r.json.task?.id);
    await openDetail(pid, true);
}

export async function assignUserToTask(projectId, taskId, userId) {
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

export async function unassignUserFromTask(taskId, userId) {
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

export async function promptEditTask(projectId, taskId) {
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
    if (!r.ok) {
        alert("Update failed: " + (r.json?.error || r.status));
        return;
    }
    alert("Task updated");
    await openDetail(projectId, true);
}

export async function deleteTask(taskId) {
    const r = await apiFetch(`/api/tasks/${taskId}`, { method: "DELETE" });
    if (!r.ok) {
        alert("Delete failed: " + (r.json?.error || r.status));
        return;
    }
    alert(r.json?.message || "Deleted");
}

export async function bulkAssign(projectId, userId) {
    const r = await apiFetch(`/api/projects/${projectId}/tasks`);
    if (!r.ok) {
        alert("Cannot load tasks: " + (r.json?.error || r.status));
        return;
    }
    const tasks = r.json.tasks || [];
    for (const t of tasks) {
        await assignUserToTask(projectId, t.id, userId);
    }
    alert("Done assigning to " + tasks.length + " tasks (attempted).");
}

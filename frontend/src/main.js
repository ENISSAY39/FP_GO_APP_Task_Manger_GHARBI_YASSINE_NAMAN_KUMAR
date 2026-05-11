import { handleLogin, handleSignup, handleLogout } from "./scripts/auth.js";
import { loadProjects, openDetail, } from "./scripts/projects.js";
import { handleCreateTask, loadTasksIntoRaw, assignUserToTask, unassignUserFromTask, promptEditTask, deleteTask, bulkAssign } from "./scripts/tasks.js";
import { handleAddMember, removeMemberByParam } from "./scripts/members.js";
import { apiFetch } from "./scripts/api.js";
import { el } from "./scripts/utils.js";

window.addEventListener("DOMContentLoaded", () => {
    document.addEventListener("click", async (ev) => {
        if (ev.target.matches("#btnLoad")) {
            await loadProjects();
        } else if (ev.target.matches("#btnClearToken")) {
            el("jwt").value = "";
        } else if (ev.target.matches("#btnLogin")) {
            await handleLogin();
        } else if (ev.target.matches("#btnSignup")) {
            await handleSignup();
        } else if (ev.target.matches("#btnLogout")) {
            await handleLogout();
            el("projectsWrap").innerHTML = "";
            closeDetail();
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
            if (!confirm("Remove user " + uid + " from project?")) return;
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
            if (!confirm("Delete task " + taskId + " ?")) return;
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


    if (el("jwt").value.trim()) {
        loadProjects();
    }
});

import { apiFetch } from "./api.js";
import { el } from "./utils.js";
import { openDetail } from "./projects.js";

export async function handleAddMember() {
    const pid = el("detailCard").dataset.projectId;
    if (!pid) {
        alert("Open a project detail first");
        return;
    }
    const uid = Number(el("newMemberId").value.trim());
    const role = el("newMemberRole").value;
    if (!uid) {
        alert("User id required (number)");
        return;
    }
    const r = await apiFetch(`/api/projects/${pid}/members`, {
        method: "POST",
        body: JSON.stringify({ user_id: uid, role })
    });
    if (!r.ok) {
        alert("Add member failed: " + (r.json?.error || r.status));
        return;
    }
    el("newMemberId").value = "";
    await openDetail(pid, true);
}

export async function removeMemberByParam(projectId, userId) {
    const r = await apiFetch(`/api/projects/${projectId}/members/${userId}`, {
        method: "DELETE"
    });
    if (!r.ok) {
        alert("Remove member failed: " + (r.json?.error || r.status));
        return;
    }
    alert("Member removed");
}

const API = "http://localhost:3000";

// ---------- AUTH ----------
async function signup() {
    let email = document.getElementById('email').value;
    let password = document.getElementById('password').value;

    let res = await fetch(`${API}/signup`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({email, password})
    });

    if (res.ok) {
        alert("Signup successful!");
        window.location = "index.html";
    } else {
        alert("Signup failed");
    }
}

async function login() {
    let email = document.getElementById('email').value;
    let password = document.getElementById('password').value;

    let res = await fetch(`${API}/login`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({email, password}),
        credentials: "include"
    });

    if (res.ok) {
        window.location = "projects.html";
    } else {
        alert("Login failed");
    }
}

// ---------- PROJECTS ----------
async function loadProjects() {
    let res = await fetch(`${API}/projects`, { credentials: "include" });
    let data = await res.json();

    let html = "";
    data.forEach(p => {
        html += `<p><a href="tasks.html?pid=${p.ID}">${p.Title}</a></p>`;
    });

    document.getElementById("projectList").innerHTML = html;
}

async function createProject() {
    let title = document.getElementById("projectTitle").value;

    let res = await fetch(`${API}/projects`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        credentials: "include",
        body: JSON.stringify({title})
    });

    if (res.ok) loadProjects();
}

// ---------- TASKS ----------
async function loadTasks() {
    let pid = new URLSearchParams(window.location.search).get("pid");

    let res = await fetch(`${API}/projects/${pid}/tasks`, {
        credentials: "include"
    });

    let data = await res.json();

    let html = "";
    data.forEach(t => {
        html += `<p>${t.Title} - ${t.Status}</p>`;
    });

    document.getElementById("taskList").innerHTML = html;
}

async function createTask() {
    let pid = new URLSearchParams(window.location.search).get("pid");

    let title = document.getElementById("taskTitle").value;
    let description = document.getElementById("taskDesc").value;

    await fetch(`${API}/projects/${pid}/tasks`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        credentials: "include",
        body: JSON.stringify({title, description})
    });

    loadTasks();
}

// Auto-run when pages load
if (window.location.pathname.includes("projects.html")) loadProjects();
if (window.location.pathname.includes("tasks.html")) loadTasks();

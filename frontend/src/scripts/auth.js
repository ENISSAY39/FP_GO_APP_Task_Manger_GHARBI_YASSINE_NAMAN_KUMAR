import { apiFetch } from "./api.js";
import { el } from "./utils.js";

export async function handleLogin() {
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

export async function handleSignup() {
    const name = el("signupName").value.trim();
    const email = el("signupEmail").value.trim();
    const password = el("signupPass").value;

    if (!name || !email || !password) {
        el("signupMsg").textContent = "All fields are required.";
        return;
    }

    el("signupMsg").textContent = "Creating account...";
    const r = await apiFetch("/api/signup", {
        method: "POST",
        body: JSON.stringify({ name, email, password })
    });

    if (!r.ok) {
        el("signupMsg").textContent = "Signup failed: " + (r.json?.error || r.status);
        return;
    }

    el("signupMsg").textContent = "Account created successfully. You can now login.";
    el("loginEmail").value = email;
    el("loginPass").value = password;
    el("signupName").value = "";
    el("signupEmail").value = "";
    el("signupPass").value = "";
}

export async function handleLogout() {
    try {
        el("loginMsg").textContent = "Logging out...";
        const r = await apiFetch("/api/logout", { method: "POST" });
        if (!r.ok) {
            el("loginMsg").textContent = "Logout failed: " + (r.json?.error || r.status);
            return;
        }
        el("jwt").value = "";
        el("loginMsg").textContent = "Logged out.";
    } catch (err) {
        el("loginMsg").textContent = "Logout error: " + (err?.message || String(err));
    }
}

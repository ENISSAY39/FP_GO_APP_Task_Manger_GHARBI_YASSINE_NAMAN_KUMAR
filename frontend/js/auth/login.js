// frontend/js/auth/login.js
import { login as apiLogin } from '../lib/api.js';

function rootPrefix() {
  return window.location.pathname.includes('/frontend/') ? '/frontend' : '';
}

document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('loginForm');
  const emailInput = document.getElementById('email');
  const passwordInput = document.getElementById('password');
  const msg = document.getElementById('msg');
  const btn = document.getElementById('btnLogin');

  if (!form) {
    console.error('login form not found');
    return;
  }

  function show(text, isError = false) {
    if (!msg) return;
    msg.textContent = text || '';
    msg.style.color = isError ? '#9b1c1c' : '#222';
  }

  form.addEventListener('submit', async (ev) => {
    ev.preventDefault();
    show('');
    try {
      if (btn) btn.disabled = true;
      const email = (emailInput?.value || '').trim();
      const password = (passwordInput?.value || '').trim();
      if (!email || !password) { show('Please fill email & password', true); return; }

      const res = await apiLogin({ email, password });

      // apiLogin stored token in localStorage if backend returned it
      show('Login successful');
      const prefix = rootPrefix();
      // redirect to manager dashboard (prefix will handle /frontend vs /)
      window.location.href = `${prefix}/pages/manager/dashboard.html`;
    } catch (err) {
      console.error('login error', err);
      show(err?.message || 'Login failed', true);
    } finally {
      if (btn) btn.disabled = false;
    }
  });
});

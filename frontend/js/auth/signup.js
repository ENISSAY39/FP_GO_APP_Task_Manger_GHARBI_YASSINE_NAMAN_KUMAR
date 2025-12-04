// frontend/js/auth/signup.js
import { signup as apiSignup, login as apiLogin, validate as apiValidate } from '../lib/api.js';

const emailEl = document.getElementById('email');
const passEl = document.getElementById('password');
const pass2El = document.getElementById('password2');
const btn = document.getElementById('btnSignup');
const err = document.getElementById('err');

function showError(msg){
  err.textContent = msg;
  err.classList.remove('hidden');
}
function hideError(){
  err.classList.add('hidden');
  err.textContent = '';
}

btn.addEventListener('click', async () => {
  hideError();
  btn.disabled = true;
  btn.textContent = 'Signing...';

  const email = emailEl.value.trim();
  const p1 = passEl.value;
  const p2 = pass2El.value;

  if (!email || !p1) {
    showError('Email and password required');
    btn.disabled = false;
    btn.textContent = 'Sign up';
    return;
  }
  if (p1.length < 6) {
    showError('Password must be at least 6 characters');
    btn.disabled = false;
    btn.textContent = 'Sign up';
    return;
  }
  if (p1 !== p2) {
    showError("Passwords do not match");
    btn.disabled = false;
    btn.textContent = 'Sign up';
    return;
  }

  try {
    // SIGNUP
    await apiSignup({ email, password: p1 });

    // LOGIN to get cookie (backend sets cookie only at login)
    await apiLogin({ email, password: p1 });

    // Validate and redirect
    const v = await apiValidate();
    if (v && v.user && v.user.is_admin) {
      window.location.href = '/frontend/pages/manager/dashboard.html';
    } else {
      window.location.href = '/frontend/pages/user/dashboard.html';
    }
  } catch (e) {
    showError(e.message || 'Signup/login failed');
  } finally {
    btn.disabled = false;
    btn.textContent = 'Sign up';
  }
});

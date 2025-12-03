// simple API helper
const API_BASE = window.location.origin.replace(/:5500$/, ':3000'); // if you're using live server at 5500 during dev
// If backend served same host: uncomment next line and comment previous:
// const API_BASE = window.location.origin;

const api = {
  get: (path) => fetch(API_BASE + path, { credentials: 'include' }),
  post: (path, body) => fetch(API_BASE + path, { method: 'POST', headers: {'Content-Type':'application/json'}, credentials: 'include', body: JSON.stringify(body) }),
  put: (path, body) => fetch(API_BASE + path, { method: 'PUT', headers: {'Content-Type':'application/json'}, credentials: 'include', body: JSON.stringify(body) }),
  del: (path) => fetch(API_BASE + path, { method: 'DELETE', credentials: 'include' }),
};

# Frontend (vanilla HTML/JS/CSS) for TaskManager

Quick steps:
1. Make sure your backend runs at: http://localhost:3000
2. Serve this `frontend/` folder with a simple static server.
   - Example (npm): `npx http-server ./frontend -c-1 -p 5173`
   - Or open `index.html` via live server extension.
3. Open `http://localhost:5173` (or the port you used).

Important notes:
- Auth uses an HttpOnly cookie set by the backend. All API calls include `credentials: 'include'`.
- On app load, the frontend calls `/api/validate` to get current user, projects and tasks.

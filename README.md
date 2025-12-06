# Backend API Documentation

This document describes the backend architecture and API endpoints for the Task Manager application built with **Go (Gin)**, **GORM**, and **JWT Authentication**.

---

## 📌 Overview

The backend handles:

* User Authentication (Signup, Login)
* Project Management (Create, List, Delete, Membership Control)
* Task Management (Create, Update, Assign, List)
* Role-Based Permissions

The project follows a modular structure:

```
controllers/
    authController.go
    projectsControllers.go
    tasksControllers.go
    permission_helpers.go
models/
initializers/
middlewares/
routes/
```

---

## 🔐 Authentication

### **POST /signup** – Create a New User

**Request Body:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secret123"
}
```

**Validations:**

* Name ≥ 2 characters
* Valid email
* Password ≥ 6 characters

**Response:**

```json
{
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

---

### **POST /login** – Login & Get JWT Token

**Request Body:**

```json
{
  "email": "john@example.com",
  "password": "secret123"
}
```

**Response:**

```json
{
  "token": "<jwt-token>",
  "user": { "id": 1, "name": "John Doe", "email": "john@example.com" }
}
```

Token contains:

* `sub`: user ID
* `exp`: expiry (24h)

---

## 📁 Project Management

All project routes require **JWT authentication**.

### **POST /projects** – Create Project

Request Body:

```json
{
  "name": "My Project",
  "description": "Project description"
}
```

The creator becomes:

* Project Owner
* Project Member with role "OWNER"

---

### **GET /projects/my** – Get All Projects of Logged-in User

Returns list of projects where the user is a member.

---

### **GET /projects/:projectId** – Get Project Details

Only accessible if user is a project member.
Loads:

* Members
* Tasks

---

### **DELETE /projects/:projectId** – Delete Project

Only owner can delete.

---

## 👥 Project Members

### **POST /projects/:projectId/members** – Add Member (OWNER Only)

Body:

```json
{
  "user_id": 5,
  "role": "MEMBER" // optional
}
```

### **DELETE /projects/:projectId/members** – Remove Member

Body:

```json
{
  "user_id": 5
}
```

Owner cannot remove themselves.

---

## 📝 Task Management

### **POST /projects/:projectId/tasks** – Create Task

Members only.

Fields:

```json
{
  "title": "Task Name",
  "description": "Optional text",
  "priority": "HIGH", // optional
  "due_date": "2025-01-20T10:00:00Z"
}
```

---

### **GET /projects/:projectId/tasks** – List Tasks

Includes task assignees.

---

### **PATCH /tasks/:taskId** – Update Task

Allowed by:

* Task creator OR
* Project owner

Partial update fields:

```json
{
  "title": "New title",
  "status": "DONE",
  "priority": "LOW"
}
```

---

## 🔑 Roles & Permission Logic

Defined in `permission_helpers.go`:

* `IsProjectOwner(projectID, userID)`
* `IsProjectMember(projectID, userID)`
* `AddProjectMember()`
* `RemoveProjectMember()`

Roles:

* **OWNER** – full permissions
* **MEMBER** – limited to tasks + viewing

---

## 🧱 Models (Summary)

* User
* Project
* Task
* ProjectMember (pivot table)

---

## 🚀 Running the Server

### Environment Variables

```
JWT_SECRET=your_secret_key
DB_URL=mysql://user:password@tcp(127.0.0.1:3306)/dbname
```

### Start Server

```
go run main.go
```

---

## 📬 Contact

For improvements or bugs, please reach out.

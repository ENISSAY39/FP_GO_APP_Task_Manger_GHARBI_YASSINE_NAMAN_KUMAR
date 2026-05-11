import './styles/app.css'
import './main'

function App() {
  return (
    <div className="wrap">

      <h1>Project Manager — Front (central UI)</h1>

      <div className="card">
        <div className="row">

          <div style={{ flex: 1 }}>
            <label>API base</label>
            <input id="apiBase" defaultValue="http://localhost:3000" />
          </div>

          <div style={{ width: "360px" }}>
            <label>Bearer token (paste or login below)</label>
            <input id="jwt" placeholder="Paste JWT here or login" />
          </div>

          <div style={{ width: "140px" }}>
            <label>&nbsp;</label>

            <div className="row">
              <button id="btnLoad">Load projects</button>
              <button id="btnClearToken" className="ghost">
                Clear
              </button>
            </div>
          </div>

        </div>

        <div style={{ marginTop: "10px" }} className="small">
          Tip: login with email/password (panel below) to auto-set token.
        </div>
      </div>

      {/* LOGIN */}
      <div className="card">

        <h2>Login</h2>

        <div style={{ display: "flex", gap: "12px", flexWrap: "wrap" }}>

          <div style={{ minWidth: "240px", flex: 1 }}>
            <label>Email</label>

            <input
              id="loginEmail"
              placeholder="alice@example.com"
              defaultValue="alice@example.com"
            />
          </div>

          <div style={{ width: "220px" }}>
            <label>Password</label>

            <input
              id="loginPass"
              placeholder="secret"
              defaultValue="secret123"
              type="password"
            />
          </div>

          <div style={{ width: "200px", alignSelf: "end" }}>
            <div className="row">
              <button id="btnLogin">Login</button>

              <button id="btnLogout" className="ghost">
                Logout
              </button>
            </div>
          </div>

        </div>

        <div
          id="loginMsg"
          className="small muted"
          style={{ marginTop: "8px" }}
        ></div>

      </div>

      {/* SIGNUP */}
      <div className="card">

        <h2>Signup</h2>

        <div style={{ display: "flex", gap: "12px", flexWrap: "wrap" }}>

          <div style={{ minWidth: "220px", flex: 1 }}>
            <label>Name</label>
            <input id="signupName" placeholder="John Doe" />
          </div>

          <div style={{ minWidth: "240px", flex: 1 }}>
            <label>Email</label>
            <input id="signupEmail" placeholder="john@example.com" />
          </div>

          <div style={{ width: "220px" }}>
            <label>Password</label>

            <input
              id="signupPass"
              type="password"
              placeholder="secret123"
            />
          </div>

          <div style={{ width: "160px", alignSelf: "end" }}>
            <button id="btnSignup">
              Create account
            </button>
          </div>

        </div>

        <div
          id="signupMsg"
          className="small muted"
          style={{ marginTop: "8px" }}
        ></div>

      </div>

      {/* CREATE PROJECT */}
      <div className="card">

        <h2>Create project</h2>

        <div style={{ display: "flex", gap: "12px", flexWrap: "wrap" }}>

          <div style={{ flex: 1 }}>
            <label>Name</label>
            <input id="pName" placeholder="Project name" />
          </div>

          <div style={{ flex: 2 }}>
            <label>Description</label>
            <input id="pDesc" placeholder="Description" />
          </div>

          <div style={{ width: "140px", alignSelf: "end" }}>
            <button id="btnCreateProject">
              Create
            </button>
          </div>

        </div>

        <div
          id="pStatus"
          className="small muted"
          style={{ marginTop: "8px" }}
        ></div>

      </div>

      {/* PROJECTS */}
      <div className="card">

        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center"
          }}
        >

          <h2>Projects</h2>

          <div className="small">
            Click "View" to open project details.
          </div>

        </div>

        <div
          id="projectsWrap"
          className="grid"
          style={{ marginTop: "10px" }}
        ></div>

      </div>
      {/* PROJECT DETAIL */}
      <div
        id="detailCard"
        className="card"
        style={{ display: "none" }}
      >

        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center"
          }}
        >

          <div>
            <h2 id="detailTitle">Project detail</h2>
            <div className="small" id="detailInfo"></div>
          </div>

          <div>
            <button id="btnRefreshDetail" className="ghost">
              Refresh
            </button>

            <button id="btnCloseDetail" className="ghost">
              Close
            </button>
          </div>

        </div>

        <div
          style={{
            display: "flex",
            gap: "12px",
            marginTop: "12px",
            flexWrap: "wrap"
          }}
        >

          {/* MEMBERS */}
          <div style={{ flex: 1, minWidth: "320px" }}>

            <h3>Members</h3>

            <div id="membersList"></div>

            <div style={{ marginTop: "8px" }}>

              <label>Add member (user id)</label>

              <div style={{ display: "flex", gap: "8px" }}>

                <input
                  id="newMemberId"
                  placeholder="user id (number)"
                />

                <select id="newMemberRole">
                  <option value="MEMBER">MEMBER</option>
                  <option value="OWNER">OWNER</option>
                </select>

                <button id="btnAddMember">
                  Add
                </button>

              </div>

            </div>

          </div>

          {/* TASKS */}
          <div style={{ flex: 1, minWidth: "360px" }}>

            <h3>Tasks</h3>

            <div id="tasksList"></div>

            <div style={{ marginTop: "8px" }}>

              <label>Create task</label>

              <input id="tTitle" placeholder="title" />

              <input
                id="tDesc"
                placeholder="description"
                style={{ marginTop: "6px" }}
              />

              <select
                id="tPriority"
                style={{ marginTop: "6px" }}
              >
                <option>LOW</option>
                <option defaultValue>MEDIUM</option>
                <option>HIGH</option>
              </select>

              <div
                style={{
                  marginTop: "8px",
                  display: "flex",
                  gap: "8px"
                }}
              >

                <button id="btnCreateTask">
                  Create task
                </button>

                <button
                  id="btnBulkAssign"
                  className="ghost"
                >
                  Assign me to all
                </button>

              </div>

            </div>

          </div>

        </div>

        <div style={{ marginTop: "12px" }}>

          <h3>Raw response</h3>

          <pre id="rawDetail"></pre>

        </div>

      </div>
      {/* RAW */}
      <div className="card">
        <h2>Raw response</h2>
        <pre id="raw"></pre>
      </div>

    </div>
  )
}

export default App
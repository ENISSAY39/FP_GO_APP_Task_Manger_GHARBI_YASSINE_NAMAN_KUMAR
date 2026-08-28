import { redirectIfAuthenticated } from '../lib/auth.js'
import { mountPage } from '../lib/mount.jsx'
import HomePage from '../pages/HomePage.jsx'

// A logged-in visitor goes straight to their projects.
if (redirectIfAuthenticated()) {
  mountPage(HomePage)
}

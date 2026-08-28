import { requireAuth } from '../lib/auth.js'
import { mountPage } from '../lib/mount.jsx'
import ProjectsPage from '../pages/ProjectsPage.jsx'

// The guard runs before React mounts, so the page never flashes to a
// logged-out visitor before redirecting.
if (requireAuth()) {
  mountPage(ProjectsPage)
}

import { requireAuth } from '../lib/auth.js'
import { mountPage } from '../lib/mount.jsx'
import ProjectPage from '../pages/ProjectPage.jsx'

if (requireAuth()) {
  mountPage(ProjectPage)
}

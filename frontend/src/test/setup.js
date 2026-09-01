// Adds the jest-dom matchers (toBeInTheDocument, toHaveTextContent, ...) to
// expect. Vitest's `globals: true` gives Testing Library its automatic cleanup
// between tests, so nothing else is needed here.
import '@testing-library/jest-dom/vitest'

import { useState } from 'react'
import AuthLayout from '../components/AuthLayout.jsx'
import Alert from '../components/ui/Alert.jsx'
import Button from '../components/ui/Button.jsx'
import Field from '../components/ui/Field.jsx'
import { signup } from '../lib/auth.js'

/** Mirrors the binding rules of the Go handler so errors show up before the call. */
function validate({ name, email, password }) {
  const errors = {}
  if (name.trim().length < 2) errors.name = 'At least 2 characters.'
  if (!/^\S+@\S+\.\S+$/.test(email.trim())) errors.email = 'Enter a valid email address.'
  if (password.length < 6) errors.password = 'At least 6 characters.'
  return errors
}

export default function SignupPage() {
  const [values, setValues] = useState({ name: '', email: '', password: '' })
  const [errors, setErrors] = useState({})
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const update = (key) => (event) => {
    const { value } = event.target
    setValues((current) => ({ ...current, [key]: value }))
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    const found = validate(values)
    setErrors(found)
    if (Object.keys(found).length > 0) return

    setSubmitting(true)
    setError('')
    try {
      // signup() creates the account then logs straight in, so the user
      // does not have to type the same credentials twice.
      await signup({
        name: values.name.trim(),
        email: values.email.trim(),
        password: values.password,
      })
      window.location.href = '/projects.html'
    } catch (signupError) {
      setError(signupError.message)
      setSubmitting(false)
    }
  }

  return (
    <AuthLayout
      title="Create your account"
      subtitle="You will be logged in right after."
      footer={
        <>
          Already registered?{' '}
          <a href="/login.html" className="text-accent-400 hover:text-accent-300">
            Log in
          </a>
        </>
      }
    >
      <form onSubmit={handleSubmit} noValidate className="space-y-4">
        <Alert tone="error">{error}</Alert>

        <Field
          label="Name"
          placeholder="Jane Doe"
          autoComplete="name"
          value={values.name}
          onChange={update('name')}
          error={errors.name}
          autoFocus
        />

        <Field
          label="Email"
          type="email"
          autoComplete="email"
          placeholder="you@example.com"
          value={values.email}
          onChange={update('email')}
          error={errors.email}
        />

        <Field
          label="Password"
          type="password"
          autoComplete="new-password"
          placeholder="••••••••"
          value={values.password}
          onChange={update('password')}
          error={errors.password}
          hint={errors.password ? undefined : 'At least 6 characters.'}
        />

        <Button type="submit" className="w-full" loading={submitting}>
          Create account
        </Button>
      </form>
    </AuthLayout>
  )
}

import { create } from 'zustand'
import { devtools } from 'zustand/middleware'

interface User {
  id: string
  email: string
  name: string
  role: 'admin' | 'user'
}

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (user: User) => void
  logout: () => void
  setLoading: (loading: boolean) => void
}

export const useAuthStore = create<AuthState>()(
  devtools(
    (set) => ({
      user: null,
      isAuthenticated: false,
      isLoading: false,
      login: (user) =>
        set(
          { user, isAuthenticated: true, isLoading: false },
          false,
          'auth/login'
        ),
      logout: () =>
        set(
          { user: null, isAuthenticated: false, isLoading: false },
          false,
          'auth/logout'
        ),
      setLoading: (isLoading) =>
        set({ isLoading }, false, 'auth/setLoading'),
    }),
    { name: 'auth-store' }
  )
)
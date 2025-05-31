import { create } from 'zustand'
import { devtools } from 'zustand/middleware'
import { goalsApi, type GoalResponse } from '../api/goals'

interface GoalState {
  currentGoal: string
  response: string | null
  isLoading: boolean
  error: string | null
  setCurrentGoal: (goal: string) => void
  submitGoal: (goal: string) => Promise<void>
  clearResponse: () => void
  clearError: () => void
}

export const useGoalStore = create<GoalState>()(
  devtools(
    (set, get) => ({
      currentGoal: '',
      response: null,
      isLoading: false,
      error: null,

      setCurrentGoal: (goal: string) =>
        set({ currentGoal: goal }, false, 'goal/setCurrentGoal'),

      submitGoal: async (goal: string) => {
        set({ isLoading: true, error: null, response: null }, false, 'goal/submitStart')
        
        try {
          const result = await goalsApi.submitGoal(goal)
          
          if (result.success && result.response) {
            set(
              { 
                response: result.response, 
                isLoading: false,
                currentGoal: ''
              },
              false,
              'goal/submitSuccess'
            )
          } else {
            set(
              { 
                error: result.error || 'Unknown error occurred',
                isLoading: false 
              },
              false,
              'goal/submitError'
            )
          }
        } catch (error) {
          set(
            { 
              error: error instanceof Error ? error.message : 'Failed to submit goal',
              isLoading: false 
            },
            false,
            'goal/submitError'
          )
        }
      },

      clearResponse: () =>
        set({ response: null }, false, 'goal/clearResponse'),

      clearError: () =>
        set({ error: null }, false, 'goal/clearError'),
    }),
    { name: 'goal-store' }
  )
)
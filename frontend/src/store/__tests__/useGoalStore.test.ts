import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useGoalStore } from '../useGoalStore'
import { goalsApi } from '../../api/goals'

// Mock the API
vi.mock('../../api/goals')
const mockGoalsApi = vi.mocked(goalsApi)

describe('useGoalStore', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    useGoalStore.setState({
      currentGoal: '',
      response: null,
      mediaRecommendations: null,
      isLoading: false,
      error: null,
    })
  })

  it('should have initial state', () => {
    const state = useGoalStore.getState()
    expect(state.currentGoal).toBe('')
    expect(state.response).toBeNull()
    expect(state.mediaRecommendations).toBeNull()
    expect(state.isLoading).toBe(false)
    expect(state.error).toBeNull()
  })

  it('should set current goal', () => {
    useGoalStore.getState().setCurrentGoal('Learn to code')
    const state = useGoalStore.getState()
    expect(state.currentGoal).toBe('Learn to code')
  })

  it('should clear response', () => {
    useGoalStore.setState({ 
      response: 'Some response',
      mediaRecommendations: {
        podcasts: [],
        streaming: [],
        books: [],
        websites: []
      }
    })
    useGoalStore.getState().clearResponse()
    const state = useGoalStore.getState()
    expect(state.response).toBeNull()
    expect(state.mediaRecommendations).toBeNull()
  })

  it('should clear error', () => {
    useGoalStore.setState({ error: 'Some error' })
    useGoalStore.getState().clearError()
    const state = useGoalStore.getState()
    expect(state.error).toBeNull()
  })

  it('should handle successful goal submission', async () => {
    const mockResponse = {
      success: true,
      response: 'Great goal! Here are some steps...',
      timestamp: '2024-01-01T00:00:00Z'
    }
    
    mockGoalsApi.submitGoal.mockResolvedValue(mockResponse)

    await useGoalStore.getState().submitGoal('Learn to play guitar')
    
    const state = useGoalStore.getState()
    expect(state.response).toBe('Great goal! Here are some steps...')
    expect(state.isLoading).toBe(false)
    expect(state.error).toBeNull()
    expect(state.currentGoal).toBe('')
  })

  it('should handle API error response', async () => {
    const mockResponse = {
      success: false,
      error: 'Invalid goal format',
      timestamp: '2024-01-01T00:00:00Z'
    }
    
    mockGoalsApi.submitGoal.mockResolvedValue(mockResponse)

    await useGoalStore.getState().submitGoal('Bad goal')
    
    const state = useGoalStore.getState()
    expect(state.error).toBe('Invalid goal format')
    expect(state.isLoading).toBe(false)
    expect(state.response).toBeNull()
  })

  it('should handle network error', async () => {
    mockGoalsApi.submitGoal.mockRejectedValue(new Error('Network error'))

    await useGoalStore.getState().submitGoal('Some goal')
    
    const state = useGoalStore.getState()
    expect(state.error).toBe('Network error')
    expect(state.isLoading).toBe(false)
    expect(state.response).toBeNull()
  })

  it('should set loading state during submission', async () => {
    let resolvePromise: (value: any) => void
    const promise = new Promise(resolve => {
      resolvePromise = resolve
    })
    
    mockGoalsApi.submitGoal.mockReturnValue(promise)

    // Start submission
    const submitPromise = useGoalStore.getState().submitGoal('Test goal')
    
    // Check loading state
    expect(useGoalStore.getState().isLoading).toBe(true)
    expect(useGoalStore.getState().error).toBeNull()
    expect(useGoalStore.getState().response).toBeNull()
    
    // Resolve the promise
    resolvePromise!({
      success: true,
      response: 'Test response',
      timestamp: '2024-01-01T00:00:00Z'
    })
    
    await submitPromise
    
    // Check final state
    expect(useGoalStore.getState().isLoading).toBe(false)
  })

  it('should handle successful goal submission with media recommendations', async () => {
    const mockResponse = {
      success: true,
      response: 'Great goal! Here are some steps...',
      mediaRecommendations: {
        podcasts: [
          { title: 'Test Podcast', platform: 'Spotify', link: 'https://spotify.com', description: 'A test podcast' }
        ],
        streaming: [
          { title: 'Test Video', platform: 'YouTube', description: 'A test video' }
        ],
        books: [
          { title: 'Test Book', link: 'https://amazon.com', description: 'A test book' }
        ],
        websites: [
          { title: 'Test Website', link: 'https://example.com', description: 'A test website' }
        ]
      },
      timestamp: '2024-01-01T00:00:00Z'
    }
    
    mockGoalsApi.submitGoal.mockResolvedValue(mockResponse)

    await useGoalStore.getState().submitGoal('Learn to play guitar')
    
    const state = useGoalStore.getState()
    expect(state.response).toBe('Great goal! Here are some steps...')
    expect(state.mediaRecommendations).toEqual(mockResponse.mediaRecommendations)
    expect(state.isLoading).toBe(false)
    expect(state.error).toBeNull()
    expect(state.currentGoal).toBe('')
  })
})
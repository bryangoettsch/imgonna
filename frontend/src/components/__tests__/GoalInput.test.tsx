import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { GoalInput } from '../GoalInput'
import { useGoalStore } from '../../store/useGoalStore'

// Mock the store
vi.mock('../../store/useGoalStore')
const mockUseGoalStore = vi.mocked(useGoalStore)

describe('GoalInput', () => {
  const mockSubmitGoal = vi.fn()
  const mockClearError = vi.fn()
  const mockClearResponse = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    
    mockUseGoalStore.mockReturnValue({
      currentGoal: '',
      response: null,
      isLoading: false,
      error: null,
      setCurrentGoal: vi.fn(),
      submitGoal: mockSubmitGoal,
      clearResponse: mockClearResponse,
      clearError: mockClearError,
    })
  })

  it('renders the goal input form', () => {
    render(<GoalInput />)
    
    expect(screen.getByText("What's your goal?")).toBeInTheDocument()
    expect(screen.getByPlaceholderText('Enter your goal here...')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Get AI Guidance' })).toBeInTheDocument()
  })

  it('updates input value when typing', async () => {
    const user = userEvent.setup()
    render(<GoalInput />)
    
    const input = screen.getByPlaceholderText('Enter your goal here...')
    await user.type(input, 'Learn to code')
    
    expect(input).toHaveValue('Learn to code')
  })

  it('shows character count', async () => {
    const user = userEvent.setup()
    render(<GoalInput />)
    
    const input = screen.getByPlaceholderText('Enter your goal here...')
    await user.type(input, 'Test goal')
    
    expect(screen.getByText('9/500')).toBeInTheDocument()
  })

  it('submits goal when form is submitted', async () => {
    const user = userEvent.setup()
    render(<GoalInput />)
    
    const input = screen.getByPlaceholderText('Enter your goal here...')
    const submitButton = screen.getByRole('button', { name: 'Get AI Guidance' })
    
    await user.type(input, 'Learn to play guitar')
    await user.click(submitButton)
    
    expect(mockSubmitGoal).toHaveBeenCalledWith('Learn to play guitar')
    expect(mockClearError).toHaveBeenCalled()
    expect(mockClearResponse).toHaveBeenCalled()
  })

  it('shows loading state when submitting', () => {
    mockUseGoalStore.mockReturnValue({
      currentGoal: '',
      response: null,
      isLoading: true,
      error: null,
      setCurrentGoal: vi.fn(),
      submitGoal: mockSubmitGoal,
      clearResponse: mockClearResponse,
      clearError: mockClearError,
    })

    render(<GoalInput />)
    
    expect(screen.getByText('Getting guidance...')).toBeInTheDocument()
    expect(screen.getByRole('button')).toBeDisabled()
  })

  it('shows response when available', () => {
    mockUseGoalStore.mockReturnValue({
      currentGoal: '',
      response: 'Great goal! Here are some steps to get started...',
      isLoading: false,
      error: null,
      setCurrentGoal: vi.fn(),
      submitGoal: mockSubmitGoal,
      clearResponse: mockClearResponse,
      clearError: mockClearError,
    })

    render(<GoalInput />)
    
    expect(screen.getByText('Your personalized guidance:')).toBeInTheDocument()
    expect(screen.getByText('Great goal! Here are some steps to get started...')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Set Another Goal' })).toBeInTheDocument()
  })

  it('shows error message when there is an error', () => {
    mockUseGoalStore.mockReturnValue({
      currentGoal: '',
      response: null,
      isLoading: false,
      error: 'Failed to process goal',
      setCurrentGoal: vi.fn(),
      submitGoal: mockSubmitGoal,
      clearResponse: mockClearResponse,
      clearError: mockClearError,
    })

    render(<GoalInput />)
    
    expect(screen.getByText('Failed to process goal')).toBeInTheDocument()
  })

  it('clears error when close button is clicked', async () => {
    const user = userEvent.setup()
    
    mockUseGoalStore.mockReturnValue({
      currentGoal: '',
      response: null,
      isLoading: false,
      error: 'Failed to process goal',
      setCurrentGoal: vi.fn(),
      submitGoal: mockSubmitGoal,
      clearResponse: mockClearResponse,
      clearError: mockClearError,
    })

    render(<GoalInput />)
    
    const closeButton = screen.getByText('×')
    await user.click(closeButton)
    
    expect(mockClearError).toHaveBeenCalled()
  })

  it('resets form when "Set Another Goal" is clicked', async () => {
    const user = userEvent.setup()
    
    mockUseGoalStore.mockReturnValue({
      currentGoal: '',
      response: 'Great goal! Here are some steps...',
      isLoading: false,
      error: null,
      setCurrentGoal: vi.fn(),
      submitGoal: mockSubmitGoal,
      clearResponse: mockClearResponse,
      clearError: mockClearError,
    })

    render(<GoalInput />)
    
    const newGoalButton = screen.getByRole('button', { name: 'Set Another Goal' })
    await user.click(newGoalButton)
    
    expect(mockClearResponse).toHaveBeenCalled()
    expect(mockClearError).toHaveBeenCalled()
  })
})
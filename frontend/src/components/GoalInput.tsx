import { useState } from 'react'
import { useGoalStore } from '../store/useGoalStore'
import { MediaRecommendations } from './MediaRecommendations'

export function GoalInput() {
  const [inputValue, setInputValue] = useState('')
  const { submitGoal, isLoading, response, mediaRecommendations, error, clearError, clearResponse } = useGoalStore()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!inputValue.trim()) {
      return
    }

    // Clear previous responses/errors
    clearError()
    clearResponse()
    
    await submitGoal(inputValue.trim())
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value)
    // Clear error when user starts typing
    if (error) {
      clearError()
    }
  }

  const handleNewGoal = () => {
    setInputValue('')
    clearResponse()
    clearError()
  }

  return (
    <div className="w-full max-w-6xl mx-auto p-6">
      <div className="text-center mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          What's your goal?
        </h1>
        <p className="text-gray-600">
          Share your goal and get personalized guidance to achieve it
        </p>
      </div>

      {!response ? (
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="relative">
            <input
              type="text"
              value={inputValue}
              onChange={handleInputChange}
              placeholder="Enter your goal here..."
              maxLength={500}
              disabled={isLoading}
              className="w-full px-4 py-3 text-lg border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none disabled:bg-gray-100 disabled:cursor-not-allowed"
            />
            <div className="text-right mt-1 text-sm text-gray-500">
              {inputValue.length}/500
            </div>
          </div>

          <button
            type="submit"
            disabled={isLoading || !inputValue.trim()}
            className="w-full py-3 px-6 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            {isLoading ? (
              <div className="flex items-center justify-center">
                <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
                Getting guidance...
              </div>
            ) : (
              'Get AI Guidance'
            )}
          </button>
        </form>
      ) : (
        <div className="space-y-6">
          <div className="bg-green-50 border border-green-200 rounded-lg p-6">
            <h3 className="text-lg font-semibold text-green-800 mb-3">
              Your personalized guidance:
            </h3>
            <div className="text-green-700 whitespace-pre-wrap leading-relaxed">
              {response}
            </div>
          </div>
          
          {mediaRecommendations && (
            <MediaRecommendations recommendations={mediaRecommendations} />
          )}
          
          <button
            onClick={handleNewGoal}
            className="w-full py-3 px-6 bg-gray-600 text-white font-medium rounded-lg hover:bg-gray-700 focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-colors"
          >
            Set Another Goal
          </button>
        </div>
      )}

      {error && (
        <div className="mt-4 bg-red-50 border border-red-200 rounded-lg p-4">
          <div className="flex">
            <div className="text-red-700">
              <strong>Error:</strong> {error}
            </div>
            <button
              onClick={clearError}
              className="ml-auto text-red-500 hover:text-red-700"
            >
              ×
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
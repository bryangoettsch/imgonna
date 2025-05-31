export interface GoalRequest {
  goal: string;
}

export interface GoalResponse {
  success: boolean;
  response?: string;
  error?: string;
  timestamp: string;
}

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const goalsApi = {
  async submitGoal(goal: string): Promise<GoalResponse> {
    const response = await fetch(`${API_BASE_URL}/api/v1/goals`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ goal }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  },
};
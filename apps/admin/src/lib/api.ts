// Types for the API
export interface Project {
  id: string;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface Environment {
  id: string;
  project_id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Flag {
  id: string;
  project_id: string;
  key: string;
  description: string;
  type: 'boolean' | 'string' | 'number' | 'json';
  created_at: string;
  updated_at: string;
}

export interface FlagValue {
  id: string;
  flag_id: string;
  env_id: string;
  value: string;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface FlagWithValues extends Flag {
  values?: FlagValue[];
}

// Request types
export interface CreateProjectRequest {
  name: string;
  description?: string;
}

export interface UpdateProjectRequest {
  name?: string;
  description?: string;
}

export interface CreateEnvironmentRequest {
  project_id: string;
  name: string;
}

export interface UpdateEnvironmentRequest {
  name?: string;
}

export interface CreateFlagRequest {
  project_id: string;
  key: string;
  description?: string;
  type: 'boolean' | 'string' | 'number' | 'json';
}

export interface UpdateFlagRequest {
  key?: string;
  description?: string;
  type?: 'boolean' | 'string' | 'number' | 'json';
}

export interface CreateFlagValueRequest {
  flag_id: string;
  env_id: string;
  value: string;
  enabled: boolean;
}

export interface UpdateFlagValueRequest {
  value?: string;
  enabled?: boolean;
}

// API client
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    const response = await fetch(url, config);

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Unknown error' }));
      throw new Error(error.error || `HTTP error! status: ${response.status}`);
    }

    // Handle 204 No Content
    if (response.status === 204) {
      return null as T;
    }

    return response.json();
  }

  // Projects
  async getProjects(): Promise<Project[]> {
    return this.request<Project[]>('/api/projects');
  }

  async getProject(id: string): Promise<Project> {
    return this.request<Project>(`/api/projects/${id}`);
  }

  async createProject(data: CreateProjectRequest): Promise<Project> {
    return this.request<Project>('/api/projects', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateProject(id: string, data: UpdateProjectRequest): Promise<Project> {
    return this.request<Project>(`/api/projects/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteProject(id: string): Promise<void> {
    return this.request<void>(`/api/projects/${id}`, {
      method: 'DELETE',
    });
  }

  // Environments
  async getEnvironments(): Promise<Environment[]> {
    return this.request<Environment[]>('/api/environments');
  }

  async getEnvironment(id: string): Promise<Environment> {
    return this.request<Environment>(`/api/environments/${id}`);
  }

  async getProjectEnvironments(projectId: string): Promise<Environment[]> {
    return this.request<Environment[]>(`/api/projects/${projectId}/environments`);
  }

  async createEnvironment(data: CreateEnvironmentRequest): Promise<Environment> {
    return this.request<Environment>('/api/environments', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateEnvironment(id: string, data: UpdateEnvironmentRequest): Promise<Environment> {
    return this.request<Environment>(`/api/environments/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteEnvironment(id: string): Promise<void> {
    return this.request<void>(`/api/environments/${id}`, {
      method: 'DELETE',
    });
  }

  // Flags
  async getProjectFlags(projectId: string, includeValues = false): Promise<FlagWithValues[]> {
    return this.request<FlagWithValues[]>(
      `/api/projects/${projectId}/flags?includeValues=${includeValues}`
    );
  }

  async getFlag(id: string): Promise<Flag> {
    return this.request<Flag>(`/api/flags/${id}`);
  }

  async createFlag(data: CreateFlagRequest): Promise<Flag> {
    return this.request<Flag>('/api/flags', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateFlag(id: string, data: UpdateFlagRequest): Promise<Flag> {
    return this.request<Flag>(`/api/flags/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteFlag(id: string): Promise<void> {
    return this.request<void>(`/api/flags/${id}`, {
      method: 'DELETE',
    });
  }

  // Flag Values
  async getFlagValues(flagId: string): Promise<FlagValue[]> {
    return this.request<FlagValue[]>(`/api/flags/${flagId}/values`);
  }

  async getEnvironmentFlags(envId: string): Promise<FlagValue[]> {
    return this.request<FlagValue[]>(`/api/environments/${envId}/flags`);
  }

  async createOrUpdateFlagValue(data: CreateFlagValueRequest): Promise<FlagValue> {
    return this.request<FlagValue>('/api/flags/values', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateFlagValue(id: string, data: UpdateFlagValueRequest): Promise<FlagValue> {
    return this.request<FlagValue>(`/api/flags/values/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteFlagValue(id: string): Promise<void> {
    return this.request<void>(`/api/flags/values/${id}`, {
      method: 'DELETE',
    });
  }
}

export const apiClient = new ApiClient();

// SSE client for real-time updates
export class SSEClient {
  private eventSource: EventSource | null = null;
  private listeners: Map<string, Set<(data: any) => void>> = new Map();

  connect(url: string = `${API_BASE_URL}/api/events`) {
    this.eventSource = new EventSource(url);

    this.eventSource.onopen = () => {
      console.log('SSE connection established');
    };

    this.eventSource.onerror = (error) => {
      console.error('SSE connection error:', error);
    };

    this.eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        const eventType = data.type || event.type;
        
        const listeners = this.listeners.get(eventType);
        if (listeners) {
          listeners.forEach(callback => callback(data));
        }
      } catch (error) {
        console.error('Error parsing SSE message:', error);
      }
    };
  }

  disconnect() {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
  }

  on(eventType: string, callback: (data: any) => void) {
    if (!this.listeners.has(eventType)) {
      this.listeners.set(eventType, new Set());
    }
    this.listeners.get(eventType)!.add(callback);
  }

  off(eventType: string, callback: (data: any) => void) {
    const listeners = this.listeners.get(eventType);
    if (listeners) {
      listeners.delete(callback);
      if (listeners.size === 0) {
        this.listeners.delete(eventType);
      }
    }
  }
}

export const sseClient = new SSEClient();

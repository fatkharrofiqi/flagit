import { create } from 'zustand';
import type { Project, Environment, FlagWithValues } from '../lib/api';

interface AppState {
  // Current project
  currentProject: Project | null;
  setCurrentProject: (project: Project | null) => void;
  
  // Projects
  projects: Project[];
  setProjects: (projects: Project[]) => void;
  addProject: (project: Project) => void;
  updateProject: (id: string, project: Partial<Project>) => void;
  removeProject: (id: string) => void;
  
  // Environments
  environments: Environment[];
  setEnvironments: (environments: Environment[]) => void;
  addEnvironment: (environment: Environment) => void;
  updateEnvironment: (id: string, environment: Partial<Environment>) => void;
  removeEnvironment: (id: string) => void;
  
  // Flags
  flags: FlagWithValues[];
  setFlags: (flags: FlagWithValues[]) => void;
  addFlag: (flag: FlagWithValues) => void;
  updateFlag: (id: string, flag: Partial<FlagWithValues>) => void;
  removeFlag: (id: string) => void;
  
  // UI state
  sidebarOpen: boolean;
  setSidebarOpen: (open: boolean) => void;
  
  // Real-time updates
  lastUpdated: number;
  setLastUpdated: () => void;
}

export const useStore = create<AppState>((set) => ({
  // Current project
  currentProject: null,
  setCurrentProject: (project) => set({ currentProject: project }),
  
  // Projects
  projects: [],
  setProjects: (projects) => set({ projects }),
  addProject: (project) => set((state) => ({ 
    projects: [...state.projects, project] 
  })),
  updateProject: (id, updatedProject) => set((state) => ({
    projects: state.projects.map(p => 
      p.id === id ? { ...p, ...updatedProject } : p
    )
  })),
  removeProject: (id) => set((state) => ({
    projects: state.projects.filter(p => p.id !== id),
    currentProject: state.currentProject?.id === id ? null : state.currentProject
  })),
  
  // Environments
  environments: [],
  setEnvironments: (environments) => set({ environments }),
  addEnvironment: (environment) => set((state) => ({ 
    environments: [...state.environments, environment] 
  })),
  updateEnvironment: (id, updatedEnvironment) => set((state) => ({
    environments: state.environments.map(e => 
      e.id === id ? { ...e, ...updatedEnvironment } : e
    )
  })),
  removeEnvironment: (id) => set((state) => ({
    environments: state.environments.filter(e => e.id !== id)
  })),
  
  // Flags
  flags: [],
  setFlags: (flags) => set({ flags }),
  addFlag: (flag) => set((state) => ({ 
    flags: [...state.flags, flag] 
  })),
  updateFlag: (id, updatedFlag) => set((state) => ({
    flags: state.flags.map(f => 
      f.id === id ? { ...f, ...updatedFlag } : f
    )
  })),
  removeFlag: (id) => set((state) => ({
    flags: state.flags.filter(f => f.id !== id)
  })),
  
  // UI state
  sidebarOpen: true,
  setSidebarOpen: (open) => set({ sidebarOpen: open }),
  
  // Real-time updates
  lastUpdated: Date.now(),
  setLastUpdated: () => set({ lastUpdated: Date.now() })
}));

import { create } from 'zustand';
import type { Task } from '../api/types';

interface AppState {
  tasks: Record<string, Task>;
  updateTask: (id: string, task: Task) => void;
}

export const useAppStore = create<AppState>((set) => ({
  tasks: {},
  updateTask: (id, task) =>
    set((s) => ({ tasks: { ...s.tasks, [id]: task } })),
}));

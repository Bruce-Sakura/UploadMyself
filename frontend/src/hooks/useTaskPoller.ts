import { getTask } from '../api/endpoints';
import type { Task } from '../api/types';

export interface PollOptions {
  /** Polling interval in ms (default: 2000) */
  interval?: number;
  /** Max polling attempts (default: 60) */
  maxAttempts?: number;
  /** Called on each poll with current task */
  onProgress?: (task: Task) => void;
  /** Called when task completes */
  onDone?: (task: Task) => void;
  /** Called when task fails */
  onFailed?: (task: Task) => void;
}

/**
 * Shared task polling utility.
 * Replaces duplicated pollTask logic across SkillClone, VoiceClone, AvatarPage.
 */
export async function pollTask(taskId: string, options: PollOptions = {}): Promise<Task | null> {
  const {
    interval = 2000,
    maxAttempts = 60,
    onProgress,
    onDone,
    onFailed,
  } = options;

  let attempts = 0;
  while (attempts < maxAttempts) {
    await new Promise((r) => setTimeout(r, interval));
    try {
      const { data: task } = (await getTask(taskId)) as { data: Task };

      if (task.status === 'done' || task.status === 'completed') {
        onDone?.(task);
        return task;
      }
      if (task.status === 'failed') {
        onFailed?.(task);
        return null;
      }

      onProgress?.(task);
    } catch {
      // continue polling on network errors
    }
    attempts++;
  }
  return null;
}

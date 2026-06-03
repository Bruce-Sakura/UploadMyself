import api from './client';
import type { Skill, Voice, Avatar, Task, UploadResult } from './types';

/* ── Skills ─────────────────────────────────────────────── */
export const createSkill = (data: { name: string; corpus: string }) =>
  api.post<Skill>('/skills', data);

export const listSkills = () => api.get<Skill[]>('/skills');

export const getSkill = (id: string) => api.get<Skill>(`/skills/${id}`);

export const processSkill = (id: string) =>
  api.post<Task>(`/skills/${id}/process`);

export const deleteSkill = (id: string) => api.delete(`/skills/${id}`);

/* ── Voices ─────────────────────────────────────────────── */
export const createVoice = (data: { name: string; audio_path: string }) =>
  api.post<Voice>('/voices', data);

export const listVoices = () => api.get<Voice[]>('/voices');

export const getVoice = (id: string) => api.get<Voice>(`/voices/${id}`);

export const trainVoice = (id: string) =>
  api.post<Task>(`/voices/${id}/train`);

export const synthesizeVoice = (id: string, text: string) =>
  api.post<{ audio_url: string }>(`/voices/${id}/synthesize`, { text });

export const deleteVoice = (id: string) => api.delete(`/voices/${id}`);

/* ── Avatars ────────────────────────────────────────────── */
export const createAvatar = (data: {
  name: string;
  type: '2d' | '3d';
  photo_path: string;
  style?: string;
}) => api.post<Avatar>('/avatars', data);

export const listAvatars = (type?: '2d' | '3d') =>
  api.get<Avatar[]>('/avatars', { params: type ? { type } : {} });

export const getAvatar = (id: string) => api.get<Avatar>(`/avatars/${id}`);

export const processAvatar = (id: string) =>
  api.post<Task>(`/avatars/${id}/process`);

export const deleteAvatar = (id: string) => api.delete(`/avatars/${id}`);

/* ── Upload ─────────────────────────────────────────────── */
export const uploadFile = (file: File) => {
  const fd = new FormData();
  fd.append('file', file);
  return api.post<UploadResult>('/upload', fd);
};

/* ── Tasks ──────────────────────────────────────────────── */
export const getTask = (id: string) => api.get<Task>(`/tasks/${id}`);

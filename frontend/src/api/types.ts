export interface Skill {
  id: string;
  name: string;
  corpus: string;
  status: string;
  result: string;
  created_at: string;
}

export interface Voice {
  id: string;
  name: string;
  audio_path: string;
  duration: number;
  status: string;
  created_at: string;
}

export interface Avatar {
  id: string;
  name: string;
  type: '2d' | '3d';
  photo_path: string;
  style: string;
  status: string;
  result: string;
  created_at: string;
}

export interface Task {
  id: string;
  type: string;
  ref_id: string;
  status: string;
  progress: number;
  error: string;
}

export interface UploadResult {
  id: string;
  filename: string;
  path: string;
  size: number;
}

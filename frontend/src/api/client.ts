import axios from 'axios';

const api = axios.create({ baseURL: '/api/v1' });

api.interceptors.response.use(
  (res) => res,
  (err) => {
    const msg = err.response?.data?.detail ?? err.message;
    console.error('[API]', msg);
    return Promise.reject(err);
  },
);

export default api;

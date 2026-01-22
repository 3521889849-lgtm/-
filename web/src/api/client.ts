import axios from 'axios';

const client = axios.create({
  baseURL: '/', // Use proxy
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

client.interceptors.response.use(
  (response) => response.data,
  (error) => {
    return Promise.reject(error);
  }
);

export default client;

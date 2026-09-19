import axios from "axios";

const client = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
});

let authToken: string | null = null;

export function setAuthToken(token: string | null) {
  authToken = token;
}

client.interceptors.request.use((config) => {
  if (authToken) {
    config.headers.Authorization = `Bearer ${authToken}`;
  }
  return config;
});

client.interceptors.response.use(
  (res) => res,
  (err) => {
    const isLoginRequest = err.config?.url?.includes("/auth/login");

    if (err.response?.status === 401 && !isLoginRequest && typeof window !== "undefined") {
      setAuthToken(null);
      window.location.href = "/login";
    }
    return Promise.reject(err);
  }
);

export default client;
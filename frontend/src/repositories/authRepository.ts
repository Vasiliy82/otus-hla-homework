const AUTH_TOKEN_KEY = 'auth_token';
const DECODED_TOKEN_KEY = 'decoded_token';
const USER_PROFILE_KEY = 'user_profile';


export const authRepository = {
  setToken: (token: string) => {
    localStorage.setItem(AUTH_TOKEN_KEY, token);
  },
  getToken: () => {
    return localStorage.getItem(AUTH_TOKEN_KEY);
  },
  removeToken: () => {
    localStorage.removeItem(AUTH_TOKEN_KEY);
  },

  // Для хранения декодированных данных токена
  setDecodedToken: (decodedToken: any) => {
    localStorage.setItem(DECODED_TOKEN_KEY, JSON.stringify(decodedToken));
  },
  getDecodedToken: () => {
    const token = localStorage.getItem(DECODED_TOKEN_KEY);
    return token ? JSON.parse(token) : null;
  },
  removeDecodedToken: () => {
    localStorage.removeItem(DECODED_TOKEN_KEY);
  },

  setProfile: (profile: any) => {
    localStorage.setItem(USER_PROFILE_KEY, JSON.stringify(profile));
  },
  getProfile: () => {
    const profile = localStorage.getItem(USER_PROFILE_KEY);
    return profile ? JSON.parse(profile) : null;
  },
  removeProfile: () => {
    localStorage.removeItem(USER_PROFILE_KEY);
  },
};
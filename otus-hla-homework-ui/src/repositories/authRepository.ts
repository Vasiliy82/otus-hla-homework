const TOKEN_KEY = 'token';
const PROFILE_KEY = 'profile';

// Репозиторий для работы с токенами и профилем пользователя
export const authRepository = {
  // Работа с токеном
  getToken: () => localStorage.getItem(TOKEN_KEY),
  setToken: (token: string) => localStorage.setItem(TOKEN_KEY, token),
  removeToken: () => localStorage.removeItem(TOKEN_KEY),

  // Работа с профилем пользователя
  getProfile: () => {
    const profile = localStorage.getItem(PROFILE_KEY);
    return profile ? JSON.parse(profile) : null;
  },
  setProfile: (profile: any) => localStorage.setItem(PROFILE_KEY, JSON.stringify(profile)),
  removeProfile: () => localStorage.removeItem(PROFILE_KEY),
};

export default authRepository;
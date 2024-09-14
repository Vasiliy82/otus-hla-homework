import { authRepository } from '../repositories/authRepository';

export const authService = {
  // Проверка аутентификации
  isAuthenticated: () => !!authRepository.getToken(),

  // Логин (сохранение токена и профиля)
  login: (token: string, profile: any) => {
    authRepository.setToken(token);
    authRepository.setProfile(profile);
  },

  // Логаут (удаление токена и профиля)
  logout: () => {
    authRepository.removeToken();
    authRepository.removeProfile();
  },

  // Получение профиля пользователя
  getProfile: () => authRepository.getProfile(),
  getToken: () => authRepository.getToken(),
};
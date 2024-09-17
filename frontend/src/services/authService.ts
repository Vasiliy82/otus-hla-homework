import * as jwtDecode from 'jwt-decode';
import { authRepository } from '../repositories/authRepository';
import api from './api';

export interface DecodedToken {
  userId: string;
  tokenId: string;
  permissions: string[];
  exp: number;
}

let interceptorId: number | null = null;

export const authService = {
  // Проверка аутентификации
  isAuthenticated: () => !!authRepository.getToken(),

  // Логин (сохранение токена и профиля)
  login: (token: string, profile: any) => {
    authRepository.setToken(token);
    authRepository.setProfile(profile);
    
    // Парсим и сохраняем данные из токена
    const decodedToken = authService.parseToken(token);
    authRepository.setDecodedToken(decodedToken);

    // Настраиваем interceptor после логина
    authService.setupAxiosInterceptor(token);
  },

  // Логаут (удаление токена и профиля)
  logout: () => {
    authRepository.removeToken();
    authRepository.removeProfile();
    authService.removeAxiosInterceptor();
  },

  // Получение профиля пользователя
  getProfile: () => authRepository.getProfile(),
  getToken: () => authRepository.getToken(),

  // Парсинг токена и извлечение нужных данных
  parseToken: (token: string): DecodedToken => {
    const decoded: any = jwtDecode.jwtDecode(token);
    /*
    This expression is not callable.
  Type 'typeof import("/home/vasiliy/Projects.local/otus-hla-homework/frontend/node_modules/jwt-decode/build/esm/index")' has no call signatures.ts(2349)
    */
    return {
      userId: decoded.sub,
      tokenId: decoded.jti,
      permissions: decoded.permissions || [],
      exp: decoded.exp,
    };
  },

  // Настройка interceptor для добавления токена в заголовки
  setupAxiosInterceptor: (token: string) => {
    interceptorId = api.interceptors.request.use((config) => {
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    }, (error) => {
      return Promise.reject(error);
    });
  },

    // Удаление interceptor (например, при логауте)
  removeAxiosInterceptor: () => {
    if (interceptorId !== null) {
      api.interceptors.request.eject(interceptorId); // Передаем идентификатор
      interceptorId = null;  // Сбрасываем идентификатор после удаления
    }
  },
};
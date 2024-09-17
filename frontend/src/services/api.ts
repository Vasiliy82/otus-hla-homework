import axios from 'axios';
import { authService } from './authService';

// Создаем экземпляр axios
const api = axios.create({
  // TODO: надо бы его брать из ENV-ов при сборке
  baseURL: 'http://localhost:9090/api', // Адрес твоего API
});

// Интерфейсы для DTO

// DTO для регистрации
export interface RegisterUserDto {
  first_name: string;
  last_name: string;
  birthdate: string;
  biography: string;
  city: string;
  username: string;
  password: string;
}

// DTO для ответа при успешной регистрации
export interface RegisterUserResponse {
  user_id: string;
}

// DTO для получения пользователя
export interface UserDto {
  id: string;
  first_name: string;
  last_name: string;
  birthdate: string;
  biography: string;
  city: string;
  username: string;
  created_at: string;
  updated_at: string;
}

// DTO для логина
export interface LoginUserDto {
  username: string;
  password: string;
}

// DTO для успешного логина
export interface LoginResponse {
  user_id: string;
  token: string;
}

// DTO для ошибок
export interface ErrorResponse {
  error: string;
  code?: string;
}


// Добавляем перехватчик для запросов
api.interceptors.request.use(
  (config) => {
    const token = authService.getToken();  // Получаем токен через authService
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;  // Добавляем заголовок Authorization
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);  // Обрабатываем ошибки
  }
);

// API методы

// Регистрация пользователя
export const registerUser = (data: RegisterUserDto) => {
  return api.post<RegisterUserResponse>('/user/register', data);
};

// Получение пользователя
export const getUserById = (userId: string) => {
  return api.get<UserDto>(`/user/get/${userId}`);
};

// Логин пользователя
export const loginUser = (data: LoginUserDto) => {
  return api.post<LoginResponse>('/login', data);
};

export default api;
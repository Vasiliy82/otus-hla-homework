import React, { useState } from 'react';
import './authForm.css';  // Стили для формы авторизации
import { getUserById } from '../services/api';  // Импортируем API метод
import { authService } from '../services/authService';
import Modal, { useModal } from './Modal';  // Импортируем модальное окно для ошибок

const AuthForm: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const { isModalVisible, modalHeader, modalMessage, showModal, closeModal } = useModal();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      // Выполняем запрос на логин
      const response = await authService.loginUser({ username, password });
      const { token } = response.data;

      // Парсим токен и получаем userId и другие данные
      const { userId } = authService.parseToken(token);

      // Получаем профиль пользователя по userId
      const profileResponse = await getUserById(userId);
      const profile = profileResponse.data;

      // Логиним пользователя и сохраняем данные (токен и профиль)
      authService.login(token, profile);

      // Перезагрузить или перенаправить на другую страницу после авторизации
      window.location.reload();
    } catch (error: any) {
      if (!error.response) {
        // Ошибка сети (например, сервер недоступен)
        showModal('Ошибка подключения', 'Не удалось подключиться к серверу. Проверьте подключение к Интернету и повторите попытку.');
      } else {
        // Ошибка от сервера (REST API)
        const apiError = error.response.data.error
          ? error.response.data.error
          : 'Неизвестная ошибка';
        showModal('Ошибка API', `(${error.response.status}) ${apiError}`);
      }
    }
  };

  return (
    <div className="auth-block">
      <form onSubmit={handleLogin}>
        <input
          type="text"
          placeholder="Имя пользователя"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Пароль"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit">Авторизация</button>
      </form>
      <div className="auth-links">
        <a href="/register">Регистрация</a>
        <a href="/forgot-password">Восстановить пароль</a>
      </div>
      <Modal isVisible={isModalVisible} onClose={closeModal} header={modalHeader} message={modalMessage} />
    </div>
  );
};

export default AuthForm;
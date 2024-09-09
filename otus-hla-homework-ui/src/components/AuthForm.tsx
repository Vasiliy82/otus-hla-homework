import React, { useState } from 'react';
import './authForm.css';  // Стили для формы авторизации
import { loginUser, getUserById} from '../services/api';  // Импортируем API методы
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
      const response = await loginUser({ username, password });
      const { token } = response.data;

      const userId = response.data.user_id;
      // Если логин успешен, получаем профиль пользователя по user_id
      const profileResponse = await getUserById(userId);
      // На всякий случай сохраним его в репозиторий (хотя сейчас оттуда используется только логин для отображения в шапке)
      const profile = profileResponse.data;

      // Логиним пользователя и сохраняем данные (токен и профиль)
      // TODO: на самом деле это не правильно, т.к. будучи не залогиненными, мы getUserById не получим!
      // TODO: сначала надо успешно авторизоваться, затем ходить за инфой в профиль
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
          placeholder="Login"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit">Login</button>
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
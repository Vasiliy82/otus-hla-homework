import React from 'react';
import './layout.css';  // Стили для layout
import { authService } from '../services/authService';
import { Outlet, useNavigate } from 'react-router-dom';
import AuthForm from './AuthForm';  // Импортируем компонент формы авторизации

const Layout: React.FC = () => {
  const isAuthenticated = authService.isAuthenticated();
  const navigate = useNavigate();

  const handleLogout = () => {
    // Удаляем токен и профиль при выходе
    authService.logout();
    navigate('/');  // Перенаправляем на страницу логина после выхода
  };

  return (
    <div className="wrapper">
      <header className="header">
        <div className="header-content">
          <div className="logo">MySocialNetwork</div>
          <input type="text" placeholder="Search..." className="search-bar" />
          {isAuthenticated && <div>Welcome, {authService.getProfile()?.username}</div>}
        </div>
      </header>

      <div className="content">
        <aside className="left-column">
          {isAuthenticated ? (
            <>
              <nav className="menu">
                <ul>
                  <li>Мой профиль</li>
                  <li>Друзья</li>
                  <li>Сообщения</li>
                  <li>Настройки</li>
                </ul>
              </nav>
              <button className="logout-button" onClick={handleLogout}>
                Выход
              </button>
            </>
          ) : (
            <AuthForm />  /* Если не авторизован, показываем форму авторизации */
          )}
        </aside>

        <main className="main-content">
          <Outlet /> {/* Контент страниц */}
        </main>
      </div>
    </div>
  );
};

export default Layout;
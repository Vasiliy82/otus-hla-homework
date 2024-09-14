import React, { useState } from 'react';
import './modal.css'; // Стили для модального окна

type ModalProps = {
  isVisible: boolean;
  onClose: () => void;
  header: string;
  message: string;
};

const Modal: React.FC<ModalProps> = ({ isVisible, onClose, header, message }) => {
  if (!isVisible) return null;

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <h2>{header}</h2>
        <p>{message}</p>
        <button onClick={onClose}>Close</button>
      </div>
    </div>
  );
};

export default Modal;

// Хук для использования модального окна на разных страницах
export const useModal = () => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalHeader, setModalHeader] = useState('');
  const [modalMessage, setModalMessage] = useState('');

  const showModal = (header: string, message: string) => {
    setModalHeader(header);
    setModalMessage(message);
    setIsModalVisible(true);
  };

  const closeModal = () => {
    setIsModalVisible(false);
  };

  return {
    isModalVisible,
    modalHeader,
    modalMessage,
    showModal,
    closeModal,
  };
};
import React from "react";

const Modal = ({ isOpen, onClose, children }) => {
  if (!isOpen) return null;

  return (
    <div className="tp-modal-overlay">
      <div className="tp-modal">
        <button className="tp-modal-close-btn" onClick={onClose}>
          &times;
        </button>
        <div className="tp-modal-content">{children}</div>
      </div>
    </div>
  );
};

export default Modal;

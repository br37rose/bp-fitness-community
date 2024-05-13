import React, { useState } from "react";
import { useDrag } from "react-dnd";

function DraggableItem({ id, content, className, onAdd }) {
  const [{ isDragging }, dragRef] = useDrag({
    type: "item",
    item: { id: id },
  });

  const [showAddButton, setShowAddButton] = useState(false);

  const handleAddClick = () => {
    onAdd(id); // Call the onAdd function passed as prop
    setShowAddButton(false); // Hide the add button after adding
  };

  return (
    <div
      ref={dragRef}
      style={{
        opacity: isDragging ? 0.5 : 1,
        cursor: "pointer",
        width: "auto",
      }}
      className={className}
      onClick={(e) => {
        setShowAddButton(true);
      }} // Show add button on click
      onMouseLeave={() => setShowAddButton(false)} // Hide add button on mouse leave
      // onMouseEnter={() => setShowAddButton(true)}
    >
      {content}
      {showAddButton && (
        <button
          className="is-small button is-light is-success"
          onClick={handleAddClick}
        >
          Add
        </button> // Show add button if state is true
      )}
    </div>
  );
}

export default DraggableItem;

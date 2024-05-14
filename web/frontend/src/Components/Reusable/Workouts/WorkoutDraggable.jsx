import React, { useState } from "react";
import { useDrag } from "react-dnd";

function WorkoutDraggableItem({ id, exercise, content, onAdd, component: Component = "div" }) {
  const [{ isDragging }, dragRef] = useDrag({
    type: "item",
    item: { id },
  });

  const [showAddButton, setShowAddButton] = useState(false);

  const handleAddClick = (e) => {
    e.stopPropagation(); // Prevent triggering onClick of the row when clicking the button
    onAdd(exercise);
    setShowAddButton(false);
  };

  return (
    <Component ref={dragRef} onClick={() => setShowAddButton(true)} style={{ opacity: isDragging ? 0.5 : 1 }}>
      {content}
      {showAddButton && (
        <button
          className="button is-small is-success"
          onClick={handleAddClick}
          style={{ marginLeft: '10px' }}
        >
          Add
        </button>
      )}
    </Component>
  );
}

export default WorkoutDraggableItem;
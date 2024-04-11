import { faBars } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useEffect, useState } from "react";

const DragSortableList = ({ items, onSortChange }) => {
  const [list, setList] = useState(items);
  const [draggedItem, setDraggedItem] = useState(null);

  useEffect(() => {
    setList(items);
    return () => {};
  }, [items]);

  const handleDragEnter = (index) => {
    if (draggedItem === null) return;
    const newList = [...list];
    const draggedItemContent = newList[draggedItem];
    newList.splice(draggedItem, 1);
    newList.splice(index, 0, draggedItemContent);
    setList(newList); // Update the list state with the new order
    setDraggedItem(index);
    onSortChange(newList);
  };

  const handleDragEnd = () => {
    setDraggedItem(null);
  };

  const handleDragStart = (index, event) => {
    event.dataTransfer.setData("text/plain", ""); // Required for Firefox
    setDraggedItem(index);
  };
  return (
    <ul>
      {list.map((item, index) => (
        <li
          key={index}
          draggable
          onDragStart={(event) => handleDragStart(index, event)}
          onDragEnter={() => handleDragEnter(index)}
          onDragEnd={handleDragEnd}
          style={{
            transition: "transform 0.4s ease-in-out",
          }}
        >
          <div className="sortable-item">
            <div className="is-flex is-justify-content-space-between ml-2 mr-2">
              {item.content}
              <div style={{ cursor: "move" }}>
                <FontAwesomeIcon icon={faBars} color="grey" />
              </div>
            </div>
          </div>
        </li>
      ))}
    </ul>
  );
};

export default DragSortableList;

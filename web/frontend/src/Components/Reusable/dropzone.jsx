import { faFileImport } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useEffect } from "react";
import { useDrop } from "react-dnd";

function DropZone({ onDrop, className, children, placeholder }) {
  const [{ isOver }, dropRef] = useDrop({
    accept: "item",
    drop: (item) => {
      onDrop(item);
    },
  });

  return (
    <div
      ref={dropRef}
      style={{
        border: isOver ? "2px dashed #333" : "2px solid #ddd",
        minHeight: "15em",
        borderRadius: "8px",
      }}
      className={className}
    >
      {placeholder && (
        <div className="mt-6">
          <div className="mt-auto is-flex is-flex-direction-column is-align-items-center">
            <div>
              <FontAwesomeIcon
                icon={faFileImport}
                color="lightgrey"
                size="5x"
              />
            </div>
            <div className="has-text-primary">drop here</div>
          </div>
        </div>
      )}
      {children}
    </div>
  );
}

export default DropZone;

import React from "react";
import WorkoutDraggableItem from "./WorkoutDraggable";

function ExerciseDisplay({ exercises, onAdd }) {
  const TableRow = ({ exercise, index }) => (
    <WorkoutDraggableItem
      key={exercise.id}
      id={exercise.id}
      onAdd={onAdd}
      exercise={exercise}
      component="tr"
      content={(
        <>
          <td>{index + 1}</td>
          <td>{exercise.name}</td>
          <td><img src={exercise.thumbnailObjectUrl} alt={exercise.name} style={{ width: '100px', height: 'auto' }} /></td>
        </>
      )}
    />
  );

  return (
    <div className="table-container" style={{ maxHeight: "400px", overflowY: "auto" }}>
      <table className="table is-fullwidth is-hoverable is-striped">
        <thead>
          <tr>
            <th>#</th>
            <th>Name</th>
            <th>Video</th>
          </tr>
        </thead>
        <tbody>
          {exercises.map((exercise, index) => <TableRow key={exercise.id} exercise={exercise} index={index} />)}
        </tbody>
      </table>
    </div>
  );
}

export default ExerciseDisplay;
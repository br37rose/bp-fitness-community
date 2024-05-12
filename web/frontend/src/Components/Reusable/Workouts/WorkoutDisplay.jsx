import React from "react";
import WorkoutDraggableItem from "./WorkoutDraggable";

function WorkoutDisplay({ workouts, onAdd }) {
  const TableRow = ({ workout, index }) => (
    <WorkoutDraggableItem
      key={workout.id}
      id={workout.id}
      onAdd={onAdd}
      workout={workout}
      component="tr"
      content={(
        <>
          <td>{index + 1}</td>
          <td>{workout.name}</td>
          <td><img src={workout.thumbnailObjectUrl} alt={workout.name} style={{ width: '100px', height: 'auto' }} /></td>
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
          {workouts.map((workout, index) => <TableRow key={workout.id} workout={workout} index={index} />)}
        </tbody>
      </table>
    </div>
  );
}

export default WorkoutDisplay;
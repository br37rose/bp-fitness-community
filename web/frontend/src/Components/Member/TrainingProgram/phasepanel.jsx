import { useState, useEffect } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCaretUp,
  faCaretDown,
  faEdit,
} from "@fortawesome/free-solid-svg-icons";

function PhasePanel({ phases, onAddWorkout, setSelectedPhase }) {
  const [activeIndex, setActiveIndex] = useState(0); // Initialize activeIndex with 0

  useEffect(() => {
    if (phases?.length > 0 && setSelectedPhase) {
      setSelectedPhase(phases[0]); // Set the selected phase to the first phase
    }
  }, [phases]); // Update selected phase when phases change

  const toggleAccordion = (index) => {
    setActiveIndex(index === activeIndex ? null : index);
    if (setSelectedPhase) {
      setSelectedPhase(phases[index]); // Update selected phase when accordion is toggled
    }
  };

  return (
    <nav className="box">
      <div className="columns is-vcentered">
        <div className="column">
          <p className="title is-4">Phases</p>
        </div>
      </div>
      {phases?.map((phase, index) => (
        <div key={index} className="accordion">
          <div
            className={`accordion-header ${
              activeIndex === index ? "is-active" : ""
            }`}
            onClick={() => toggleAccordion(index)}
          >
            <p>{phase.name}</p>
            <span className="icon">
              {activeIndex === index ? (
                <FontAwesomeIcon icon={faCaretUp} />
              ) : (
                <FontAwesomeIcon icon={faCaretDown} />
              )}
            </span>
          </div>
          {/* <div
            className={`accordion-body ${
              activeIndex === index ? "is-active" : ""
            }`}
          >
            <button
              className="button is-success is-light is-small"
              onClick={() => onAddWorkout(phase)}
            >
              <FontAwesomeIcon icon={faEdit} className="mr-1" />
              Workouts
            </button>
            Render other phase details here
          </div> */}
        </div>
      ))}
    </nav>
  );
}

export default PhasePanel;

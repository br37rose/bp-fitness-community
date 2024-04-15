import { useState } from "react";
import PropTypes from "prop-types";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCaretDown, faCaretUp } from "@fortawesome/free-solid-svg-icons";

const Accordion = ({ head, content, isOpenByDefault }) => {
  const [isOpen, setIsOpen] = useState(isOpenByDefault);

  const toggleAccordion = () => {
    setIsOpen(!isOpen);
  };

  return (
    <div className="accordion">
      <div className="accordion-header" onClick={toggleAccordion}>
        {head}
        <span className="accordion-icon">
          {isOpen ? (
            <FontAwesomeIcon icon={faCaretUp} />
          ) : (
            <FontAwesomeIcon icon={faCaretDown} />
          )}
        </span>
      </div>
      {isOpen && <div className="accordion-content">{content}</div>}
    </div>
  );
};

Accordion.propTypes = {
  head: PropTypes.node.isRequired,
  content: PropTypes.node.isRequired,
  isOpenByDefault: PropTypes.bool,
};

Accordion.defaultProps = {
  isOpenByDefault: false,
};

export default Accordion;

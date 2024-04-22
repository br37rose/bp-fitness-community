import React, { useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeftLong,
  faArrowRight,
  faArrowRightLong,
  faCheck,
} from "@fortawesome/free-solid-svg-icons";
import { faBars } from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { onHamburgerClickedState, currentUserState } from "../../../AppState";

const Title = ({ text, subtitle, quizTime }) => (
  <div>
    <h1 className="title is-size-2 has-text-centered has-text-white">{text}</h1>
    {subtitle && (
      <h2 className="subtitle is-size-3 has-text-centered has-text-white">
        {subtitle}
      </h2>
    )}
    <p className="is-size-7 has-text-centered has-text-white pb-3">
      {quizTime}
    </p>
  </div>
);

const Card = ({ imgUrl, id, selected, buttonText, card, onSelect }) => (
  <div className="column">
    <div
      className={`card ${selected ? "is-selected" : ""}`}
      style={{
        maxWidth: "300px",
        margin: "auto",
        borderRadius: "20px",
        overflow: "hidden",
        backgroundColor: selected ? "#3273dc" : "#4a4a4a", // Example: change background color when selected
        cursor: "pointer",
        boxShadow: selected ? "0 0 0 2px #3273dc" : "none", // Adding a box shadow instead of border for a "selected" effect
      }}
      onClick={() => onSelect(id)}
    >
      <div className="card-image">
        <figure className="image is-4by3">
          <img
            src={imgUrl}
            alt="User"
            style={{ objectFit: "cover", height: "100%" }}
          />
        </figure>
      </div>
      <div className="button is-info is-fullwidth">
        {buttonText}:&nbsp;{card.type === "age" ? card.ageRange : ""}&nbsp;
        <span className="icon">
          &nbsp;
          <FontAwesomeIcon icon={faArrowRight} />
        </span>
      </div>
    </div>
  </div>
);

const SelectableOption = ({ option, isSelected, onSelect }) => (
  <div className="field">
    <button
      id={option}
      type="button"
      name={option}
      aria-pressed={isSelected}
      onClick={() => onSelect(option)}
      className={`button is-large is-fullwidth ${
        isSelected ? "is-primary" : "is-dark"
      }`}
    >
      <span>{option}</span>
      {isSelected && (
        <span className="icon is-small">
          <FontAwesomeIcon icon={faCheck} />
        </span>
      )}
    </button>
  </div>
);

export { Title, Card, SelectableOption };

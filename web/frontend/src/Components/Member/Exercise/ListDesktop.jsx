import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEye, faFire, faDroplet, faGenderless } from '@fortawesome/free-solid-svg-icons';
import {
    EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER,
    EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL
}  from "../../../Constants/App";

function MemberExerciseListDesktop({ description,
  name,
  thumbnailUrl,
  thumbnailObjectUrl,
  thumbnailType,
  gender,
  id }) {

  return (
    <div>
      <div className="columns is-flex is-justify-content-space-between pt-5 pb-4 border-bottom">
        {/* <img className="is-radiusless image is-1by3" src={thumbnailUrl} alt={name} /> */}
        <figure class="image is-256x256" >
          {thumbnailType === EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER && <img src={thumbnailObjectUrl} alt={name} style={{ borderRadius: "10px" }} />}
          {thumbnailType === EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL && <img src={thumbnailUrl} alt={name} style={{ borderRadius: "10px" }} />}
        </figure>
        <div className="column is-7 is-flex is-flex-wrap-wrap is-align-content-space-between ">
          <div className="">
            <h4 className="is-size-5 has-text-weight-bold mb-2">{name}</h4>
            <p>{description}</p>
            {/* <p>1. Hook your heels into the wall with your hamstrings. 2. Tuck your tailbone between your knees. 3. Push your lower back into the floor. 4. Exhale fully, dropping your ribs down as far as they'll go. 5. Inhale without losing the ribs-down position and repeat.</p> */}
          </div>
          <div className="is-flex-desktop mt-4 is-justify-content-end">
            <span className="is-flex is-align-items-center mr-5 mt-m-3">

              <h5 className="has-text-weight-semibold is-size-6"><FontAwesomeIcon className="fas" icon={faFire} />&nbsp;Ground-Based Exercises</h5>
            </span>

            <span className="is-flex is-align-items-center mr-5 mt-m-3">

              <h5 className="has-text-weight-semibold is-size-6"><FontAwesomeIcon className="fas" icon={faDroplet} />&nbsp;Warmups & Mobility Fillers</h5>
            </span>

            <span className="is-flex is-align-items-center mt-m-3">
              <h5 className="has-text-weight-semibold is-size-6"><FontAwesomeIcon className="fas" icon={faGenderless} />&nbsp;{gender}</h5>
            </span>
          </div>

        </div>
        <div className="column is-2 has-text-right">
          <Link
            to={`/exercise/${id}`}
            className="button"
            type="button"
          >
            <FontAwesomeIcon
              className="mdi"
              icon={faEye}
            />
            &nbsp;View
          </Link>
        </div>
      </div>
      <hr />
    </div>
  );
}

export default MemberExerciseListDesktop;

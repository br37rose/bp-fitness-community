import React, { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faMessage,
  faHeartbeat,
  faRankingStar,
  faStar,
  faChartLine,
  faLeaf,
  faTrophy,
  faHandHolding,
  faVideoCamera,
  faDumbbell,
  faUsers,
  faBuilding,
  faTachometer,
  faUserCircle,
  faSignOut,
  faCalendarPlus,
  faQuestionCircle,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { onHamburgerClickedState, currentUserState } from "../../AppState";
import { faWhatsapp } from "@fortawesome/free-brands-svg-icons";

export default (props) => {
  ////
  //// Global State
  ////
  const [onHamburgerClicked, setOnHamburgerClicked] = useRecoilState(
    onHamburgerClickedState
  );
  const [currentUser] = useRecoilState(currentUserState);

  ////
  //// Local State
  ////

  const [showLogoutWarning, setShowLogoutWarning] = useState(false);

  ////
  //// Events
  ////

  // Do nothing.

  ////
  //// Rendering.
  ////

  //-------------//
  // CASE 1 OF 3 //
  //-------------//

  // Get the current location and if we are at specific URL paths then we
  // will not render this component.
  const ignorePathsArr = [
    "/",
    "/register",
    "/register-step-1",
    "/register-step-2",
    "/register-successful",
    "/index",
    "/login",
    "/login/2fa",
    "/login/2fa/step-1",
    "/login/2fa/step-2",
    "/login/2fa/step-3",
    "/logout",
    "/verify",
    "/forgot-password",
    "/password-reset",
    "/terms",
    "/privacy",
  ];
  const location = useLocation();
  var arrayLength = ignorePathsArr.length;
  for (var i = 0; i < arrayLength; i++) {
    // console.log(location.pathname, "===", ignorePathsArr[i], " EQUALS ", location.pathname === ignorePathsArr[i]); // For debugging purposes only.
    if (location.pathname === ignorePathsArr[i]) {
      return null;
    }
  }

  //-------------//
  // CASE 2 OF 3 //
  //-------------//

  if (currentUser === null) {
    console.log("No current user detected, hiding siedbard menu.");
    return null;
  }

  //-------------//
  //   Whatsapp  //
  //-------------//
  const WhatsAppGroupButton = () => {
    const whatsappGroupUrl = "https://chat.whatsapp.com/FdVtOFCoUN48PUy5E4AcKe";

    const handleClick = () => {
      // Open WhatsApp Group link in a new tab
      window.open(whatsappGroupUrl, "_blank", "noopener,noreferrer");
    };
    return (
      <Link className="has-text-success is-size-6" onClick={handleClick}>
        <FontAwesomeIcon icon={faWhatsapp} />
        &nbsp;BP8 Community
      </Link>
    );
  };

  return (
    <>
      <div class={`modal ${showLogoutWarning ? "is-active" : ""}`}>
        <div class="modal-background"></div>
        <div class="modal-card">
          <header class="modal-card-head">
            <p class="is-size-4 modal-card-title m-0">Are you sure?</p>
            <button
              class="delete"
              aria-label="close"
              onClick={(e) => setShowLogoutWarning(false)}
            ></button>
          </header>
          <p class="modal-card-body">
            You are about to log out of the system and you'll need to log in
            again next time. Are you sure you want to continue?
          </p>
          <footer class="modal-card-foot">
            <Link class="button is-success" to={`/logout`}>
              Yes
            </Link>
            <button class="button" onClick={(e) => setShowLogoutWarning(false)}>
              No
            </button>
          </footer>
        </div>
      </div>
      {/*
                ---------------------
                ADMINISTRATOR (ROOT)
                ---------------------
            */}
      {currentUser.role === 1 && (
        <>
          <p>Not implemeneted yet</p>
        </>
      )}

      {/*
                ---------------------
                ADMINISTRATOR (STAFF)
                ---------------------
            */}
      {(currentUser.role === 2 || currentUser.role === 3) && (
        <div
          className={`column is-one-fifth has-background-black ${
            onHamburgerClicked ? "" : "is-hidden"
          }`}
        >
          <aside class="menu p-4">
            <p class="menu-label has-text-grey-light">Menu</p>
            <ul class="menu-list">
              <li>
                <a
                  href="/admin/dashboard"
                  class={`has-text-grey-light ${
                    location.pathname.includes("dashboard") &&
                    !location.pathname.includes("trainer") &&
                    "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faTachometer} />
                  &nbsp;Dashboard
                </a>
              </li>
              <li>
                <a
                  href="/admin/members"
                  class={`has-text-grey-light ${
                    location.pathname.includes("member") && "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;Members
                </a>
              </li>
              <li>
                <a
                  href="/admin/exercises"
                  class={`has-text-grey-light ${
                    location.pathname.includes("exercises") && "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faDumbbell} />
                  &nbsp;Exercises
                </a>
              </li>
              <li>
                <a
                  href="/admin/video-collections"
                  class={`has-text-grey-light ${
                    location.pathname.includes("video-collections") &&
                    "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faVideoCamera} />
                  &nbsp;Video Collections
                </a>
              </li>
              <li>
                <a
                  href="/admin/workouts"
                  class={`has-text-grey-light ${
                    location.pathname.includes("workouts") && "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faDumbbell} />
                  &nbsp;Workouts
                </a>
              </li>
              <li>
                <a
                  href="/admin/training-program"
                  class={`has-text-grey-light ${
                    location.pathname.includes("training-program") &&
                    "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faCalendarPlus} />
                  &nbsp;Training Program
                </a>
              </li>
              {/*
                                Videos, Exercises, Workouts, Programs, Fitness Plans, Nutrition Plans, Social Media Feed
                            */}
            </ul>

            {currentUser.role === 2 && (
              <>
                <p class="menu-label has-text-grey-light">System</p>
                <ul class="menu-list">
                  <li>
                    <a
                      href={`/admin/organization`}
                      class={`has-text-grey-light ${
                        location.pathname.includes("organization") &&
                        "is-active"
                      }`}
                    >
                      <FontAwesomeIcon className="fas" icon={faBuilding} />
                      &nbsp;Organization
                    </a>
                  </li>
                  <li>
                    <a
                      href="/admin/video-categories"
                      class={`has-text-grey-light ${
                        location.pathname.includes("video-categor") &&
                        "is-active"
                      }`}
                    >
                      <FontAwesomeIcon className="fas" icon={faDumbbell} />
                      &nbsp;Video Categories
                    </a>
                  </li>
                  <li>
                    <a
                      href={`/admin/offers`}
                      class={`has-text-grey-light ${
                        location.pathname.includes("offer") && "is-active"
                      }`}
                    >
                      <FontAwesomeIcon className="fas" icon={faHandHolding} />
                      &nbsp;Offers
                    </a>
                  </li>
                  <li>
                    <a
                      href={`/admin/questions`}
                      class={`has-text-grey-light ${
                        location.pathname.includes("questions") && "is-active"
                      }`}
                    >
                      <FontAwesomeIcon
                        className="fas"
                        icon={faQuestionCircle}
                      />
                      &nbsp;Onboarding Questions
                    </a>
                  </li>
                </ul>
              </>
            )}

            <p class="menu-label has-text-grey-light">Account</p>
            <ul class="menu-list">
              <li>
                <a
                  href={`/account`}
                  class={`has-text-grey-light ${
                    location.pathname.includes("account") && "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faUserCircle} />
                  &nbsp;Account
                </a>
              </li>

              <li>
                <a
                  onClick={(e) => setShowLogoutWarning(true)}
                  class={`has-text-grey-light ${
                    location.pathname.includes("logout") && "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faSignOut} />
                  &nbsp;Sign Off
                </a>
              </li>
            </ul>
            <p class="menu-label has-text-grey-light">Connect with us on</p>
            <WhatsAppGroupButton />
          </aside>
        </div>
      )}
      {/*
                ---------------------
                MEMBER (REGULAR USERS)
                ---------------------
            */}
      {currentUser.role === 4 && (
        <div
          className={`column is-one-fifth has-background-black ${
            onHamburgerClicked ? "" : "is-hidden"
          }`}
        >
          <aside class="menu p-4">
            <p class="menu-label has-text-grey-light">Menu</p>
            <ul class="menu-list">
              <li>
                <a
                  href="/dashboard"
                  class={`has-text-grey-light ${
                    location.pathname.includes("dashboard") && "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faTachometer} />
                  &nbsp;Dashboard
                </a>
              </li>
              <li>
                <a
                  href="/exercises"
                  class={`has-text-grey-light ${
                    location.pathname.includes("exercises") && "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faDumbbell} />
                  &nbsp;Exercises
                </a>
              </li>
              <li>
                <a
                  href="/video-categories"
                  class={`has-text-grey-light ${
                    location.pathname.includes("video-categories") &&
                    "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faVideoCamera} />
                  &nbsp;Videos
                </a>
              </li>
              <li>
                <a
                  href="/fitness-plans"
                  class={`has-text-grey-light ${
                    location.pathname.includes("fitness-plan") && "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faTrophy} />
                  &nbsp;Fitness Plan
                </a>
              </li>
              <li>
                <a
                  href="/nutrition-plans"
                  class={`has-text-grey-light ${
                    location.pathname.includes("nutrition-plan") && "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faLeaf} />
                  &nbsp;Nutrition Plan
                </a>
              </li>
              <li>
                <a
                  href="/workouts"
                  class={`has-text-grey-light ${
                    location.pathname.includes("workouts") && "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faDumbbell} />
                  &nbsp;Workouts
                </a>
              </li>
              <li>
                <a
                  href="/training-program"
                  class={`has-text-grey-light ${
                    location.pathname.includes("training-program") &&
                    "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faCalendarPlus} />
                  &nbsp;Training Program
                </a>
              </li>
              <li>
                <a
                  href="/biometrics"
                  class={`has-text-grey-light ${
                    location.pathname.includes("/biometrics") &&
                    !location.pathname.includes("/leaderboard") &&
                    !location.pathname.includes("/summary") &&
                    !location.pathname.includes("/history") &&
                    "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faHeartbeat} />
                  &nbsp;Biometrics
                </a>
                <ul>
                  <li>
                    <a
                      href="/biometrics/leaderboard/global"
                      class={`has-text-grey-light ${
                        location.pathname.includes("/biometrics/leaderboard") &&
                        "is-active"
                      }`}
                    >
                      <FontAwesomeIcon className="fas" icon={faRankingStar} />
                      &nbsp;Leaderboard
                    </a>
                  </li>
                  <li>
                    <a
                      href="/biometrics/summary"
                      class={`has-text-grey-light ${
                        location.pathname.includes("/biometrics/summary") &&
                        "is-active"
                      }`}
                    >
                      <FontAwesomeIcon className="fas" icon={faStar} />
                      &nbsp;My Summary
                    </a>
                  </li>
                  <li>
                    <a
                      href="/biometrics/history/tableview"
                      class={`has-text-grey-light ${
                        location.pathname.includes("/biometrics/history/") &&
                        "is-active"
                      }`}
                    >
                      <FontAwesomeIcon className="fas" icon={faChartLine} />
                      &nbsp;My History
                    </a>
                  </li>
                </ul>
              </li>
              {/*
                                Videos, Exercises, Workouts, Programs, Fitness Plans, Nutrition Plans, Social Media Feed
                            */}
            </ul>
            <p class="menu-label has-text-grey-light">Account</p>
            <ul class="menu-list">
              <li>
                <a
                  href={`/account`}
                  class={`has-text-grey-light ${
                    location.pathname.includes("account") && "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faUserCircle} />
                  &nbsp;Account
                </a>
              </li>
              <li>
                <a
                  onClick={(e) => setShowLogoutWarning(true)}
                  class={`has-text-grey-light ${
                    location.pathname.includes("logout") && "is-active"
                  }`}
                >
                  <FontAwesomeIcon className="fas" icon={faSignOut} />
                  &nbsp;Sign Off
                </a>
              </li>
            </ul>
            <p class="menu-label has-text-grey-light">Connect with us on</p>
            <WhatsAppGroupButton />
          </aside>
        </div>
      )}
    </>
  );
};

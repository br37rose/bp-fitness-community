import React, { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBars,
  faRightFromBracket,
  faTachometer,
  faTasks,
  faSignOut,
  faUserCircle,
  faUsers,
  faBuilding,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { onHamburgerClickedState, currentUserState } from "../../AppState";

function Topbar() {
  ////
  //// Global State
  ////
  const [onHamburgerClicked, setOnHamburgerClicked] = useRecoilState(
    onHamburgerClickedState,
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
    "/terms-of-service",
    "/privacy-policy",
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
    console.log("No current user detected, hiding sidebar menu.");
    return null;
  }

  //-------------//
  // CASE 3 OF 3 //
  //-------------//

  // Render the following component GUI
  return (
    <div>
      <nav
        class="navbar has-background-black"
        role="navigation"
        aria-label="main navigation"
      >
        <div class="navbar-brand">
          {currentUser && (
            <>
              {currentUser.role === 1 && <>TODO - SaaSify</>}
              {currentUser.role === 2 && (
                <a
                  class="navbar-item"
                  href="/admin/dashboard"
                  style={{ color: "white" }}
                >
                  <img
                    src="/static/logo.png"
                    style={{
                      maxWidth: "8rem",
                      maxHeight: "3rem",
                      width: "100%", // This ensures the image scales within the given maxWidth and maxHeight
                      height: "auto", // This maintains the aspect ratio of the image
                    }}
                    alt="BP8_Fitness_Logo"
                  />
                </a>
              )}
              {currentUser.role === 3 && (
                <a
                  class="navbar-item"
                  href="/dashboard"
                  style={{ color: "white" }}
                >
                  <img
                    src="/static/logo.png"
                    style={{
                      maxWidth: "8rem",
                      maxHeight: "3rem",
                      width: "100%", // This ensures the image scales within the given maxWidth and maxHeight
                      height: "auto", // This maintains the aspect ratio of the image
                    }}
                    alt="BP8_Fitness_Logo"
                  />
                </a>
              )}
              {currentUser.role === 4 && (
                <a
                  class="navbar-item"
                  href="/dashboard"
                  style={{ color: "white" }}
                >
                  <img
                    src="/static/logo.png"
                    style={{
                      maxWidth: "8rem",
                      maxHeight: "3rem",
                      width: "100%", // This ensures the image scales within the given maxWidth and maxHeight
                      height: "auto", // This maintains the aspect ratio of the image
                    }}
                    alt="BP8_Fitness_Logo"
                  />
                </a>
              )}
            </>
          )}
          <a
            role="button"
            class="navbar-burger has-text-white"
            aria-label="menu"
            aria-expanded="false"
            data-target="navbarBasicExample"
            onClick={(e) => setOnHamburgerClicked(!onHamburgerClicked)}
          >
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
          </a>
        </div>
        <div id="navbarBasicExample" class="navbar-menu has-text-white">
          <div class="navbar-left">
            <div class="navbar-item">
              <div
                class="buttons p-3"
                onClick={(e) => setOnHamburgerClicked(!onHamburgerClicked)}
              >
                <FontAwesomeIcon className="fas has-text-white" icon={faBars} />
              </div>
            </div>
          </div>
        </div>
      </nav>
      <div class={`modal ${showLogoutWarning ? "is-active" : ""}`}>
        <div class="modal-background"></div>
        <div class="modal-card">
          <header class="modal-card-head">
            <p class="modal-card-title">Are you sure?</p>
            <button
              class="delete"
              aria-label="close"
              onClick={(e) => setShowLogoutWarning(false)}
            ></button>
          </header>
          <section class="modal-card-body">
            You are about to log out of the system and you'll need to log in
            again next time. Are you sure you want to continue?
          </section>
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
    </div>
  );
}

export default Topbar;

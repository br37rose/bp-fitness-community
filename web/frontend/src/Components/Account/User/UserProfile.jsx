import { useState, useEffect } from "react";
import { Link, Navigate, useLocation } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faEllipsis,
  faArrowLeft,
  faUserCircle,
  faGauge,
  faKey,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { getAccountDetailAPI } from "../../../API/Account";
import FormErrorBox from "../../Reusable/FormErrorBox";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import Layout from "../../Menu/Layout";
import UserDetail from "./UserDetail";
import UserInfo from "./UserInfo";
import ActivityStream from "./ActivityStream";
import AccountTagList from "../DetailForTags";
import AccountFriendList from "../Friend/Friend";
import AccountWearableTechLaunchpad from "../WearableTech/Launchpad";
import AccountSubscriptionDetailAndCancel from "../Subscription/Subscription";
import AccountMoreLaunchpad from "../More/Launchpad";
import Survey from "../Survey";

function UserProfile(props) {
  ////
  ////
  ////

  const location = useLocation();
  const activeTabProp = location.state?.activeTabProp || "detail";
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);
  const [activeTab, setActiveTab] = useState(activeTabProp);

  ////
  //// Event handling.
  ////

  //

  ////
  //// API.
  ////

  function onAccountDetailSuccess(response) {
    console.log("onAccountDetailSuccess: Starting...");
    setCurrentUser(response);
  }

  function onAccountDetailError(apiErr) {
    console.log("onAccountDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onAccountDetailDone() {
    console.log("onAccountDetailDone: Starting...");
    setFetching(false);
  }

  ////
  //// BREADCRUMB
  ////

  const generateBreadcrumbItemLink = (currentUser) => {
    let dashboardLink;
    switch (currentUser.role) {
      case 1:
        dashboardLink = "/root/dashboard";
        break;
      case 2:
        dashboardLink = "/admin/dashboard";
        break;
      case 3:
        dashboardLink = "/trainer/dashboard";
        break;
      case 4:
        dashboardLink = "/dashboard";
        break;
      default:
        dashboardLink = "/"; // Default or error handling
        break;
    }
    return dashboardLink;
  };

  const breadcrumbItems = {
    items: [
      {
        text: "Dashboard",
        link: generateBreadcrumbItemLink(currentUser),
        isActive: false,
        icon: faGauge,
      },
      { text: "Account", link: "/account", icon: faUserCircle, isActive: true },
    ],
    mobileBackLinkItems: {
      link: generateBreadcrumbItemLink(currentUser),
      text: "Back to Dashboard",
      icon: faArrowLeft,
    },
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      setFetching(true);
      setErrors({});
      getAccountDetailAPI(
        onAccountDetailSuccess,
        onAccountDetailError,
        onAccountDetailDone
      );
    }

    return () => {
      mounted = false;
    };
  }, []);

  const renderTabContent = () => {
    if (!activeTab) {
      // Handle the case where activeTab is null or undefined
      return <UserInfo {...currentUser} />;
    }

    switch (activeTab) {
      case "detail":
        return <UserInfo {...currentUser} />;
      case "tags":
        return <AccountTagList />;
      case "friends":
        return <AccountFriendList />;
      case "wearableTech":
        return <AccountWearableTechLaunchpad />;
      case "survey":
        return <Survey />;
      case "subscription":
        return <AccountSubscriptionDetailAndCancel />;
      case "more":
        return <AccountMoreLaunchpad />;
      default:
        return <UserInfo {...currentUser} />;
    }
  };

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <Layout breadcrumbItems={breadcrumbItems}>
      <div class="box">
        {/* Title + Options */}
        {currentUser && (
          <div class="columns">
            <div class="column">
              <p class="title is-4">
                <FontAwesomeIcon className="fas" icon={faUserCircle} />
                &nbsp;Account
              </p>
            </div>
            <div class="column has-text-right-tablet has-text-right-desktop">
              <Link
                class="button is-success is-hidden-touch"
                to={"/account/change-password"}
              >
                <FontAwesomeIcon className="mdi" icon={faKey} />
                &nbsp;Change Password
              </Link>
              <Link
                to={"/account/change-password"}
                class="button is-success is-fullwidth is-hidden-desktop"
              >
                <FontAwesomeIcon className="mdi" icon={faKey} />
                &nbsp;Change Password
              </Link>
            </div>
          </div>
        )}
        <FormErrorBox errors={errors} />

        {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

        {isFetching ? (
          <PageLoadingContent displayMessage={"Please wait..."} />
        ) : (
          <>
            {currentUser && (
              <div>
                <div className="columns is-mobile is-multiline">
                  {/* Left side box items */}
                  <div className="column column is-one-quarter-desktop is-full-tablet is-full-mobile">
                    <div className="mb-5">
                      <UserDetail {...currentUser} />
                    </div>
                    <ActivityStream {...currentUser} />
                  </div>

                  {/* Right side box items */}
                  <div className="column column is-three-quarters-desktop is-full-tablet is-full-mobile">
                    {/* Tab Navigation */}
                    <div className="box">
                      <div class="tabs is-centered is-boxed">
                        <ul>
                          <li
                            className={
                              activeTab === "detail" ? "is-active" : ""
                            }
                          >
                            <a onClick={() => setActiveTab("detail")}>
                              <strong>Detail</strong>
                            </a>
                          </li>
                          <li
                            className={activeTab === "tags" ? "is-active" : ""}
                          >
                            <a onClick={() => setActiveTab("tags")}>
                              <strong>Tags</strong>
                            </a>
                          </li>
                          <li
                            className={
                              activeTab === "friends" ? "is-active" : ""
                            }
                          >
                            <a onClick={() => setActiveTab("friends")}>
                              <strong>Friends</strong>
                            </a>
                          </li>
                          <li
                            className={
                              activeTab === "wearableTech" ? "is-active" : ""
                            }
                          >
                            <a onClick={() => setActiveTab("wearableTech")}>
                              <strong>Wearable Tech</strong>
                            </a>
                          </li>
                          <li
                            className={
                              activeTab === "survey" ? "is-active" : ""
                            }
                          >
                            <a onClick={() => setActiveTab("survey")}>
                              <strong>Survey</strong>
                            </a>
                          </li>
                          <li
                            className={
                              activeTab === "subscription" ? "is-active" : ""
                            }
                          >
                            <a onClick={() => setActiveTab("subscription")}>
                              <strong>Subscription</strong>
                            </a>
                          </li>
                          <li
                            className={activeTab === "more" ? "is-active" : ""}
                          >
                            <a onClick={() => setActiveTab("more")}>
                              <strong>
                                More&nbsp;
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faEllipsis}
                                />
                              </strong>
                            </a>
                          </li>
                        </ul>
                      </div>
                      <div>{renderTabContent()}</div>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </Layout>
  );
}

export default UserProfile;

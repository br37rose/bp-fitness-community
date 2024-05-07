import React, {useState, useEffect} from "react";
import {Link, useNavigate, useSearchParams} from "react-router-dom";
import Scroll from "react-scroll";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
	faArrowUpRightFromSquare,
	faHeartPulse,
	faArrowRight,
	faArrowLeft,
} from "@fortawesome/free-solid-svg-icons";
import {useRecoilState} from "recoil";

import RedirectURL from "../../../Hooks/RedirectURL";
import {getAccountDetailAPI} from "../../../API/Account";
import {
	getGoogleFitRegistrationURLAPI,
	postFitBitAppCreateSimulatorAPI,
} from "../../../API/Wearable";
import {
	topAlertMessageState,
	topAlertStatusState,
	currentUserState,
} from "../../../AppState";

function AccountWearableTechLaunchpad() {
	////
	//// URL Parameters.
	////
	const navigate = useNavigate();

	const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140

	// DEVELOPERS NOTE:
	// This url parameter is set to `true` from the backend when the oAuth2.0
	// authorization was successfull between BP8 Fitness Community and Google.
	// Use this variable to notify the user they have successfully registered
	// their Google Fit with us.
	const wasGoogleFitRegistered = searchParams.get("google_fit_registered");

	////
	//// Global state.
	////

	const [topAlertMessage, setTopAlertMessage] =
		useRecoilState(topAlertMessageState);
	const [topAlertStatus, setTopAlertStatus] =
		useRecoilState(topAlertStatusState);
	const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

	////
	//// Component states.
	////

	const [errors, setErrors] = useState({});
	const [isFetching, setFetching] = useState(false);
	const [forceURL, setForceURL] = useState("");

	////
	//// Event handling.
	////
	const handleNavigateToAccount = () => {
		// e.preventDefault();
		navigate("/account", {state: {activeTabProp: "wearableTech"}});
	};

	useEffect(() => {
		console.log("calling handleNavigateToAccount");
		handleNavigateToAccount();
	}, []);

	const onRegisterClick = (e) => {
		e.preventDefault();
		console.log("onRegisterClick: Clicked");
		setFetching(true);
		setErrors({});
		getGoogleFitRegistrationURLAPI(
			onRegistrationSuccess,
			onRegistrationError,
			onRegistrationDone
		);
	};

	const onCreateSimulator = (e) => {
		setFetching(true);
		setErrors({});
		postFitBitAppCreateSimulatorAPI(
			currentUser.id,
			"random",
			onCreateSimulatorSuccess,
			onCreateSimulatorError,
			onCreateSimulatorDone
		);
	};

	////
	//// API.
	////

	// --- Detail --- //

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

	// --- Simulator --- //

	function onCreateSimulatorSuccess(response) {
		console.log("onCreateSimulatorSuccess: Starting...");
		window.location.reload();
	}

	function onCreateSimulatorError(apiErr) {
		console.log("onCreateSimulatorError: Starting...");
		setErrors(apiErr);

		// The following code will cause the screen to scroll to the top of
		// the page. Please see ``react-scroll`` for more information:
		// https://github.com/fisshy/react-scroll
		var scroll = Scroll.animateScroll;
		scroll.scrollToTop();
	}

	function onCreateSimulatorDone() {
		console.log("onCreateSimulatorDone: Starting...");
		setFetching(false);
	}

	// --- Register --- //

	function onRegistrationSuccess(response) {
		console.log("onRegistrationSuccess: Starting...");
		setForceURL(response.url);
	}

	function onRegistrationError(apiErr) {
		console.log("onRegistrationError: Starting...");
		setErrors(apiErr);

		// The following code will cause the screen to scroll to the top of
		// the page. Please see ``react-scroll`` for more information:
		// https://github.com/fisshy/react-scroll
		var scroll = Scroll.animateScroll;
		scroll.scrollToTop();
	}

	function onRegistrationDone() {
		console.log("onRegistrationDone: Starting...");
		setFetching(false);
	}

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

	////
	//// Component rendering.
	////

	if (forceURL !== "") {
		return <RedirectURL url={forceURL} />;
	}

	return (
		<div>
			{/*
          DEVELOPERS NOTE:
          THIS IS IMPORTANT AS THE BACKEND SPECIFIES WHETHER THE REGISTRATION
          WAS SUCCESSFUL OR NOT AND WE NEED TO LET THE USER KNOW THE STATUS.
      */}
			{wasGoogleFitRegistered !== undefined &&
				wasGoogleFitRegistered !== null &&
				wasGoogleFitRegistered !== "" && (
					<>
						{wasGoogleFitRegistered === "true" ? (
							<>
								<article class="message is-success">
									<div class="message-body">
										You have successfully registered your{" "}
										<strong>Google Fit</strong> with us!
									</div>
								</article>
							</>
						) : (
							<>
								<article class="message is-danger">
									<div class="message-body">
										Registered <strong>Google Fit</strong> was unsuccessfuly,
										please try again.
									</div>
								</article>
							</>
						)}
					</>
				)}

			<div className="columns">
				<div className="column">
					{/* Subtitle */}
					<p class="title is-6">
						<FontAwesomeIcon className="fas" icon={faHeartPulse} />
						&nbsp;Google Fit Fitness Tracker
					</p>
					<hr />

					{/* Empty list */}
					{currentUser.primaryHealthTrackingDeviceType === 0 && (
						<section className="hero has-background-white-ter">
							<div className="hero-body">
								<p className="title">
									<FontAwesomeIcon className="fas" icon={faHeartPulse} />
									&nbsp;No Connection
								</p>
								<p className="subtitle">
									Your Google Fit fitness tracker is not connected with us.{" "}
									<b>
										<Link onClick={onRegisterClick}>
											Click here&nbsp;
											<FontAwesomeIcon
												className="mdi"
												icon={faArrowUpRightFromSquare}
											/>{" "}
										</Link>
									</b>{" "}
									to get started by registering your device and let us read the
									latest your biometrics. We currently we will extract the
									following data from your device:
									<div className="content">
										<ul>
											<li>Activity and exercise</li>
											<li>Heart rate</li>
										</ul>
										<p>
											{" "}
											However, we ask for numerous other biometrics so when new
											features come around, we will support those device types.
											Please accept those as well.
										</p>

                    {/* DEVELOPERS ONLY */}

										{/* <p>
											<i>
												DEVELOPERS ONLY:{" "}
												<b>
													<Link onClick={onCreateSimulator}>
														Click here&nbsp;
														<FontAwesomeIcon
															className="mdi"
															icon={faArrowRight}
														/>{" "}
													</Link>
												</b>{" "}
												attach a Google Fit simulator with fake data.
											</i>
										</p> */}
									</div>
								</p>
							</div>
						</section>
					)}

					{/* Google Fit */}
					{currentUser.primaryHealthTrackingDeviceType === 1 && (
						<>
							{currentUser.primaryHealthTrackingDeviceRequiresLoginAgain ? (
								<section className="hero has-background-white-ter">
									<div className="hero-body">
										<p className="title">
											<FontAwesomeIcon className="fas" icon={faHeartPulse} />
											&nbsp;Authentication Required
										</p>
										<p className="subtitle">
											Your Google Fit fitness tracker requires you to login
											again.{" "}
											<b>
												<Link onClick={onRegisterClick}>
													Click here&nbsp;
													<FontAwesomeIcon
														className="mdi"
														icon={faArrowUpRightFromSquare}
													/>{" "}
												</Link>
											</b>{" "}
											to login again and meet the requirements of Google.
										</p>
									</div>
								</section>
							) : (
								<section className="hero has-background-white-ter">
									<div className="hero-body">
										<p className="title">
											<FontAwesomeIcon className="fas" icon={faHeartPulse} />
											&nbsp;Google Fit Connected
										</p>
										<p className="subtitle">
											Your Google Fit fitness tracker is connected with us - you
											are done!
										</p>
										<p>
											<i>
												If for any reason you need to login again then{" "}
												<b>
													<Link onClick={onRegisterClick}>
														click here&nbsp;
														<FontAwesomeIcon
															className="mdi"
															icon={faArrowUpRightFromSquare}
														/>{" "}
													</Link>
												</b>{" "}
												to get redirected to Google's authentication portal.
											</i>
										</p>
									</div>
								</section>
							)}
						</>
					)}
				</div>
			</div>

			<div class="columns pt-5">
				<div class="column is-half">
					<Link class="button is-hidden-touch" to={"/dashboard"}>
						<FontAwesomeIcon className="fas" icon={faArrowLeft} />
						&nbsp;Back to Dashboard
					</Link>
					<Link class="button is-fullwidth is-hidden-desktop" to={"/dashboard"}>
						<FontAwesomeIcon className="fas" icon={faArrowLeft} />
						&nbsp;Back to Dashboard
					</Link>
				</div>
				<div class="column is-half has-text-right">
					{currentUser.primaryHealthTrackingDeviceType === 0 && (
						<>
							<Link
								class="button is-success is-hidden-touch"
								onClick={onRegisterClick}>
								Register With Google Fit&nbsp;
								<FontAwesomeIcon
									className="fas"
									icon={faArrowUpRightFromSquare}
								/>
							</Link>
							<Link
								class="button is-success is-fullwidth is-hidden-desktop"
								onClick={onRegisterClick}>
								Register My Google Fit&nbsp;
								<FontAwesomeIcon
									className="fas"
									icon={faArrowUpRightFromSquare}
								/>
							</Link>
						</>
					)}
				</div>
			</div>
		</div>
	);
}

export default AccountWearableTechLaunchpad;

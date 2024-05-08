import React, {useState, useEffect, useRef} from "react";
import Scroll from "react-scroll";
import moment from "moment";
import {useRecoilState} from "recoil";

import FormErrorBox from "../../Reusable/FormErrorBox";
import {getMySummaryAPI} from "../../../API/Biometric";
import {
	topAlertMessageState,
	topAlertStatusState,
	currentUserState,
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import Slideshow from "./Slider";

const ActivityStream = () => {
	const [topAlertMessage, setTopAlertMessage] =
		useRecoilState(topAlertMessageState);
	const [topAlertStatus, setTopAlertStatus] =
		useRecoilState(topAlertStatusState);
	const [currentUser] = useRecoilState(currentUserState);

	const [errors, setErrors] = useState({});
	const [datum, setDatum] = useState(null);
	const [isFetching, setFetching] = useState(false);

	const getDatum = (user) => {
		if (!user || user.primaryHealthTrackingDeviceType === 0) {
			console.log("user does not have a device, prevented pulling data.");
			return;
		}

		setFetching(true);
		setErrors({});
		let params = new Map();
		params.set("user_id", user.id);

		getMySummaryAPI(params, onSummarySuccess, onSummaryError, onSummaryDone);
	};

	const onSummarySuccess = (response) => {
		setDatum(response);
		console.log("onSummarySuccess:", response);
	};

	const onSummaryError = (apiErr) => {
		setErrors(apiErr);
		Scroll.animateScroll.scrollToTop();
		console.error("onSummaryError:", apiErr);
	};

	const onSummaryDone = () => {
		setFetching(false);
	};

	useEffect(() => {
		window.scrollTo(0, 0);
		getDatum(currentUser);

		return () => {};
	}, [currentUser]);

	const formatFunction = (timeframe) => (item) => {
		switch (timeframe) {
			case "hours":
				return moment(item.end).format("ha");
			case "week":
				return moment(item.end).format("ddd");
			case "month":
				return moment(item.end).format("MMM D");
			case "year":
				return moment(item.end).format("MMM YYYY");
			default:
				return moment(item.end).format("ha");
		}
	};

	const transformData = (data, label, text, timeframe, mode) => {
		if (!data) return null;

		const formatter = formatFunction(timeframe);
		const dataset = {
			label: label,
			data: data.map(
				(item) =>
					({
						1: item.count,
						2: item.count,
						3: item.average,
						4: item.sum,
					}[mode])
			),
			borderColor: "#E1BD67",
			backgroundColor: "#ffffff",
			type: mode === 2 ? "bar" : "line",
			order: 1,
		};

		return {
			text: text,
			labels: data.map(formatter),
			datasets: [dataset],
		};
	};

	return (
		<div className="card has-background-primary">
			<div className="card-content">
				<p className="title is-6 has-text-centered">Recent Trends</p>
				{isFetching ? (
					<PageLoadingContent displayMessage="Loading..." />
				) : (
					<>
						{datum ? (
							<Slideshow datum={datum} transformData={transformData} />
						) : (
							<FormErrorBox errors={errors} />
						)}
					</>
				)}
			</div>
		</div>
	);
};

export default ActivityStream;

// Datum not working

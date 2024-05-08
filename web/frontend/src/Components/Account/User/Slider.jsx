import React, {useState, useEffect} from "react";
import BarChart from "../../Reusable/Charts/Bar";

const Slideshow = ({datum, transformData}) => {
	const [activeSlide, setActiveSlide] = useState(0);

	// Define the data for each slide
	const slides = [
		{
			data: datum?.heartRateThisIsoWeekData,
			label: "Heart Rate",
			text: "Heart Rate - Week",
			timeframe: "week",
			mode: 3,
		},
		{
			data: datum?.stepCountDeltaThisIsoWeekData,
			label: "Steps Count",
			text: "Steps Count - Week",
			timeframe: "week",
			mode: 4,
		},
		{
			data: datum?.caloriesBurnedThisIsoWeekData,
			label: "Calories Burned",
			text: "Calories Burned - Week",
			timeframe: "hours",
			mode: 4,
		},
		{
			data: datum?.distanceDeltaThisIsoWeekData, // Assuming 'heartRateThisDayData' contains distance data
			label: "Distance",
			text: "Distance - Week",
			timeframe: "hours",
			mode: 4,
		},
	];

	useEffect(() => {
		// Change slide every 5 seconds
		const interval = setInterval(() => {
			setActiveSlide((currentSlide) => (currentSlide + 1) % slides.length);
		}, 5000);

		return () => clearInterval(interval);
	}, [slides.length]);

	return (
		<div>
			{slides.map((slide, index) => (
				<div
					key={index}
					style={{display: index === activeSlide ? "block" : "none"}}>
					<BarChart
						data={transformData(
							slide.data,
							slide.label,
							slide.text,
							slide.timeframe,
							slide.mode
						)}
					/>
				</div>
			))}
		</div>
	);
};

export default Slideshow;
